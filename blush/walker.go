package blush

import (
	"io"
	"os"

	"github.com/arsham/blush/tools"
	"github.com/pkg/errors"
)

// Walker implements io.ReadCloser. It walks into the Paths and outputs the
// lines of each files.
type Walker struct {
	Paths     []string
	Recursive bool
	r         []io.ReadCloser
	mr        io.Reader
}

// NewWalker returns an error if any of given files are not found. It ignores
// any files that cannot be read or opened.
func NewWalker(paths []string, recursive bool) (*Walker, error) {
	w := &Walker{Paths: paths, Recursive: recursive}
	files, err := tools.Files(w.Recursive, w.Paths...)
	if err != nil {
		return nil, errors.Wrap(err, "NewWalker")
	}
	w.r = make([]io.ReadCloser, 0)
	var r []io.Reader
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
func (w *Walker) Read(b []byte) (n int, err error) {
	return w.mr.Read(b)
}

// Close closes all files opened by this Walker.
func (w *Walker) Close() (err error) {
	for _, r := range w.r {
		err = r.Close()
	}
	return
}
