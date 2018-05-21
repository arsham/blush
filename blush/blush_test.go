package blush_test

import (
	"strings"
	"testing"

	"github.com/arsham/blush/blush"
)

func TestNewErrors(t *testing.T) {
	tcs := []struct {
		name   string
		input  string
		errTxt string
	}{
		{"empty", "", blush.ErrNoInput.Error()},
		{"one missing file", "/doesnotexist", "/doesnotexist"},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			g, err := blush.New(tc.input)
			if g != nil {
				t.Errorf("New(): g = %v, want nil", g)
			}
			if err == nil {
				t.Error("New(): err = nil, want error")
				return
			}
			if !strings.Contains(err.Error(), tc.errTxt) {
				t.Errorf("want %s in %s", tc.errTxt, err.Error())
			}
		})
	}
}

func TestNewFindFiles(t *testing.T) {
	tcs := []struct {
		name  string
		input string
		count int
	}{
		{"path", "/", 1},
		{"path duplicate", "/ /", 1},
		{"path trailing spaces", "/         ", 1},
		{"path with prefix", "something else /", 1},
		{"paths", "/ /dev", 2},
		{"file", "/dev/null", 1},
		{"file duplicate", "/dev/null /dev/null", 1},
		{"file trailing spaces", "/dev/null      ", 1},
		{"file with prefix", "something else /dev/null", 1},
		{"files", "/dev/null /dev/zero", 2},
		{"file and path", "/dev/null /dev", 2},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			g, err := blush.New(tc.input)
			if err != nil {
				t.Errorf("New(): err = %s, want nil", err)
			}
			if g == nil {
				t.Error("New(): g = nil, want *Blush")
				return
			}
			if len(g.Paths) != tc.count {
				t.Errorf("len(g.Paths) = %d, want %d", len(g.Paths), tc.count)
			}
		})
	}
}

func TestNewOtherArgs(t *testing.T) {
	g, err := blush.New("/")
	if g.Sensitive {
		t.Error("g.Sensitive = true, want false")
	}
	if g.Recursive {
		t.Error("g.Recursive = true, want false")
	}
	if err != nil {
		t.Errorf("New(): err = %s, want nil", err)
	}
	tcs := []string{
		"-s /",
		"-s -s /",
		"-s -R /",
		"-R -s /",
		"aaa -s /",
		"-s aaa -s /",
		"aaa -s -s /",
	}
	for _, tc := range tcs {
		g, err := blush.New(tc)
		if err != nil {
			t.Errorf("New(%s): err = %s, want nil", tc, err)
		}
		if g == nil {
			t.Errorf("New(%s): g = nil, want *Blush", tc)
			continue
		}
		if !g.Sensitive {
			t.Errorf("%s: g.Sensitive = false, want true", tc)
		}
	}
	tcs = []string{
		"-R /",
		"-R -R /",
		"-s -R /",
		"-R -s /",
		"aaa -R /",
		"-R aaa -R /",
		"aaa -R -R /",
	}
	for _, tc := range tcs {
		g, err := blush.New(tc)
		if err != nil {
			t.Errorf("New(%s): err = %s, want nil", tc, err)
		}
		if g == nil {
			t.Errorf("New(%s): g = nil, want *Blush", tc)
			continue
		}
		if !g.Recursive {
			t.Errorf("%s: g.Recursive = false, want true", tc)
		}
	}

	g, err = blush.New("-R -s /")
	if err != nil {
		t.Errorf("New(): err = %s, want nil", err)
	}
	if g == nil {
		t.Error("New(): g = nil, want *Blush")
		return
	}
	if !g.Recursive {
		t.Error("g.Recursive = false, want true")
	}
	if !g.Sensitive {
		t.Error("g.Sensitive = false, want true")
	}
}

func argsEqual(a, b []blush.Arg) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	isIn := func(a blush.Arg, haystack []blush.Arg) bool {
		for _, b := range haystack {
			af, bf := a.Find.(blush.Exact), b.Find.(blush.Exact)
			if a.Colour == b.Colour && string(af) == string(bf) {
				return true
			}
		}
		return false
	}

	for _, item := range b {
		if !isIn(item, a) {
			return false
		}
	}
	return true
}

func TestNewColourArgs(t *testing.T) {
	aaa := blush.Exact("aaa")
	bbb := blush.Exact("bbb")
	tcs := []struct {
		name  string
		input string
		want  []blush.Arg
	}{
		{"empty", "/", []blush.Arg{}},
		{"1-no colour", "aaa /", []blush.Arg{
			blush.Arg{Colour: blush.DefaultColour, Find: aaa},
		}},
		{"1-colour", "-b aaa /", []blush.Arg{
			blush.Arg{Colour: blush.FgBlue, Find: aaa},
		}},
		{"1-colour long", "--blue aaa /", []blush.Arg{
			blush.Arg{Colour: blush.FgBlue, Find: aaa},
		}},
		{"2-no colour", "aaa bbb /", []blush.Arg{
			blush.Arg{Colour: blush.DefaultColour, Find: aaa},
			blush.Arg{Colour: blush.DefaultColour, Find: bbb},
		}},
		{"2-colour", "-b aaa bbb /", []blush.Arg{
			blush.Arg{Colour: blush.FgBlue, Find: aaa},
			blush.Arg{Colour: blush.FgBlue, Find: bbb},
		}},
		{"2-two colours", "-b aaa -g bbb /", []blush.Arg{
			blush.Arg{Colour: blush.FgBlue, Find: aaa},
			blush.Arg{Colour: blush.FgGreen, Find: bbb},
		}},
		{"red", "-r aaa --red bbb /", []blush.Arg{
			blush.Arg{Colour: blush.FgRed, Find: aaa},
			blush.Arg{Colour: blush.FgRed, Find: bbb},
		}},
		{"green", "-g aaa --green bbb /", []blush.Arg{
			blush.Arg{Colour: blush.FgGreen, Find: aaa},
			blush.Arg{Colour: blush.FgGreen, Find: bbb},
		}},
		{"blue", "-b aaa --blue bbb /", []blush.Arg{
			blush.Arg{Colour: blush.FgBlue, Find: aaa},
			blush.Arg{Colour: blush.FgBlue, Find: bbb},
		}},
		{"white", "-w aaa --white bbb /", []blush.Arg{
			blush.Arg{Colour: blush.FgWhite, Find: aaa},
			blush.Arg{Colour: blush.FgWhite, Find: bbb},
		}},
		{"black", "-bl aaa --black bbb /", []blush.Arg{
			blush.Arg{Colour: blush.FgBlack, Find: aaa},
			blush.Arg{Colour: blush.FgBlack, Find: bbb},
		}},
		{"cyan", "-cy aaa --cyan bbb /", []blush.Arg{
			blush.Arg{Colour: blush.FgCyan, Find: aaa},
			blush.Arg{Colour: blush.FgCyan, Find: bbb},
		}},
		{"magenta", "-mg aaa --magenta bbb /", []blush.Arg{
			blush.Arg{Colour: blush.FgMagenta, Find: aaa},
			blush.Arg{Colour: blush.FgMagenta, Find: bbb},
		}},
		{"yellow", "-yl aaa --yellow bbb /", []blush.Arg{
			blush.Arg{Colour: blush.FgYellow, Find: aaa},
			blush.Arg{Colour: blush.FgYellow, Find: bbb},
		}},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			g, err := blush.New(tc.input)
			if err != nil {
				t.Errorf("New(): err = %s, want nil", err)
			}
			if g == nil {
				t.Error("New(): g = nil, want *Blush")
			}
			if !argsEqual(g.Args, tc.want) {
				t.Errorf("(%s): g.Args = %v, want %v", tc.input, g.Args, tc.want)
			}
		})
	}
}

// testing everything
func TestNew(t *testing.T) {}
