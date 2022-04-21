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

	"github.com/alecthomas/assert"
	"github.com/arsham/blush/blush"
	"github.com/arsham/blush/internal/reader"
)

func TestWriteToErrors(t *testing.T) {
	t.Parallel()
	w := &bytes.Buffer{}
	e := errors.New("something")
	nn := 10
	bw := &badWriter{
		writeFunc: func([]byte) (int, error) {
			return nn, e
		},
	}
	getReader := func() io.ReadCloser {
		return io.NopCloser(bytes.NewBufferString("something"))
	}
	tcs := []struct {
		name    string
		b       *blush.Blush
		writer  io.Writer
		wantN   int
		wantErr error
	}{
		{"no input", &blush.Blush{}, w, 0, reader.ErrNoReader},
		{"no writer", &blush.Blush{Reader: getReader()}, nil, 0, blush.ErrNoWriter},
		{"bad writer", &blush.Blush{Reader: getReader(), NoCut: true}, bw, nn, e},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tc.b.Finders = []blush.Finder{blush.NewExact("", blush.NoColour)}
			n, err := tc.b.WriteTo(tc.writer)
			assert.Error(t, err)
			assert.EqualValues(t, tc.wantN, n)
			assert.EqualError(t, err, tc.wantErr.Error())
		})
	}
}

func TestWriteToNoMatch(t *testing.T) {
	t.Parallel()
	pwd, err := os.Getwd()
	assert.NoError(t, err)
	location := path.Join(pwd, "testdata")
	r, err := reader.NewMultiReader(reader.WithPaths([]string{location}, true))
	assert.NoError(t, err)
	b := &blush.Blush{
		Reader:  r,
		Finders: []blush.Finder{blush.NewExact("SHOULDNOTFINDTHISONE", blush.NoColour)},
	}
	buf := &bytes.Buffer{}
	n, err := b.WriteTo(buf)
	assert.NoError(t, err)
	assert.Zero(t, n)
	assert.Zero(t, buf.Len())
}

func TestWriteToMatchNoColourPlain(t *testing.T) {
	t.Parallel()
	match := "TOKEN"
	pwd, err := os.Getwd()
	assert.NoError(t, err)
	location := path.Join(pwd, "testdata")
	r, err := reader.NewMultiReader(reader.WithPaths([]string{location}, true))
	assert.NoError(t, err)

	b := &blush.Blush{
		Reader:  r,
		Finders: []blush.Finder{blush.NewExact(match, blush.NoColour)},
	}

	buf := &bytes.Buffer{}
	n, err := b.WriteTo(buf)
	assert.NoError(t, err)
	assert.NotZero(t, buf.Len())
	assert.EqualValues(t, buf.Len(), n)

	assert.Contains(t, buf.String(), match)
	assert.NotContains(t, buf.String(), leaveMeHere)
}

func TestWriteToMatchColour(t *testing.T) {
	t.Parallel()
	match := blush.Colourise("TOKEN", blush.Blue)
	pwd, err := os.Getwd()
	assert.NoError(t, err)
	location := path.Join(pwd, "testdata")
	r, err := reader.NewMultiReader(reader.WithPaths([]string{location}, true))
	assert.NoError(t, err)
	b := &blush.Blush{
		Reader:  r,
		Finders: []blush.Finder{blush.NewExact("TOKEN", blush.Blue)},
	}

	buf := &bytes.Buffer{}
	n, err := b.WriteTo(buf)
	assert.NoError(t, err)
	assert.NotZero(t, buf.Len())
	assert.EqualValues(t, buf.Len(), n)

	assert.Contains(t, buf.String(), match)
	assert.NotContains(t, buf.String(), leaveMeHere)
}

func TestWriteToMatchCountColour(t *testing.T) {
	t.Parallel()
	pwd, err := os.Getwd()
	assert.NoError(t, err)

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
			assert.NoError(t, err)

			match := blush.Colourise(tc.name, blush.Red)
			b := &blush.Blush{
				Reader:  r,
				Finders: []blush.Finder{blush.NewExact(tc.name, blush.Red)},
			}

			buf := &bytes.Buffer{}
			n, err := b.WriteTo(buf)
			assert.NoError(t, err)
			assert.EqualValues(t, buf.Len(), n)
			count := strings.Count(buf.String(), match)
			assert.EqualValues(t, tc.count, count)
			assert.NotContains(t, buf.String(), leaveMeHere)
		})
	}
}

func TestWriteToMultiColour(t *testing.T) {
	t.Parallel()
	two := blush.Colourise("TWO", blush.Magenta)
	three := blush.Colourise("THREE", blush.Red)
	pwd, err := os.Getwd()
	assert.NoError(t, err)
	location := path.Join(pwd, "testdata")
	r, err := reader.NewMultiReader(reader.WithPaths([]string{location}, true))
	assert.NoError(t, err)
	b := &blush.Blush{
		Reader: r,
		Finders: []blush.Finder{
			blush.NewExact("TWO", blush.Magenta),
			blush.NewExact("THREE", blush.Red),
		},
	}

	buf := &bytes.Buffer{}
	n, err := b.WriteTo(buf)
	assert.NoError(t, err)
	assert.NotZero(t, buf.Len())
	assert.EqualValues(t, buf.Len(), n)
	count := strings.Count(buf.String(), two)
	assert.EqualValues(t, 2*3, count)
	count = strings.Count(buf.String(), three)
	assert.EqualValues(t, 3*3, count)
	if strings.Contains(buf.String(), leaveMeHere) {
		t.Errorf("didn't expect to see %s", leaveMeHere)
	}
}

func TestWriteToMultiColourColourMode(t *testing.T) {
	t.Parallel()
	two := blush.Colourise("TWO", blush.Magenta)
	three := blush.Colourise("THREE", blush.Red)
	pwd, err := os.Getwd()
	assert.NoError(t, err)
	location := path.Join(pwd, "testdata")
	r, err := reader.NewMultiReader(reader.WithPaths([]string{location}, true))
	assert.NoError(t, err)
	b := &blush.Blush{
		Reader: r,
		NoCut:  true,
		Finders: []blush.Finder{
			blush.NewExact("TWO", blush.Magenta),
			blush.NewExact("THREE", blush.Red),
		},
	}

	buf := &bytes.Buffer{}
	n, err := b.WriteTo(buf)
	assert.NoError(t, err)

	assert.NotZero(t, buf.Len())
	assert.EqualValues(t, buf.Len(), n)
	count := strings.Count(buf.String(), two)
	assert.EqualValues(t, 2*3, count)
	count = strings.Count(buf.String(), three)
	assert.EqualValues(t, 3*3, count)
	count = strings.Count(buf.String(), leaveMeHere)
	assert.EqualValues(t, 1, count)
}

func TestWriteToMultipleMatchInOneLine(t *testing.T) {
	t.Parallel()
	line1 := "this is an example\n"
	line2 := "someone should find this line\n"
	input1 := bytes.NewBuffer([]byte(line1))
	input2 := bytes.NewBuffer([]byte(line2))
	r := io.NopCloser(io.MultiReader(input1, input2))
	match := fmt.Sprintf(
		"someone %s find %s line",
		blush.Colourise("should", blush.Red),
		blush.Colourise("this", blush.Magenta),
	)
	out := &bytes.Buffer{}

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
	assert.EqualValues(t, match, example)
}

func TestBlushClosesReader(t *testing.T) {
	t.Parallel()
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
	assert.NoError(t, err)
	assert.True(t, called, "didn't close the reader")
}

func TestBlushReadOneStream(t *testing.T) {
	t.Parallel()
	input := bytes.NewBuffer([]byte("one two three four"))
	match := blush.NewExact("three", blush.Blue)
	r := io.NopCloser(input)
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
			assert.True(t, errors.Is(err, tc.wantErr))
			assert.EqualValues(t, tc.wantLen, n)
			assert.EqualValues(t, tc.wantP, tc.p)
		})
	}
}

func TestBlushReadTwoStreams(t *testing.T) {
	t.Parallel()
	b1 := []byte("one for all\n")
	b2 := []byte("all for one\n")
	input1 := bytes.NewBuffer(b1)
	input2 := bytes.NewBuffer(b2)
	match := blush.NewExact("one", blush.Blue)
	r := io.NopCloser(io.MultiReader(input1, input2))
	b := &blush.Blush{
		Finders: []blush.Finder{match},
		Reader:  r,
	}
	defer b.Close()

	buf := &bytes.Buffer{}
	n, err := buf.ReadFrom(b)
	assert.NoError(t, err)
	expectLen := len(b1) + len(b2) - len("one")*2 + len(match.String())*2
	assert.EqualValues(t, expectLen, n)
	expectStr := fmt.Sprintf("%s%s",
		strings.Replace(string(b1), "one", match.String(), 1),
		strings.Replace(string(b2), "one", match.String(), 1),
	)
	assert.EqualValues(t, expectStr, buf.String())
}

func TestBlushReadHalfWay(t *testing.T) {
	t.Parallel()
	b1 := []byte("one for all\n")
	b2 := []byte("all for one\n")
	input1 := bytes.NewBuffer(b1)
	input2 := bytes.NewBuffer(b2)
	match := blush.NewExact("one", blush.Blue)
	r := io.NopCloser(io.MultiReader(input1, input2))
	b := &blush.Blush{
		Finders: []blush.Finder{match},
		Reader:  r,
	}
	p := make([]byte, len(b1))
	_, err := b.Read(p)
	assert.NoError(t, err)
	n, err := b.Read(p)
	assert.Len(t, b1, n)
	assert.NoError(t, err)
}

func TestBlushReadOnClosed(t *testing.T) {
	t.Parallel()
	b1 := []byte("one for all\n")
	b2 := []byte("all for one\n")
	input1 := bytes.NewBuffer(b1)
	input2 := bytes.NewBuffer(b2)
	match := blush.NewExact("one", blush.Blue)
	r := io.NopCloser(io.MultiReader(input1, input2))
	b := &blush.Blush{
		Finders: []blush.Finder{match},
		Reader:  r,
	}
	p := make([]byte, len(b1))
	_, err := b.Read(p)
	assert.NoError(t, err)
	err = b.Close()
	assert.NoError(t, err)

	n, err := b.Read(p)
	assert.True(t, errors.Is(err, blush.ErrClosed))
	assert.Zero(t, n)
}

func TestBlushReadLongOneLineText(t *testing.T) {
	t.Parallel()
	head := strings.Repeat("a", 10000)
	tail := strings.Repeat("b", 10000)
	input := bytes.NewBuffer([]byte(head + " FINDME " + tail))
	match := blush.NewExact("FINDME", blush.Blue)
	r := io.NopCloser(input)
	b := &blush.Blush{
		Finders: []blush.Finder{match},
		Reader:  r,
	}
	p := make([]byte, 20)
	_, err := b.Read(p)
	assert.NoError(t, err)
	err = b.Close()
	assert.NoError(t, err)
	n, err := b.Read(p)
	assert.True(t, errors.Is(err, blush.ErrClosed))
	assert.Zero(t, n)
}

func TestPrintName(t *testing.T) {
	t.Parallel()
	line1 := "line one\n"
	line2 := "line two\n"
	r1 := io.NopCloser(bytes.NewBuffer([]byte(line1)))
	r2 := io.NopCloser(bytes.NewBuffer([]byte(line2)))
	name1 := "reader1"
	name2 := "reader2"
	r, err := reader.NewMultiReader(
		reader.WithReader(name1, r1),
		reader.WithReader(name2, r2),
	)
	assert.NoError(t, err)
	b := blush.Blush{
		Reader:       r,
		Finders:      []blush.Finder{blush.NewExact("line", blush.NoColour)},
		WithFileName: true,
	}
	buf := &bytes.Buffer{}
	n, err := b.WriteTo(buf)
	assert.NoError(t, err)
	total := len(line1+line2+name1+name2) + len(blush.Separator)*2
	assert.EqualValues(t, total, n)

	s := strings.Split(buf.String(), "\n")
	assert.Contains(t, s[0], name1)
	assert.Contains(t, s[1], name2)
}

// testing stdin should not print the name
func TestStdinPrintName(t *testing.T) {
	t.Parallel()
	input := "line one"
	oldStdin := os.Stdin
	f, err := ioutil.TempFile("", "blush_stdin")
	assert.NoError(t, err)
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
	buf := &bytes.Buffer{}
	_, err = b.WriteTo(buf)
	assert.NoError(t, err)

	assert.EqualValues(t, input, buf.String())
	assert.NotContains(t, buf.String(), f.Name())
}

func TestPrintFileName(t *testing.T) {
	t.Parallel()
	p, err := ioutil.TempDir("", "blush_name")
	assert.NoError(t, err)
	f1, err := ioutil.TempFile(p, "blush_name")
	assert.NoError(t, err)
	f2, err := ioutil.TempFile(p, "blush_name")
	assert.NoError(t, err)
	defer func() {
		if err = os.RemoveAll(p); err != nil {
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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r, err := reader.NewMultiReader(
				reader.WithPaths([]string{p}, false),
			)
			assert.NoError(t, err)
			b := blush.Blush{
				Reader:       r,
				Finders:      []blush.Finder{blush.NewExact("line", blush.NoColour)},
				WithFileName: tc.withFilename,
			}
			buf := &bytes.Buffer{}
			n, err := b.WriteTo(buf)
			assert.NoError(t, err)
			assert.EqualValues(t, tc.wantLen, n)
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
	t.Parallel()
	var (
		ret   []byte
		p     = make([]byte, 2)
		count int
		input = "one two three four\nfive six three seven"
	)
	match := blush.NewExact("three", blush.Blue)
	r := io.NopCloser(bytes.NewBufferString(input))
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
		if errors.Is(err, io.EOF) {
			ret = append(ret, p...)
			break
		}
		assert.NoError(t, err)
		ret = append(ret, p...)
	}
	if c := strings.Count(string(ret), "three"); c != 2 {
		t.Errorf("count %s = %d, want %d", "three", c, 2)
	}
	for _, s := range []string{"one", "two", "three", "four", "five", "six", "seven"} {
		assert.Contains(t, string(ret), s)
	}
}

func TestReadMiddleOfMatch(t *testing.T) {
	t.Parallel()
	var (
		search = "aa this aa"
		match  = blush.NewExact("this", blush.Blue)
		p      = make([]byte, (len(search)+len(match.String()))/2)
		ret    []byte
	)
	r := io.NopCloser(bytes.NewBufferString(search))
	b := &blush.Blush{
		Finders: []blush.Finder{match},
		Reader:  r,
	}
	for i := 0; i < 2; i++ {
		_, err := b.Read(p)
		if err != nil && !errors.Is(err, io.EOF) {
			t.Fatal(err)
		}

		ret = append(ret, p...)
	}
	assert.Contains(t, string(ret), "this")
}

func TestReadComplete(t *testing.T) {
	t.Parallel()
	input := "123456789"
	match := blush.NewExact("1", blush.NoColour)
	r := io.NopCloser(bytes.NewBufferString(input))
	b := &blush.Blush{
		Finders: []blush.Finder{match},
		Reader:  r,
	}
	p := make([]byte, 10)
	n, err := b.Read(p)
	assert.Len(t, input, n)
	assert.True(t, errors.Is(err, io.EOF))
	assert.EqualValues(t, input, string(bytes.Trim(p, "\x00")))
	p = make([]byte, 4)
	n, err = b.Read(p)
	assert.Zero(t, n)
	assert.True(t, errors.Is(err, io.EOF))
	assert.Empty(t, string(bytes.Trim(p, "\x00")))
}

func TestReadPartComplete(t *testing.T) {
	t.Parallel()
	input := "123456789"
	match := blush.NewExact("1", blush.NoColour)
	r := io.NopCloser(bytes.NewBufferString(input))
	b := &blush.Blush{
		Finders: []blush.Finder{match},
		Reader:  r,
	}
	p := make([]byte, 3)
	n, err := b.Read(p)
	assert.EqualValues(t, 3, n)
	assert.NoError(t, err)
	assert.EqualValues(t, string(bytes.Trim(p, "\x00")), "123")

	p = make([]byte, 6)
	n, err = b.Read(p)
	assert.NoError(t, err)
	assert.EqualValues(t, 6, n)
	assert.EqualValues(t, string(bytes.Trim(p, "\x00")), "456789")
}

func TestReadPartPartOver(t *testing.T) {
	t.Parallel()
	input := "123456789"
	match := blush.NewExact("1", blush.NoColour)
	r := io.NopCloser(bytes.NewBufferString(input))
	b := &blush.Blush{
		Finders: []blush.Finder{match},
		Reader:  r,
	}
	p := make([]byte, 3)
	n, err := b.Read(p)
	assert.NoError(t, err)
	assert.EqualValues(t, 3, n)
	assert.EqualValues(t, string(bytes.Trim(p, "\x00")), "123")

	p = make([]byte, 3)
	n, err = b.Read(p)
	assert.EqualValues(t, 3, n)
	assert.NoError(t, err)
	assert.EqualValues(t, string(bytes.Trim(p, "\x00")), "456")

	p = make([]byte, 10)
	n, err = b.Read(p)
	assert.EqualValues(t, 3, n)
	assert.True(t, errors.Is(err, io.EOF))
	assert.EqualValues(t, string(bytes.Trim(p, "\x00")), "789")
}

func TestReadMultiLine(t *testing.T) {
	t.Parallel()
	input := "line1\nline2\nline3\nline4\n"
	match := blush.NewExact("l", blush.NoColour)
	r := io.NopCloser(bytes.NewBufferString(input))
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
			assert.True(t, errors.Is(err, tc.wantErr))
			assert.EqualValues(t, tc.wantLen, n)
			assert.EqualValues(t, tc.want, string(bytes.Trim(p, "\x00")))
		})
	}
}

func TestReadWriteToMode(t *testing.T) {
	t.Parallel()
	p := make([]byte, 1)
	r := io.NopCloser(bytes.NewBufferString("input"))
	b := &blush.Blush{
		Finders: []blush.Finder{blush.NewExact("", blush.NoColour)},
		Reader:  r,
	}
	_, err := b.Read(p)
	assert.NoError(t, err)
	_, err = b.WriteTo(&bytes.Buffer{})
	assert.True(t, errors.Is(err, blush.ErrReadWriteMix))

	b = &blush.Blush{
		Finders: []blush.Finder{blush.NewExact("", blush.NoColour)},
		Reader:  r,
	}
	_, err = b.WriteTo(&bytes.Buffer{})
	assert.NoError(t, err)

	_, err = b.Read(p)
	assert.True(t, errors.Is(err, blush.ErrReadWriteMix))
}
