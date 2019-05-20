package fs

import (
	"github.com/mvisonneau/automount/pkg/exec"
)

// GetFSType
func GetFSType(device string) (string, error) {
	cmd := "blkid"
	args := []string{
		"-o",
		"value",
		"-s",
		"TYPE",
		device,
	}

	stdout, _, status, err := exec.Exec(cmd, args, "")
	if err != nil && status != 2 {
		return "", err
	}

	return stdout, nil
}

// CreateFS
func CreateFS(device, type, label string) error {
	cmd := "mkfs." + type
	args := []string{
		device,
		"-L",
		label,
	}

	if _, err := exec.Exec(cmd, args, ""); err != nil {
		return err
	}

	return nil
}
