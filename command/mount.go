package command

import (
	"fmt"
	"math"
	"os"

	"github.com/jaypipes/ghw"
	"github.com/mvisonneau/automount/pkg/fs"
	"github.com/mvisonneau/automount/pkg/lvm"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	label = "automount"
)

// Mount actually formats, then mount a FS
func Mount(ctx *cli.Context) error {
	if err := configure(ctx); err != nil {
		return cli.NewExitError(err, 1)
	}

	fsType := ctx.GlobalString("fstype")
	mode := os.FileMode(ctx.GlobalInt("mountpoint-mode"))

	var devices *fs.Devices

	mountPoint := ctx.Args().First()
	if mountPoint == "" {
		return exit(fmt.Errorf("You must provide the mountpoint path"), 1)
	}

	log.Info("Parsing /etc/fstab")
	mounts, err := fs.GetMounts()
	if err != nil {
		return err
	}
	log.Infof("Found %d entries in /etc/fstab", len(*mounts))

	if len(ctx.GlobalStringSlice("device")) > 0 {
		log.Infof("Attempting to mount '%v' block device(s) at '%s'", ctx.GlobalStringSlice("device"), mountPoint)
		for _, d := range ctx.GlobalStringSlice("device") {
			device := &fs.Device{Path: d}
			if exists, err := device.Exists(); err != nil || !exists {
				return exit(fmt.Errorf("%s does not exist or is not a block device", device.Path), 1)
			}
			*devices = append(*devices, device)
		}
	} else {
		log.Infof("No device specified, looking up available ones")
		devices, err = findAvailableDevices(mounts, ctx.GlobalBool("use-formatted-devices"))
		if err != nil {
			exit(fmt.Errorf("Error whilst looking for available devices %v", err), 1)
		}
	}

	if len(*devices) == 0 {
		log.Warnf("No available device found, exiting..")
		return exit(nil, 0)
	}

	var device *fs.Device
	if ctx.GlobalBool("use-lvm") {
		device, err = createLogicalVolume(devices, ctx.GlobalBool("use-all-devices"))
		if err != nil {
			return exit(err, 1)
		}
	} else {
		device = devices.First()
		log.Infof("Using block device : '%s'", device.Path)
	}

	isDeviceFormatted := false
	if deviceFsType, _ := device.GetFSType(); len(deviceFsType) == 0 {
		log.Infof("%s is not formatted, will format it.", device.Path)
	} else {
		if deviceFsType == fsType {
			log.Infof("%s is formatted to '%s' as expected, continuing..", device.Path, fsType)
			isDeviceFormatted = true
		} else if ctx.GlobalBool("use-formatted-devices") {
			log.Infof("%s is not configured but already formatted in %v, will reformat in %v", device.Path, deviceFsType, fsType)
		} else {
			return exit(fmt.Errorf("Cannot mount device '%s' (%s) as %s", device.Path, deviceFsType, fsType), 1)
		}
	}

	if !isDeviceFormatted {
		log.Infof("Formatting device %s to %s", device.Path, fsType)
		if err := device.CreateFS(fsType, label); err != nil {
			return exit(err, 1)
		}
	}

	// Create the mount point directory and ensure permissions
	log.Infof("Ensuring that mountpoint %s exists with correct permissions (%d)", mountPoint, mode)
	directory := fs.Directory{Path: mountPoint}
	directory.EnsureExists(mode)

	// Check if it is already configured in /etc/fstab otherwise amend the config
	mount := &fs.Mount{
		Spec:    device.Path,
		File:    directory.Path,
		VfsType: fsType,
		MntOps:  map[string]string{"defaults": ""},
		Freq:    0,
		PassNo:  0,
	}

	if !mounts.Exists(mount.File) {
		log.Infof("%s is not configured within fstab, appending configuration", device.Path)
		mounts.Add(mount)

		log.Info("Writing configuration to /etc/fstab")
		if err := mounts.WriteFstab(); err != nil {
			return exit(err, 1)
		}
	}

	// Mount it
	log.Infof("Attempting to mount %s to %s", device.Path, mountPoint)
	if err := mount.Mount(); err != nil {
		return exit(err, 1)
	}
	log.Infof("Mounted!")

	return exit(nil, 0)
}

func findAvailableDevices(mounts *fs.Mounts, reuseFormattedAndUnmounted bool) (*fs.Devices, error) {
	devices := &fs.Devices{}
	block, err := ghw.Block()
	if err != nil {
		return nil, fmt.Errorf("Error getting block storage info: %v", err)
	}

	log.Infof("Found %v disk(s), total size of %v GB", len(block.Disks), math.Ceil(float64(block.TotalPhysicalBytes/1024/1024/1024)))
	for _, disk := range block.Disks {
		device := &fs.Device{Path: fmt.Sprintf("/dev/%v", disk.Name)}
		if mounts.Exists(device.Path) {
			log.Infof("%s is already mounted, skipping..", device.Path)
			continue
		}

		if len(disk.Partitions) > 0 {
			log.Infof("/dev/%v has partitions, skipping", disk.Name)
			continue
		}

		if deviceFsType, _ := device.GetFSType(); len(deviceFsType) > 0 {
			if !reuseFormattedAndUnmounted {
				log.Infof("%s is formatted (%v), skipping..", device.Path, deviceFsType)
				continue
			}
		}

		log.Infof("%s is available", device.Path)
		*devices = append(*devices, device)
	}

	return devices, nil
}

func createLogicalVolume(devices *fs.Devices, useAllDevices bool) (*fs.Device, error) {
	log.Infof("Using LVM for managing the partitions")

	if !fs.IsLvmAvailable() {
		return nil, fmt.Errorf("LVM is not available on the OS")
	}

	log.Debugf("LVM: getting current state")
	l, err := lvm.New()
	if err != nil {
		return nil, err
	}

	pvs := lvm.PhysicalVolumes{}
	for _, device := range *devices {
		log.Debugf("LVM: creating physical volume on %s", device.Path)

		pv, err := l.CreatePhysicalVolume(device.Path)
		if err != nil {
			return nil, fmt.Errorf("LVM: Error creating physical volume on %s : %v", device.Path, err)
		}
		pvs = append(pvs, pv)

		if !useAllDevices {
			break
		}
	}

	log.Debugf("LVM: creating volume group")
	vg, err := l.CreateVolumeGroup("automount", pvs, []string{"automount"})
	if err != nil {
		return nil, fmt.Errorf("LVM: Error creating volume group : %v", err)
	}

	log.Debugf("LVM: creating logical volume")
	lv, err := l.CreateLogicalVolume("automount", vg, 0, []string{"automount"})
	if err != nil {
		return nil, fmt.Errorf("LVM: Error creating logical volume %v", err)
	}

	log.Infof("physical volume, volume group and logical volume created, using this as a device")
	return &fs.Device{Path: lv.Path()}, nil
}
