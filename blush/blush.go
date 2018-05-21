package blush

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/arsham/blush/tools"
	"github.com/pkg/errors"
)

// Blush has a slice of given regexp, matching paths, and operation
// configurations.
type Blush struct {
	Args      []Arg
	Paths     []string
	Sensitive bool // case sensitivity
	Recursive bool
	Colouring bool
}

// New returns an error if `a` is empty, or there is no files found. You should
// remove the application name otherwise it will be accounted as an expression.
func New(input string) (*Blush, error) {
	var (
		g       = &Blush{}
		ok      bool
		remains string
	)
	if input == "" {
		return nil, ErrNoInput
	}
	if remains, ok = hasArg(input, "-C"); ok {
		g.Colouring = true
	}
	if remains, ok = hasArg(remains, "-i"); ok {
		g.Sensitive = true
	}
	if remains, ok = hasArg(remains, "-R"); ok {
		g.Recursive = true
	}

	remains, p, err := files(remains)
	if err != nil {
		return nil, errors.Wrap(err, "provided files")
	}
	g.Paths = p

	g.Args = getArgs(remains)
	return g, nil
}

// WriteTo writes matches to w.
func (b Blush) Write(w io.Writer) error {
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

// files starts from the end and removes any file matches it finds and returns
// them.
func files(input string) (remaining string, p []string, err error) {
	var (
		foundOne bool
		counter  int
		ret      []string
	)
	input = strings.Trim(input, " ")
	tokens := strings.Split(input, " ")
	sort.Slice(tokens, func(i, j int) bool {
		return i > j
	})
	for _, t := range tokens {
		if inStringSlice(t, p) {
			continue
		}
		if _, err := os.Stat(t); err == nil {
			foundOne = true
			p = append(p, t)
			counter++
			continue
		}
		if !foundOne {
			return input, nil, fmt.Errorf("%s not found", t)
		}
		ret = append(ret, t)
	}

	//We have reversed it. We need to return back in the same order.
	sort.Slice(ret, func(i, j int) bool {
		return i > j
	})
	remaining = strings.Join(ret, " ")
	return remaining, p, nil
}

func inStringSlice(s string, haystack []string) bool {
	for _, a := range haystack {
		if a == s {
			return true
		}
	}
	return false
}
