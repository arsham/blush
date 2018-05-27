package blush

import (
	"io"
	"os"

	"github.com/arsham/blush/internal/tools"
	"github.com/pkg/errors"
)

// MultiReader is an io.MultiReader which also implements io.Closer.
type MultiReader struct {
	r  []io.ReadCloser
	mr io.Reader
}

// NewMultiReader keeps all readers in memory in order to be able to close
// them when the Close() method is called.
func NewMultiReader(readers ...io.ReadCloser) *MultiReader {
	r := make([]io.Reader, len(readers))
	mr := &MultiReader{
		r: make([]io.ReadCloser, len(readers)),
	}
	for i, rd := range readers {
		mr.r[i] = rd
		r[i] = rd
	}
	mr.mr = io.MultiReader(r...)
	return mr
}

// NewMultiReaderFromPaths returns an error if any of given files are not found.
// It ignores any files that cannot be read or opened.
func NewMultiReaderFromPaths(paths []string, recursive bool) (*MultiReader, error) {
	readers := make([]io.ReadCloser, 0)
	files, err := tools.Files(recursive, paths...)
	if err != nil {
		return nil, errors.Wrap(err, "NewMultiReaderFromPaths")
	}
	for _, f := range files {
		file, err := os.Open(f)
		if err != nil {
			continue
		}
		readers = append(readers, file)
	}
	return NewMultiReader(readers...), nil
}

// Read returns an io.MultiReader created by all files found in the Paths.
func (w *MultiReader) Read(b []byte) (n int, err error) {
	return w.mr.Read(b)
}

// Close closes all files opened by this MultiReader.
func (w *MultiReader) Close() (err error) {
	for _, r := range w.r {
		err = r.Close()
	}
	return
}
