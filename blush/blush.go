package blush

import (
	"bufio"
	"fmt"
	"io"
	"sync"

	"github.com/arsham/blush/internal/reader"
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

	oncePrepare  sync.Once
	onceTransfer sync.Once
	scanner      *bufio.Reader
	readLineCh   chan []byte
	readCh       chan byte
}

// Read creates a goroutine on first invocation to read from the underlying
// reader. It is considerably slower than WriteTo as it reads the bytes one by
// one in order to produce the results, therefore you should use WriteTo
// directly or use io.Copy() on blush.
func (b *Blush) Read(p []byte) (n int, err error) {
	if b.closed {
		return 0, ErrClosed
	}
	b.oncePrepare.Do(func() {
		err = b.prepare()
	})
	b.onceTransfer.Do(func() {
		go b.transfer()
	})
	for n = 0; n < cap(p); n++ {
		select {
		case c, ok := <-b.readCh:
			if !ok {
				return n, io.EOF
			}
			p[n] = c
		}
	}
	return n, err
}

// WriteTo writes matches to w. It returns an error if the writer is nil or
// there are not paths defined or there is no files found in the Reader.
func (b *Blush) WriteTo(w io.Writer) (int64, error) {
	if b.closed {
		return 0, ErrClosed
	}
	b.oncePrepare.Do(func() {
		b.prepare()
	})
	var total int
	if w == nil {
		return 0, ErrNoWriter
	}
	if b.Reader == nil {
		return 0, reader.ErrNoReader
	}
	for line := range b.readLineCh {
		if n, err := fmt.Fprintf(w, "%s", line); err != nil {
			return int64(n), err
		}
		total += len(line)
	}
	return int64(total), nil
}

func (b *Blush) prepare() error {
	if b.Reader == nil {
		return reader.ErrNoReader
	}
	b.scanner = bufio.NewReader(b.Reader)
	b.readLineCh = make(chan []byte, 50)
	b.readCh = make(chan byte, 1000)
	go b.readLines()
	return nil
}

func (b *Blush) decorate(input string) (string, bool) {
	str, ok := lookInto(b.Finders, input)
	if ok || b.NoCut {
		var prefix string
		if b.WithFileName {
			prefix = fileName(b.Reader)
		}
		return prefix + str, true
	}
	return "", false
}

func (b *Blush) readLines() {
	for {
		line, err := b.scanner.ReadString('\n')
		if line, ok := b.decorate(line); ok {
			b.readLineCh <- []byte(line)
		}
		if err != nil {
			break
		}
	}
	close(b.readLineCh)
}

func (b *Blush) transfer() {
	for line := range b.readLineCh {
		for _, c := range line {
			b.readCh <- c
		}
	}
	close(b.readCh)
}

// Close closes the reader and returns whatever error it returns.
func (b *Blush) Close() error {
	b.closed = true
	return b.Reader.Close()
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
