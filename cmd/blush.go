package cmd

import (
	"log"
	"os"
	"sort"
	"strings"

	"github.com/arsham/blush/blush"
	"github.com/pkg/errors"
)

// These variables are provided to support the tests.
var (
	FatalErr = func(s string) {
		log.Fatal(s)
	}
)

// Main reads the provided arguments from the command line and creates a
// blush.Blush instance. If there is any error, it will terminate the
// application with error code `1`, otherwise it calls the Write() method of
// Blush and exits with `0`.
func Main() {
	b, err := GetBlush(os.Args)
	if err != nil {
		FatalErr(err.Error())
		return
	}

	b.Write(os.Stdout)
}

// GetBlush returns an error if no arguments are provided or it can't find all
// the passed files. Files should be last arguments, otherwise they are counted
// as matching strings. We are not using the usual flag package because it
// cannot handle variables in the args.
func GetBlush(input []string) (b *blush.Blush, err error) {
	var ok bool
	b = &blush.Blush{}
	if len(input) == 1 {
		return nil, ErrNoInput
	}
	remains, p, err := files(input[1:])
	if err != nil {
		return nil, err
	}
	if remains, ok = hasArg(remains, "-C"); ok {
		b.Colouring = true
	}
	if remains, ok = hasArg(remains, "-i"); ok {
		b.Insensitive = true
	}
	if remains, ok = hasArg(remains, "-R"); ok {
		b.Recursive = true
	}
	b.Paths = p
	b.Args = getArgs(remains)
	return
}

// files starts from the end of the slice and removes any file it finds and
// returns them in p.
func files(input []string) (remaining []string, p []string, err error) {
	var (
		foundOne bool
		counter  int
		ret      []string
	)
	sort.Slice(input, func(i, j int) bool {
		return i > j
	})
	for _, t := range input {
		t = strings.Trim(t, " ")
		if t == "" || inStringSlice(t, p) {
			continue
		}
		if _, err := os.Stat(t); err == nil {
			foundOne = true
			p = append(p, t)
			counter++
			continue
		}
		if !foundOne {
			return input, nil, errors.Wrap(ErrFileNotFound, t)
		}
		ret = append(ret, t)
	}

	//We have reversed it. We need to return back in the same order.
	sort.Slice(ret, func(i, j int) bool {
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

// hasArg removes the `arg` argument and returns the remaining []string.
func hasArg(input []string, arg string) ([]string, bool) {
	for i, a := range input {
		if a == arg {
			return append(input[:i], input[i+1:]...), true
		}
	}
	return input, false
}

func getArgs(input []string) []blush.Arg {
	var (
		lastColour = blush.DefaultColour
		ret        []blush.Arg
	)
	for _, token := range input {
		if strings.HasPrefix(token, "-") {
			lastColour = colorFromArg(token)
			continue
		}
		a := blush.Arg{
			Colour: lastColour,
			Find:   blush.NewLocator(token),
		}
		ret = append(ret, a)
	}
	return ret
}

func colorFromArg(arg string) blush.Colour {
	switch arg {
	case "-r", "--red":
		return blush.FgRed
	case "-b", "--blue":
		return blush.FgBlue
	case "-g", "--green":
		return blush.FgGreen
	case "-bl", "--black":
		return blush.FgBlack
	case "-w", "--white":
		return blush.FgWhite
	case "-cy", "--cyan":
		return blush.FgCyan
	case "-mg", "--magenta":
		return blush.FgMagenta
	case "-yl", "--yellow":
		return blush.FgYellow
	}
	return blush.DefaultColour
}
