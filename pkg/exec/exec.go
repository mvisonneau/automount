package exec

import (
	"bytes"
	"fmt"
	"os/exec"
	"syscall"
)

// CommandInfo is used to configure a command and it's arguments
type CommandInfo struct {
	Command, Dir string
	Args         []string
	Result       *CommandResult
}

// CommandResult handles information about the command execution
type CommandResult struct {
	Stdout, Stderr string
	Status         int
}

// Exec executes a given command with specific arguments
func (c *CommandInfo) Exec() (err error) {
	// Initialize output
	c.Result = &CommandResult{}

	// Instanciate exec handler
	cmd := exec.Cmd{}

	cmd.Path, err = exec.LookPath(c.Command)
	if err != nil {
		return
	}

	cmd.Args = append([]string{cmd.Path}, c.Args...)

	if c.Dir != "" {
		cmd.Dir = c.Dir
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
				c.Result.Status = cmdStatus.ExitStatus()
				err = fmt.Errorf("exit code: %d", cmdStatus.ExitStatus())
			}
		} else {
			err = fmt.Errorf("error: %v", err)
		}
	}

	c.Result.Stdout = outBuf.String()
	c.Result.Stderr = errBuf.String()
	return
}
