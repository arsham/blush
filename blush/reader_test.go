package blush_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"path"
	"strings"
	"testing"

	"github.com/arsham/blush/blush"
)

func TestWithReader(t *testing.T) {
	m, err := blush.NewMultiReader(blush.WithReader("name", nil))
	if err == nil {
		t.Error("err = nil, want error")
	}
	if m != nil {
		t.Errorf("m = %v, want nil", m)
	}

	r := ioutil.NopCloser(new(bytes.Buffer))
	m, err = blush.NewMultiReader(blush.WithReader("name", r))
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}
	if m == nil {
		t.Error("m = nil, want *blush.MultiReader")
	}

	m, err = blush.NewMultiReader(blush.WithReader("", r))
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}
	if m == nil {
		t.Error("m = nil, want *blush.MultiReader")
	}
}

func TestWithReaderMultipleReadersClose(t *testing.T) {
	var called []string
	input1 := "afmBEswIRYosG7"
	input2 := "UbMFeIFjvAhdA3sdT"
	r1 := nopCloser{
		Reader: bytes.NewBufferString(input1),
		closeFunc: func() error {
			called = append(called, "r1")
			return nil
		},
	}
	r2 := nopCloser{
		Reader: bytes.NewBufferString(input2),
		closeFunc: func() error {
			called = append(called, "r2")
			return nil
		},
	}
	m, err := blush.NewMultiReader(blush.WithReader("r1", r1), blush.WithReader("r2", r2))
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}
	if m == nil {
		t.Fatal("m = nil, want *blush.MultiReader")
	}

	b := make([]byte, 100)
	_, err = m.Read(b)

	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}
	if string(bytes.Trim(b, "\x00")) != input1 {
		t.Errorf("b = %s, want %s", b, input1)
	}
	_, err = m.Read(b)
	if !inStringSlice("r1", called) {
		t.Error("m.Close() didn't close r1")
	}
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}
	if string(bytes.Trim(b, "\x00")) != input2 {
		t.Errorf("b = %s, want %s", b, input2)
	}
	_, err = m.Read(b)
	if err != io.EOF {
		t.Errorf("err = %v, want io.EOF", err)
	}
	if !inStringSlice("r2", called) {
		t.Error("m.Close() didn't close r2")
	}
}

func TestWithReaderMultipleReadersError(t *testing.T) {
	r := nopCloser{
		Reader: new(bytes.Buffer),
		closeFunc: func() error {
			return nil
		},
	}
	m, err := blush.NewMultiReader(blush.WithReader("r", r), nil)
	if err == nil {
		t.Error("err = nil, want error")
	}
	if m != nil {
		t.Errorf("m = %v, want nil", m)
	}
}

func TestWithPathsError(t *testing.T) {
	tcs := []struct {
		name  string
		input []string
	}{
		{"nil", nil},
		{"empty", []string{}},
		{"empty string", []string{""}},
		{"not found", []string{"nomansland2987349237"}},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			input := blush.WithPaths(tc.input, true)
			m, err := blush.NewMultiReader(input)
			if err == nil {
				t.Error("NewMultiReader(WithPaths): err = nil, want error")
			}
			if m != nil {
				t.Errorf("NewMultiReader(WithPaths): m = %v, want nil", m)
			}
		})
	}
}

func TestNewMultiReaderWithPaths(t *testing.T) {
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
	m, err := blush.NewMultiReader(blush.WithPaths(dirs, false))
	if err != nil {
		t.Errorf("NewMultiReader(WithPaths(): err = %v, want nil", err)
	}
	if m == nil {
		t.Fatal("NewMultiReader(WithPaths(): m = nil, want *blush.MultiReader")
	}
	if err = m.Close(); err != nil {
		t.Error(err)
	}
}

func TestMultiReaderReadOneReader(t *testing.T) {
	input := "3wAgvZ4bSfQYawl5OEEg"
	r := ioutil.NopCloser(bytes.NewBufferString(input))
	m, err := blush.NewMultiReader(blush.WithReader("r", r))
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}
	if m == nil {
		t.Fatal("m = nil, want *blush.MultiReader")
	}
	b := make([]byte, len(input))
	n, err := m.Read(b)
	if err != nil {
		t.Errorf("Read(): err = %v, want nil", err)
	}
	if n != len(input) {
		t.Errorf("Read(): n = %d, want %d", n, len(input))
	}
	if string(b) != input {
		t.Errorf("b = %s, want %s", b, input)
	}

	n, err = m.Read(b)
	if err != io.EOF {
		t.Errorf("Read(): err = %v, want io.EOF", err)
	}
	if n != 0 {
		t.Errorf("Read(): n = %d, want %d", n, 0)
	}
}

func TestMultiReaderReadZeroBytes(t *testing.T) {
	input := "3wAgvZ4bSfQYawl5OEEg"
	r := ioutil.NopCloser(bytes.NewBufferString(input))
	m, err := blush.NewMultiReader(blush.WithReader("r", r))
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}
	if m == nil {
		t.Fatal("m = nil, want *blush.MultiReader")
	}
	b := make([]byte, 0)
	n, err := m.Read(b)
	if err != nil {
		t.Errorf("Read(): err = %v, want nil", err)
	}
	if n != 0 {
		t.Errorf("Read(): n = %d, want %d", n, 0)
	}
	if string(b) != "" {
		t.Errorf("b = %s, want %s", b, "")
	}
}

func TestMultiReaderReadOneReaderMoreSpace(t *testing.T) {
	input := "3wAgvZ4bSfQYawl5OEEg"
	r := ioutil.NopCloser(bytes.NewBufferString(input))
	m, err := blush.NewMultiReader(blush.WithReader("r", r))
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}
	if m == nil {
		t.Fatal("m = nil, want *blush.MultiReader")
	}
	b := make([]byte, len(input)+1)
	n, err := m.Read(b)
	if err != nil {
		t.Errorf("Read(): err = %v, want nil", err)
	}
	if n != len(input) {
		t.Errorf("Read(): n = %d, want %d", n, len(input))
	}
	if string(bytes.Trim(b, "\x00")) != input {
		t.Errorf("b = %s, want %s", b, input)
	}
}

func TestMultiReaderReadMultipleReaders(t *testing.T) {
	input := []string{"P5tyugWXFn", "b8YbUO7pMX3G8j4Bi"}
	r1 := ioutil.NopCloser(bytes.NewBufferString(input[0]))
	r2 := ioutil.NopCloser(bytes.NewBufferString(input[1]))
	m, err := blush.NewMultiReader(
		blush.WithReader("r1", r1),
		blush.WithReader("r2", r2),
	)
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}
	if m == nil {
		t.Fatal("m = nil, want *blush.MultiReader")
	}
	tcs := []struct {
		name    string
		b       []byte
		wantErr error
		wantLen int
		wantOut string
	}{
		{"r1", make([]byte, len(input[0])), nil, len(input[0]), input[0]},
		{"r2", make([]byte, len(input[1])), nil, len(input[1]), input[1]},
		{"nothing left", make([]byte, 10), io.EOF, 0, ""},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			n, err := m.Read(tc.b)
			if err != tc.wantErr {
				t.Errorf("Read(): err = %v, want %v", err, tc.wantErr)
			}
			if n != tc.wantLen {
				t.Errorf("Read(): n = %d, want %d", n, tc.wantLen)
			}
			if string(bytes.Trim(tc.b, "\x00")) != tc.wantOut {
				t.Errorf("tc.b = `%b`, want `%b`", tc.b, []byte(tc.wantOut))
			}
		})
	}
}

func TestMultiReaderNames(t *testing.T) {
	input := []string{"Mw0mxekLYOpXaKl8PVT", "1V5MjHUXYTPChW"}
	r1 := ioutil.NopCloser(bytes.NewBufferString(input[0]))
	r2 := ioutil.NopCloser(bytes.NewBufferString(input[1]))
	m, err := blush.NewMultiReader(
		blush.WithReader("r1", r1),
		blush.WithReader("r2", r2),
	)
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}
	if m == nil {
		t.Fatal("m = nil, want *blush.MultiReader")
	}
	b := make([]byte, 100)
	tcs := []struct {
		name    string
		wantErr error
	}{
		{"r1", nil},
		{"r2", nil},
		{"", io.EOF},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			_, err := m.Read(b)
			if err != tc.wantErr {
				t.Errorf("m.Read(): err = %v, want %v", err, tc.wantErr)
			}
			if m.String() != tc.name {
				t.Errorf("m.String() = %s, want %s", m.String(), tc.name)
			}
		})
	}
}

func TestNewMultiReaderWithPathsRead(t *testing.T) {
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
	w, err := blush.NewMultiReader(blush.WithPaths(dirs, false))
	if err != nil {
		t.Fatalf("NewMultiReaderFromPaths(): err = %v, want nil", err)
	}
	if w == nil {
		t.Fatal("NewMultiReaderFromPaths(): w = nil, want *blush.MultiReader")
	}
	defer func() {
		if err = w.Close(); err != nil {
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

func TestNewMultiReaderRecursive(t *testing.T) {
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
	w, err := blush.NewMultiReader(blush.WithPaths([]string{base}, true))
	if err != nil {
		t.Fatalf("NewMultiReaderFromPaths(): err = %v, want nil - %v", err, base)
	}
	if w == nil {
		t.Fatal("NewMultiReaderFromPaths(): w = nil, want *blush.MultiReader")
	}
	defer func() {
		if err = w.Close(); err != nil {
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

func TestNewMultiReaderNonRecursive(t *testing.T) {
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
	w, err := blush.NewMultiReader(blush.WithPaths([]string{base}, false))
	if err != nil {
		t.Fatalf("NewMultiReaderFromPaths(): err = %v, want nil - %v", err, base)
	}
	if w == nil {
		t.Fatal("NewMultiReaderFromPaths(): w = nil, want *blush.MultiReader")
	}
	defer func() {
		if err = w.Close(); err != nil {
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
