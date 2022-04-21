package cmd_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/arsham/blush/blush"
	"github.com/bouk/monkey"
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

func newStdFile(t *testing.T, name string) *stdFile {
	t.Helper()
	f, err := ioutil.TempFile("", name)
	if err != nil {
		t.Fatal(err)
	}
	sf := &stdFile{f}
	t.Cleanup(func() {
		f.Close()
		os.Remove(f.Name())
	})
	return sf
}

func setup(t *testing.T, args string) (stdout, stderr *stdFile) {
	t.Helper()
	oldArgs := os.Args
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	stdout = newStdFile(t, "stdout")
	stderr = newStdFile(t, "stderr")
	os.Stdout = stdout.f
	os.Stderr = stderr.f

	os.Args = []string{"blush"}
	if len(args) > 1 {
		os.Args = append(os.Args, strings.Split(args, " ")...)
	}
	fatalPatch := monkey.Patch(log.Fatal, func(msg ...interface{}) {
		fmt.Fprintln(os.Stderr, msg)
	})
	fatalfPatch := monkey.Patch(log.Fatalf, func(format string, v ...interface{}) {
		fmt.Fprintf(os.Stderr, format, v...)
	})

	t.Cleanup(func() {
		os.Args = oldArgs
		os.Stdout = oldStdout
		os.Stderr = oldStderr
		fatalPatch.Unpatch()
		fatalfPatch.Unpatch()
	})
	return stdout, stderr
}

func getPipe(t *testing.T) *os.File {
	t.Helper()
	file, err := ioutil.TempFile("", "blush_pipe")
	assert.NoError(t, err)
	name := file.Name()
	rmFile := func() {
		err := os.Remove(name)
		assert.NoError(t, err)
	}
	file.Close()
	rmFile()
	file, err = os.OpenFile(name, os.O_CREATE|os.O_RDWR, os.ModeCharDevice|os.ModeDevice)
	assert.NoError(t, err)
	return file
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
