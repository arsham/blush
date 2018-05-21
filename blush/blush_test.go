package blush_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/arsham/blush/blush"
)

// In the testdata folder, there are three files. In each file there are 1 ONE,
// 2 TWO, 3 THREE and 4 FOURs.

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
		"-i /",
		"-i -i /",
		"-i -R /",
		"-R -i /",
		"aaa -i /",
		"-i aaa -i /",
		"aaa -i -i /",
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
		"-i -R /",
		"-R -i /",
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

	g, err = blush.New("-R -i /")
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

func TestWriteErrors(t *testing.T) {
	dir, err := ioutil.TempDir("", "blush")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			t.Errorf("could not remove the folder: %s", dir)
		}
	}()
	l, err := blush.New("SHOULDNOTFINDTHISONE " + dir)
	if err != nil {
		t.Fatal(err)
	}
	buf := new(bytes.Buffer)
	err = l.Write(buf)
	if err == nil {
		t.Error("err = nil, want error")
	}

	// Creating a file, letting Blush register it and then we remove it just
	// before we attempt to read. It should throw an error.
	name := path.Join(dir, "something")
	_, err = os.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	l, err = blush.New("SHOULDNOTFINDTHISONE " + dir)
	if err != nil {
		t.Fatal(err)
	}
	err = os.Remove(name)
	if err != nil {
		t.Fatal(err)
	}
	err = l.Write(buf)
	if err == nil {
		t.Error("err = nil, want error")
	}
}

func TestWriteNoMatch(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	location := path.Join(pwd, "testdata")
	l, err := blush.New("SHOULDNOTFINDTHISONE " + location)
	if err != nil {
		t.Fatal(err)
	}
	buf := new(bytes.Buffer)
	err = l.Write(buf)
	if buf.Len() > 0 {
		t.Errorf("buf.Len() = %d, want 0", buf.Len())
	}
	if err != nil {
		t.Errorf("err = %v, want %v", err, nil)
	}
}

func TestWriteMatchNoColourPlain(t *testing.T) {
	match := "TOKEN"
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	location := path.Join(pwd, "testdata")
	input := fmt.Sprintf("-R %s %s", match, location)
	l, err := blush.New(input)
	if err != nil {
		t.Fatal(err)
	}
	buf := new(bytes.Buffer)
	err = l.Write(buf)
	if buf.Len() == 0 {
		t.Errorf("buf.Len() = %d, want > 0", buf.Len())
	}
	if err != nil {
		t.Errorf("err = %v, want %v", err, nil)
	}
	if !strings.Contains(buf.String(), match) {
		t.Errorf("want `%s` in `%s`", match, buf.String())
	}
	if strings.Contains(buf.String(), "[38;5;") {
		t.Errorf("didn't expect colouring: `%s`", buf.String())
	}
}

func TestWriteMatchColour(t *testing.T) {
	match := blush.Colourise("TOKEN", blush.FgBlue)
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	location := path.Join(pwd, "testdata")
	l, err := blush.New("-b TOKEN " + location)
	if err != nil {
		t.Fatal(err)
	}
	buf := new(bytes.Buffer)
	err = l.Write(buf)
	if buf.Len() == 0 {
		t.Errorf("buf.Len() = %d, want > 0", buf.Len())
	}
	if err != nil {
		t.Errorf("err = %v, want %v", err, nil)
	}
	if !strings.Contains(buf.String(), match) {
		t.Errorf("want `%s` in `%s`", match, buf.String())
	}
}

func TestWriteMatchCountColour(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	location := path.Join(pwd, "testdata")
	tcs := []struct {
		name      string
		recursive bool
		count     int
	}{
		{"ONE", false, 1},
		{"ONE", true, 3 * 1},
		{"TWO", false, 2},
		{"TWO", true, 3 * 2},
		{"THREE", false, 3},
		{"THREE", true, 3 * 3},
		{"FOUR", false, 4},
		{"FOUR", true, 3 * 4},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			match := blush.Colourise(tc.name, blush.FgRed)
			input := fmt.Sprintf("-r %s %s", tc.name, location)
			if tc.recursive {
				input = "-R " + input
			}
			l, err := blush.New(input)
			if err != nil {
				t.Error(err)
				return
			}
			buf := new(bytes.Buffer)
			err = l.Write(buf)
			if err != nil {
				t.Errorf("err = %v, want %v", err, nil)
			}
			count := strings.Count(buf.String(), match)
			if count != tc.count {
				t.Errorf("count = %d, want %d", count, tc.count)
			}
		})
	}
}

func TestWriteMultiColour(t *testing.T) {
	two := blush.Colourise("TWO", blush.FgMagenta)
	three := blush.Colourise("THREE", blush.FgRed)
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	location := path.Join(pwd, "testdata")
	input := fmt.Sprintf("-R -mg TWO -r THREE %s", location)
	l, err := blush.New(input)
	if err != nil {
		t.Fatal(err)
	}
	buf := new(bytes.Buffer)
	err = l.Write(buf)
	if buf.Len() == 0 {
		t.Errorf("buf.Len() = %d, want > 0", buf.Len())
	}
	if err != nil {
		t.Errorf("err = %v, want %v", err, nil)
	}
	count := strings.Count(buf.String(), two)
	if count != 2*3 {
		t.Errorf("count = %d, want %d", count, 2*3)
	}
	count = strings.Count(buf.String(), three)
	if count != 3*3 {
		t.Errorf("count = %d, want %d", count, 3*3)
	}
}
