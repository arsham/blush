package blush_test

import "io"

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

type badWriter struct {
	writeFunc func([]byte) (int, error)
}

func (b *badWriter) Write(p []byte) (int, error) { return b.writeFunc(p) }
