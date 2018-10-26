package main

import (
	"flag"
	"os"
	"time"
)

// ParseFlags parses the supplied arguments to the settings.
func ParseFlags(args []string, settings *Settings) (err error) {
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)

	settings.SysfsPath, err = GetDeviceDirectory()
	if err != nil {
		return err
	}

	var config string
	var interval string

	fs.StringVar(&config, "config", "", "The path to the config file")
	fs.StringVar(&config, "c", "", "The path to the config file (shorthand)")
	fs.StringVar(&interval, "interval", "5s", "The interval at which the daemon executes.")
	fs.StringVar(&settings.SysfsPath, "sysfs", settings.SysfsPath, "The path to the SYSFS device.")

	fs.Uint("draghys", DefaultDragHysteresis, "Drag Hysteresis (how hard it is to drag with Z-axis pressed).")
	fs.Uint("thresh", DefaultThreshold, "Minimum value for a Z-axis press.")
	fs.Uint("upthresh", DefaultUpThreshold, "Used to generate a 'click' on Z-axis.")
	fs.Uint("ztime", DefaultZTime, "How sharp of a press.")
	fs.Uint("reach", DefaultReach, "Backup for Z-axis press.")
	fs.Uint("jenks", DefaultJenks, "Minimum curvature for double click.")
	fs.Uint("drifttime", DefaultDriftTime, "How long a 'hands off' condition must last for drift correction to occur.")
	fs.Uint("speed", DefaultSpeed, "Speed of TP Cursor.")
	fs.Uint("sensitivity", DefaultSensitivity, "Sensitivity.")
	fs.Uint("inertia", DefaultInertia, "Negative Inertia.")
	fs.Uint("mindrag", DefaultMinDrag, "Minimum amount of force needed to trigger dragging.")
	fs.Bool("pts", DefaultPressToSelect, "If press-to-select should be active.")
	fs.Bool("skipback", DefaultSkipback, "Suppress movement after drag release.")
	fs.Bool("extdev", DefaultExternalDevice, "Disable external device.")

	fs.BoolVar(&settings.Daemon, "daemon", false, "Run as a daemon")
	fs.BoolVar(&settings.Daemon, "d", false, "Run as a daemon (shorthand)")

	fs.Parse(args[1:])

	settings.Interval, err = time.ParseDuration(interval)
	if err != nil {
		return err
	}

	if config != "" {
		settings.Path = config
		err = settings.ReadYAML(config)
		if err != nil {
			return
		}
	}

	fs.Visit(func(f *flag.Flag) {
		v := f.Value.(flag.Getter).Get()
		switch f.Name {
		case "draghys":
			settings.Values.DragHysteresis = uint8(v.(uint))
		case "thresh":
			settings.Values.Threshold = uint8(v.(uint))
		case "upthresh":
			settings.Values.UpThreshold = uint8(v.(uint))
		case "ztime":
			settings.Values.ZTime = uint8(v.(uint))
		case "reach":
			settings.Values.Reach = uint8(v.(uint))
		case "jenks":
			settings.Values.Jenks = uint8(v.(uint))
		case "drifttime":
			settings.Values.DriftTime = uint8(v.(uint))
		case "speed":
			settings.Values.Speed = uint8(v.(uint))
		case "sensitivity":
			settings.Values.Sensitivity = uint8(v.(uint))
		case "inertia":
			settings.Values.Inertia = uint8(v.(uint))
		case "mindrag":
			settings.Values.MinDrag = uint8(v.(uint))
		case "pts":
			settings.Values.PressToSelect = v.(bool)
		case "skipback":
			settings.Values.Skipback = v.(bool)
		case "extdev":
			settings.Values.ExternalDevice = v.(bool)
		}
	})
	return
}

func main() {
	var err error
	settings := NewSettings()

	err = ParseFlags(os.Args, settings)
	if err != nil {
		panic(err)
	}
	if settings.Daemon {
		d := NewSettingsDaemon(settings)
		if err = Run(d); err != nil {
			panic(err)
		}
	} else {
		rw := NewSettingsReaderWriter(settings.SysfsPath)
		if err = rw.Set(settings); err != nil {
			panic(err)
		}
	}
}
