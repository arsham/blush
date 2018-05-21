package blush

import (
	"strings"
)

// Arg contains a pair of colour name and corresponding regexp.
type Arg struct {
	Colour Colour
	Find   Locator
}

func getArgs(input string) []Arg {
	var (
		lastColour = DefaultColour
		ret        []Arg
	)
	input = strings.Trim(input, " ")
	if input == "" {
		return ret
	}
	tokens := strings.Split(input, " ")
	for _, token := range tokens {
		if strings.HasPrefix(token, "-") {
			lastColour = colorFromArg(token)
			continue
		}
		a := Arg{
			Colour: lastColour,
			Find:   NewLocator(token),
		}
		ret = append(ret, a)
	}
	return ret
}

// hasArg removes the `arg` argument and returns the remaining string.
func hasArg(input, arg string) (string, bool) {
	for _, a := range []string{arg + " ", " " + arg, arg} {
		if strings.Contains(input, a) {
			return strings.Replace(input, a, " ", -1), true
		}
	}
	return input, false
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
