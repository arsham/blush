package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
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
		{name: "colour american", input: []string{"--color"}, colour: true},
		{name: "colour short", input: []string{"-C"}, colour: true},
		{name: "no filename", input: []string{"-h"}, noFilename: true},
		{name: "no filename long", input: []string{"--no-filename"}, noFilename: true},
		{name: "rec", input: []string{"-R"}, recursive: true},
		{name: "ins", input: []string{"-i"}, insensitive: true},
		{name: "ins rec", input: []string{"-i", "-R"}, insensitive: true, recursive: true},
		{name: "rec ins", input: []string{"-R", "-i"}, insensitive: true, recursive: true},
		{name: "rec ins nofile", input: []string{"-R", "-i", "-h"},
			insensitive: true, recursive: true, noFilename: true},
		{name: "nofile rec ins", input: []string{"-h", "-R", "-i"},
			insensitive: true, recursive: true, noFilename: true},
		{name: "nofile rec ins colour", input: []string{"-h", "-R", "-i", "-C"},
			insensitive: true, recursive: true, noFilename: true, colour: true},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			a, err := newArgs(tc.input...)
			if tc.wantErr != nil && err != tc.wantErr {
				t.Errorf("newArgs(%v): err = %v, want %v", tc.input, err, tc.wantErr)
			}
			if err != nil {
				if a != nil {
					t.Errorf("newArgs(%v): a = %v, want nil", tc.input, a)
				}
				return
			}
			if a.colour != tc.colour {
				t.Errorf("a.colour = %t, want %t", a.colour, tc.colour)
			}
			if a.noFilename != tc.noFilename {
				t.Errorf("a.noFilename = %t, want %t", a.noFilename, tc.noFilename)
			}
			if a.recursive != tc.recursive {
				t.Errorf("a.recursive = %t, want %t", a.recursive, tc.recursive)
			}
			if a.insensitive != tc.insensitive {
				t.Errorf("a.insensitive = %t, want %t", a.insensitive, tc.insensitive)
			}
		})
	}
}

func TestArgsPipe(t *testing.T) {
	_, cleanup := getPipe(t)
	defer cleanup()
	a, err := newArgs("something")
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}
	if a == nil {
		t.Fatal("a = nil, want *args")
	}
	if !a.stdin {
		t.Errorf("a.stdin = %t, want %t", a.stdin, true)
	}
}

func TestArgsPaths(t *testing.T) {
	dir, err := ioutil.TempDir("", "blush_main")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err = os.RemoveAll(dir); err != nil {
			t.Error(err)
		}
	}()
	f1, err := ioutil.TempFile(dir, "main")
	if err != nil {
		t.Fatal(err)
	}
	f2, err := ioutil.TempFile(dir, "main")
	if err != nil {
		t.Fatal(err)
	}

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
		{"many prefixes",
			[]string{"--#7ff", "main", "-g", "package", "-r", "a", path.Join(dir, "*")},
			[]string{path.Join(dir, "*")}, false,
		},
		{"many prefixes star",
			[]string{"--#7ff", "main", "-g", "package", "-r", "a", dir + "*"},
			[]string{dir + "*"}, false,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			a, err := newArgs(tc.input...)
			if tc.wantErr {
				if err == nil {
					t.Error("err = nil, want error")
				}
				return
			}
			if a == nil {
				t.Fatal("a = nil, want *args")
			}
			if !stringSliceEq(a.paths, tc.wantPaths) {
				t.Errorf("files(%v): a.paths = %v, want %v", tc.input, a.paths, tc.wantPaths)
			}
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
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			_, cleanup := getPipe(t)
			defer cleanup()
			a, err := newArgs(tc.input...)
			if err != nil {
				t.Fatal(err)
			}
			ok := a.hasArg(tc.args...)
			if !stringSliceEq(a.remaining, tc.want) {
				t.Errorf("a.hasArg(%v, %s): a.remaining = %v, want %v", tc.input, tc.args, a.remaining, tc.want)
			}
			if ok != tc.wantOk {
				t.Errorf("a.hasArg(%v, %s): ok = %v, want %v", tc.input, tc.args, ok, tc.wantOk)
			}
		})
	}
}
