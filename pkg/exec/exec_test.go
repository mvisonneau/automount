package exec

import (
	"testing"
)

func TestExecReturnCodeSuccess(t *testing.T) {
	if _, _, status, err := Exec("true", []string{}, ""); err != nil {
		t.Fatalf("Errored: %v", err)
	} else {
		if status != 0 {
			t.Fatalf("Expected exit code: 0")
		}
	}
}

func TestExecReturnCodeFailure(t *testing.T) {
	if _, _, status, err := Exec("false", []string{}, ""); err == nil {
		t.Fatalf("Expected an error")
	} else {
		if status != 1 {
			t.Fatal("Expected exit code: 1")
		}
	}
}

func TestExecStdout(t *testing.T) {
	if stdout, stderr, status, err := Exec("bash", []string{"-c", "echo foo"}, ""); err != nil {
		t.Fatalf("Error: %v", err)
	} else {
		if stderr != "" {
			t.Fatal("Expected stdout to be empty")
		}
		if stdout != "foo\n" {
			t.Fatalf("Expected 'foo' to be printed to stdout. Got: %s", stdout)
		}
		if status != 0 {
			t.Fatalf("Expected exit code: 0")
		}
	}
}

func TestExecStderr(t *testing.T) {
	if stdout, stderr, status, err := Exec("bash", []string{"-c", "echo foo 1>&2"}, ""); err != nil {
		t.Fatalf("Error: %v", err)
	} else {
		if stdout != "" {
			t.Fatal("Expected stdout to be empty")
		}
		if stderr != "foo\n" {
			t.Fatalf("Expected 'foo' to be printed to stderr. Got: %s", stderr)
		}
		if status != 0 {
			t.Fatalf("Expected exit code: 0")
		}
	}
}

func TestExecFailedToStart(t *testing.T) {
	if stdout, stderr, status, err := Exec("foo", []string{}, ""); err == nil {
		t.Fatal("Expected the 'foo' command would not run")
	} else {
		if status != 0 {
			t.Fatal("Return code should be uninitialized i.e. 0")
		}
		if stdout != "" || stderr != "" {
			t.Fatal("Both stdout and stderr should be empty")
		}
	}
}
