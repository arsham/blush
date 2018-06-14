package blush_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/arsham/blush/blush"
	"github.com/arsham/blush/internal/reader"
)

func TestWriteToErrors(t *testing.T) {
	w := new(bytes.Buffer)
	e := errors.New("something")
	nn := 10
	bw := &badWriter{
		writeFunc: func([]byte) (int, error) {
			return nn, e
		},
	}
	getReader := func() io.ReadCloser {
		return ioutil.NopCloser(bytes.NewBufferString("something"))
	}
	tcs := []struct {
		name    string
		b       *blush.Blush
		writer  io.Writer
		wantN   int
		wantErr string
	}{
		{"no input", &blush.Blush{}, w, 0, reader.ErrNoReader.Error()},
		{"no writer", &blush.Blush{Reader: getReader()}, nil, 0, blush.ErrNoWriter.Error()},
		{"bad writer", &blush.Blush{Reader: getReader(), NoCut: true}, bw, nn, e.Error()},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			n, err := tc.b.WriteTo(tc.writer)
			if err == nil {
				t.Error("New(): err = nil, want error")
				return
			}
			if int(n) != tc.wantN {
				t.Errorf("l.WriteTo(): n = %d, want %d", n, tc.wantN)
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("want `%s` in `%s`", tc.wantErr, err.Error())
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
	r, err := reader.NewMultiReader(reader.WithPaths([]string{location}, true))
	if err != nil {
		t.Fatal(err)
	}
	b := &blush.Blush{
		Reader:  r,
		Finders: []blush.Finder{blush.NewExact("SHOULDNOTFINDTHISONE", blush.NoColour)},
	}
	buf := new(bytes.Buffer)
	n, err := b.WriteTo(buf)
	if err != nil {
		t.Errorf("err = %v, want %v", err, nil)
	}
	if n != 0 {
		t.Errorf("b.WriteTo(): n = %d, want %d", n, 0)
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
	r, err := reader.NewMultiReader(reader.WithPaths([]string{location}, true))
	if err != nil {
		t.Fatal(err)
	}
	b := &blush.Blush{
		Reader:  r,
		Finders: []blush.Finder{blush.NewExact(match, blush.NoColour)},
	}

	buf := new(bytes.Buffer)
	n, err := b.WriteTo(buf)
	if err != nil {
		t.Errorf("err = %v, want %v", err, nil)
	}
	if buf.Len() == 0 {
		t.Errorf("buf.Len() = %d, want > 0", buf.Len())
	}
	if int(n) != buf.Len() {
		t.Errorf("b.WriteTo(): n = %d, want %d", int(n), buf.Len())
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
	match := blush.Colourise("TOKEN", blush.Blue)
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	location := path.Join(pwd, "testdata")
	r, err := reader.NewMultiReader(reader.WithPaths([]string{location}, true))
	if err != nil {
		t.Fatal(err)
	}
	b := &blush.Blush{
		Reader:  r,
		Finders: []blush.Finder{blush.NewExact("TOKEN", blush.Blue)},
	}

	buf := new(bytes.Buffer)
	n, err := b.WriteTo(buf)
	if err != nil {
		t.Errorf("err = %v, want %v", err, nil)
	}
	if buf.Len() == 0 {
		t.Errorf("buf.Len() = %d, want > 0", buf.Len())
	}
	if int(n) != buf.Len() {
		t.Errorf("b.WriteTo(): n = %d, want %d", int(n), buf.Len())
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
			r, err := reader.NewMultiReader(reader.WithPaths([]string{location}, tc.recursive))
			if err != nil {
				t.Fatal(err)
			}

			match := blush.Colourise(tc.name, blush.Red)
			b := &blush.Blush{
				Reader:  r,
				Finders: []blush.Finder{blush.NewExact(tc.name, blush.Red)},
			}

			buf := new(bytes.Buffer)
			n, err := b.WriteTo(buf)
			if err != nil {
				t.Errorf("b.WriteTo(): err = %v, want %v", err, nil)
			}
			if int(n) != buf.Len() {
				t.Errorf("b.WriteTo(): n = %d, want %d", int(n), buf.Len())
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
	two := blush.Colourise("TWO", blush.Magenta)
	three := blush.Colourise("THREE", blush.Red)
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	location := path.Join(pwd, "testdata")
	r, err := reader.NewMultiReader(reader.WithPaths([]string{location}, true))
	if err != nil {
		t.Fatal(err)
	}
	b := &blush.Blush{
		Reader: r,
		Finders: []blush.Finder{
			blush.NewExact("TWO", blush.Magenta),
			blush.NewExact("THREE", blush.Red),
		},
	}

	buf := new(bytes.Buffer)
	n, err := b.WriteTo(buf)
	if err != nil {
		t.Errorf("err = %v, want %v", err, nil)
	}
	if buf.Len() == 0 {
		t.Errorf("buf.Len() = %d, want > 0", buf.Len())
	}
	if int(n) != buf.Len() {
		t.Errorf("b.WriteTo(): n = %d, want %d", int(n), buf.Len())
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
	two := blush.Colourise("TWO", blush.Magenta)
	three := blush.Colourise("THREE", blush.Red)
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	location := path.Join(pwd, "testdata")
	r, err := reader.NewMultiReader(reader.WithPaths([]string{location}, true))
	if err != nil {
		t.Fatal(err)
	}
	b := &blush.Blush{
		Reader: r,
		NoCut:  true,
		Finders: []blush.Finder{
			blush.NewExact("TWO", blush.Magenta),
			blush.NewExact("THREE", blush.Red),
		},
	}

	buf := new(bytes.Buffer)
	n, err := b.WriteTo(buf)
	if err != nil {
		t.Errorf("err = %v, want %v", err, nil)
	}
	if buf.Len() == 0 {
		t.Errorf("buf.Len() = %d, want > 0", buf.Len())
	}
	if int(n) != buf.Len() {
		t.Errorf("b.WriteTo(): n = %d, want %d", int(n), buf.Len())
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
		blush.Colourise("should", blush.Red),
		blush.Colourise("this", blush.Magenta),
	)
	out := new(bytes.Buffer)

	b := &blush.Blush{
		Reader: r,
		Finders: []blush.Finder{
			blush.NewExact("this", blush.Magenta),
			blush.NewExact("should", blush.Red),
		},
	}

	b.WriteTo(out)
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
	b := &blush.Blush{
		Reader: w,
	}
	err := b.Close()
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}
	if !called {
		t.Error("didn't close the reader")
	}
}

func TestBlushReadOneStream(t *testing.T) {
	input := bytes.NewBuffer([]byte("one two three four"))
	match := blush.NewExact("three", blush.Blue)
	r := ioutil.NopCloser(input)
	b := &blush.Blush{
		Finders: []blush.Finder{match},
		Reader:  r,
	}
	defer b.Close()
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
		{"four", make([]byte, len(" four")), nil, len(" four"), " four"},
		{"empty", emptyP, io.EOF, 0, string(emptyP)},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			n, err := b.Read(tc.p)
			if err != tc.wantErr {
				t.Errorf("err = %v, want %v", err, tc.wantErr)
			}
			if n != tc.wantLen {
				t.Errorf("b.Read(): n = %d, want %d", n, tc.wantLen)
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
	match := blush.NewExact("one", blush.Blue)
	r := ioutil.NopCloser(io.MultiReader(input1, input2))
	b := &blush.Blush{
		Finders: []blush.Finder{match},
		Reader:  r,
	}
	defer b.Close()

	buf := new(bytes.Buffer)
	n, err := buf.ReadFrom(b)
	if err != nil {
		t.Error(err)
	}
	expectLen := len(b1) + len(b2) - len("one")*2 + len(match.String())*2
	if int(n) != expectLen {
		t.Errorf("b.Read(): n = %d, want %d", n, expectLen)
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
	match := blush.NewExact("one", blush.Blue)
	r := ioutil.NopCloser(io.MultiReader(input1, input2))
	b := &blush.Blush{
		Finders: []blush.Finder{match},
		Reader:  r,
	}
	p := make([]byte, len(b1))
	_, err := b.Read(p)
	if err != nil {
		t.Error(err)
	}
	n, err := b.Read(p)
	if n != len(b1) {
		t.Errorf("b.Read(): n = %d, want %d", n, len(b1))
	}
	if err != nil {
		t.Errorf("b.Read(): err = %v, want %v", err, nil)
	}
}

func TestBlushReadOnClosed(t *testing.T) {
	b1 := []byte("one for all\n")
	b2 := []byte("all for one\n")
	input1 := bytes.NewBuffer(b1)
	input2 := bytes.NewBuffer(b2)
	match := blush.NewExact("one", blush.Blue)
	r := ioutil.NopCloser(io.MultiReader(input1, input2))
	b := &blush.Blush{
		Finders: []blush.Finder{match},
		Reader:  r,
	}
	p := make([]byte, len(b1))
	_, err := b.Read(p)
	if err != nil {
		t.Error(err)
	}
	err = b.Close()
	if err != nil {
		t.Fatal(err)
	}
	n, err := b.Read(p)
	if n != 0 {
		t.Errorf("b.Read(): n = %d, want 0", n)
	}
	if err != blush.ErrClosed {
		t.Errorf("b.Read(): err = %v, want %v", err, blush.ErrClosed)
	}
}

func TestBlushReadLongOneLineText(t *testing.T) {
	head := strings.Repeat("a", 10000)
	tail := strings.Repeat("b", 10000)
	input := bytes.NewBuffer([]byte(head + " FINDME " + tail))
	match := blush.NewExact("FINDME", blush.Blue)
	r := ioutil.NopCloser(input)
	b := &blush.Blush{
		Finders: []blush.Finder{match},
		Reader:  r,
	}
	p := make([]byte, 20)
	_, err := b.Read(p)
	if err != nil {
		t.Error(err)
	}
	err = b.Close()
	if err != nil {
		t.Fatal(err)
	}
	n, err := b.Read(p)
	if n != 0 {
		t.Errorf("b.Read(): n = %d, want 0", n)
	}
	if err != blush.ErrClosed {
		t.Errorf("b.Read(): err = %v, want %v", err, blush.ErrClosed)
	}
}

func TestPrintName(t *testing.T) {
	line1 := "line one\n"
	line2 := "line two\n"
	r1 := ioutil.NopCloser(bytes.NewBuffer([]byte(line1)))
	r2 := ioutil.NopCloser(bytes.NewBuffer([]byte(line2)))
	name1 := "reader1"
	name2 := "reader2"
	r, err := reader.NewMultiReader(
		reader.WithReader(name1, r1),
		reader.WithReader(name2, r2),
	)
	if err != nil {
		t.Fatal(err)
	}
	b := blush.Blush{
		Reader:       r,
		Finders:      []blush.Finder{blush.NewExact("line", blush.NoColour)},
		WithFileName: true,
	}
	buf := new(bytes.Buffer)
	n, err := b.WriteTo(buf)
	if err != nil {
		t.Fatal(err)
	}
	total := len(line1+line2+name1+name2) + len(blush.Separator)*2
	if int(n) != total {
		t.Errorf("total reads = %d, want %d", n, total)
	}
	s := strings.Split(buf.String(), "\n")
	if !strings.Contains(s[0], name1) {
		t.Errorf("want `%s` in `%s`", name1, s[0])
	}
	if !strings.Contains(s[1], name2) {
		t.Fatalf("want `%s` in `%s`", name2, s[1])
	}
}

// testing stdin should not print the name
func TestStdinPrintName(t *testing.T) {
	input := "line one"
	oldStdin := os.Stdin
	f, err := ioutil.TempFile("", "blush_stdin")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err = os.Remove(f.Name()); err != nil {
			t.Error(err)
		}
		os.Stdin = oldStdin
	}()
	os.Stdin = f
	f.WriteString(input)
	f.Seek(0, 0)
	b := blush.Blush{
		Reader:       f,
		Finders:      []blush.Finder{blush.NewExact("line", blush.NoColour)},
		WithFileName: true,
	}
	buf := new(bytes.Buffer)
	_, err = b.WriteTo(buf)
	if err != nil {
		t.Fatal(err)
	}
	if buf.String() != input {
		t.Errorf("buf.String() = `%s`, want `%s`", buf.String(), input)
	}
	if strings.Contains(buf.String(), f.Name()) {
		t.Errorf("buf.String() = `%s`, don't want `%s` in it", buf.String(), f.Name())
	}
}

func TestPrintFileName(t *testing.T) {
	path, err := ioutil.TempDir("", "blush_name")
	if err != nil {
		t.Fatal(err)
	}
	f1, err := ioutil.TempFile(path, "blush_name")
	if err != nil {
		t.Fatal(err)
	}
	f2, err := ioutil.TempFile(path, "blush_name")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err = os.RemoveAll(path); err != nil {
			t.Error(err)
		}
	}()
	line1 := "line one\n"
	line2 := "line two\n"
	f1.WriteString(line1)
	f2.WriteString(line2)
	tcs := []struct {
		name          string
		withFilename  bool
		wantLen       int
		wantFilenames bool
	}{
		{"with filename", true, len(line1+line2+f1.Name()+f2.Name()) + len(blush.Separator)*2, true},
		{"without filename", false, len(line1 + line2), false},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			r, err := reader.NewMultiReader(
				reader.WithPaths([]string{path}, false),
			)
			if err != nil {
				t.Fatal(err)
			}
			b := blush.Blush{
				Reader:       r,
				Finders:      []blush.Finder{blush.NewExact("line", blush.NoColour)},
				WithFileName: tc.withFilename,
			}
			buf := new(bytes.Buffer)
			n, err := b.WriteTo(buf)
			if err != nil {
				t.Fatal(err)
			}
			if int(n) != tc.wantLen {
				t.Errorf("total reads = %d, want %d", n, tc.wantLen)
			}
			notStr := "not"
			if tc.wantFilenames {
				notStr = ""
			}
			if strings.Contains(buf.String(), f1.Name()) != tc.wantFilenames {
				t.Errorf("want `%s` %s in `%s`", f1.Name(), notStr, buf.String())
			}
			if strings.Contains(buf.String(), f2.Name()) != tc.wantFilenames {
				t.Errorf("want `%s` %s in `%s`", f2.Name(), notStr, buf.String())
			}
		})
	}
}

// reading with a small byte slice until the read is done.
func TestReadContiniously(t *testing.T) {
	var (
		ret   []byte
		p     = make([]byte, 2)
		count int
		input = "one two three four\nfive six three seven"
	)
	match := blush.NewExact("three", blush.Blue)
	r := ioutil.NopCloser(bytes.NewBufferString(input))
	b := &blush.Blush{
		Finders: []blush.Finder{match},
		Reader:  r,
	}
	for {
		if count > len(input) { // more that required
			t.Errorf("didn't finish after %d reads: len = %d", count, len(input))
			break
		}
		count++
		_, err := b.Read(p)
		if err == io.EOF {
			ret = append(ret, p...)
			break
		}
		if err != nil {
			t.Error(err)
		}
		ret = append(ret, p...)
	}
	if c := strings.Count(string(ret), "three"); c != 2 {
		t.Errorf("count %s = %d, want %d", "three", c, 2)
	}
	for _, s := range []string{"one", "two", "three", "four", "five", "six", "seven"} {
		if !strings.Contains(string(ret), s) {
			t.Errorf("`%s` not found in `%s`", s, ret)
		}
	}
}

func TestReadMiddleOfMatch(t *testing.T) {
	var (
		search = "aa this aa"
		match  = blush.NewExact("this", blush.Blue)
		p      = make([]byte, (len(search)+len(match.String()))/2)
		ret    []byte
	)
	r := ioutil.NopCloser(bytes.NewBufferString(search))
	b := &blush.Blush{
		Finders: []blush.Finder{match},
		Reader:  r,
	}
	for i := 0; i < 2; i++ {
		_, err := b.Read(p)
		if err != nil && err != io.EOF {
			t.Fatal(err)
		}

		ret = append(ret, p...)
	}
	if !strings.Contains(string(ret), "this") {
		t.Errorf("`%s` not found in `%s`", "this", ret)
	}
}

func TestReadComplete(t *testing.T) {
	input := "123456789"
	match := blush.NewExact("1", blush.NoColour)
	r := ioutil.NopCloser(bytes.NewBufferString(input))
	b := &blush.Blush{
		Finders: []blush.Finder{match},
		Reader:  r,
	}
	p := make([]byte, 10)
	n, err := b.Read(p)
	if n != len(input) {
		t.Errorf("n = %d, want %d", n, len(input))
	}
	if err != io.EOF {
		t.Errorf("err = %v, want %v", err, io.EOF)
	}
	if string(bytes.Trim(p, "\x00")) != input {
		t.Errorf("p = `%v: %s`, want `%v: %s`", p, p, input, input)
	}
	p = make([]byte, 4)
	n, err = b.Read(p)
	if n != 0 {
		t.Errorf("n = %d, want %d", n, 0)
	}
	if err != io.EOF {
		t.Errorf("err = %v, want %v", err, io.EOF)
	}
	if string(bytes.Trim(p, "\x00")) != "" {
		t.Errorf("p = `%v: %s`, want ``", p, p)
	}
}

func TestReadPartComplete(t *testing.T) {
	input := "123456789"
	match := blush.NewExact("1", blush.NoColour)
	r := ioutil.NopCloser(bytes.NewBufferString(input))
	b := &blush.Blush{
		Finders: []blush.Finder{match},
		Reader:  r,
	}
	p := make([]byte, 3)
	n, err := b.Read(p)
	if n != 3 {
		t.Errorf("n = %d, want %d", n, 3)
	}
	if err != nil {
		t.Errorf("err = %v, want %v", err, nil)
	}
	if string(bytes.Trim(p, "\x00")) != "123" {
		t.Errorf("p = `%v: %s`, want `%v: %s`", p, p, "123", "123")
	}
	p = make([]byte, 6)
	n, err = b.Read(p)
	if n != 6 {
		t.Errorf("n = %d, want %d", n, 6)
	}
	if err != nil {
		t.Errorf("err = %v, want %v", err, nil)
	}
	if string(bytes.Trim(p, "\x00")) != "456789" {
		t.Errorf("p = `%v: %s`, want `%v: %s`", p, p, "456789", "456789")
	}
}

func TestReadPartPartOver(t *testing.T) {
	input := "123456789"
	match := blush.NewExact("1", blush.NoColour)
	r := ioutil.NopCloser(bytes.NewBufferString(input))
	b := &blush.Blush{
		Finders: []blush.Finder{match},
		Reader:  r,
	}
	p := make([]byte, 3)
	n, err := b.Read(p)
	if n != 3 {
		t.Errorf("n = %d, want %d", n, 3)
	}
	if err != nil {
		t.Errorf("err = %v, want %v", err, nil)
	}
	if string(bytes.Trim(p, "\x00")) != "123" {
		t.Errorf("p = `%v: %s`, want `%v: %s`", p, p, "123", "123")
	}
	p = make([]byte, 3)
	n, err = b.Read(p)
	if n != 3 {
		t.Errorf("n = %d, want %d", n, 3)
	}
	if err != nil {
		t.Errorf("err = %v, want %v", err, nil)
	}
	if string(bytes.Trim(p, "\x00")) != "456" {
		t.Errorf("p = `%v: %s`, want `%v: %s`", p, p, "456", "456")
	}
	p = make([]byte, 10)
	n, err = b.Read(p)
	if n != 3 {
		t.Errorf("n = %d, want %d", n, 3)
	}
	if err != io.EOF {
		t.Errorf("err = %v, want %v", err, io.EOF)
	}
	if string(bytes.Trim(p, "\x00")) != "789" {
		t.Errorf("p = `%v: %s`, want `%v: %s`", p, p, "789", "789")
	}
}

func TestReadMultiLine(t *testing.T) {
	input := "line1\nline2\nline3\nline4\n"
	match := blush.NewExact("l", blush.NoColour)
	r := ioutil.NopCloser(bytes.NewBufferString(input))
	b := &blush.Blush{
		Finders: []blush.Finder{match},
		Reader:  r,
	}

	tcs := []struct {
		name    string
		length  int
		want    string
		wantLen int
		wantErr error
	}{
		{"line1", 5, "line1", 5, nil},
		{"\nli", 3, "\nli", 3, nil},
		{"ne2\nline", 8, "ne2\nline", 8, nil},
		{"3\nline4", 7, "3\nline4", 7, nil},
		{"\n", 1, "\n", 1, nil},
		{"finish", 10, "", 0, io.EOF},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			p := make([]byte, tc.length)
			n, err := b.Read(p)
			if err != tc.wantErr {
				t.Errorf("err = %v, want %v", err, tc.wantErr)
			}
			if n != tc.wantLen {
				t.Errorf("n = %d, want %d", n, tc.wantLen)
			}
			if string(bytes.Trim(p, "\x00")) != tc.want {
				t.Errorf("p = `%v: %s`, want `%v: %s`", p, p, tc.want, tc.want)
			}
		})
	}
}
