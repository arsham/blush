package cmd

import "errors"

// Errors regarding application start up.
var (
	ErrNoInput      = errors.New("no input provided")
	ErrNoFilesFound = errors.New("no files found")
)
