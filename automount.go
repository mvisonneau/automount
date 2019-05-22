package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/jaypipes/ghw"
	"github.com/mvisonneau/automount/pkg/fs"
)

var (
	fsType                  = "ext4"
	mountPoint              = "/mnt/test"
	permissions os.FileMode = 0755
)

const (
	label = "formatted_by_automount"
)

func main() {
	if user, err := user.Current(); err != nil {
		panic(fmt.Errorf("Unable to determine current user"))
	} else if user.Uid != "0" {
		panic(fmt.Errorf("You have to run this function as root"))
	}

	block, err := ghw.Block()

	if err != nil {
		fmt.Printf("Error getting block storage info: %v", err)
	}

	fmt.Printf("Found %v disk(s), total size of %v bytes\n", len(block.Disks), block.TotalPhysicalBytes)

	// Parse current fstab
	mounts, err := fs.GetMounts()
	if err != nil {
		panic(err)
	}

	for _, disk := range block.Disks {
		if len(disk.Partitions) > 0 {
			fmt.Printf("/dev/%v has partitions, skipping\n", disk.Name)
			continue
		}

		d := fs.Device{Path: fmt.Sprintf("/dev/%v", disk.Name)}
		m := fs.Directory{Path: mountPoint}
		deviceFsType, _ := d.GetFSType()
		if len(deviceFsType) > 0 {
			fmt.Printf("%v is formatted (%v), skipping\n", d.Path, deviceFsType)
			continue
		}

		fmt.Printf("%v is available, formatting to %s and mounting to %s\n", d.Path, fsType, mountPoint)
		if err := d.CreateFS(fsType, label); err != nil {
			panic(err)
		}

		// Create the mount point directory and ensure permissions
		m.EnsureExists(permissions)

		// Check if it is already configured in /etc/fstab otherwise amend the config
		mount := &fs.Mount{
			Spec:    d.Path,
			File:    m.Path,
			VfsType: fsType,
			MntOps:  map[string]string{"defaults": ""},
			Freq:    0,
			PassNo:  0,
		}

		if !mounts.Exists(mount.File) {
			mounts.Add(mount)

			if err := mounts.WriteFstab(); err != nil {
				panic(err)
			}
		}

		// Mount it
		if err := mount.Mount(); err != nil {
			panic(err)
		}
	}
}
