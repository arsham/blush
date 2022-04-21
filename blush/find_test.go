package blush_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/arsham/blush/blush"
)

type colourer interface {
	Colour() blush.Colour
}

// nolint:misspell // it's ok.
func TestNewLocatorColours(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		name   string
		colour string
		want   blush.Colour
	}{
		{"default", "", blush.DefaultColour},
		{"garbage", "sdsds", blush.DefaultColour},
		{"blue", "blue", blush.Blue},
		{"blue short", "b", blush.Blue},
		{"red", "red", blush.Red},
		{"red short", "r", blush.Red},
		{"blue", "blue", blush.Blue},
		{"blue short", "b", blush.Blue},
		{"green", "green", blush.Green},
		{"green short", "g", blush.Green},
		{"black", "black", blush.Black},
		{"black short", "bl", blush.Black},
		{"white", "white", blush.White},
		{"white short", "w", blush.White},
		{"cyan", "cyan", blush.Cyan},
		{"cyan short", "cy", blush.Cyan},
		{"magenta", "magenta", blush.Magenta},
		{"magenta short", "mg", blush.Magenta},
		{"yellow", "yellow", blush.Yellow},
		{"yellow short", "yl", blush.Yellow},
		{"no colour", "no-colour", blush.NoColour},
		{"no colour american", "no-color", blush.NoColour},

		{"hash 000", "#000", blush.Colour{Foreground: blush.RGB{R: 0, G: 0, B: 0}, Background: blush.NoRGB}},
		{"hash 666", "#666", blush.Colour{Foreground: blush.RGB{R: 102, G: 102, B: 102}, Background: blush.NoRGB}},
		{"hash 000000", "#000000", blush.Colour{Foreground: blush.RGB{R: 0, G: 0, B: 0}, Background: blush.NoRGB}},
		{"hash 666666", "#666666", blush.Colour{Foreground: blush.RGB{R: 102, G: 102, B: 102}, Background: blush.NoRGB}},
		{"hash FFF", "#FFF", blush.Colour{Foreground: blush.RGB{R: 255, G: 255, B: 255}, Background: blush.NoRGB}},
		{"hash fff", "#fff", blush.Colour{Foreground: blush.RGB{R: 255, G: 255, B: 255}, Background: blush.NoRGB}},
		{"hash ffffff", "#ffffff", blush.Colour{Foreground: blush.RGB{R: 255, G: 255, B: 255}, Background: blush.NoRGB}},
		{"hash ababAB", "#ababAB", blush.Colour{Foreground: blush.RGB{R: 171, G: 171, B: 171}, Background: blush.NoRGB}},
		{"hash hhhhhh", "#hhhhhh", blush.DefaultColour},
		{"hash aaaaaaa", "#aaaaaaa", blush.DefaultColour},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			l := blush.NewLocator(tc.colour, "aaa", false)
			c, ok := l.(colourer)
			assert.True(t, ok)
			assert.Equal(t, tc.want, c.Colour())
		})
	}
}

func TestNewLocatorExact(t *testing.T) {
	t.Parallel()
	l := blush.NewLocator("", "aaa", false)
	assert.IsType(t, blush.Exact{}, l)
	l = blush.NewLocator("", "*aaa", false)
	assert.IsType(t, blush.Exact{}, l)
}

func TestNewLocatorIexact(t *testing.T) {
	t.Parallel()
	l := blush.NewLocator("", "aaa", true)
	assert.IsType(t, blush.Iexact{}, l)
	l = blush.NewLocator("", "*aaa", true)
	assert.IsType(t, blush.Iexact{}, l)
}

func TestNewLocatorRx(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		name    string
		input   string
		matches []string
	}{
		{"empty", "^$", []string{""}},
		{"starts with", "^aaa", []string{"aaa", "aaa sss"}},
		{"ends with", "aaa$", []string{"aaa", "sss aaa"}},
		{"with star", "blah blah.*", []string{"blah blah", "aa blah blah aa"}},
		{"with curly brackets", "a{3}", []string{"aaa", "aa aaa aa"}},
		{"with brackets", "[ab]", []string{"kjhadf", "kjlrbrlkj", "sdbsdha"}},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			l := blush.NewLocator("", tc.input, false)
			assert.IsType(t, blush.Rx{}, l)
			l = blush.NewLocator("", tc.input, false)
			assert.IsType(t, blush.Rx{}, l)
		})
	}
}

func TestNewLocatorRxColours(t *testing.T) {
	t.Parallel()
	rx := blush.NewLocator("b", "a{3}", false)
	want := "this " + blush.Colourise("aaa", blush.Blue) + "meeting"
	got, ok := rx.Find("this aaameeting")
	assert.Equal(t, want, got)
	assert.True(t, ok)
}

func TestExactNotFound(t *testing.T) {
	t.Parallel()
	l := blush.NewExact("nooooo", blush.NoColour)
	got, ok := l.Find("yessss")
	assert.Empty(t, got)
	assert.False(t, ok)
}

func TestExactFind(t *testing.T) {
	t.Parallel()
	l := blush.NewExact("nooooo", blush.NoColour)
	got, ok := l.Find("yessss")
	assert.Empty(t, got)
	assert.False(t, ok)

	tcs := []struct {
		name   string
		search string
		colour blush.Colour
		input  string
		want   string
		wantOk bool
	}{
		{"exact no colour", "aaa", blush.NoColour, "aaa", "aaa", true},
		{"exact not found", "aaaa", blush.NoColour, "aaa", "", false},
		{"some parts no colour", "aaa", blush.NoColour, "bb aaa bb", "bb aaa bb", true},
		{"exact blue", "aaa", blush.Blue, "aaa", blush.Colourise("aaa", blush.Blue), true},
		{"some parts blue", "aaa", blush.Blue, "bb aaa bb", "bb " + blush.Colourise("aaa", blush.Blue) + " bb", true},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			l := blush.NewExact(tc.search, tc.colour)
			got, ok := l.Find(tc.input)
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantOk, ok)
		})
	}
}

func TestRxNotFound(t *testing.T) {
	t.Parallel()
	l := blush.NewRx(regexp.MustCompile("no{5}"), blush.NoColour)
	got, ok := l.Find("yessss")
	assert.Empty(t, got)
	assert.False(t, ok)
}

func TestRxFind(t *testing.T) {
	t.Parallel()
	l := blush.NewRx(regexp.MustCompile("no{5}"), blush.NoColour)
	got, ok := l.Find("yessss")
	assert.Empty(t, got)
	assert.False(t, ok)

	tcs := []struct {
		name   string
		search string
		colour blush.Colour
		input  string
		want   string
		wantOk bool
	}{
		{"exact no colour", "(^aaa$)", blush.NoColour, "aaa", "aaa", true},
		{"exact not found", "(^aa$)", blush.NoColour, "aaa", "", false},
		{"some parts no colour", "(aaa)", blush.NoColour, "bb aaa bb", "bb aaa bb", true},
		{"some parts not matched", "(Aaa)", blush.NoColour, "bb aaa bb", "", false},
		{"exact blue", "(aaa)", blush.Blue, "aaa", blush.Colourise("aaa", blush.Blue), true},
		{"some parts blue", "(aaa)", blush.Blue, "bb aaa bb", "bb " + blush.Colourise("aaa", blush.Blue) + " bb", true},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			l := blush.NewRx(regexp.MustCompile(tc.search), tc.colour)
			got, ok := l.Find(tc.input)
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantOk, ok)
			assert.Equal(t, tc.colour, l.Colour())
		})
	}
}

func TestIexactNotFound(t *testing.T) {
	t.Parallel()
	l := blush.NewIexact("nooooo", blush.NoColour)
	got, ok := l.Find("yessss")
	assert.Empty(t, got)
	assert.False(t, ok)
}

func TestIexact(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		name   string
		search string
		colour blush.Colour
		input  string
		want   string
		wantOk bool
	}{
		{"exact no colour", "aaa", blush.NoColour, "aaa", "aaa", true},
		{"exact not found", "aaaa", blush.NoColour, "aaa", "", false},
		{"i exact no colour", "AAA", blush.NoColour, "aaa", "aaa", true},
		{"some parts no colour", "aaa", blush.NoColour, "bb aaa bb", "bb aaa bb", true},
		{"i some parts no colour", "AAA", blush.NoColour, "bb aaa bb", "bb aaa bb", true},
		{"exact blue", "aaa", blush.Blue, "aaa", blush.Colourise("aaa", blush.Blue), true},
		{"i exact blue", "AAA", blush.Blue, "aaa", blush.Colourise("aaa", blush.Blue), true},
		{"some parts blue", "aaa", blush.Blue, "bb aaa bb", "bb " + blush.Colourise("aaa", blush.Blue) + " bb", true},
		{"i some parts blue", "AAA", blush.Blue, "bb aaa bb", "bb " + blush.Colourise("aaa", blush.Blue) + " bb", true},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			l := blush.NewIexact(tc.search, tc.colour)
			got, ok := l.Find(tc.input)
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantOk, ok)
			assert.Equal(t, tc.colour, l.Colour())
		})
	}
}

func TestRxInsensitiveFind(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		name   string
		search string
		colour string
		input  string
		want   string
		wantOk bool
	}{
		{"exact no colour", "^AAA$", "no-colour", "aaa", "aaa", true},
		{"exact not found", "^AA$", "no-colour", "aaa", "", false},
		{"some words no colour", `AAA*`, "no-colour", "bb aaa bb", "bb aaa bb", true},
		{"exact blue", "^AAA$", "b", "aaa", blush.Colourise("aaa", blush.Blue), true},
		{"default colour", "^AAA$", "", "aaa", blush.Colourise("aaa", blush.DefaultColour), true},
		{"some words blue", "AAA?", "b", "bb aaa bb", "bb " + blush.Colourise("aaa", blush.Blue) + " bb", true},
		{"some words blue long", "AAA?", "blue", "bb aaa bb", "bb " + blush.Colourise("aaa", blush.Blue) + " bb", true},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			l := blush.NewLocator(tc.colour, tc.search, true)
			assert.IsType(t, blush.Rx{}, l)

			got, ok := l.Find(tc.input)
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantOk, ok)
		})
	}

	rx := blush.NewLocator("b", "A{3}", true)
	want := "this " + blush.Colourise("aaa", blush.Blue) + "meeting"
	got, ok := rx.Find("this aaameeting")
	assert.Equal(t, want, got)
	assert.True(t, ok)
}

func TestColourGroup(t *testing.T) {
	t.Parallel()
	for i := 0; i < 100; i++ {
		l := blush.NewLocator(fmt.Sprintf("b%d", i), "aaa", false)
		e, ok := l.(blush.Exact)
		assert.True(t, ok)
		assert.Equal(t, blush.FgBlue, e.Colour().Foreground)
		assert.NotEqual(t, blush.NoRGB, e.Colour().Background)
		assert.NotEqual(t, e.Colour().Background, e.Colour().Foreground)
	}
}

func TestColourNewGroup(t *testing.T) {
	t.Parallel()
	l := blush.NewLocator("b1", "aaa", false)
	assert.IsType(t, blush.Exact{}, l)

	e := l.(blush.Exact)
	c1 := e.Colour()
	l = blush.NewLocator("b1", "aaa", false)
	e = l.(blush.Exact)
	assert.EqualValues(t, c1, e.Colour())

	l = blush.NewLocator("b2", "aaa", false)
	e = l.(blush.Exact)
	assert.Equal(t, blush.FgBlue, e.Colour().Foreground)
	assert.NotEqual(t, c1, e.Colour())
	assert.NotEqual(t, c1.Background, e.Colour().Background)
}
