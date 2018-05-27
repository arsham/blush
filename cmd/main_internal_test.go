package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func getPipe(t *testing.T) (*os.File, func()) {
	oldStdin := os.Stdin

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
	os.Stdin = file
	return file, func() {
		os.Stdin = oldStdin
		rmFile()
	}
}

func TestGetReaderPipe(t *testing.T) {
	pipe, cleanup := getPipe(t)
	defer cleanup()
	_, r, err := getReader(nil)
	if err != nil {
		t.Fatal(err)
	}
	if r != pipe {
		t.Errorf("r = %v, want %v", r, pipe)
	}
}

func TestGetReaderNoFiles(t *testing.T) {
	_, r, err := getReader([]string{"/nomansland"})
	if err == nil {
		t.Error("err = nil, want error")
	}
	if r != nil {
		t.Errorf("r = %v, want nil", r)
	}
}

func TestGetReaderNewMultiReaderFromPathsError(t *testing.T) {
	_, r, err := getReader([]string{""})
	if err == nil {
		t.Error("err = nil, want error")
	}
	if r != nil {
		t.Errorf("r = %v, want nil", r)
	}
}

func stringSliceEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestFiles(t *testing.T) {
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
		name          string
		input         []string
		wantRemaining []string
		wantP         []string
		wantErr       bool
	}{
		{"not found", []string{"nowhere"}, []string{}, []string{}, true},
		{"only a file", []string{f1.Name()}, []string{}, []string{f1.Name()}, false},
		{"two files", []string{f1.Name(), f2.Name()}, []string{}, []string{f1.Name(), f2.Name()}, false},
		{"arg between two files", []string{f1.Name(), "-a", f2.Name()}, []string{f1.Name(), "-a", f2.Name()}, []string{}, true},
		{"prefix file", []string{"a", f1.Name()}, []string{"a"}, []string{f1.Name()}, false},
		{"prefix arg file", []string{"-r", f1.Name()}, []string{"-r", f1.Name()}, []string{}, true},
		{"file matches but is an argument", []string{"-r", f1.Name(), f2.Name()}, []string{"-r", f1.Name()}, []string{f2.Name()}, false},
		{
			"star dir",
			[]string{path.Join(dir, "*")},
			[]string{},
			[]string{path.Join(dir, "*")},
			false,
		},
		{
			"stared dir",
			[]string{dir + "*"},
			[]string{},
			[]string{dir + "*"},
			false,
		},
		{
			"many prefixes",
			[]string{"--#7ff", "main", "-g", "package", "-r", "a", path.Join(dir, "*")},
			[]string{"--#7ff", "main", "-g", "package", "-r", "a"},
			[]string{path.Join(dir, "*")},
			false,
		},
		{
			"many prefixes star",
			[]string{"--#7ff", "main", "-g", "package", "-r", "a", dir + "*"},
			[]string{"--#7ff", "main", "-g", "package", "-r", "a"},
			[]string{dir + "*"},
			false,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			remaining, p, err := paths(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Error("err = nil, want error")
				}
				return
			}
			if !stringSliceEq(remaining, tc.wantRemaining) {
				t.Errorf("files(%v): remaining = %v, want %v", tc.input, remaining, tc.wantRemaining)
			}
			if !stringSliceEq(p, tc.wantP) {
				t.Errorf("files(%v): p = %v, want %v", tc.input, p, tc.wantP)
			}
		})
	}
}

func TestHasArgs(t *testing.T) {
	tcs := []struct {
		input  []string
		arg    string
		want   []string
		wantOk bool
	}{
		{[]string{}, "", []string{}, false},
		{[]string{}, "a", []string{}, false},
		{[]string{"a"}, "-a", []string{"a"}, false},
		{[]string{"-a"}, "-a", []string{}, true},
		{[]string{"-a", "-b"}, "-a", []string{"-b"}, true},
		{[]string{"-a", "-c", "-b"}, "-c", []string{"-a", "-b"}, true},
		{[]string{"-a", "-c", "-b"}, "-d", []string{"-a", "-c", "-b"}, false},
	}
	for i, tc := range tcs {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			got, ok := hasArg(tc.input, tc.arg)
			if !stringSliceEq(got, tc.want) {
				t.Errorf("hasArg(%v, %s): got = %v, want %v", tc.input, tc.arg, got, tc.want)
			}
			if ok != tc.wantOk {
				t.Errorf("hasArg(%v, %s): ok = %v, want %v", tc.input, tc.arg, ok, tc.wantOk)
			}
		})
	}
}
