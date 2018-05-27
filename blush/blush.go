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
