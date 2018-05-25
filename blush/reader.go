package blush

import (
	"io"
	"os"

	"github.com/arsham/blush/internal/tools"
	"github.com/pkg/errors"
)

// MultiReadCloser is an io.MultiReader which also implements io.Closer.
type MultiReadCloser struct {
	r  []io.ReadCloser
	mr io.Reader
}

// NewMultiReadCloser returns an error if any of given files are not found. It
// ignores any files that cannot be read or opened.
func NewMultiReadCloser(paths []string, recursive bool) (*MultiReadCloser, error) {
	var r []io.Reader
	files, err := tools.Files(recursive, paths...)
	if err != nil {
		return nil, errors.Wrap(err, "NewMultiReadCloser")
	}
	w := &MultiReadCloser{
		r: make([]io.ReadCloser, 0),
	}
	for _, f := range files {
		file, err := os.Open(f)
		if err != nil {
			continue
		}

		w.r = append(w.r, file)
		r = append(r, file)
	}
	w.mr = io.MultiReader(r...)
	return w, nil
}

// Read returns an io.MultiReader created by all files found in the Paths.
func (w *MultiReadCloser) Read(b []byte) (n int, err error) {
	return w.mr.Read(b)
}

// Close closes all files opened by this MultiReadCloser.
func (w *MultiReadCloser) Close() (err error) {
	for _, r := range w.r {
		err = r.Close()
	}
	return
}
