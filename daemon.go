package trackpoint

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Daemon is daemon that does stuff.
type Daemon interface {
	DoStuff(stop <-chan bool) error
}

// Run starts the daemon.
func Run(d Daemon) error {
	var (
		stop = make(chan bool, 1)
		e    = make(chan error, 1)
		s    = make(chan os.Signal, 1)
	)

	signal.Notify(s, syscall.SIGHUP,
		syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGQUIT, syscall.SIGKILL)

	go func() { e <- d.DoStuff(stop) }()

	for {
		select {
		case signal := <-s:
			log.Printf("received %v", signal)
			stop <- true
		case err := <-e:
			return err
		}
	}
}

// SettingsDaemon is a simple daemon implementation.
type SettingsDaemon struct {
	*sync.RWMutex
	rw        *SettingsReaderWriter
	Settings  *Settings
	SysfsPath string
}

// NewSettingsDaemon creates a new daemon.
func NewSettingsDaemon(settings *Settings) (d *SettingsDaemon) {
	return &SettingsDaemon{
		RWMutex:  &sync.RWMutex{},
		rw:       NewSettingsReaderWriter(settings.SysfsPath),
		Settings: settings,
	}
}

func (d *SettingsDaemon) watchSettings(stop <-chan bool) (changed chan bool, errors chan error) {
	changed = make(chan bool, 1)
	errors = make(chan error, 1)

	watcher, err := fsnotify.NewWatcher()
	log.Print("Created watcher")
	if err != nil {
		errors <- err
		close(changed)
		close(errors)
		return
	}

	closeWatcher := func(err error) {
		err2 := watcher.Close()
		if err2 != nil {
			log.Println("error closing watcher", err)
		}
		if err != nil {
			errors <- err
		} else if err2 != nil {
			errors <- err2
		}
		close(changed)
		close(errors)
	}

	err = watcher.Add(d.Settings.Path)
	if err != nil {
		closeWatcher(err)
		return
	}

	go func() {
		for {
			select {
			case <-stop:
				log.Println("watcher received stop")
				closeWatcher(nil)
				return
			case err = <-watcher.Errors:
				log.Println("error watching file", err)
				closeWatcher(err)
				return
			case event := <-watcher.Events:
				switch op := event.Op; op {
				case fsnotify.Create, fsnotify.Rename, fsnotify.Write:
					if path := event.Name; path == d.Settings.Path {
						err = d.onSettingsChange()
						if err != nil {
							closeWatcher(err)
							return
						}
						changed <- true
					}
				}
			}
		}
	}()
	return
}

func (d *SettingsDaemon) onSettingsChange() error {
	log.Println("settings file changed")
	return d.refreshSettings()
}

func (d *SettingsDaemon) refreshSettings() error {
	d.Lock()
	defer d.Unlock()
	log.Printf("refreshing settings")
	return d.Settings.ReadYAML(d.Settings.Path)
}

// DoStuff does the stuff.
func (d *SettingsDaemon) DoStuff(stop <-chan bool) (err error) {
	err = d.applySettings()
	if err != nil {
		log.Print(err)
		err = nil
	}

	if d.Settings.Path == "" {
		for {
			select {
			case <-stop:
				log.Printf("received stop")
				return
			case <-time.After(d.Settings.Interval):
				d.applySettingsNoError()
			}
		}
	} else {

		stop2 := make(chan bool, 1)
		changed, errors := d.watchSettings(stop2)
		changed = DebounceBool(2*time.Second, changed)

		for {
			select {
			case <-stop:
				log.Printf("received stop")
				stop2 <- true
				err = nil
				return
			case err = <-errors:
				return
			case <-changed:
				d.applySettingsNoError()
			case <-time.After(d.Settings.Interval):
				d.applySettingsNoError()
			}
		}
	}

}

func (d *SettingsDaemon) applySettings() error {
	d.RLock()
	defer d.RUnlock()
	return d.rw.Set(d.Settings)
}

func (d *SettingsDaemon) applySettingsNoError() {
	err := d.applySettings()
	if err != nil {
		log.Print(err)
	}
}

func (d *SettingsDaemon) applySetting(key string) error {
	d.RLock()
	defer d.RUnlock()
	return d.rw.SetValue(key, d.Settings.Get(key))
}
