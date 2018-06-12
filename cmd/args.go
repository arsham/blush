package cmd

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/arsham/blush/blush"
)

// Note that hasArgs, setFinders and setPaths methods of args are designed to
// shrink the input as they go. Therefore the order of calls matters in some
// cases.
type args struct {
	colour      bool
	noFilename  bool
	recursive   bool
	insensitive bool
	stdin       bool
	paths       []string
	matches     []string
	remaining   []string
	finders     []blush.Finder
}

func newArgs(input ...string) (*args, error) {
	a := &args{
		matches:   make([]string, 0),
		remaining: input,
	}
	if a.hasArgs("--help") {
		return nil, errShowHelp
	}
	a.recursive = a.hasArgs("-R")
	a.colour = a.hasArgs("-C", "--colour", "--color")
	a.noFilename = a.hasArgs("-h", "--no-filename")
	a.insensitive = a.hasArgs("-i")

	if stat, _ := os.Stdin.Stat(); (stat.Mode() & os.ModeCharDevice) == 0 {
		a.stdin = true
	} else if err := a.setPaths(); err != nil {
		return nil, err
	}
	a.setFinders()
	return a, nil
}

// hasArgs removes any occurring `args` argument.
func (a *args) hasArgs(args ...string) (found bool) {
	remains := a.remaining[:]
LOOP:
	for _, arg := range args {
		for i, ar := range a.remaining {
			if ar == arg {
				remains = append(remains[:i], remains[i+1:]...)
				found = true
				if len(remains) == 0 {
					break LOOP
				}
			}
		}
	}
	a.remaining = remains
	return found
}

// setPaths starts from the end of the slice and removes any paths/globs/files
// it finds and put them in the paths property.
func (a *args) setPaths() error {
	var (
		foundOne bool
		counter  int
		p, ret   []string
		input    = a.remaining
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
			return err
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
		return ErrNoFilesFound
	}

	// We have reversed it. We need to return back in the same order.
	sort.SliceStable(ret, func(i, j int) bool {
		return i > j
	})
	// to keep the original user's preference.
	sort.SliceStable(p, func(i, j int) bool {
		return i > j
	})
	a.remaining = ret
	a.paths = p
	return nil
}

func (a *args) setFinders() {
	var lastColour string
	a.finders = make([]blush.Finder, 0)
	for _, token := range a.remaining {
		if strings.HasPrefix(token, "-") {
			lastColour = strings.TrimLeft(token, "-")
			continue
		}
		l := blush.NewLocator(lastColour, token, a.insensitive)
		a.finders = append(a.finders, l)
	}
}

func inStringSlice(s string, haystack []string) bool {
	for _, a := range haystack {
		if a == s {
			return true
		}
	}
	return false
}
