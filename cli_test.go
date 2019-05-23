package main

import (
	"testing"
)

func TestRunCli(t *testing.T) {
	c := runCli()
	if c.Name != "automount" {
		t.Fatalf("Expected c.Name to be automount, got '%v'", c.Name)
	}
}
