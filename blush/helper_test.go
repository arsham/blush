package blush_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/arsham/blush/blush"
)

// this file contains helpers for all tests in this package.

// In the testdata folder, there are three files. In each file there are 1 ONE,
// 2 TWO, 3 THREE and 4 FOURs. There is a line containing `LEAVEMEHERE` which
// does not collide with any of these numbers.

var leaveMeHere = "LEAVEMEHERE"

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

// this function reads everything in `w` and returns the length of contents.
// Particularly useful for comparing the WriteTo and Write returning length.
func walkerLen(walker *blush.Walker) (int64, error) {
	w, err := blush.NewWalker(walker.Paths, walker.Recursive)
	if err != nil {
		return 0, err
	}
	buf := new(bytes.Buffer)
	return buf.ReadFrom(w)
}
