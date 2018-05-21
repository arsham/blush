package blush

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Blush has a slice of given regexp, matching paths, and operation
// configurations.
type Blush struct {
	Args      []Arg
	Paths     []string
	Sensitive bool // case sensitivity
	Recursive bool
}

// New returns an error if `a` is empty, or there is no files found. You should
// remove the application name otherwise it will be accounted as an expression.
func New(input string) (*Blush, error) {
	var ok bool
	if input == "" {
		return nil, ErrNoInput
	}
	remains, p, err := files(input)
	if err != nil {
		return nil, err
	}
	g := &Blush{
		Paths: p,
	}

	if remains, ok = hasArg(remains, "-s"); ok {
		g.Sensitive = true
	}
	if remains, ok = hasArg(remains, "-R"); ok {
		g.Recursive = true
	}

	g.Args = getArgs(remains)
	return g, nil
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
