package reader_test

import (
	"io"
	"os"
	"path"
	"testing"

	"github.com/alecthomas/assert"
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

func setup(t *testing.T, input []testCase) []string {
	t.Helper()
	dir := t.TempDir()
	ret := make([]string, len(input))
	for i, d := range input {
		name := path.Join(dir, d.name)
		base := path.Dir(name)
		err := os.MkdirAll(base, os.ModePerm)
		assert.NoError(t, err)
		f, err := os.Create(name)
		assert.NoError(t, err)
		f.WriteString(d.content)
		f.Close()
		ret[i] = base
	}

	return ret
}

func inSlice(niddle string, haystack []string) bool {
	for _, s := range haystack {
		if s == niddle {
			return true
		}
	}
	return false
}
