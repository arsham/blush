package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/alecthomas/assert"
)

func TestArgs(t *testing.T) {
	tcs := []struct {
		name        string
		input       []string
		wantErr     error
		colour      bool
		noFilename  bool
		recursive   bool
		insensitive bool
	}{
		{name: "help", input: []string{"--help"}, wantErr: errShowHelp},
		{name: "colour", input: []string{"--colour"}, colour: true},
		{name: "colour and help", input: []string{"--colour", "--help"}, wantErr: errShowHelp},
		{name: "colour american", input: []string{"--colour"}, colour: true},
		{name: "colour short", input: []string{"-C"}, colour: true},
		{name: "no filename", input: []string{"-h"}, noFilename: true},
		{name: "no filename long", input: []string{"--no-filename"}, noFilename: true},
		{name: "rec", input: []string{"-R"}, recursive: true},
		{name: "ins", input: []string{"-i"}, insensitive: true},
		{name: "ins rec", input: []string{"-i", "-R"}, insensitive: true, recursive: true},
		{name: "rec ins", input: []string{"-R", "-i"}, insensitive: true, recursive: true},
		{
			name: "rec ins nofile", input: []string{"-R", "-i", "-h"},
			insensitive: true, recursive: true, noFilename: true,
		},
		{
			name: "nofile rec ins", input: []string{"-h", "-R", "-i"},
			insensitive: true, recursive: true, noFilename: true,
		},
		{
			name: "nofile rec ins colour", input: []string{"-h", "-R", "-i", "-C"},
			insensitive: true, recursive: true, noFilename: true, colour: true,
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			a, err := newArgs(tc.input...)
			if tc.wantErr != nil {
				assert.True(t, errors.Is(err, tc.wantErr))
			}
			if err != nil {
				assert.Nil(t, a)
				return
			}
			assert.Equal(t, tc.colour, a.colour)
			assert.Equal(t, tc.noFilename, a.noFilename)
			assert.Equal(t, tc.recursive, a.recursive)
			assert.Equal(t, tc.insensitive, a.insensitive)
		})
	}
}

func TestArgsPipe(t *testing.T) {
	getPipe(t)
	a, err := newArgs("something")
	assert.NoError(t, err)
	assert.NotNil(t, a)
	assert.True(t, a.stdin)
}

func TestArgsPaths(t *testing.T) {
	dir, err := ioutil.TempDir("", "blush_main")
	assert.NoError(t, err)
	defer func() {
		err := os.RemoveAll(dir)
		assert.NoError(t, err)
	}()
	f1, err := ioutil.TempFile(dir, "main")
	assert.NoError(t, err)
	f2, err := ioutil.TempFile(dir, "main")
	assert.NoError(t, err)

	tcs := []struct {
		name      string
		input     []string
		wantPaths []string
		wantErr   bool
	}{
		{"not found", []string{"nowhere"}, []string{}, true},
		{"only a file", []string{f1.Name()}, []string{f1.Name()}, false},
		{"two files", []string{f1.Name(), f2.Name()}, []string{f1.Name(), f2.Name()}, false},
		{"arg between two files", []string{f1.Name(), "-a", f2.Name()}, []string{}, true},
		{"prefix file", []string{"a", f1.Name()}, []string{f1.Name()}, false},
		{"prefix arg file", []string{"-r", f1.Name()}, []string{}, true},
		{"file matches but is an argument", []string{"-r", f1.Name(), f2.Name()}, []string{f2.Name()}, false},
		{"star dir", []string{path.Join(dir, "*")}, []string{path.Join(dir, "*")}, false},
		{"stared dir", []string{dir + "*"}, []string{dir + "*"}, false},
		{
			"many prefixes",
			[]string{"--#7ff", "main", "-g", "package", "-r", "a", path.Join(dir, "*")},
			[]string{path.Join(dir, "*")},
			false,
		},
		{
			"many prefixes star",
			[]string{"--#7ff", "main", "-g", "package", "-r", "a", dir + "*"},
			[]string{dir + "*"},
			false,
		},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			a, err := newArgs(tc.input...)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NotNil(t, a)
			assert.True(t, stringSliceEq(a.paths, tc.wantPaths))
		})
	}
}

func TestArgsHasArgs(t *testing.T) {
	tcs := []struct {
		input  []string
		args   []string
		want   []string
		wantOk bool
	}{
		{[]string{}, []string{""}, []string{}, false},
		{[]string{}, []string{"-a"}, []string{}, false},
		{[]string{}, []string{"-a", "-a"}, []string{}, false},
		{[]string{"a"}, []string{"-a"}, []string{"a"}, false},
		{[]string{"a"}, []string{"-a", "-a"}, []string{"a"}, false},
		{[]string{"-a"}, []string{"-a"}, []string{}, true},
		{[]string{"-a"}, []string{"-a", "-a"}, []string{}, true},
		{[]string{"-a", "-b"}, []string{"-a"}, []string{"-b"}, true},
		{[]string{"-a", "-c", "-b"}, []string{"-c"}, []string{"-a", "-b"}, true},
		{[]string{"-a", "-c", "-b"}, []string{"-d"}, []string{"-a", "-c", "-b"}, false},
	}
	for i, tc := range tcs {
		tc := tc
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			getPipe(t)
			a, err := newArgs(tc.input...)
			assert.NoError(t, err)
			ok := a.hasArgs(tc.args...)
			assert.True(t, stringSliceEq(a.remaining, tc.want))
			assert.EqualValues(t, tc.wantOk, ok)
		})
	}
}
