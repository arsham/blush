package blush_test

import (
	"strings"
	"testing"

	"github.com/arsham/blush/blush"
)

func TestColourise(t *testing.T) {
	input := "nswkgjTusmxWoiZLhZOGBG"
	got := blush.Colourise(input, blush.NoColour)
	if !strings.Contains(got, input) {
		t.Errorf("want `%s` in `%s`", input, got)
	}
	if strings.Contains(got, "[38;") {
		t.Errorf("don't want `%v` in `%v`", []byte("[38;"), []byte(got))
	}
	if strings.Contains(got, "[48;") {
		t.Errorf("don't want `%v` in `%v`", []byte("[48;"), []byte(got))
	}
	if strings.Contains(got, "\033[0m") {
		t.Errorf("don't want `%v` in `%v`", []byte("\033[0m"), []byte(got))
	}

	c := blush.Colour{
		Foreground: blush.FgGreen,
		Background: blush.FgRed,
	}
	got = blush.Colourise(input, c)
	if !strings.Contains(got, input) {
		t.Errorf("want `%v` in `%v`", []byte(input), []byte(got))
	}
	if !strings.Contains(got, "[38;") {
		t.Errorf("want `%v` in `%v`", []byte("[38;"), []byte(got))
	}
	if !strings.Contains(got, "[48;") {
		t.Errorf("want `%v` in `%v`", []byte("[48;"), []byte(got))
	}
	if strings.Count(got, "\033[0m") != 1 {
		t.Errorf("want unformat to appear once, got %d", strings.Count(got, "\033[0m"))
	}

	c = blush.Colour{
		Foreground: blush.NoRGB,
		Background: blush.FgRed,
	}
	got = blush.Colourise(input, c)
	if !strings.Contains(got, input) {
		t.Errorf("want `%v` in `%v`", []byte(input), []byte(got))
	}
	if strings.Contains(got, "[38;") {
		t.Errorf("don't want `%v` in `%v`", []byte("[38;"), []byte(got))
	}
	if !strings.Contains(got, "[48;") {
		t.Errorf("want `%v` in `%v`", []byte("[48;"), []byte(got))
	}
	if strings.Count(got, "\033[0m") != 1 {
		t.Errorf("want unformat string to appear once, got %d", strings.Count(got, "\033[0m"))
	}
}
