package cmd

import "errors"

// Errors regarding application start up.
var (
	ErrNoInput      = errors.New("no input provided")
	ErrFileNotFound = errors.New("file not found")
)
