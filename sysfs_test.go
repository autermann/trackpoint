package main

import (
	"testing"
)

func TestGetDeviceDirectory(t *testing.T) {
	path, err := GetDeviceDirectory()
	if err != nil {
		t.Fatal(err)
	}
	expected := "/sys/devices/platform/i8042/serio1/serio2"
	if path != expected {
		t.Fatalf("expected %v, got %v", expected, path)
	}
}
