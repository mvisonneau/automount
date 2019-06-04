package fs

import (
	"fmt"
	"os"
)

// Directory is being used to handle information
type Directory struct {
	Path string
}

// EnsureExists is going to create a folder if it does not already exist
func (d *Directory) EnsureExists(mode os.FileMode) error {
	exists, err := d.Exists()
	if err != nil {
		return err
	}

	if !exists {
		return d.Create(mode)
	}

	dm, err := d.GetMode()
	if err != nil {
		return err
	}

	if dm != mode {
		return d.SetMode(mode)
	}

	return nil
}

// Create actually creates the directory on the filesystem
func (d *Directory) Create(mode os.FileMode) error {
	return os.MkdirAll(d.Path, mode)
}

// Delete a directory from the filesystem
func (d *Directory) Delete() error {
	return os.RemoveAll(d.Path)
}

// Exists returns if a directory exists or not
func (d *Directory) Exists() (bool, error) {
	fi, err := os.Lstat(d.Path)

	if os.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	if !fi.IsDir() {
		return false, fmt.Errorf("%s exists but is not a directory", d.Path)
	}

	return true, nil
}

// GetMode returns the filemode of the directory
func (d *Directory) GetMode() (os.FileMode, error) {
	file, err := os.Open(d.Path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	f, err := file.Stat()
	if err != nil {
		return 0, err
	}

	return f.Mode().Perm(), nil
}

// SetMode configures the mode of the directory
func (d *Directory) SetMode(mode os.FileMode) error {
	return os.Chmod(d.Path, mode)
}
