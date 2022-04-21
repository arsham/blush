package cmd

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/alecthomas/assert"
)

func getPipe(t *testing.T) *os.File {
	t.Helper()
	oldStdin := os.Stdin

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
	os.Stdin = file
	t.Cleanup(func() {
		os.Stdin = oldStdin
		rmFile()
	})
	return file
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
			a, err := newArgs(tc.input...)
			if tc.wantErr {
				if err == nil {
					t.Error("err = nil, want error")
				}
				return
			}
			if !stringSliceEq(a.remaining, tc.wantRemaining) {
				t.Errorf("files(%v): a.remaining = %v, want %v", tc.input, a.remaining, tc.wantRemaining)
			}
			if !stringSliceEq(a.paths, tc.wantP) {
				t.Errorf("files(%v): a.paths = %v, want %v", tc.input, a.paths, tc.wantP)
			}
		})
	}
}
