package blush_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
	"testing"

	"github.com/arsham/blush/blush"
)

func BenchmarkColourise(b *testing.B) {
	var (
		shortStr = "jNa8SZ1RPM"
		longStr  = strings.Repeat("ZL5B2kNexCcTPvf9 ", 50)
	)

	// nf = not found
	// cl = coloured
	// ln = long
	bcs := []struct {
		name   string
		input  string
		colour blush.Colour
	}{
		{"nc", shortStr, blush.NoColour},
		{"fg", shortStr, blush.Red},
		{"bg", shortStr, blush.Colour{Foreground: blush.NoRGB, Background: blush.FgRed}},
		{"both cl", shortStr, blush.Colour{Foreground: blush.FgRed, Background: blush.FgRed}},
		{"ln nc", longStr, blush.NoColour},
		{"ln fg", longStr, blush.Colour{Foreground: blush.FgRed, Background: blush.NoRGB}},
		{"ln bg", longStr, blush.Colour{Foreground: blush.NoRGB, Background: blush.FgRed}},
		{"ln both cl", longStr, blush.Colour{Foreground: blush.FgRed, Background: blush.FgRed}},
	}
	for _, bc := range bcs {
		b.Run(bc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				blush.Colourise(bc.input, bc.colour)
			}
		})
	}
}

func BenchmarkNewLocator(b *testing.B) {
	var l blush.Finder
	bcs := []struct {
		name        string
		input       string
		insensitive bool
	}{
		{"plain", "aaa", false},
		{"asterisk", "*aaa", false},
		{"i plain", "aaa", true},
		{"i asterisk", "*aaa", true},
		{"rx empty", "^$", false},
		{"irx empty", "^$", true},
		{"rx starts with", "^aaa", false},
		{"irx starts with", "^aaa", true},
		{"rx ends with", "aaa$", false},
		{"irx ends with", "aaa$", true},
		{"rx with star", "blah blah.*", false},
		{"irx with star", "blah blah.*", true},
		{"rx with curly brackets", "a{3}", false},
		{"irx with curly brackets", "a{3}", true},
		{"rx with brackets", "[ab]", false},
		{"irx with brackets", "[ab]", true},
	}
	for _, bc := range bcs {
		b.Run(bc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				l = blush.NewLocator("", "aaa", bc.insensitive)
				_ = l
			}
		})
	}
}

func BenchmarkFind(b *testing.B) {
	var (
		find         = "FIND"
		blob         = "you should FIND this \n"
		blobLong     = strings.Repeat(blob, 50)
		notFound     = "YWzvnLKGyU "
		notFoundLong = strings.Repeat(notFound, 50)
		got          string
		ok           bool
	)
	// nf = not found
	// nc = no colour
	// cl = coloured
	// ln = long
	bcs := []struct {
		name  string
		l     blush.Finder
		input string
		finds bool
	}{
		{"E nc", blush.NewExact(find, blush.NoColour), blob, true},
		{"E cl", blush.NewExact(find, blush.Blue), blob, true},
		{"E nf", blush.NewExact(find, blush.NoColour), notFound, false},
		{"E nf cl", blush.NewExact(find, blush.Blue), notFound, false},
		{"E ln nc", blush.NewExact(find, blush.NoColour), blobLong, true},
		{"E ln cl", blush.NewExact(find, blush.Blue), blobLong, true},
		{"E ln nf", blush.NewExact(find, blush.NoColour), notFoundLong, false},
		{"E ln nf cl", blush.NewExact(find, blush.Blue), notFoundLong, false},
		{"IE nc", blush.NewIexact(find, blush.NoColour), blob, true},
		{"IE cl", blush.NewIexact(find, blush.Blue), blob, true},
		{"IE nf", blush.NewIexact(find, blush.NoColour), notFound, false},
		{"IE nf cl", blush.NewIexact(find, blush.Blue), notFound, false},
		{"IE ln nc", blush.NewIexact(find, blush.NoColour), blobLong, true},
		{"IE ln cl", blush.NewIexact(find, blush.Blue), blobLong, true},
		{"IE ln nf", blush.NewIexact(find, blush.NoColour), notFoundLong, false},
		{"IE ln nf cl", blush.NewIexact(find, blush.Blue), notFoundLong, false},
		{"rx nc", blush.NewRx(regexp.MustCompile(find), blush.NoColour), blob, true},
		{"rx cl", blush.NewRx(regexp.MustCompile(find), blush.Blue), blob, true},
		{"rx nf", blush.NewRx(regexp.MustCompile(find), blush.NoColour), notFound, false},
		{"rx nf cl", blush.NewRx(regexp.MustCompile(find), blush.Blue), notFound, false},
		{"rx ln nc", blush.NewRx(regexp.MustCompile(find), blush.NoColour), blobLong, true},
		{"rx ln cl", blush.NewRx(regexp.MustCompile(find), blush.Blue), blobLong, true},
		{"rx ln nf", blush.NewRx(regexp.MustCompile(find), blush.NoColour), notFoundLong, false},
		{"rx ln nf cl", blush.NewRx(regexp.MustCompile(find), blush.Blue), notFoundLong, false},
		{"irx nc", blush.NewRx(regexp.MustCompile("(?i)"+find), blush.NoColour), blob, true},
		{"irx cl", blush.NewRx(regexp.MustCompile("(?i)"+find), blush.Blue), blob, true},
		{"irx nf", blush.NewRx(regexp.MustCompile("(?i)"+find), blush.NoColour), notFound, false},
		{"irx nf cl", blush.NewRx(regexp.MustCompile("(?i)"+find), blush.Blue), notFound, false},
		{"irx ln nc", blush.NewRx(regexp.MustCompile("(?i)"+find), blush.NoColour), blobLong, true},
		{"irx ln cl", blush.NewRx(regexp.MustCompile("(?i)"+find), blush.Blue), blobLong, true},
		{"irx ln nf", blush.NewRx(regexp.MustCompile("(?i)"+find), blush.NoColour), notFoundLong, false},
		{"irx ln nf cl", blush.NewRx(regexp.MustCompile("(?i)"+find), blush.Blue), notFoundLong, false},
	}
	for _, bc := range bcs {
		b.Run(bc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				got, ok = bc.l.Find(bc.input)
				if ok != bc.finds {
					b.Fail()
				}
				_ = got
			}
		})
	}
}

// lnr = long reader
// ml = multi
// rds = readers
// mdl = middle
type benchCase struct {
	name   string
	reader func() io.Reader
	match  []blush.Finder
	length int
}

func BenchmarkBlush(b *testing.B) {
	readLarge := func() io.Reader {
		input := bytes.Repeat([]byte("one two three four\n"), 200)
		v := make([]io.Reader, 100)
		for i := 0; i < 100; i++ {
			v[i] = bytes.NewBuffer(input)
		}
		return io.MultiReader(v...)
	}

	readMedium := func() io.Reader {
		input := bytes.Repeat([]byte("one two three four\n"), 200)
		return io.MultiReader(bytes.NewBuffer(input), bytes.NewBuffer(input))
	}

	multiReader100 := func() io.Reader {
		input := []byte("one two three four\n")
		v := make([]io.Reader, 100)
		for i := 0; i < 100; i++ {
			v[i] = bytes.NewBuffer(input)
		}
		return io.MultiReader(v...)
	}

	readMiddle := func() io.Reader {
		p := bytes.Repeat([]byte("one two three four"), 100)
		return bytes.NewBuffer(bytes.Join([][]byte{p, p}, []byte(" MIDDLE ")))
	}

	bcs := []benchCase{
		{
			"short-10",
			func() io.Reader {
				return bytes.NewBuffer([]byte("one two three four"))
			},
			[]blush.Finder{blush.NewExact("three", blush.Blue)},
			10,
		},
		{
			"short-1000",
			func() io.Reader {
				return bytes.NewBuffer([]byte("one two three four"))
			},
			[]blush.Finder{blush.NewExact("three", blush.Blue)},
			1000,
		},
		{
			"lnr-10",
			func() io.Reader {
				return bytes.NewBuffer(bytes.Repeat([]byte("one two three four"), 200))
			},
			[]blush.Finder{blush.NewExact("three", blush.Blue)},
			10,
		},
		{
			"lnr-1000",
			func() io.Reader {
				return bytes.NewBuffer(bytes.Repeat([]byte("one two three four"), 200))
			},
			[]blush.Finder{blush.NewExact("three", blush.Blue)},
			1000,
		},
		{
			"lnr ml lines-10",
			func() io.Reader {
				return bytes.NewBuffer(bytes.Repeat([]byte("one two three four\n"), 200))
			},
			[]blush.Finder{blush.NewExact("three", blush.Blue)},
			10,
		},
		{
			"lnr ml lines-1000",
			func() io.Reader {
				return bytes.NewBuffer(bytes.Repeat([]byte("one two three four\n"), 200))
			},
			[]blush.Finder{blush.NewExact("three", blush.Blue)},
			1000,
		},
		{
			"lnr in mdl-10",
			readMiddle,
			[]blush.Finder{blush.NewExact("MIDDLE", blush.Blue)},
			10,
		},
		{
			"lnr in mdl-1000",
			readMiddle,
			[]blush.Finder{blush.NewExact("MIDDLE", blush.Blue)},
			1000,
		},
		{
			"two rds-10",
			func() io.Reader {
				input := []byte("one two three four\n")
				return io.MultiReader(bytes.NewBuffer(input), bytes.NewBuffer(input))
			},
			[]blush.Finder{blush.NewExact("three", blush.Blue)},
			10,
		},
		{
			"two rds-1000",
			func() io.Reader {
				input := []byte("one two three four\n")
				return io.MultiReader(bytes.NewBuffer(input), bytes.NewBuffer(input))
			},
			[]blush.Finder{blush.NewExact("three", blush.Blue)},
			1000,
		},
		{
			"ln two rds-10",
			readMedium,
			[]blush.Finder{blush.NewExact("three", blush.Blue)},
			10,
		},
		{
			"ln two rds-1000",
			readMedium,
			[]blush.Finder{blush.NewExact("three", blush.Blue)},
			1000,
		},
		{
			"100 rds-10",
			multiReader100,
			[]blush.Finder{blush.NewExact("three", blush.Blue)},
			10,
		},
		{
			"100 rds-1000",
			multiReader100,
			[]blush.Finder{blush.NewExact("three", blush.Blue)},
			1000,
		},
		{
			"ln 100 rds-10",
			readLarge,
			[]blush.Finder{blush.NewExact("three", blush.Blue)},
			10,
		},
		{
			"ln 100 rds-1000",
			readLarge,
			[]blush.Finder{blush.NewExact("three", blush.Blue)},
			1000,
		},
	}
	for _, bc := range bcs {
		b.Run("read_"+bc.name, func(b *testing.B) {
			benchmarkRead(b, bc)
		})
		b.Run("writeTo_"+bc.name, func(b *testing.B) {
			benchmarkWriteTo(b, bc)
		})
	}
}

func benchmarkRead(b *testing.B, bc benchCase) {
	p := make([]byte, bc.length)
	for i := 0; i < b.N; i++ {
		input := ioutil.NopCloser(bc.reader())
		bl := &blush.Blush{
			Finders: bc.match,
			Reader:  input,
		}
		for {
			_, err := bl.Read(p)
			if err != nil {
				if err != io.EOF {
					b.Errorf("err = %v", err)
				}
				break
			}
		}
	}
}

func benchmarkWriteTo(b *testing.B, bc benchCase) {
	for i := 0; i < b.N; i++ {
		input := ioutil.NopCloser(bc.reader())
		bl := &blush.Blush{
			Finders: bc.match,
			Reader:  input,
		}
		buf := new(bytes.Buffer)
		n, err := io.Copy(buf, bl)
		if n == 0 {
			b.Errorf("b.Read(): n = 0, want some read")
		}
		if err != nil {
			b.Error(err)
		}
	}
}
