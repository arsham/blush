package blush

import (
	"bufio"
	"fmt"
	"io"
)

// Blush has a slice of given regexp, matching paths, and operation
// configuration. If NoCut is true, the unmatched lines are printed as well.
type Blush struct {
	Locator []ColourLocator
	Reader  io.ReadCloser
	NoCut   bool
}

// ColourLocator contains a pair of colour name and corresponding matcher.
type ColourLocator struct {
	Locator
	Colour Colour
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
	if err := b.find(w, b.Reader); err != nil {
		return err
	}
	return nil
}

func (b Blush) find(w io.Writer, file io.Reader) error {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		lineWritten := false
		for _, a := range b.Locator {
			s, ok := a.Find(line, a.Colour)
			if ok {
				fmt.Fprintf(w, "%s\n", s)
				lineWritten = true
			}
		}
		if !lineWritten && b.NoCut {
			fmt.Fprintf(w, "%s\n", line)
		}
	}
	return nil
}
