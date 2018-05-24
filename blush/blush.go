package blush

import (
	"bufio"
	"fmt"
	"io"
)

// Blush has a slice of given regexp, matching paths, and operation
// configuration. If NoCut is true, the unmatched lines are printed as well.
type Blush struct {
	Locator []Locator
	Reader  io.ReadCloser
	NoCut   bool
}

// WriteTo writes matches to w. It returns an error if the writer is nil or
// there are not paths defined or there is no files found in the Reader.
func (b Blush) Write(w io.Writer) error {
	if w == nil {
		return ErrNoWriter
	}
	if b.Reader == nil {
		return ErrNoInput
	}
	b.find(w, b.Reader)
	return nil
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
		for _, a := range b.Locator {
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
}
