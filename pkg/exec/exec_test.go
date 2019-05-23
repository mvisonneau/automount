package exec

import (
	"testing"
)

func TestExecReturnCodeSuccess(t *testing.T) {
	c := &CommandInfo{
		Command: "true",
		Args:    []string{},
		Dir:     "/",
	}

	if err := c.Exec(); err != nil {
		t.Fatalf("Errored: %v", err)
	} else {
		if c.Result.Status != 0 {
			t.Fatalf("Expected exit code: 0")
		}
	}
}

func TestExecReturnCodeFailure(t *testing.T) {
	c := &CommandInfo{
		Command: "false",
		Args:    []string{},
	}

	if err := c.Exec(); err == nil {
		t.Fatalf("Expected an error")
	} else {
		if c.Result.Status != 1 {
			t.Fatal("Expected exit code: 1")
		}
	}
}

func TestExecStdout(t *testing.T) {
	c := &CommandInfo{
		Command: "bash",
		Args: []string{
			"-c",
			"echo foo",
		},
	}

	if err := c.Exec(); err != nil {
		t.Fatalf("Error: %v", err)
	} else {
		if c.Result.Stderr != "" {
			t.Fatalf("Expected stderr to be empty, got %s", c.Result.Stderr)
		}
		if c.Result.Stdout != "foo\n" {
			t.Fatalf("Expected 'foo' to be printed to stdout, got: %s", c.Result.Stdout)
		}
		if c.Result.Status != 0 {
			t.Fatalf("Expected exit code 0, got %d", c.Result.Status)
		}
	}
}

func TestExecStderr(t *testing.T) {
	c := &CommandInfo{
		Command: "bash",
		Args: []string{
			"-c",
			"echo foo 1>&2",
		},
	}

	if err := c.Exec(); err != nil {
		t.Fatalf("Error: %v", err)
	} else {
		if c.Result.Stdout != "" {
			t.Fatalf("Expected stdout to be empty, got %s", c.Result.Stdout)
		}
		if c.Result.Stderr != "foo\n" {
			t.Fatalf("Expected 'foo' to be printed to stderr, got: %s", c.Result.Stderr)
		}
		if c.Result.Status != 0 {
			t.Fatalf("Expected exit code 0, got %d", c.Result.Status)
		}
	}
}

func TestExecFailedToStart(t *testing.T) {
	c := &CommandInfo{
		Command: "foo",
		Args:    []string{},
	}

	if err := c.Exec(); err == nil {
		t.Fatal("Expected the 'foo' command would not run")
	} else {
		if c.Result.Status != 0 {
			t.Fatal("Return code should be uninitialized i.e. 0")
		}
		if c.Result.Stdout != "" || c.Result.Stderr != "" {
			t.Fatalf("Both stdout and stderr should be empty, got stdout:'%s', stderr:'%s'", c.Result.Stdout, c.Result.Stderr)
		}
	}
}
