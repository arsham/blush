package blush

import (
	"fmt"
)

// Some stock colours.
var (
	NoColour  = Colour{-1, -1, -1}
	FgRed     = Colour{255, 0, 0}
	FgBlue    = Colour{0, 0, 255}
	FgGreen   = Colour{0, 255, 0}
	FgBlack   = Colour{0, 0, 0}
	FgWhite   = Colour{255, 255, 255}
	FgCyan    = Colour{0, 255, 255}
	FgMagenta = Colour{255, 0, 255}
	FgYellow  = Colour{255, 255, 0}
)

// DefaultColour is a no colour. There will be no colouring when used.
var DefaultColour = NoColour

// Colour is a RGB colour scheme. R, G and B should be between 0 and 255.
type Colour struct {
	R, G, B int
}

// Colourise wraps the `input` between colours for terminals.
func Colourise(input string, c Colour) string {
	return fmt.Sprintf("%s%s%s", format(c), input, unformat())
}

func format(c Colour) string {
	return fmt.Sprintf("\033[38;5;%dm", colour(c.R, c.G, c.B))
}

func unformat() string {
	return "\033[0m"
}

func colour(red, green, blue int) int {
	return 16 + baseColor(red, 36) + baseColor(green, 6) + baseColor(blue, 1)
}

func baseColor(value int, factor int) int {
	return int(6*float64(value)/256) * factor
}
