package fs

import (
	"fmt"
	"os"
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
	if os.Getenv("SKIP_PRIVILEGED") == "true" {
		t.Skip("skipping testing in non-privileged environment")
	}

	d := &Device{Path: generateRandomRAMDevicePath()}

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
	if os.Getenv("SKIP_PRIVILEGED") == "true" {
		t.Skip("skipping testing in non-privileged environment")
	}

	d := &Device{Path: generateRandomRAMDevicePath()}
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
	if os.Getenv("SKIP_PRIVILEGED") == "true" {
		t.Skip("skipping testing in non-privileged environment")
	}

	d := &Device{Path: generateRandomRAMDevicePath()}
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

func TestDeviceCreateFSError(t *testing.T) {
	d := &Device{Path: generateRandomRAMDevicePath()}
	if err := d.CreateFS("foo", "bar"); err == nil {
		t.Fatal("Expected to get an error")
	}
}

func TestDeviceGetFSInfoError(t *testing.T) {
	d := &Device{Path: generateRandomRAMDevicePath()}
	if _, err := d.getFSInfo("foo"); err == nil {
		t.Fatal("Expected to get an error")
	}
}

func TestDeviceExists(t *testing.T) {
	if os.Getenv("SKIP_PRIVILEGED") == "true" {
		t.Skip("skipping testing in non-privileged environment")
	}

	d := &Device{Path: generateRandomRAMDevicePath()}
	if exists, err := d.Exists(); err != nil || exists {
		t.Fatalf("Expected device '%s' to not exist", d.Path)
	}

	if err := createRAMDisk(d.Path, testFSType, testFSLabel, testFSSize); err != nil {
		t.Fatalf("Errored: %v", err)
	}
	defer deleteRAMDisk(d.Path)

	if exists, err := d.Exists(); err != nil || !exists {
		t.Fatalf("Expected device '%s' to exist : %v", d.Path, err)
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

	return c.Exec()
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

	return c.Exec()
}

func generateRandomRAMDevicePath() string {
	r, err := random.GenerateString(2)
	if err != nil {
		panic(err)
	}
	return "/dev/ram" + r
}
