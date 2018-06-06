package reader_test

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

// this file contains helpers for all tests in this package.

type nopCloser struct {
	io.Reader
	closeFunc func() error
}

func (n nopCloser) Close() error { return n.closeFunc() }

type testCase struct {
	name    string
	content string
}

func setup(t *testing.T, input []testCase) ([]string, func()) {
	dir, err := ioutil.TempDir("", "blush_walker")
	if err != nil {
		t.Fatal(err)
	}
	ret := make([]string, len(input))
	for i, d := range input {
		name := path.Join(dir, d.name)
		base := path.Dir(name)
		err = os.MkdirAll(base, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
		f, err := os.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		f.WriteString(d.content)
		f.Close()
		ret[i] = base
	}
	return ret, func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Error(err)
		}
	}
}

func inStringSlice(niddle string, haystack []string) bool {
	for _, s := range haystack {
		if s == niddle {
			return true
		}
	}
	return false
}
