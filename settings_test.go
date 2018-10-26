package trackpoint

import (
	"errors"
	"testing"
)

var defaults = map[string]string{
	"draghys":         "255",
	"thresh":          "8",
	"upthresh":        "255",
	"ztime":           "38",
	"sensitivity":     "128",
	"inertia":         "6",
	"speed":           "97",
	"reach":           "10",
	"mindrag":         "20",
	"jenks":           "135",
	"drift_time":      "5",
	"press_to_select": "0",
	"skipback":        "0",
	"ext_dev":         "0",
}

var keys = []string{
	"draghys",
	"thresh",
	"upthresh",
	"ztime",
	"sensitivity",
	"inertia",
	"speed",
	"reach",
	"mindrag",
	"jenks",
	"drift_time",
	"press_to_select",
	"skipback",
	"ext_dev",
}

func checkDefaults(s *Values, t *testing.T) {
	if s.DragHysteresis != DefaultDragHysteresis {
		t.Fatalf("Expected %v, got %v", DefaultDragHysteresis, s.DragHysteresis)
	}
	if s.Threshold != DefaultThreshold {
		t.Fatalf("Expected %v, got %v", DefaultThreshold, s.Threshold)
	}
	if s.UpThreshold != DefaultUpThreshold {
		t.Fatalf("Expected %v, got %v", DefaultUpThreshold, s.UpThreshold)
	}
	if s.ZTime != DefaultZTime {
		t.Fatalf("Expected %v, got %v", DefaultZTime, s.ZTime)
	}
	if s.Sensitivity != DefaultSensitivity {
		t.Fatalf("Expected %v, got %v", DefaultSensitivity, s.Sensitivity)
	}
	if s.Inertia != DefaultInertia {
		t.Fatalf("Expected %v, got %v", DefaultInertia, s.Inertia)
	}
	if s.Speed != DefaultSpeed {
		t.Fatalf("Expected %v, got %v", DefaultSpeed, s.Speed)
	}
	if s.Reach != DefaultReach {
		t.Fatalf("Expected %v, got %v", DefaultReach, s.Reach)
	}
	if s.MinDrag != DefaultMinDrag {
		t.Fatalf("Expected %v, got %v", DefaultMinDrag, s.MinDrag)
	}
	if s.Jenks != DefaultJenks {
		t.Fatalf("Expected %v, got %v", DefaultJenks, s.Jenks)
	}
	if s.DriftTime != DefaultDriftTime {
		t.Fatalf("Expected %v, got %v", DefaultDriftTime, s.DriftTime)
	}
	if s.PressToSelect != DefaultPressToSelect {
		t.Fatalf("Expected %v, got %v", DefaultPressToSelect, s.PressToSelect)
	}
	if s.Skipback != DefaultSkipback {
		t.Fatalf("Expected %v, got %v", DefaultSkipback, s.Skipback)
	}
	if s.ExternalDevice != DefaultExternalDevice {
		t.Fatalf("Expected %v, got %v", DefaultExternalDevice, s.ExternalDevice)
	}
}

func TestSettingsNew(t *testing.T) {
	s := NewSettings()
	checkDefaults(s.Values, t)
}

func TestSettings_SetDefaults(t *testing.T) {
	s2 := NewSettings()
	s := s2.Values
	s.DragHysteresis = 0
	s.Threshold = 0
	s.UpThreshold = 0
	s.ZTime = 0
	s.Sensitivity = 0
	s.Inertia = 0
	s.Speed = 0
	s.Reach = 0
	s.MinDrag = 0
	s.Jenks = 0
	s.DriftTime = 0
	s.PressToSelect = true
	s.Skipback = true
	s.ExternalDevice = true
	s.SetDefaults()
	checkDefaults(s, t)
}

func TestSettings_ReadYAML(t *testing.T) {
	s2 := NewSettings()
	s := s2.Values
	s.DragHysteresis = 0
	s.Threshold = 0
	s.UpThreshold = 0
	s.ZTime = 0
	s.Sensitivity = 0
	s.Inertia = 0
	s.Speed = 0
	s.Reach = 0
	s.MinDrag = 0
	s.Jenks = 0
	s.DriftTime = 0
	s.PressToSelect = true
	s.Skipback = true
	s.ExternalDevice = true
	err := s2.ReadYAML("./trackpoint.yml")
	if err != nil {
		t.Fatal(err)
	}
	checkDefaults(s, t)

	if s2.ReadYAML("./trackpoint-asdf.yml") == nil {
		t.Fatal("expected error")
	}
}

func TestSettings_ForEach(t *testing.T) {
	s := NewSettings()
	e := errors.New("dummy error")
	err := s.ForEach(func(key string, value string) (err error) {
		if key == "press_to_select" {
			err = e
		}
		return
	})
	if err != e {
		t.Fatalf("Expected %v, got %v", e, err)
	}
	err = s.ForEach(func(key string, value string) (err error) {
		if key == "thresh" {
			err = e
		}
		return
	})
	if err != e {
		t.Fatalf("Expected %v, got %v", e, err)
	}
	err = s.ForEach(func(key string, value string) (err error) {
		if key == "skipback" {
			err = e
		}
		return
	})
	if err != e {
		t.Fatalf("Expected %v, got %v", e, err)
	}
}

func TestSettings_ToStringMap(t *testing.T) {

	s := NewSettings()
	s.Values.PressToSelect = true
	actual := s.ToStringMap()
	if len(defaults) != len(actual) {
		t.Fatalf("expected length of %v, got %v", len(defaults), len(actual))
	}

	for key, value := range actual {
		if key == "press_to_select" {
			if value != "1" {
				t.Fatalf("Expected %v, got %v", defaults[key], value)
			}
		} else if defaults[key] != value {
			t.Fatalf("Expected %v, got %v", defaults[key], value)
		}
	}

}

func TestSettings_Keys(t *testing.T) {
	expectedMap := make(map[string]bool)
	for _, v := range keys {
		expectedMap[v] = true
	}
	k := NewSettings().Keys()
	if len(expectedMap) != len(k) {
		t.Fatalf("expected keys with length %v, got %v", len(expectedMap), len(keys))
	}
	for _, v := range k {
		if !expectedMap[v] {
			t.Fatalf("unexpected key %v", v)
		}
	}
}

func TestSettings_Get(t *testing.T) {
	s := NewSettings()
	for key, value := range defaults {
		if actual := s.Get(key); actual != value {
			t.Fatalf("expected %v, got %v", value, actual)
		}
	}

	if actual := s.Get("someBogusKey"); actual != "" {
		t.Fatalf("expected %v, got %v", nil, actual)
	}
}
