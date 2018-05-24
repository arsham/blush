package blush_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/arsham/blush/blush"
)

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

func TestNewWalkerError(t *testing.T) {
	dirs := []string{"nomansland2987349237"}
	w, err := blush.NewWalker(dirs, false)
	if err == nil {
		t.Error("NewWalker(): err = nil, want error")
	}
	if w != nil {
		t.Errorf("NewWalker(): w = %v, want nil", w)
	}
}

func TestNewWalker(t *testing.T) {
	var (
		c1 = "VJSNS5IeLCtEB"
		c2 = "kkNL8vGNJn"
		c3 = "o6Ubb5Taj"
	)
	input := []testCase{
		{"a/a.txt", c1},
		{"a/b.txt", c2},
		{"ab.txt", c3},
	}

	dirs, cleanup := setup(t, input)
	defer cleanup()
	w, err := blush.NewWalker(dirs, false)
	if err != nil {
		t.Fatalf("NewWalker(): err = %v, want nil", err)
	}
	if w == nil {
		t.Fatal("NewWalker(): w = nil, want *blush.Walker")
	}
	defer func() {
		if err := w.Close(); err != nil {
			t.Error(err)
		}
	}()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(w)
	if err != nil {
		t.Error(err)
	}
	for _, s := range []string{c1, c2, c3} {
		if !strings.Contains(buf.String(), s) {
			t.Errorf("`%s` not found in `%s`", s, buf.String())
		}
	}
}

func TestNewWalkerRecursive(t *testing.T) {
	var (
		c1 = "1JQey4agQ3w9pqg3"
		c2 = "7ToNRMgsOAR6A"
		c3 = "EtOkn9C5zoH0Dla2rF9"
	)
	input := []testCase{
		{"a/a.txt", c1},
		{"a/b.txt", c2},
		{"a/b/c.txt", c3},
	}

	dirs, cleanup := setup(t, input)
	defer cleanup()
	base := path.Join(path.Dir(dirs[0]), "a")
	w, err := blush.NewWalker([]string{base}, true)
	if err != nil {
		t.Fatalf("NewWalker(): err = %v, want nil - %v", err, base)
	}
	if w == nil {
		t.Fatal("NewWalker(): w = nil, want *blush.Walker")
	}
	defer func() {
		if err := w.Close(); err != nil {
			t.Error(err)
		}
	}()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(w)
	if err != nil {
		t.Error(err)
	}
	for _, s := range []string{c1, c2, c3} {
		if !strings.Contains(buf.String(), s) {
			t.Errorf("`%s` not found in `%s`", s, buf.String())
		}
	}
}

func TestNewWalkerNonRecursive(t *testing.T) {
	var (
		c1 = "DRAjfSq2y"
		c2 = "ht3xCIQ"
		c3 = "jPqPoAbMNb"
	)
	input := []testCase{
		{"a/a.txt", c1},
		{"a/b.txt", c2},
		{"a/b/c.txt", c3},
	}

	dirs, cleanup := setup(t, input)
	defer cleanup()
	base := path.Join(path.Dir(dirs[0]), "a")
	w, err := blush.NewWalker([]string{base}, false)
	if err != nil {
		t.Fatalf("NewWalker(): err = %v, want nil - %v", err, base)
	}
	if w == nil {
		t.Fatal("NewWalker(): w = nil, want *blush.Walker")
	}
	defer func() {
		if err := w.Close(); err != nil {
			t.Error(err)
		}
	}()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(w)
	if err != nil {
		t.Error(err)
	}
	for _, s := range []string{c1, c2} {
		if !strings.Contains(buf.String(), s) {
			t.Errorf("`%s` not found in `%s`", s, buf.String())
		}
	}
	if strings.Contains(buf.String(), c3) {
		t.Errorf("`%s` should not be found in `%s`", c3, buf.String())
	}
}
