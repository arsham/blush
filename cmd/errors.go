package cmd

import "errors"

var (
	// ErrNoInput is returned when the application doesn't receive any files as the
	// last arguments or a stream of inputs from shell's pipe.
	ErrNoInput = errors.New("no input provided")

	// ErrNoFilesFound is returned when the files pattern passed to the application
	// doesn't match any existing files.
	ErrNoFilesFound = errors.New("no files found")
)
