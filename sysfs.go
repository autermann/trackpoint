package trackpoint

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// TrackPointName is the string the device name has to contain to be a TrackPoint.
	TrackPointName = "TrackPoint"
	// SysfsBaseDir is the base directory to search for the TrackPoint device.
	SysfsBaseDir = "/sys/devices/platform/i8042"
)

var (
	// ErrDeviceDirNotFound indicates that the TrackPoint device was not found.
	ErrDeviceDirNotFound = errors.New("device directory not found")
)

// GetDeviceDirectory get the device directory of the TrackPoint in the SYS FS.
func GetDeviceDirectory() (result string, err error) {
	err = RetryWait(1*time.Second, func(attempt uint) (bool, error) {
		errFileFound := errors.New("file found")
		err = filepath.Walk(SysfsBaseDir,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() && info.Name() == "name" {
					bytes, err := ioutil.ReadFile(path)
					if err != nil {
						return err
					}
					if strings.Contains(string(bytes), TrackPointName) {
						result = path
						return errFileFound
					}
				}
				return nil
			})
		if err == errFileFound {
			err = nil
			result = filepath.Dir(filepath.Dir(filepath.Dir(result)))
		} else if err == nil {
			err = ErrDeviceDirNotFound
		}
		return attempt < 10, err
	})
	return result, err
}
