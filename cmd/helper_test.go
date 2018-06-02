package cmd_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/arsham/blush/blush"
	"github.com/arsham/blush/cmd"
)

var leaveMeHere = "LEAVEMEHERE"

type stdFile struct {
	f *os.File
}

func (s *stdFile) String() string {
	s.f.Seek(0, 0)
	buf := new(bytes.Buffer)
	buf.ReadFrom(s.f)
	return buf.String()
}
func (s *stdFile) Close() error {
	return s.f.Close()
}

func newStdFile(t *testing.T, name string) (*stdFile, func()) {
	f, err := ioutil.TempFile("", name)
	if err != nil {
		t.Fatal(err)
	}
	sf := &stdFile{f}
	return sf, func() {
		f.Close()
		os.Remove(f.Name())
	}
}

func setup(t *testing.T, args string) (stdout, stderr *stdFile, cleanup func()) {
	oldArgs := os.Args
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	oldFatalErr := cmd.FatalErr

	stdout, outCleanup := newStdFile(t, "stdout")
	stderr, errCleanup := newStdFile(t, "stderr")
	os.Stdout = stdout.f
	os.Stderr = stderr.f

	os.Args = []string{"blush"}
	if len(args) > 1 {
		os.Args = append(os.Args, strings.Split(args, " ")...)
	}
	cmd.FatalErr = func(s error) {
		fmt.Fprintf(os.Stderr, "%s\n", s)
	}

	cleanup = func() {
		outCleanup()
		errCleanup()
		os.Args = oldArgs
		os.Stdout = oldStdout
		os.Stderr = oldStderr
		cmd.FatalErr = oldFatalErr
	}
	return stdout, stderr, cleanup
}

func getPipe(t *testing.T) (*os.File, func()) {
	file, err := ioutil.TempFile("", "blush_pipe")
	if err != nil {
		t.Fatal(err)
	}
	name := file.Name()
	rmFile := func() {
		if err = os.Remove(name); err != nil {
			t.Error(err)
		}
	}
	file.Close()
	rmFile()
	file, err = os.OpenFile(name, os.O_CREATE|os.O_RDWR, os.ModeCharDevice|os.ModeDevice)
	if err != nil {
		t.Fatal(err)
	}
	return file, rmFile
}

func argsEqual(a, b []blush.Finder) bool {
	isIn := func(a blush.Finder, haystack []blush.Finder) bool {
		for _, b := range haystack {
			if reflect.DeepEqual(a, b) {
				return true
			}
		}
		return false
	}

	for _, item := range b {
		if !isIn(item, a) {
			return false
		}
	}
	return true
}
