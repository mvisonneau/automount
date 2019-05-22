package fs

import (
	"fmt"
	"testing"

	"github.com/mvisonneau/automount/pkg/exec"
	"github.com/mvisonneau/automount/pkg/random"
)

const (
	testFSType         = "ext2"
	testFSLabel        = "foo"
	testFSSize  uint32 = 8192
)

func TestDeviceGetFSType(t *testing.T) {
	var d = &Device{Path: generateRandomRAMDevicePath()}

	if err := createRAMDisk(d.Path, testFSType, testFSLabel, testFSSize); err != nil {
		t.Fatalf("Errored: %v", err)
	}
	defer deleteRAMDisk(d.Path)

	if fsType, err := d.GetFSType(); err != nil {
		t.Fatalf("Errored: %v", err)
	} else {
		if fsType != testFSType {
			t.Fatalf("Expected fsType to be '%s', got '%s'", testFSType, fsType)
		}
	}
}

func TestDeviceGetFSLabel(t *testing.T) {
	var d = &Device{Path: generateRandomRAMDevicePath()}
	if err := createRAMDisk(d.Path, testFSType, testFSLabel, testFSSize); err != nil {
		t.Fatalf("Errored: %v", err)
	}
	defer deleteRAMDisk(d.Path)

	if fsLabel, err := d.GetFSLabel(); err != nil {
		t.Fatalf("Errored: %v", err)
	} else {
		if fsLabel != testFSLabel {
			t.Fatalf("Expected fsLabel to be '%s', got '%s'", testFSLabel, fsLabel)
		}
	}
}

func TestDeviceCreateFS(t *testing.T) {
	var d = &Device{Path: generateRandomRAMDevicePath()}
	if err := createRAMDisk(d.Path, testFSType, testFSLabel, testFSSize); err != nil {
		t.Fatalf("Errored: %v", err)
	}
	defer deleteRAMDisk(d.Path)

	// Reformat to ext4
	d.CreateFS("ext4", "bar")

	if fsType, err := d.GetFSType(); err != nil {
		t.Fatalf("Errored: %v", err)
	} else {
		if fsType != "ext4" {
			t.Fatalf("Expected fsType to be 'ext4', got '%s'", fsType)
		}
	}

	if fsLabel, err := d.GetFSLabel(); err != nil {
		t.Fatalf("Errored: %v", err)
	} else {
		if fsLabel != "bar" {
			t.Fatalf("Expected fsLabel to be 'bar', got '%s'", fsLabel)
		}
	}
}

// createRamDisk creates a block device in RAM
func createRAMDisk(path, fsType, fsLabel string, size uint32) error {
	c := exec.CommandInfo{
		Command: "mkfs",
		Args: []string{
			"-q",
			"-t",
			fsType,
			"-L",
			fsLabel,
			path,
			fmt.Sprint(size),
		},
	}

	if err := c.Exec(); err != nil {
		return err
	}

	return nil
}

// deleteRamDisk delete a ram block device
func deleteRAMDisk(path string) error {
	c := exec.CommandInfo{
		Command: "rm",
		Args: []string{
			"-f",
			path,
		},
	}

	if err := c.Exec(); err != nil {
		return err
	}

	return nil
}

func generateRandomRAMDevicePath() string {
	r, err := random.GenerateString(2)
	if err != nil {
		panic(err)
	}
	return "/dev/ram" + r
}
