// Package cmd bootstraps the application.
//
// Main() reads the provided arguments from the command line and creates a
// blush.Blush instance. If there is any error, it will terminate the
// application with os.Exit(1), otherwise it calls the Write() method of Blush
// and exits normally.
//
// GetBlush() returns an error if no arguments are provided or it can't find all
// the passed files. Files should be last arguments, otherwise they are counted
// as matching strings. If there is no file passed, the input should come in
// from Stdin as a pipe.
//
// hasArg(input []string, arg string) function looks for arg in input and if it
// finds it, it removes it and returns the remaining slice with a boolean to
// tell it was found.
//
// Notes
//
// We are not using the usual flag package because it cannot handle variables in
// the args and continues grouping of passed arguments.
package cmd
