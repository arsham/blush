package blush

import "errors"

// ErrNoWriter is returned if a nil object is passed to the WriteTo method.
var ErrNoWriter = errors.New("no output defined")

// ErrClosed is returned if the reader is closed and you try to read from it.
var ErrClosed = errors.New("reader already closed")
