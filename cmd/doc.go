// Package cmd bootstraps the application.
//
// Main() reads the provided arguments from the command line and creates a
// blush.Blush instance. If there is any error, it will terminate the
// application with os.Exit(1), otherwise it then uses io.Copy() to write to
// standard output and exits normally.
//
// GetBlush() returns an error if no arguments are provided or it can't find all
// the passed files. Files should be last arguments, otherwise they are counted
// as matching strings. If there is no file passed, the input should come in
// from Stdin as a pipe.
//
// hasArgs(args ...string) function looks for args in input and if it finds it,
// it removes it and put the rest in the remaining slice.
//
// Notes
//
// We are not using the usual flag package because it cannot handle variables in
// the args and continues grouping of passed arguments.
package cmd
