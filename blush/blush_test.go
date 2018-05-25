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

func TestWriteToErrors(t *testing.T) {
	w := new(bytes.Buffer)
	r := ioutil.NopCloser(new(bytes.Buffer))
	tcs := []struct {
		name   string
		b      *blush.Blush
		writer io.Writer
		errTxt string
	}{
		{"no input", &blush.Blush{}, w, blush.ErrNoReader.Error()},
		{"no writer", &blush.Blush{Reader: r}, nil, blush.ErrNoWriter.Error()},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			n, err := tc.b.WriteTo(tc.writer)
			if err == nil {
				t.Error("New(): err = nil, want error")
				return
			}
			if n != 0 {
				t.Errorf("l.WriteTo(): n = %d, want 0", n)
			}
			if !strings.Contains(err.Error(), tc.errTxt) {
				t.Errorf("want `%s` in `%s`", tc.errTxt, err.Error())
			}
		})
	}
}

func TestWriteToNoMatch(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	location := path.Join(pwd, "testdata")
	r, err := blush.NewMultiReadCloser([]string{location}, true)
	if err != nil {
		t.Fatal(err)
	}
	l := &blush.Blush{
		Reader:  r,
		Finders: []blush.Finder{blush.NewExact("SHOULDNOTFINDTHISONE", blush.NoColour)},
	}
	buf := new(bytes.Buffer)
	n, err := l.WriteTo(buf)
	if err != nil {
		t.Errorf("err = %v, want %v", err, nil)
	}
	if n != 0 {
		t.Errorf("l.WriteTo(): n = %d, want %d", n, 0)
	}
	if buf.Len() > 0 {
		t.Errorf("buf.Len() = %d, want 0", buf.Len())
	}
}

func TestWriteToMatchNoColourPlain(t *testing.T) {
	match := "TOKEN"
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	location := path.Join(pwd, "testdata")
	r, err := blush.NewMultiReadCloser([]string{location}, true)
	if err != nil {
		t.Fatal(err)
	}
	l := &blush.Blush{
		Reader:  r,
		Finders: []blush.Finder{blush.NewExact(match, blush.NoColour)},
	}

	buf := new(bytes.Buffer)
	n, err := l.WriteTo(buf)
	if err != nil {
		t.Errorf("err = %v, want %v", err, nil)
	}
	if buf.Len() == 0 {
		t.Errorf("buf.Len() = %d, want > 0", buf.Len())
	}
	if int(n) != buf.Len() {
		t.Errorf("l.WriteTo(): n = %d, want %d", int(n), buf.Len())
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

func TestWriteToMatchColour(t *testing.T) {
	match := blush.Colourise("TOKEN", blush.FgBlue)
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	location := path.Join(pwd, "testdata")
	r, err := blush.NewMultiReadCloser([]string{location}, true)
	if err != nil {
		t.Fatal(err)
	}
	l := &blush.Blush{
		Reader:  r,
		Finders: []blush.Finder{blush.NewExact("TOKEN", blush.FgBlue)},
	}

	buf := new(bytes.Buffer)
	n, err := l.WriteTo(buf)
	if err != nil {
		t.Errorf("err = %v, want %v", err, nil)
	}
	if buf.Len() == 0 {
		t.Errorf("buf.Len() = %d, want > 0", buf.Len())
	}
	if int(n) != buf.Len() {
		t.Errorf("l.WriteTo(): n = %d, want %d", int(n), buf.Len())
	}
	if !strings.Contains(buf.String(), match) {
		t.Errorf("want `%s` in `%s`", match, buf.String())
	}
	if strings.Contains(buf.String(), leaveMeHere) {
		t.Errorf("didn't expect to see %s", leaveMeHere)
	}
}

func TestWriteToMatchCountColour(t *testing.T) {
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
			r, err := blush.NewMultiReadCloser([]string{location}, tc.recursive)
			if err != nil {
				t.Fatal(err)
			}

			match := blush.Colourise(tc.name, blush.FgRed)
			l := &blush.Blush{
				Reader:  r,
				Finders: []blush.Finder{blush.NewExact(tc.name, blush.FgRed)},
			}

			buf := new(bytes.Buffer)
			n, err := l.WriteTo(buf)
			if err != nil {
				t.Errorf("l.WriteTo(): err = %v, want %v", err, nil)
			}
			if int(n) != buf.Len() {
				t.Errorf("l.WriteTo(): n = %d, want %d", int(n), buf.Len())
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

func TestWriteToMultiColour(t *testing.T) {
	two := blush.Colourise("TWO", blush.FgMagenta)
	three := blush.Colourise("THREE", blush.FgRed)
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	location := path.Join(pwd, "testdata")
	r, err := blush.NewMultiReadCloser([]string{location}, true)
	if err != nil {
		t.Fatal(err)
	}
	l := &blush.Blush{
		Reader: r,
		Finders: []blush.Finder{
			blush.NewExact("TWO", blush.FgMagenta),
			blush.NewExact("THREE", blush.FgRed),
		},
	}

	buf := new(bytes.Buffer)
	n, err := l.WriteTo(buf)
	if err != nil {
		t.Errorf("err = %v, want %v", err, nil)
	}
	if buf.Len() == 0 {
		t.Errorf("buf.Len() = %d, want > 0", buf.Len())
	}
	if int(n) != buf.Len() {
		t.Errorf("l.WriteTo(): n = %d, want %d", int(n), buf.Len())
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

func TestWriteToMultiColourColourMode(t *testing.T) {
	two := blush.Colourise("TWO", blush.FgMagenta)
	three := blush.Colourise("THREE", blush.FgRed)
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	location := path.Join(pwd, "testdata")
	r, err := blush.NewMultiReadCloser([]string{location}, true)
	if err != nil {
		t.Fatal(err)
	}
	l := &blush.Blush{
		Reader: r,
		NoCut:  true,
		Finders: []blush.Finder{
			blush.NewExact("TWO", blush.FgMagenta),
			blush.NewExact("THREE", blush.FgRed),
		},
	}

	buf := new(bytes.Buffer)
	n, err := l.WriteTo(buf)
	if err != nil {
		t.Errorf("err = %v, want %v", err, nil)
	}
	if buf.Len() == 0 {
		t.Errorf("buf.Len() = %d, want > 0", buf.Len())
	}
	if int(n) != buf.Len() {
		t.Errorf("l.WriteTo(): n = %d, want %d", int(n), buf.Len())
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

func TestWriteToMultipleMatchInOneLine(t *testing.T) {
	line1 := "this is an example\n"
	line2 := "someone should find this line\n"
	input1 := bytes.NewBuffer([]byte(line1))
	input2 := bytes.NewBuffer([]byte(line2))
	r := ioutil.NopCloser(io.MultiReader(input1, input2))
	match := fmt.Sprintf(
		"someone %s find %s line",
		blush.Colourise("should", blush.FgRed),
		blush.Colourise("this", blush.FgMagenta),
	)
	out := new(bytes.Buffer)

	l := &blush.Blush{
		Reader: r,
		Finders: []blush.Finder{
			blush.NewExact("this", blush.FgMagenta),
			blush.NewExact("should", blush.FgRed),
		},
	}

	l.WriteTo(out)
	lines := strings.Split(out.String(), "\n")
	example := lines[1]
	if strings.Contains(example, "is an example") {
		example = lines[0]
	}
	if example != match {
		t.Errorf("example = %s, want %s", example, match)
	}
}

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
	l := &blush.Blush{
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

func TestBlushReadOneStream(t *testing.T) {
	input := bytes.NewBuffer([]byte("one two three four"))
	match := blush.NewExact("three", blush.FgBlue)
	r := ioutil.NopCloser(input)
	l := &blush.Blush{
		Finders: []blush.Finder{match},
		Reader:  r,
	}
	defer l.Close()
	emptyP := make([]byte, 10)
	tcs := []struct {
		name    string
		p       []byte
		wantErr error
		wantLen int
		wantP   string
	}{
		{"one", make([]byte, len("one ")), nil, len("one "), "one "},
		{"two", make([]byte, len("two ")), nil, len("two "), "two "},
		{"three", make([]byte, len(match.String())), nil, len(match.String()), match.String()},
		{"four", make([]byte, len(" four\n")), nil, len(" four\n"), " four\n"}, // there is always a new line after each reader.
		{"empty", emptyP, io.EOF, 0, string(emptyP)},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			n, err := l.Read(tc.p)
			if err != tc.wantErr {
				t.Error(err)
			}
			if n != tc.wantLen {
				t.Errorf("l.Read(): n = %d, want %d", n, tc.wantLen)
			}
			if string(tc.p) != tc.wantP {
				t.Errorf("p = `%s`, want `%s`", tc.p, tc.wantP)
			}
		})
	}
}

func TestBlushReadTwoStreams(t *testing.T) {
	b1 := []byte("one for all\n")
	b2 := []byte("all for one\n")
	input1 := bytes.NewBuffer(b1)
	input2 := bytes.NewBuffer(b2)
	match := blush.NewExact("one", blush.FgBlue)
	r := ioutil.NopCloser(io.MultiReader(input1, input2))
	l := &blush.Blush{
		Finders: []blush.Finder{match},
		Reader:  r,
	}
	defer l.Close()

	buf := new(bytes.Buffer)
	n, err := buf.ReadFrom(l)
	if err != nil {
		t.Error(err)
	}
	expectLen := len(b1) + len(b2) - len("one")*2 + len(match.String())*2
	if int(n) != expectLen {
		t.Errorf("l.Read(): n = %d, want %d", n, expectLen)
	}
	expectStr := fmt.Sprintf("%s%s",
		strings.Replace(string(b1), "one", match.String(), 1),
		strings.Replace(string(b2), "one", match.String(), 1),
	)
	if buf.String() != expectStr {
		t.Errorf("buf.String() = %s, want %s", buf.String(), expectStr)
	}
}

func TestBlushReadHalfWay(t *testing.T) {
	b1 := []byte("one for all\n")
	b2 := []byte("all for one\n")
	input1 := bytes.NewBuffer(b1)
	input2 := bytes.NewBuffer(b2)
	match := blush.NewExact("one", blush.FgBlue)
	r := ioutil.NopCloser(io.MultiReader(input1, input2))
	l := &blush.Blush{
		Finders: []blush.Finder{match},
		Reader:  r,
	}
	p := make([]byte, len(b1))
	_, err := l.Read(p)
	if err != nil {
		t.Error(err)
	}
	n, err := l.Read(p)
	if n != len(b1) {
		t.Errorf("l.Read(): n = %d, want %d", n, len(b1))
	}
	if err != nil {
		t.Errorf("l.Read(): err = %v, want %v", err, nil)
	}
}

func TestBlushReadOnClosed(t *testing.T) {
	b1 := []byte("one for all\n")
	b2 := []byte("all for one\n")
	input1 := bytes.NewBuffer(b1)
	input2 := bytes.NewBuffer(b2)
	match := blush.NewExact("one", blush.FgBlue)
	r := ioutil.NopCloser(io.MultiReader(input1, input2))
	l := &blush.Blush{
		Finders: []blush.Finder{match},
		Reader:  r,
	}
	p := make([]byte, len(b1))
	_, err := l.Read(p)
	if err != nil {
		t.Error(err)
	}
	err = l.Close()
	if err != nil {
		t.Fatal(err)
	}
	n, err := l.Read(p)
	if n != 0 {
		t.Errorf("l.Read(): n = %d, want 0", n)
	}
	if err != blush.ErrClosed {
		t.Errorf("l.Read(): err = %v, want %v", err, blush.ErrClosed)
	}
}

func TestPrintFileName(t *testing.T) {
	t.Skip("not implemented")
}
