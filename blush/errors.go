package blush

import "errors"

// Errors for creating a new Blush object.
var (
	ErrNoReader = errors.New("no input specified")
	ErrNoInput  = errors.New("no input")
	ErrNoWriter = errors.New("no output defined")
)
