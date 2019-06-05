package lvm

import (
	"testing"
)

func TestLVMPath(t *testing.T) {
	lv := &LogicalVolume{
		Name:        "bar",
		VolumeGroup: &VolumeGroup{Name: "foo"},
	}

	if lv.Path() != "/dev/foo/bar" {
		t.Fatalf("Expected to get '/dev/foo/bar', got %s", lv.Path())
	}
}

func TestRemoveQuotes(t *testing.T) {
	if s := removeQuotes("f'o'o"); s != "foo" {
		t.Fatalf("Expected to get 'foo', got %s", s)
	}
}
