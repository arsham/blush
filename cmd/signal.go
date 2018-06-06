package cmd

import (
	"os"
	"os/signal"
	"syscall"
)

// WaitForSignal calls exit with code 130 if receives an SIGINT or SIGTERM, 0 if
// SIGPIPE, and 1 otherwise.
func WaitForSignal(sig chan os.Signal, exit func(int)) {
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGPIPE)
	go func() {
		s := <-sig
		switch s {
		case syscall.SIGINT, syscall.SIGTERM:
			exit(130) // Ctrl+c
		case syscall.SIGPIPE:
			exit(0)
		}
		exit(1)
	}()
}
