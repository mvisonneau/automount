package fs

import (
	"testing"
)

func TestIsCommandAvailable(t *testing.T) {
	if isCommandAvailable("foo") {
		t.Fatalf("Command foo should not be available")
	}
}
