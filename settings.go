package trackpoint

import (
	"errors"
	"io/ioutil"
	"reflect"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"
)

// Settings are the configurable TrackPoint settings.
type Settings struct {
	Path      string        // Path is the path to the settings
	SysfsPath string        `yaml:"sysfs"`    // SysfsPath is the path to the SYSFS device.
	Values    *Values       `yaml:"values"`   // Values are the trackpoint properties.
	Daemon    bool          `yaml:"daemon"`   // Daemon lets the tool act as a daemon.
	Interval  time.Duration `yaml:"interval"` // Interval is the interval at which the daemon executes.
}

// Values are the configurable values.
type Values struct {
	DragHysteresis uint8 `yaml:"draghys" trackpoint:"draghys"`                 // Drag Hysteresis (how hard it is to drag with Z-axis pressed).
	Threshold      uint8 `yaml:"thresh" trackpoint:"thresh"`                   // Minimum value for a Z-axis press.
	UpThreshold    uint8 `yaml:"upthresh" trackpoint:"upthresh"`               // Used to generate a 'click' on Z-axis.
	ZTime          uint8 `yaml:"ztime" trackpoint:"ztime"`                     // How sharp of a press.
	Sensitivity    uint8 `yaml:"sensitivity" trackpoint:"sensitivity"`         // Sensitivity.
	Inertia        uint8 `yaml:"inertia" trackpoint:"inertia"`                 // Negative Inertia.
	Speed          uint8 `yaml:"speed" trackpoint:"speed"`                     // Speed of TP Cursor.
	Reach          uint8 `yaml:"reach" trackpoint:"reach"`                     // Backup for Z-axis press.
	MinDrag        uint8 `yaml:"mindrag" trackpoint:"mindrag"`                 // Minimum amount of force needed to trigger dragging.
	Jenks          uint8 `yaml:"jenks" trackpoint:"jenks"`                     // Minimum curvature for double click.
	DriftTime      uint8 `yaml:"drift_time" trackpoint:"drift_time"`           // How long a 'hands off' condition must last for drift correction to occur.
	PressToSelect  bool  `yaml:"press_to_select" trackpoint:"press_to_select"` // Press to Select.
	Skipback       bool  `yaml:"skipback" trackpoint:"skipback"`               // Suppress movement after drag release.
	ExternalDevice bool  `yaml:"ext_dev" trackpoint:"ext_dev"`                 // Disable external device.

}

// NewSettings creates a new Settings with default values.
func NewSettings() *Settings {
	s := &Settings{Values: &Values{}}
	s.Values.SetDefaults()
	return s
}

// SetDefaults resets the settings.
func (s *Values) SetDefaults() {
	s.DragHysteresis = DefaultDragHysteresis
	s.Threshold = DefaultThreshold
	s.UpThreshold = DefaultUpThreshold
	s.ZTime = DefaultZTime
	s.Sensitivity = DefaultSensitivity
	s.Inertia = DefaultInertia
	s.Speed = DefaultSpeed
	s.Reach = DefaultReach
	s.MinDrag = DefaultMinDrag
	s.Jenks = DefaultJenks
	s.DriftTime = DefaultDriftTime
	s.PressToSelect = DefaultPressToSelect
	s.Skipback = DefaultSkipback
	s.ExternalDevice = DefaultExternalDevice
}

// ReadYAML reads a YAML file into the settings.
func (s *Settings) ReadYAML(path string) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bytes, s)
}

// Get gets the value of the key.
func (s *Settings) Get(k string) string {
	var value string
	found := errors.New("found")
	s.ForEach(func(k2, v string) error {
		if k == k2 {
			value = v
			return found
		}
		return nil
	})
	return value
}

// Keys gets the keys of the settings.
func (s *Settings) Keys() []string {
	t := reflect.TypeOf(*s.Values)
	keys := make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		keys[i] = t.Field(i).Tag.Get("trackpoint")
	}
	return keys
}

// ForEach iterates over the settings.
func (s *Settings) ForEach(fn func(key, value string) error) (err error) {
	v := reflect.ValueOf(*s.Values)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("trackpoint")
		switch f.Type.Kind() {
		case reflect.Uint8:
			err = fn(tag, strconv.FormatUint(v.Field(i).Uint(), 10))
		case reflect.Bool:
			if v.Field(i).Bool() {
				err = fn(tag, "1")
			} else {
				err = fn(tag, "0")
			}
		}
		if err != nil {
			return
		}
	}
	return
}

// ToStringMap convert the settings to a map.
func (s *Settings) ToStringMap() map[string]string {
	m := make(map[string]string)
	err := s.ForEach(func(key string, value string) error {
		m[key] = value
		return nil
	})
	if err != nil {
		panic(err)
	}
	return m
}
