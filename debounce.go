package main

import (
	"time"

	"github.com/cheekybits/genny/generic"
)

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "DebounceType=bool"

// DebounceType is the chan type of the
type DebounceType generic.Type

// DebounceDebounceType debounces the input chan.
func DebounceDebounceType(interval time.Duration, input chan DebounceType) chan DebounceType {
	output := make(chan DebounceType)

	go func() {
		var i DebounceType
		var ok bool

		for {
			i, ok = <-input

			if !ok {
				close(output)
				return
			}
		F:
			for {
				select {
				case i, ok = <-input:
					if !ok {
						close(output)
						return
					}
				case <-time.After(interval):
					output <- i
					break F
				}
			}
		}
	}()
	return output
}
