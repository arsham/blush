package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/arsham/blush/blush"
	"github.com/arsham/blush/internal/reader"
)

// Main reads the provided arguments from the command line and creates a
// blush.Blush instance. It then uses io.Copy() to write to standard output.
func Main() {
	b, err := GetBlush(os.Args)
	switch err {
	case nil:
	case errShowHelp:
		fmt.Println(Usage)
		return
	default:
		log.Fatalf("%s\n%s", err, Help)
		return // this return statement should be here to support tests.
	}
	defer func() {
		if err := b.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	sig := make(chan os.Signal, 1)
	WaitForSignal(sig, os.Exit)
	if _, err := io.Copy(os.Stdout, b); err != nil {
		log.Fatal(err)
	}
}

// GetBlush returns an error if no arguments are provided or it can't find all
// the passed files in the input.
//
// Note
//
// The first argument will be dropped as it will be the application's name.
func GetBlush(input []string) (*blush.Blush, error) {
	var (
		ok         bool
		noCut      bool
		noFileName bool
		remaining  []string
		err        error
		r          io.ReadCloser
	)
	if len(input) == 1 {
		return nil, ErrNoInput
	}
	if remaining, ok = hasArg(input[1:], "--help"); ok {
		return nil, errShowHelp
	}

	remaining, r, err = getReader(remaining)
	if err != nil {
		return nil, err
	}
	if remaining, ok = hasArg(remaining, "-C", "--colour", "--color"); ok {
		noCut = true
	}
	if remaining, ok = hasArg(remaining, "-h", "--no-filename"); ok {
		noFileName = true
	}
	finders := getFinders(remaining)
	return &blush.Blush{
		Finders:      finders,
		Reader:       r,
		NoCut:        noCut,
		WithFileName: !noFileName,
	}, nil
}

// getReader returns os.Stdin if it is piped to the program, otherwise looks for
// files.
func getReader(input []string) (remaining []string, r io.ReadCloser, err error) {
	var (
		recursive bool
		ok        bool
	)
	if remaining, ok = hasArg(input, "-R"); ok {
		recursive = true
	}
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		return remaining, os.Stdin, nil
	}

	remaining, p, err := paths(input)
	if err != nil {
		return nil, nil, err
	}
	w, err := reader.NewMultiReader(reader.WithPaths(p, recursive))
	if err != nil {
		return nil, nil, err
	}
	return remaining, w, nil
}

// paths starts from the end of the slice and removes any paths/globs/files it
// finds and returns them in p.
func paths(input []string) (remaining []string, p []string, err error) {
	var (
		foundOne bool
		counter  int
		ret      []string
	)
	// going backwards from the end.
	sort.SliceStable(input, func(i, j int) bool {
		return i > j
	})

	// I don't like this label, but if we replace the `switch` statement with a
	// regular if-then-else clause, it gets ugly and doesn't show its
	// intentions. Order of cases in this switch matters.
LOOP:
	for i, t := range input {
		t = strings.Trim(t, " ")
		if t == "" || inStringSlice(t, p) {
			continue
		}

		m, err := filepath.Glob(t)
		if err != nil {
			return nil, nil, err
		}
		switch {
		case len(input) > i+1 && strings.HasPrefix(input[i+1], "-"):
			// In this case, the previous input was a flag argument, therefore
			// it might have been a colouring command. That is why we are
			// ignoring this item.
			ret = append(ret, input[i:]...)
			break LOOP
		case len(m) > 0:
			foundOne = true
			p = append(p, t)
			counter++
		case foundOne:
			// there is already a pattern found so we stop here.
			ret = append(ret, input[i:]...)
			break LOOP
		}
	}
	if !foundOne {
		return input, nil, ErrNoFilesFound
	}

	// We have reversed it. We need to return back in the same order.
	sort.SliceStable(ret, func(i, j int) bool {
		return i > j
	})
	// to keep the original user's preference.
	sort.SliceStable(p, func(i, j int) bool {
		return i > j
	})
	return ret, p, nil
}

func inStringSlice(s string, haystack []string) bool {
	for _, a := range haystack {
		if a == s {
			return true
		}
	}
	return false
}

func getFinders(input []string) []blush.Finder {
	var (
		lastColour  string
		ret         []blush.Finder
		insensitive bool
		ok          bool
	)
	if input, ok = hasArg(input, "-i"); ok {
		insensitive = true
	}
	for _, token := range input {
		if strings.HasPrefix(token, "-") {
			lastColour = strings.TrimLeft(token, "-")
			continue
		}
		a := blush.NewLocator(lastColour, token, insensitive)
		ret = append(ret, a)
	}
	return ret
}

// hasArg removes any occurring `args` argument and returns the remaining.
func hasArg(input []string, args ...string) (remains []string, found bool) {
	remains = input[:]
	for _, arg := range args {
		for i, a := range input {
			if a == arg {
				remains = append(remains[:i], remains[i+1:]...)
				found = true
				if len(remains) == 0 {
					return remains, found
				}
			}
		}
	}
	return remains, found
}
