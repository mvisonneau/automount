package fs

import (
	"github.com/mvisonneau/automount/pkg/exec"
	log "github.com/sirupsen/logrus"
)

// IsLvmAvailable tells if lvm is avaialable
func IsLvmAvailable() bool {
	return isCommandAvailable("lvm")
}

// IsMdadmAvailable tells if mdadm is avaialable
func IsMdadmAvailable() bool {
	return isCommandAvailable("mdadm")
}

// IsLsblkAvailable tells if lsblk is avaialable
func IsLsblkAvailable() bool {
	return isCommandAvailable("lsblk")
}

// IsBlkidAvailable tells if blkid is avaialable
func IsBlkidAvailable() bool {
	return isCommandAvailable("blkid")
}

func isCommandAvailable(cmd string) bool {
	c := exec.CommandInfo{
		Command: "which",
		Args: []string{
			cmd,
		},
		Result: &exec.CommandResult{},
	}

	err := c.Exec()
	if err != nil || c.Result.Status != 0 {
		log.Debugf("command %s : %v | exit code %d", cmd, err, c.Result.Status)
		return false
	}

	log.Debugf("command %s available at %s", cmd, c.Result.Stdout)
	return true
}
