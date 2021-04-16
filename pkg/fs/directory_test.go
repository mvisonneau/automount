package fs

import (
	"testing"

	"github.com/mvisonneau/automount/pkg/random"
)

const (
	testFileMode = 0o755
)

func TestDirectoryCreate(t *testing.T) {
	d := &Directory{Path: generateRandomTmpPath()}
	if err := d.Create(testFileMode); err != nil {
		t.Fatalf("Errored: %v", err)
	}
	defer d.Delete()

	if exists, err := d.Exists(); err != nil {
		t.Fatalf("Errored: %v", err)
	} else {
		if !exists {
			t.Fatal("Directory should be existing but isn't")
		}
	}
}

func TestDirectoryDelete(t *testing.T) {
	d := &Directory{Path: generateRandomTmpPath()}
	if err := d.Create(testFileMode); err != nil {
		t.Fatalf("Errored: %v", err)
	}

	if err := d.Delete(); err != nil {
		t.Fatalf("Error while trying to delete : %v", err)
	}

	if exists, err := d.Exists(); err != nil {
		t.Fatalf("Errored: %v", err)
	} else {
		if exists {
			t.Fatal("Directory should not exist anymore")
		}
	}
}

func TestDirectoryEnsureExists(t *testing.T) {
	d := &Directory{Path: generateRandomTmpPath()}
	if err := d.EnsureExists(testFileMode); err != nil {
		t.Fatalf("Errored: %v", err)
	}
	defer d.Delete()

	if exists, err := d.Exists(); err != nil {
		t.Fatalf("Errored: %v", err)
	} else {
		if !exists {
			t.Fatal("Directory should be existing but isn't")
		}
	}

	if err := d.SetMode(0o755); err != nil {
		t.Fatalf("Errored: %v", err)
	}

	if err := d.EnsureExists(testFileMode); err != nil {
		t.Fatalf("Errored: %v", err)
	} else {
		if m, err := d.GetMode(); err != nil {
			t.Fatalf("Errored: %v", err)
		} else {
			if m != testFileMode {
				t.Fatalf("Directory mode should be %v, got %v", testFileMode, m)
			}
		}
	}
}

func TestDirectoryGetMode(t *testing.T) {
	d := &Directory{Path: generateRandomTmpPath()}
	if err := d.Create(testFileMode); err != nil {
		t.Fatalf("Errored: %v", err)
	}

	defer d.Delete()

	if m, err := d.GetMode(); err != nil {
		t.Fatalf("Errored: %v", err)
	} else {
		if m != testFileMode {
			t.Fatalf("Directory mode should be %v, got %v", testFileMode, m)
		}
	}
}

func TestDirectorySetMode(t *testing.T) {
	d := &Directory{Path: generateRandomTmpPath()}
	if err := d.Create(testFileMode); err != nil {
		t.Fatalf("Errored: %v", err)
	}

	defer d.Delete()

	if err := d.SetMode(0o755); err != nil {
		t.Fatalf("Errored: %v", err)
	}

	if m, err := d.GetMode(); err != nil {
		t.Fatalf("Errored: %v", err)
	} else {
		if m != 0o755 {
			t.Fatalf("Directory mode should be %v, got %v", 0o755, m)
		}
	}
}

func generateRandomTmpPath() string {
	r, err := random.GenerateString(8)
	if err != nil {
		panic(err)
	}
	return "/tmp/" + r
}
