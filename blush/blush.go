package blush

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/arsham/blush/tools"
	"github.com/pkg/errors"
)

// Blush has a slice of given regexp, matching paths, and operation
// configuration.
type Blush struct {
	Args        []Arg
	Paths       []string
	Insensitive bool
	Recursive   bool
	Colouring   bool
}

// Arg contains a pair of colour name and corresponding matcher.
type Arg struct {
	Colour Colour
	Find   Locator
}

// WriteTo writes matches to w. It returns an error if the writer is nil or
// there are not paths defined or there is no files found in the Paths.
func (b Blush) Write(w io.Writer) error {
	if w == nil {
		return ErrNoWriter
	}
	if b.Paths == nil {
		return ErrNoFiles
	}

	files, err := tools.Files(b.Recursive, b.Paths...)
	if err != nil {
		return errors.Wrap(err, "write")
	}

	for _, f := range files {
		if err := b.find(w, f); err != nil {
			return errors.Wrap(err, f)
		}
	}
	return nil
}

func (b Blush) find(w io.Writer, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, path)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		lineWritten := false
		for _, a := range b.Args {
			s, ok := a.Find.Find(line, a.Colour)
			if ok {
				fmt.Fprintf(w, "%s\n", s)
				lineWritten = true
			}
		}
		if !lineWritten && b.Colouring {
			fmt.Fprintf(w, "%s\n", line)
		}
	}
	return nil
}
