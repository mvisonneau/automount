package fs

import (
	"fmt"
	"os"
)

// Device is being used to handle device information
type Directory struct {
	Path string
}

func (d *Directory) EnsureExists(mode os.FileMode) error {
	if exists, err := d.Exists(); err != nil {
		return err
	} else {
		if !exists {
			return d.Create(mode)
		}

		if dm, err := d.GetMode(); err != nil {
			return err
		} else {
			if dm != mode {
				return d.SetMode(mode)
			}
		}
	}

	return nil
}

func (d *Directory) Create(mode os.FileMode) error {
	return os.MkdirAll(d.Path, mode)
}

func (d *Directory) Delete() error {
	return os.RemoveAll(d.Path)
}

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

func (d *Directory) SetMode(mode os.FileMode) error {
	return os.Chmod(d.Path, mode)
}
