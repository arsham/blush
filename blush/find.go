package blush

import (
	"fmt"
	"regexp"
	"strings"
)

var isRegExp = regexp.MustCompile(`[\^\$\.\{\}\[\]\*\?]`)

// Finder is a strategy to find texts based on a plain text or regexp logic. If
// Find finds the string, it will decorate it with the given Colour. If the
// colour is zero, it doesn't decorate and works as a regular grep.
type Finder interface {
	Find(string) (string, bool)
	Colour() Colour
}

// NewLocator returns a `rx` object if the `search` is a valid regexp, otherwise
// it returns a plain locator. If insensitive is true, the match will be case
// insensitive. The `colour` argument can be in short form (b) or long form
// (blue). If it cannot find the colour, it will fallback to DefaultColour. The
// colour also can be in hex format, which should be started with a pound sign
// (#666).
func NewLocator(colour, search string, insensitive bool) Finder {
	c := colorFromArg(colour)
	if !isRegExp.Match([]byte(search)) {
		if insensitive {
			return NewIexact(search, c)
		}
		return NewExact(search, c)
	}

	decore := fmt.Sprintf("(%s)", search)
	if insensitive {
		decore = fmt.Sprintf("(?i)%s", decore)
		if o, err := regexp.Compile(decore); err == nil {
			return NewRx(o, c)
		}
		return NewIexact(search, c)
	}

	if o, err := regexp.Compile(decore); err == nil {
		return NewRx(o, c)
	}
	return NewExact(search, c)
}

// Exact looks for the exact word in the string.
type Exact struct {
	s      string
	colour Colour
}

// NewExact returns a new instance of the Exact.
func NewExact(s string, c Colour) Exact {
	return Exact{
		s:      s,
		colour: c,
	}
}

// Find looks for the exact string.
func (e Exact) Find(input string) (string, bool) {
	if strings.Contains(input, e.s) {
		return e.colourise(input, e.colour), true
	}
	return "", false
}

func (e Exact) colourise(input string, c Colour) string {
	if c == NoColour {
		return input
	}
	return strings.Replace(input, e.s, Colourise(e.s, c), -1)
}

// Colour returns the Colour property.
func (e Exact) Colour() Colour {
	return e.colour
}

func (e Exact) String() string {
	return e.colourise(e.s, e.colour)
}

// Iexact is like Exact but case insensitive.
type Iexact struct {
	s      string
	colour Colour
}

// NewIexact returns a new instance of the Iexact.
func NewIexact(s string, c Colour) Iexact {
	return Iexact{
		s:      s,
		colour: c,
	}
}

// Find looks for the exact string.
func (i Iexact) Find(input string) (string, bool) {
	if strings.Contains(strings.ToLower(input), strings.ToLower(i.s)) {
		return i.colourise(input, i.colour), true
	}
	return "", false
}

func (i Iexact) colourise(input string, c Colour) string {
	if c == NoColour {
		return input
	}
	index := strings.Index(strings.ToLower(input), strings.ToLower(i.s))
	end := len(i.s) + index
	match := input[index:end]
	return strings.Replace(input, match, Colourise(match, c), -1)
}

// Colour returns the Colour property.
func (i Iexact) Colour() Colour {
	return i.colour
}

func (i Iexact) String() string {
	return i.colourise(i.s, i.colour)
}

// Rx is the regexp implementation of the Locator.
type Rx struct {
	*regexp.Regexp
	colour Colour
}

// NewRx returns a new instance of the Rx.
func NewRx(r *regexp.Regexp, c Colour) Rx {
	return Rx{
		Regexp: r,
		colour: c,
	}
}

// Find looks for the string matching `r` regular expression..
func (r Rx) Find(input string) (string, bool) {
	if r.MatchString(input) {
		return r.colourise(input, r.colour), true
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

// Colour returns the Colour property.
func (r Rx) Colour() Colour {
	return r.colour
}
