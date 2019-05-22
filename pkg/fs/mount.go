package fs

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	fstab "github.com/deniswernert/go-fstab"
)

// Mount is only used for aliasing purposes
type Mount fstab.Mount

// Mounts is used for searching capabilities over all mounts
type Mounts map[string][]*fstab.Mount

// Mount makes a syscall to mount a mountpoint which must be present in the fstab
func (m *Mount) Mount() error {
	return syscall.Mount(m.Spec, m.File, m.VfsType, 0, "")
}

// Unmount makes a syscall to umount a specific path
func (m *Mount) Unmount() error {
	return syscall.Unmount(m.File, 0)
}

// IsMounted allow you to check if a given path has a specific mount
func IsMounted(mountPoint string) (bool, error) {
	mntF, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return false, err
	}

	scanner := bufio.NewScanner(mntF)

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 5 {
			if mountPoint == fields[4] {
				return true, nil
			}
		}
	}

	return false, nil
}

// GetMounts returns fstab entries in a serialized format
func GetMounts() (*Mounts, error) {
	mounts := make(Mounts)
	foundMounts, err := fstab.ParseSystem()
	if err != nil {
		return &mounts, err
	}

	for _, m := range foundMounts {
		mounts[m.Spec] = append(mounts[m.Spec], m)
	}

	return &mounts, nil
}

// Exists validates that at least 1 mount is present in the fstab for a specific device
func (mounts *Mounts) Exists(device string) bool {
	if _, ok := (*mounts)[device]; ok {
		return true
	}

	return false
}

// Get returns a seralized mount information
func (mounts *Mounts) Get(device string) ([]*fstab.Mount, error) {
	if !mounts.Exists(device) {
		return nil, fmt.Errorf("Cannot find fstab entries for '%v' device", device)
	}

	return (*mounts)[device], nil
}

// Add appends a Mount to a Mounts object
func (mounts *Mounts) Add(mount *Mount) {
	(*mounts)[mount.Spec] = append((*mounts)[mount.Spec], (*fstab.Mount)(mount))
}

// Remove deletes a Mount from a Mounts object
func (mounts *Mounts) Remove(mount *Mount) {
	if mounts.Exists(mount.Spec) {
		delete(*mounts, mount.Spec)
	}
}

// WriteFstab replace the content of the fstab with Mounts values
func (mounts *Mounts) WriteFstab() error {
	mnts := fstab.Mounts{}
	for _, m := range *mounts {
		for _, e := range m {
			mnts = append(mnts, e)
		}
	}

	// Open /etc/fstab file and truncate before appending
	f, err := os.OpenFile("/etc/fstab", os.O_TRUNC|os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.WriteString(mnts.String()); err != nil {
		return err
	}

	return nil
}
