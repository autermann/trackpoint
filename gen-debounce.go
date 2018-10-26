// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package trackpoint

import "time"

// DebounceBool debounces the input chan.
func DebounceBool(interval time.Duration, input chan bool) chan bool {
	output := make(chan bool, 1)

	go func() {
		execute := false
		for {
			var i bool
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
