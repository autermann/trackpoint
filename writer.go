package main

import (
	"errors"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

var (
	// ErrReadValueIsNotWrittenValue indicates that the written value is not the read value.
	ErrReadValueIsNotWrittenValue = errors.New("read value is not written value")
)

// SettingsReaderWriter reads and writes settings.
type SettingsReaderWriter struct {
	SysfsPath           string        // SysfsPath is the SYS FS path of the device.
	MaxWriteAttempts    uint          // MaxWriteAttempts is the maximum number of attempts to write a file.
	WriteTimeout        time.Duration // WriteTimeout is the time a single write may take.
	TimeBetweenAttempts time.Duration // TimeBetweenAttempts is the time between write attempts
}

// NewSettingsReaderWriter creates a new SettingsReaderWriter.
func NewSettingsReaderWriter(path string) *SettingsReaderWriter {
	return &SettingsReaderWriter{
		SysfsPath:           path,
		MaxWriteAttempts:    10,
		WriteTimeout:        3 * time.Second,
		TimeBetweenAttempts: 10 * time.Second,
	}
}

// Set writes the settings.
func (t *SettingsReaderWriter) Set(settings *Settings) error {
	return RetryWait(t.TimeBetweenAttempts, func(attempt uint) (bool, error) {
		log.Printf("writing (attempt %v)", attempt)
		err := settings.ForEach(t.SetValue)
		if err != nil {
			log.Print(err)
			return attempt < t.MaxWriteAttempts, err
		}
		return false, nil
	})
}

// SetValue sets the value for a key.
func (t *SettingsReaderWriter) SetValue(key, value string) error {

	switch written, err := t.GetValue(key); {
	case err != nil:
		return err
	case written == value:
		log.Printf("%15v: is already %3v", key, value)
		// early return, no value change
		return nil
	}

	log.Printf("%15v: setting to %3v", key, value)

	path := filepath.Join(t.SysfsPath, key)
	if err := t.writeValue(path, value); err != nil {
		return err
	}

	switch written, err := t.GetValue(key); {
	case err != nil:
		return err
	case written != value:
		return ErrReadValueIsNotWrittenValue
	default:
		return nil
	}
}

func (t *SettingsReaderWriter) writeValue(path, value string) error {
	fd, err := syscall.Open(path, syscall.O_APPEND|syscall.O_WRONLY, syscall.O_SYNC)
	if err != nil {
		return err
	}

	if err := syscall.SetNonblock(fd, true); err != nil {
		return err
	}

	defer func() {
		if e := syscall.Close(fd); e != nil {
			log.Print(e)
		}
	}()

	return Timeout(t.WriteTimeout, func() error {
		bytes := []byte(value)
		switch n, err := syscall.Write(fd, []byte(value)); {
		case err != nil:
			return err
		case n != len(bytes):
			return errors.New("could not write all bytes")
		default:
			return nil
		}
	})
}

// GetValue gets the value for a key.
func (t *SettingsReaderWriter) GetValue(key string) (string, error) {
	bytes, err := ioutil.ReadFile(filepath.Join(t.SysfsPath, key))
	return strings.TrimSpace(string(bytes)), err
}
