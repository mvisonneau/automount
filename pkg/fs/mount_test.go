package fs

import (
	"fmt"
	"os"
	"testing"

	"github.com/mvisonneau/automount/pkg/random"
)

func TestIsMounted(t *testing.T) {
	if mounted, err := IsMounted("/"); err != nil {
		t.Fatalf("Errored: %v", err)
	} else {
		if !mounted {
			t.Fatalf("Expect / to be mounted")
		}
	}

	s, _ := random.GenerateString(8)
	randomPath := fmt.Sprintf("/mnt/%s", s)
	if mounted, err := IsMounted(randomPath); err != nil {
		t.Fatalf("Errored: %v", err)
	} else {
		if mounted {
			t.Fatalf("Expect %v to not be mounted", randomPath)
		}
	}
}

func TestMount(t *testing.T) {
	// Unfortunately drone.io does not support updating worker capabilities..
	if os.Getenv("DRONE") == "true" {
		return
	}

	s, _ := random.GenerateString(8)
	randomPath := fmt.Sprintf("/mnt/%s", s)
	if mounts, err := GetMounts(); err != nil {
		t.Fatalf("Errored: %v", err)
	} else {
		m := &Mount{
			Spec:    "tmpfs",
			File:    randomPath,
			VfsType: "tmpfs",
			MntOps:  map[string]string{"defaults": ""},
			Freq:    0,
			PassNo:  0,
		}

		mounts.Add(m)

		if err := mounts.WriteFstab(); err != nil {
			t.Fatalf("Errored: %v", err)
		}

		d := &Directory{Path: randomPath}
		if err := d.EnsureExists(0777); err != nil {
			t.Fatalf("Errored: %v", err)
		}
		defer d.Delete()

		if err := m.Mount(); err != nil {
			t.Fatalf("Errored: %v", err)
		}
		defer m.Unmount()

		if mounted, err := IsMounted(randomPath); err != nil {
			t.Fatalf("Errored: %v", err)
		} else {
			if !mounted {
				t.Fatalf("Expected %s to be mounted", randomPath)
			}
		}

		// Cleanup
		mounts.Remove(m)

		if err := mounts.WriteFstab(); err != nil {
			t.Fatalf("Errored: %v", err)
		}
	}
}
