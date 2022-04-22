package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/arsham/blush/blush"
	"github.com/arsham/blush/internal/reader"
)

// Main reads the provided arguments from the command line and creates a
// blush.Blush instance.
func Main() {
	b, err := GetBlush(os.Args)
	if errors.Is(err, errShowHelp) {
		fmt.Println(Usage)
		return
	}
	if err != nil {
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
		log.Print(err)
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
		r   io.ReadCloser = os.Stdin
		a   *args
		err error
	)
	if len(input) == 1 {
		return nil, ErrNoInput
	}
	if a, err = newArgs(input[1:]...); err != nil {
		return nil, err
	}
	if !a.stdin {
		r, err = reader.NewMultiReader(reader.WithPaths(a.paths, a.recursive))
		if err != nil {
			return nil, err
		}
	}
	return &blush.Blush{
		Finders:      a.finders,
		Reader:       r,
		Drop:         a.cut,
		WithFileName: !a.noFilename,
	}, nil
}
