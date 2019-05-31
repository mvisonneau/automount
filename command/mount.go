package command

import (
	"fmt"
	"math"
	"os"

	"github.com/jaypipes/ghw"
	"github.com/mvisonneau/automount/pkg/fs"
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

	var device *fs.Device
	isDeviceFormatted := false
	fsType := ctx.GlobalString("fstype")

	mountPoint := ctx.Args().First()
	if mountPoint == "" {
		return exit(fmt.Errorf("You must provide the mountpoint path"), 1)
	}

	mode := os.FileMode(ctx.GlobalInt("mountpoint-mode"))

	// Parse current fstab
	log.Info("Parsing /etc/fstab")
	mounts, err := fs.GetMounts()
	if err != nil {
		return exit(err, 1)
	}
	log.Infof("Found %d entries in /etc/fstab", len(*mounts))

	if ctx.GlobalString("device") != "auto" {
		log.Infof("Attempting to mount '%s' block device at '%s'", ctx.GlobalString("device"), mountPoint)
		device = &fs.Device{Path: ctx.GlobalString("device")}
		if exists, err := device.Exists(); err != nil || !exists {
			return exit(fmt.Errorf("%s does not exist or is not a block device", device.Path), 1)
		}
		if deviceFsType, _ := device.GetFSType(); len(deviceFsType) == 0 {
			log.Infof("%s is not formatted, will format it.", device.Path)
		} else {
			if deviceFsType == fsType {
				log.Infof("%s is formatted to '%s' as expected, continuing..", device.Path, fsType)
				isDeviceFormatted = true
			} else {
				return exit(fmt.Errorf("Cannot mount device '%s' (%s) as %s", device.Path, deviceFsType, fsType), 1)
			}
		}
	} else {
		log.Infof("No device specified, trying to find one automatically")

		block, err := ghw.Block()
		if err != nil {
			return exit(fmt.Errorf("Error getting block storage info: %v", err), 1)
		}

		log.Infof("Found %v disk(s), total size of %v GB", len(block.Disks), math.Ceil(float64(block.TotalPhysicalBytes/1024/1024/1024)))

		for _, disk := range block.Disks {
			if len(disk.Partitions) > 0 {
				log.Infof("/dev/%v has partitions, skipping", disk.Name)
				continue
			}

			foundDevice := &fs.Device{Path: fmt.Sprintf("/dev/%v", disk.Name)}
			if mounts.Exists(foundDevice.Path) {
				log.Infof("%s is already mounted, skipping..", foundDevice.Path)
				continue
			}

			if deviceFsType, _ := foundDevice.GetFSType(); len(deviceFsType) > 0 {
				if !ctx.GlobalBool("reuse-formatted-devices") {
					log.Infof("%s is formatted (%v), skipping..", foundDevice.Path, deviceFsType)
					continue
				}

				if deviceFsType != fsType {
					log.Infof("%s is not configured but already formatted in %v, will reformat in %v", deviceFsType, fsType)
				} else {
					isDeviceFormatted = true
				}
			}

			device = foundDevice
			log.Infof("%s is available, picking this one!", device.Path)
			break
		}
	}

	if device == nil {
		log.Warnf("No available device found, exiting")
		return exit(nil, 0)
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
