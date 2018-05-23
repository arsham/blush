package blush

import "errors"

// Errors for creating a new Blush object.
var (
	ErrNoFiles  = errors.New("no files")
	ErrNoInput  = errors.New("no input")
	ErrNoWriter = errors.New("no output defined")
)
