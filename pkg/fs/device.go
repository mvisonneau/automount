package fs

import (
	"fmt"
	"os"
	"strings"

	"github.com/mvisonneau/automount/pkg/exec"
)

const (
	blkidFSType  = "TYPE"
	blkidFSLabel = "LABEL"
)

// Device is being used to handle device information
type Device struct {
	Path string
}

// GetFSType returns the type of the filesystem of a given device
func (d *Device) GetFSType() (string, error) {
	return d.getFSInfo(blkidFSType)
}

// GetFSLabel returns the label of the filesystem of a given device
func (d *Device) GetFSLabel() (string, error) {
	return d.getFSInfo(blkidFSLabel)
}

// Exists validates that the path configured exists and points to a block device
func (d *Device) Exists() (bool, error) {
	fi, err := os.Lstat(d.Path)

	if os.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	if fi.Mode() != os.ModeDevice {
		return false, fmt.Errorf("%s exists but is not a device", d.Path)
	}

	return true, nil
}

// CreateFS creates a filesystem on a given device
func (d *Device) CreateFS(fsType, label string) error {
	c := exec.CommandInfo{
		Command: "mkfs." + fsType,
		Args: []string{
			d.Path,
			"-L",
			label,
		},
	}

	if err := c.Exec(); err != nil {
		return err
	}

	return nil
}

func (d *Device) getFSInfo(kind string) (string, error) {
	c := exec.CommandInfo{
		Command: "blkid",
		Args: []string{
			"-o",
			"value",
			"-s",
			kind,
			d.Path,
		},
		Result: &exec.CommandResult{},
	}

	err := c.Exec()
	if err != nil && c.Result.Status != 2 {
		return "", err
	}

	return strings.Trim(c.Result.Stdout, "\n"), nil
}
