// Package blush reads from a given io.Reader line by line and looks for
// patterns.
//
// Blush struct has a Reader property which can be Stdin in case of it being
// shell's pipe, or any type that implements io.ReadCloser. If NoCut is set to
// true, it will show all lines despite being not matched.
//
// The hex number should be in 3 or 6 part format (#aaaaaa or #aaa) and each
// part will be translated to a number value between 0 and 255 when creating the
// Colour instance. If any of hex parts are not between 00 and ff, it creates
// the DefaultColour value.
//
// Important Notes
//
// The Read() method could be slow in case of huge inspections. It is
// recommended to avoid it and use WriteTo() instead; io.Copy() can take care of
// that for you.
//
// When WriteTo() is called with an unavailable or un-writeable writer, there
// will be no further checks until it tries to write into it. If the Write
// encounters any errors regarding writes, it will return the amount if writes
// and stops its search.
//
// There always will be a newline after each read.
package blush

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sync"
)

// Blush reads from Reader and matches against all Finders. If NoCut is true,
// any unmatched lines are printed as well.
type Blush struct {
	Finders []Finder
	Reader  io.ReadCloser
	NoCut   bool
	closed  bool

	once sync.Once
	buf  *bytes.Buffer
}

// Read will search the whole tree and keeps the results in a buffer and uses
// the buffer to write to p. Any error that happens during the construction of
// this buffer will be returned immediately and closes the object for further
// reads.
func (b *Blush) Read(p []byte) (n int, err error) {
	if b.closed {
		return 0, ErrClosed
	}
	b.once.Do(func() {
		b.buf = new(bytes.Buffer)
		if _, er := b.WriteTo(b.buf); er != nil {
			err = er
		}
	})
	if err != nil {
		b.closed = true
		return
	}
	return b.buf.Read(p)
}

// WriteTo writes matches to w. It returns an error if the writer is nil or
// there are not paths defined or there is no files found in the Reader. Please
// read documentations for ErrNoWriter.
func (b *Blush) WriteTo(w io.Writer) (int64, error) {
	if w == nil {
		return 0, ErrNoWriter
	}
	if b.Reader == nil {
		return 0, ErrNoReader
	}
	return b.search(w), nil
}

// Close closes the reader and returns whatever error it returns.
func (b *Blush) Close() error {
	b.closed = true
	return b.Reader.Close()
}

func (b *Blush) search(w io.Writer) int64 {
	var total int
	scanner := bufio.NewScanner(b.Reader)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		var foundStr string
		line := scanner.Text()
		for _, a := range b.Finders {
			if s, ok := a.Find(line); ok {
				line = s
				foundStr = line
			}
		}
		if foundStr != "" {
			line = foundStr
		}
		if foundStr != "" || b.NoCut {
			if n, err := fmt.Fprintf(w, "%s\n", line); err != nil {
				return int64(n)
			}
			total += len(line) + 1 // new-line of each line is added here
		}
	}
	return int64(total)
}
