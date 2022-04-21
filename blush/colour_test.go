package blush_test

import (
	"strings"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/arsham/blush/blush"
)

func TestColourise(t *testing.T) {
	t.Parallel()
	input := "nswkgjTusmxWoiZLhZOGBG"
	got := blush.Colourise(input, blush.NoColour)
	assert.Contains(t, input, got)
	assert.NotContains(t, "[38;", got)
	assert.NotContains(t, "[48;", got)
	assert.NotContains(t, "\033[0m", got)

	c := blush.Colour{
		Foreground: blush.FgGreen,
		Background: blush.FgRed,
	}
	got = blush.Colourise(input, c)
	assert.Contains(t, got, input)
	assert.Contains(t, got, "[38;")
	assert.Contains(t, got, "[48;")
	assert.EqualValues(t, 1, strings.Count(got, "\033[0m"))

	c = blush.Colour{
		Foreground: blush.NoRGB,
		Background: blush.FgRed,
	}
	got = blush.Colourise(input, c)
	assert.Contains(t, got, input)
	assert.NotContains(t, got, "[38;")
	assert.Contains(t, got, "[48;")
	assert.EqualValues(t, 1, strings.Count(got, "\033[0m"))
}
