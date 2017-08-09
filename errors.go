package degob

import (
	"errors"
	"fmt"
)

// Error is a custom error type that hopefully gives a little more information
// about where your error was in the Gob.
type Error struct {
	Processed uint64
	Err       error
	RawGob    []byte
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: %v after processing %d bytes", e.Err, e.Processed)
}

func errGen(s string) func(uint64, []byte) *Error {
	return func(b uint64, gob []byte) *Error {
		return &Error{
			Processed: b,
			Err:       errors.New(s),
			RawGob:    gob,
		}
	}
}

var (
	genericError = func(err error, b uint64, gob []byte) *Error {
		return &Error{
			Processed: b,
			Err:       err,
			RawGob:    gob,
		}
	}
	errUintTooBig        = errGen("found uint claiming to be longer than 8 bytes")
	errDuplicateType     = errGen("duplicate type found")
	errUnknownDelta      = errGen("found unexpected delta value when trying to decode type")
	errCorruptCommonType = errGen("bad field number for CommonType")
	errBadString         = errGen("failed to decode string")
)
