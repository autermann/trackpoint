package trackpoint

import (
	"errors"
	"time"
)

var (
	// ErrTimeout indicates that an operation timed out.
	ErrTimeout = errors.New("timeout")
)

// Retry retries to call fn until it succeeds.
func Retry(fn func(attempt uint) (bool, error)) error {
	return RetryWait(0, fn)
}

// RetryWait retries to call fn until it succeeds.
func RetryWait(wait time.Duration, fn func(attempt uint) (bool, error)) error {
	var err error
	var cont bool
	attempt := uint(1)
	for {
		cont, err = fn(attempt)
		if !cont || err == nil {
			break
		}
		if int64(wait) > 0 {
			time.Sleep(wait)
		}
		attempt++
	}
	return err
}

// Timeout wraps a function to timeout after a given duration.
func Timeout(timeout time.Duration, fn func() error) error {
	c := make(chan error, 1)
	go func() { c <- fn() }()
	timer := time.NewTimer(timeout)
	select {
	case err := <-c:
		timer.Stop()
		return err
	case <-timer.C:
		return ErrTimeout
	}
}

// WrapError runs a function in a goroutine and offers a channel for the returned error.
func WrapError(fn func() error) <-chan error {
	c := make(chan error, 1)
	go func() {
		err := fn()
		if err != nil {
			c <- err
		}
	}()
	return c
}
