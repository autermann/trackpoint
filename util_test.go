package trackpoint

import (
	"errors"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	times := 5
	Retry(func(attempt uint) (bool, error) {
		times--
		if times > 0 {
			return true, errors.New("retry")
		}
		return false, nil
	})
	if times != 0 {
		t.Fatal("expected to be called 5 times")
	}
	times = 5
	Retry(func(attempt uint) (bool, error) {
		times--
		return true, nil
	})
	if times != 4 {
		t.Fatal("expected to be called once")
	}

	times = 5
	Retry(func(attempt uint) (bool, error) {
		times--
		return false, errors.New("retry")
	})
	if times != 4 {
		t.Fatal("expected to be called once")
	}

	times = 5
	RetryWait(1*time.Second, func(attempt uint) (bool, error) {
		times--
		if times > 0 {
			return true, errors.New("retry")
		}
		return false, nil
	})
	if times != 0 {
		t.Fatal("expected to be called 5 times")
	}
	times = 5
	RetryWait(1*time.Second, func(attempt uint) (bool, error) {
		times--
		return true, nil
	})
	if times != 4 {
		t.Fatal("expected to be called once")
	}

	times = 5
	RetryWait(1*time.Second, func(attempt uint) (bool, error) {
		times--
		return false, errors.New("retry")
	})
	if times != 4 {
		t.Fatal("expected to be called once")
	}

}

func TestTimeout(t *testing.T) {
	var err error
	err = Timeout(time.Second, func() error {
		time.Sleep(time.Second * 2)
		return nil
	})
	if err != ErrTimeout {
		t.Fatal("expected timeout")
	}
	err = Timeout(time.Second*2, func() error {
		time.Sleep(time.Second)
		return nil
	})
	if err != nil {
		t.Fatal("expected no error")
	}
	e := errors.New("dummy")
	err = Timeout(time.Second*2, func() error {
		time.Sleep(time.Second)
		return e
	})
	if err != e {
		t.Fatalf("expected %v, got %v", e, err)
	}
	err = Timeout(time.Second, func() error {
		time.Sleep(time.Second * 2)
		return e
	})
	if err != ErrTimeout {
		t.Fatalf("expected %v, got %v", ErrTimeout, err)
	}
}

func TestWrapError(t *testing.T) {
	var errc <-chan error
	result := false
	errc = WrapError(func() error {
		time.Sleep(time.Second)
		result = true
		return nil
	})

	select {
	case <-errc:
		t.Fatalf("expected no error")
	case <-time.After(2 * time.Second):
		if !result {
			t.Fatal("expected fn to run")
		}
	}

	result = false
	expectedError := errors.New("expected")
	errc = WrapError(func() error {
		time.Sleep(time.Second)
		result = true
		return expectedError
	})
	select {
	case err := <-errc:
		if !result {
			t.Fatal("expected fn to run")
		}
		if err != expectedError {
			t.Fatalf("expected %v, got %v", expectedError, err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("expected error")
	}

}
