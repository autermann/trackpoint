package trackpoint

import "time"
import "github.com/cheekybits/genny/generic"

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "DebounceType=bool"

// DebounceType is the chan type of the
type DebounceType generic.Type

// DebounceDebounceType debounces the input chan.
func DebounceDebounceType(interval time.Duration, input chan DebounceType) chan DebounceType {
	output := make(chan DebounceType, 1)

	go func() {
		execute := false
		var i DebounceType
		for {
			select {
			case i, more := <-input:
				execute = true
				if !more {
					output <- i
					close(output)
					return
				}
			case <-time.After(interval):
				if execute {
					output <- i
					execute = false
				}
			}
		}
	}()
	return output
}
