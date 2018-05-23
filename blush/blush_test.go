package blush_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/arsham/blush/blush"
)

// In the testdata folder, there are three files. In each file there are 1 ONE,
// 2 TWO, 3 THREE and 4 FOURs. There is a line containing `LEAVEMEHERE` which
// does not collide with any of these numbers.

var leaveMeHere = "LEAVEMEHERE"

func TestWriteErrors(t *testing.T) {
	w := new(bytes.Buffer)
	tcs := []struct {
		name   string
		b      blush.Blush
		writer io.Writer
		errTxt string
	}{
		{"empty", blush.Blush{}, w, blush.ErrNoFiles.Error()},
		{"no writer", blush.Blush{}, nil, blush.ErrNoWriter.Error()},
		{"no files", blush.Blush{Paths: []string{"/doesnotexist"}}, w, "doesnotexist"},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.b.Write(tc.writer)
			if err == nil {
				t.Error("New(): err = nil, want error")
				return
			}
			if !strings.Contains(err.Error(), tc.errTxt) {
				t.Errorf("want `%s` in `%s`", tc.errTxt, err.Error())
			}
		})
	}

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
	l := blush.Blush{
		Paths: []string{"SHOULDNOTFINDTHISONE " + dir},
	}
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
	l = blush.Blush{
		Paths: []string{"SHOULDNOTFINDTHISONE " + dir},
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
	l := blush.Blush{
		Paths: []string{location},
		Args: []blush.Arg{
			blush.Arg{Find: blush.Exact("SHOULDNOTFINDTHISONE")},
		},
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
	l := blush.Blush{
		Recursive: true,
		Paths:     []string{location},
		Args: []blush.Arg{
			blush.Arg{
				Colour: blush.NoColour,
				Find:   blush.Exact(match),
			},
		},
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
	if strings.Contains(buf.String(), leaveMeHere) {
		t.Errorf("didn't expect to see %s", leaveMeHere)
	}
}

func TestWriteMatchColour(t *testing.T) {
	match := blush.Colourise("TOKEN", blush.FgBlue)
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	location := path.Join(pwd, "testdata")
	l := blush.Blush{
		Paths: []string{location},
		Args: []blush.Arg{
			blush.Arg{
				Colour: blush.FgBlue,
				Find:   blush.Exact("TOKEN"),
			},
		},
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
	if strings.Contains(buf.String(), leaveMeHere) {
		t.Errorf("didn't expect to see %s", leaveMeHere)
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
			l := blush.Blush{
				Paths:     []string{location},
				Recursive: tc.recursive,
				Args: []blush.Arg{
					blush.Arg{
						Colour: blush.FgRed,
						Find:   blush.Exact(tc.name),
					},
				},
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
			if strings.Contains(buf.String(), leaveMeHere) {
				t.Errorf("didn't expect to see %s", leaveMeHere)
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
	l := blush.Blush{
		Paths:     []string{location},
		Recursive: true,
		Args: []blush.Arg{
			blush.Arg{
				Colour: blush.FgMagenta,
				Find:   blush.Exact("TWO"),
			},
			blush.Arg{
				Colour: blush.FgRed,
				Find:   blush.Exact("THREE"),
			},
		},
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
	if strings.Contains(buf.String(), leaveMeHere) {
		t.Errorf("didn't expect to see %s", leaveMeHere)
	}
}

func TestWriteMultiColourColourMode(t *testing.T) {
	two := blush.Colourise("TWO", blush.FgMagenta)
	three := blush.Colourise("THREE", blush.FgRed)
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	location := path.Join(pwd, "testdata")
	l := blush.Blush{
		Paths:     []string{location},
		Recursive: true,
		Colouring: true,
		Args: []blush.Arg{
			blush.Arg{
				Colour: blush.FgMagenta,
				Find:   blush.Exact("TWO"),
			},
			blush.Arg{
				Colour: blush.FgRed,
				Find:   blush.Exact("THREE"),
			},
		},
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
	count = strings.Count(buf.String(), leaveMeHere)
	if count != 1 {
		t.Errorf("count = %d, want to see `%s` exactly %d times", count, leaveMeHere, 1)
	}
}
