package blush

import "errors"

var (
	// ErrNoWriter is returned if a nil object is passed to the WriteTo method.
	ErrNoWriter = errors.New("no writer defined")

	// ErrNoFinder is returned if there is no finder passed to Blush.
	ErrNoFinder = errors.New("no finders defined")

	// ErrClosed is returned if the reader is closed and you try to read from
	// it.
	ErrClosed = errors.New("reader already closed")

	// ErrReadWriteMix is returned when the Read and WriteTo are called on the
	// same object.
	ErrReadWriteMix = errors.New("you cannot mix Read and WriteTo calls")
)
