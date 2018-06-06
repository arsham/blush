package cmd_test

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/arsham/blush/cmd"
)

func TestCaptureSignals(t *testing.T) {
	tcs := []struct {
		name   string
		signal os.Signal
		code   int
	}{
		{"SIGINT", syscall.SIGINT, 130},
		{"SIGTERM", syscall.SIGTERM, 130},
		{"SIGPIPE", syscall.SIGPIPE, 0},
		{"other", syscall.Signal(-1), 1},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			sig := make(chan os.Signal, 1)
			code := make(chan int)
			exit := func(c int) {
				code <- c
			}
			cmd.WaitForSignal(sig, exit)
			sig <- tc.signal
			select {
			case code := <-code:
				if code != tc.code {
					t.Errorf("exit code = %d, want %d", code, tc.code)
				}
			case <-time.After(500 * time.Millisecond):
				t.Error("exit function wasn't called")
			}
		})
	}
}
