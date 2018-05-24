package blush

import (
	"bufio"
	"fmt"
	"io"
)

// Blush has a slice of given regexp, matching paths, and operation
// configuration. If NoCut is true, the unmatched lines are printed as well.
type Blush struct {
	Locator []Locator
	Reader  io.ReadCloser
	NoCut   bool
}

// WriteTo writes matches to w. It returns an error if the writer is nil or
// there are not paths defined or there is no files found in the Reader.
func (b Blush) Write(w io.Writer) error {
	if w == nil {
		return ErrNoWriter
	}
	if b.Reader == nil {
		return ErrNoInput
	}
	if err := b.find(w, b.Reader); err != nil {
		return err
	}
	return nil
}

func (b Blush) find(w io.Writer, file io.Reader) error {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		var foundStr string
		line := scanner.Text()
		for _, a := range b.Locator {
			if s, ok := a.Find(line); ok {
				line = s
				foundStr = line
			}
		}
		if foundStr != "" {
			fmt.Fprintf(w, "%s\n", foundStr)
		} else if b.NoCut {
			fmt.Fprintf(w, "%s\n", line)
		}
	}
	return nil
}

func colorFromArg(arg string) Colour {
	switch arg {
	case "-r", "--red":
		return FgRed
	case "-b", "--blue":
		return FgBlue
	case "-g", "--green":
		return FgGreen
	case "-bl", "--black":
		return FgBlack
	case "-w", "--white":
		return FgWhite
	case "-cy", "--cyan":
		return FgCyan
	case "-mg", "--magenta":
		return FgMagenta
	case "-yl", "--yellow":
		return FgYellow
	}
	return DefaultColour
}
