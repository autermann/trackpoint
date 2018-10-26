package main

import (
	"testing"
)

func TestParseFlags(t *testing.T) {
	s := NewSettings()
	v := s.Values
	var (
		daemon bool
		config string
		err    error
	)

	err = ParseFlags([]string{"trackpoint"}, s)
	if err != nil {
		t.Fatal(err)
	}
	if daemon {
		t.Fatalf("expected %v, got %v", false, daemon)
	}
	if config != "" {
		t.Fatalf("expected %v, got %v", "", config)
	}

	checkDefaults(s.Values, t)
	args := []string{
		"trackpoint",
		"--daemon",
		"--skipback",
		"--extdev",
		"--config", "./trackpoint.yml",
		"--draghys", "0",
		"--thresh", "0",
		"--upthresh", "0",
		"--ztime", "0",
		"--reach", "0",
		"--jenks", "0",
		"--drifttime", "0",
		"--speed", "0",
		"--sensitivity", "0",
		"--inertia", "0",
		"--mindrag", "0",
		"--pts", "0",
	}

	err = ParseFlags(args, s)

	if err != nil {
		t.Fatal(err)
	}
	if config != "./trackpoint.yml" {
		t.Fatal("config not read")
	}
	if !daemon {
		t.Fatal("daemon not read")
	}
	if v.DragHysteresis != 0 {
		t.Fatal("draghys not read")
	}
	if v.Threshold != 0 {
		t.Fatal("thresh not read")
	}
	if v.UpThreshold != 0 {
		t.Fatal("upthresh not read")
	}
	if v.ZTime != 0 {
		t.Fatal("ztime not read")
	}
	if v.Sensitivity != 0 {
		t.Fatal("sensitivity not read")
	}
	if v.Inertia != 0 {
		t.Fatal("inertia not read")
	}
	if v.Speed != 0 {
		t.Fatal("speed not read")
	}
	if v.Reach != 0 {
		t.Fatal("reach not read")
	}
	if v.MinDrag != 0 {
		t.Fatal("mindrag not read")
	}
	if v.Jenks != 0 {
		t.Fatal("jenks not read")
	}
	if v.DriftTime != 0 {
		t.Fatal("drifttime not read")
	}
	if v.PressToSelect != true {
		t.Fatal("pts not read")
	}
	if v.Skipback != true {
		t.Fatal("skipback not read")
	}
	if v.ExternalDevice != true {
		t.Fatal("extdev")
	}

	args = []string{"trackpoint", "--config", "./trackpoint-asdf.yml"}

	err = ParseFlags(args, s)
	if err == nil {
		t.Fatal("expected error")
	}

}
