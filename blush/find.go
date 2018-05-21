package blush

import (
	"fmt"
	"regexp"
	"strings"
)

// Locator is a strategy to find texts based on a plain text or regexp logic. If
// Find finds the string, it will decorate it with the given Colour. If the
// colour is zero, it doesn't decorate and works as a regular grep.
type Locator interface {
	Find(string, Colour) (string, bool)
}

var isRegExp = regexp.MustCompile(`[\^\$\.\{\}\[\]\*]`)

// NewLocator returns a `rx` object id the `input` is a valid regexp, otherwise
// it returns a plain locator.
func NewLocator(input string) Locator {
	if !isRegExp.Match([]byte(input)) {
		return Exact(input)
	}
	dec := fmt.Sprintf("(%s)", input)
	if o, err := regexp.Compile(dec); err == nil {
		return Rx{o}
	}
	return Exact(input)
}

// Exact looks for the exact word in the string.
type Exact string

// Find looks for the exact string.
func (e Exact) Find(input string, c Colour) (string, bool) {
	if strings.Contains(input, string(e)) {
		return e.colourise(input, c), true
	}
	return "", false
}

func (e Exact) colourise(input string, c Colour) string {
	if c == NoColour {
		return input
	}
	return strings.Replace(input, string(e), Colourise(string(e), c), -1)
}

// Rx is the regexp implementation of the Locator.
type Rx struct {
	*regexp.Regexp
}

// Find looks for the string matching `r` regular expression..
func (r Rx) Find(input string, c Colour) (string, bool) {
	if r.MatchString(input) {
		return r.colourise(input, c), true
	}
	return "", false
}

func (r Rx) colourise(input string, c Colour) string {
	if c == NoColour {
		return input
	}
	repl := fmt.Sprintf("%s$1%s", format(c), unformat())
	return r.ReplaceAllString(input, repl)
}
