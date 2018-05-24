// Package blush reads from a given io.Reader line by line and looks for
// patterns.
//
// Blush struct has a Reader property which can be Stdin in case of it being
// shell's pipe, or any type that implements io.ReadCloser. If NoCut is set to
// true, it will show all lines even if they don't match.
package blush

import (
	"bufio"
	"fmt"
	"io"
	"sync"
	"time"
)

// Blush has a slice of given regexp, matching paths, and operation
// configuration. If NoCut is true, the unmatched lines are printed as well.
type Blush struct {
	Finders []Finder
	Reader  io.ReadCloser
	NoCut   bool
	once    sync.Once
	res     chan results
}

// WriteTo writes matches to w. It returns an error if the writer is nil or
// there are not paths defined or there is no files found in the Reader.
func (b Blush) WriteTo(w io.Writer) (n int64, err error) {
	if w == nil {
		return 0, ErrNoWriter
	}
	if b.Reader == nil {
		return 0, ErrNoInput
	}
	n = b.find(w)
	return
}

// Close closes the reader and returns whatever error it returns.
func (b Blush) Close() error {
	return b.Reader.Close()
}

func (b Blush) find(w io.Writer, file io.Reader) {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		var foundStr string
		line := scanner.Text()
		n += int64(len(line)) + 1 // new-line of each line
		for _, a := range b.Finders {
			if s, ok := a.Find(line); ok {
				line = s
				foundStr = line
			}
		}
		if foundStr != "" {
			fmt.Fprintf(w, "%s\n", foundStr)
		} else if b.NoCut {
			fmt.Fprintf(w, "%s\n", line)
		}
	}
	return
}
