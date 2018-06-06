package blush

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/arsham/blush/internal/reader"
	"github.com/pkg/errors"
)

const (
	// Separator string between name of the reader and the contents.
	Separator = ": "
)

// Blush reads from Reader and matches against all Finders. If NoCut is true,
// any unmatched lines are printed as well. If WithFileName is true, blush will
// write the filename before it writes the output.
type Blush struct {
	Finders      []Finder
	Reader       io.ReadCloser
	NoCut        bool
	WithFileName bool
	closed       bool

	once sync.Once // used in Read() for loading everything in to the buffer.
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
		if e := b.Close(); e != nil {
			err = errors.Wrap(err, e.Error())
		}
		return 0, err
	}
	return b.buf.Read(p)
}

// WriteTo writes matches to w. It returns an error if the writer is nil or
// there are not paths defined or there is no files found in the Reader.
func (b *Blush) WriteTo(w io.Writer) (int64, error) {
	if w == nil {
		return 0, ErrNoWriter
	}
	if b.Reader == nil {
		return 0, reader.ErrNoReader
	}
	return b.search(w)
}

// Close closes the reader and returns whatever error it returns.
func (b *Blush) Close() error {
	b.closed = true
	return b.Reader.Close()
}

func (b *Blush) search(w io.Writer) (int64, error) {
	var total int
	scanner := bufio.NewScanner(b.Reader)
	scanner.Split(bufio.ScanLines)
	max := bufio.MaxScanTokenSize * 120
	buf := make([]byte, max)
	scanner.Buffer(buf, max)
	for scanner.Scan() {
		line := scanner.Text()
		n, err := b.processLine(w, line)
		if err != nil {
			return int64(n), err
		}
		total += n
	}
	return int64(total), nil
}

func (b *Blush) processLine(w io.Writer, line string) (int, error) {
	var total int
	str, ok := lookInto(b.Finders, line)
	if ok || b.NoCut {
		var prefix string
		if b.WithFileName {
			prefix = fileName(b.Reader)
			total += len(prefix)
		}
		if n, err := fmt.Fprintf(w, "%s%s\n", prefix, str); err != nil {
			return n, err
		}
		total += len(str) + 1 // new-line is added here (\n above)
	}
	return total, nil
}

// lookInto returns a new decorated line if any of the Finders decorate it, or
// the given line as it is.
func lookInto(f []Finder, line string) (string, bool) {
	var found bool
	for _, a := range f {
		if s, ok := a.Find(line); ok {
			line = s
			found = true
		}
	}
	return line, found
}

// fileName returns an empty string if it could not query the fileName from r.
func fileName(r io.Reader) string {
	type namer interface {
		Name() string
	}
	if o, ok := r.(namer); ok {
		return o.Name() + Separator
	}
	return ""
}
