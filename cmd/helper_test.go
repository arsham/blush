package cmd_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/arsham/blush/blush"
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

	t.Cleanup(func() {
		os.Args = oldArgs
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	})
	return stdout, stderr
}

func getPipe(t *testing.T) *os.File {
	t.Helper()
	file, err := ioutil.TempFile("", "blush_pipe")
	assert.NoError(t, err)
	name := file.Name()
	file.Close()

	t.Cleanup(func() {
		err = os.Remove(name)
		assert.NoError(t, err)
	})

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
