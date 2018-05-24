package blush_test

import (
	"bytes"
	"fmt"
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
	r := ioutil.NopCloser(new(bytes.Buffer))
	tcs := []struct {
		name   string
		b      blush.Blush
		writer io.Writer
		errTxt string
	}{
		{"no input", blush.Blush{}, w, blush.ErrNoInput.Error()},
		{"no writer", blush.Blush{Reader: r}, nil, blush.ErrNoWriter.Error()},
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
}

func TestWriteNoMatch(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	location := path.Join(pwd, "testdata")
	w, err := blush.NewWalker([]string{location}, true)
	if err != nil {
		t.Fatal(err)
	}

	l := blush.Blush{
		Reader:  w,
		Locator: []blush.Locator{blush.NewExact("SHOULDNOTFINDTHISONE", blush.NoColour)},
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
	w, err := blush.NewWalker([]string{location}, true)
	if err != nil {
		t.Fatal(err)
	}
	l := blush.Blush{
		Reader:  w,
		Locator: []blush.Locator{blush.NewExact(match, blush.NoColour)},
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
	w, err := blush.NewWalker([]string{location}, true)
	if err != nil {
		t.Fatal(err)
	}
	l := blush.Blush{
		Reader:  w,
		Locator: []blush.Locator{blush.NewExact("TOKEN", blush.FgBlue)},
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
			location := path.Join(pwd, "testdata")
			w, err := blush.NewWalker([]string{location}, tc.recursive)
			if err != nil {
				t.Fatal(err)
			}
			match := blush.Colourise(tc.name, blush.FgRed)
			l := blush.Blush{
				Reader:  w,
				Locator: []blush.Locator{blush.NewExact(tc.name, blush.FgRed)},
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
	w, err := blush.NewWalker([]string{location}, true)
	if err != nil {
		t.Fatal(err)
	}
	l := blush.Blush{
		Reader: w,
		Locator: []blush.Locator{
			blush.NewExact("TWO", blush.FgMagenta),
			blush.NewExact("THREE", blush.FgRed),
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
	w, err := blush.NewWalker([]string{location}, true)
	if err != nil {
		t.Fatal(err)
	}
	l := blush.Blush{
		Reader: w,
		NoCut:  true,
		Locator: []blush.Locator{
			blush.NewExact("TWO", blush.FgMagenta),
			blush.NewExact("THREE", blush.FgRed),
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

func TestMultipleMatchInOneLine(t *testing.T) {
	line1 := "this is an example\n"
	line2 := "someone should find this line\n"
	input1 := bytes.NewBuffer([]byte(line1))
	input2 := bytes.NewBuffer([]byte(line2))
	w := ioutil.NopCloser(io.MultiReader(input1, input2))
	match := fmt.Sprintf(
		"someone %s find %s line",
		blush.Colourise("should", blush.FgRed),
		blush.Colourise("this", blush.FgMagenta),
	)
	out := new(bytes.Buffer)

	l := blush.Blush{
		Reader: w,
		Locator: []blush.Locator{
			blush.NewExact("this", blush.FgMagenta),
			blush.NewExact("should", blush.FgRed),
		},
	}

	l.Write(out)
	lines := strings.Split(out.String(), "\n")
	example := lines[1]
	if strings.Contains(example, "is an example") {
		example = lines[0]
	}
	if example != match {
		t.Errorf("example = %s, want %s", example, match)
	}
}

type nopCloser struct {
	io.Reader
	closeFunc func() error
}

func (n nopCloser) Close() error { return n.closeFunc() }

func TestBlushClosesReader(t *testing.T) {
	var called bool
	input := bytes.NewBuffer([]byte("DwgQnpvro5bVvrRwBB"))
	w := nopCloser{
		Reader: input,
		closeFunc: func() error {
			called = true
			return nil
		},
	}
	l := blush.Blush{
		Reader: w,
	}
	err := l.Close()
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}
	if !called {
		t.Error("didn't close the reader")
	}
}

func TestPrintFileName(t *testing.T) {
	t.Skip("not implemented")
}
