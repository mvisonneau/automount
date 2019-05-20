package exec

import (
	"bytes"
	"fmt"
	"os/exec"
	"syscall"
)

// Exec executes a given command with specific arguments
func Exec(command string, args []string, dir string) (stdout, stderr string, status int, err error) {
	stdout = ""
	stderr = ""
	status = 0

	cmd := exec.Cmd{}

	cmd.Path, err = exec.LookPath(command)
	if err != nil {
		return
	}

	cmd.Args = append([]string{cmd.Path}, args...)

	if dir != "" {
		cmd.Dir = dir
	}

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	if err = cmd.Start(); err != nil {
		return
	}

	if cmdErr := cmd.Wait(); cmdErr != nil {
		if exitErr, ok := cmdErr.(*exec.ExitError); ok {
			if cmdStatus, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				status = cmdStatus.ExitStatus()
				err = fmt.Errorf("exit code: %d", cmdStatus.ExitStatus())
			}
		} else {
			err = fmt.Errorf("error: %v", err)
		}
	}

	stdout = outBuf.String()
	stderr = errBuf.String()
	return
}
