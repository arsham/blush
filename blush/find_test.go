package blush_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/arsham/blush/blush"
)

type colourer interface {
	Colour() blush.Colour
}

func TestNewLocatorColours(t *testing.T) {
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
		t.Run(tc.name, func(t *testing.T) {
			l := blush.NewLocator(tc.colour, "aaa", false)
			c, ok := l.(colourer)
			if !ok {
				t.Fatalf("%v does not implement Colour() method", l)
			}
			if c.Colour() != tc.want {
				t.Errorf("%s: c.Colour() = %#v, want %#v", tc.colour, c.Colour(), tc.want)
			}
		})
	}
}

func TestNewLocatorExact(t *testing.T) {
	l := blush.NewLocator("", "aaa", false)
	if _, ok := l.(blush.Exact); !ok {
		t.Errorf("l = %T, want *blush.Exact", l)
	}
	l = blush.NewLocator("", "*aaa", false)
	if _, ok := l.(blush.Exact); !ok {
		t.Errorf("l = %T, want *blush.Exact", l)
	}
}

func TestNewLocatorIexact(t *testing.T) {
	l := blush.NewLocator("", "aaa", true)
	if _, ok := l.(blush.Iexact); !ok {
		t.Errorf("l = %T, want *blush.Iexact", l)
	}
	l = blush.NewLocator("", "*aaa", true)
	if _, ok := l.(blush.Iexact); !ok {
		t.Errorf("l = %T, want *blush.Iexact", l)
	}
}

func TestNewLocatorRx(t *testing.T) {
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
		t.Run(tc.name, func(t *testing.T) {
			l := blush.NewLocator("", tc.input, false)
			if _, ok := l.(blush.Rx); !ok {
				t.Errorf("l = %T, want *blush.Rx", l)
			}
			l = blush.NewLocator("", tc.input, false)
			if _, ok := l.(blush.Rx); !ok {
				t.Errorf("l = %T, want *blush.Rx", l)
			}
		})
	}
}

func TestNewLocatorRxColours(t *testing.T) {
	rx := blush.NewLocator("b", "a{3}", false)
	want := "this " + blush.Colourise("aaa", blush.Blue) + "meeting"
	got, ok := rx.Find("this aaameeting")
	if got != want {
		t.Errorf("got = `%s`, want `%s`", got, want)
	}
	if !ok {
		t.Error("ok = false, want true")
	}
}

func TestExactNotFound(t *testing.T) {
	l := blush.NewExact("nooooo", blush.NoColour)
	got, ok := l.Find("yessss")
	if got != "" {
		t.Errorf("got = %s, want `%s`", got, "")
	}
	if ok {
		t.Error("ok = true, want false")
	}
}

func TestExactFind(t *testing.T) {
	l := blush.NewExact("nooooo", blush.NoColour)
	got, ok := l.Find("yessss")
	if got != "" {
		t.Errorf("got = %s, want `%s`", got, "")
	}
	if ok {
		t.Error("ok = true, want false")
	}

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
		t.Run(tc.name, func(t *testing.T) {
			l := blush.NewExact(tc.search, tc.colour)
			got, ok := l.Find(tc.input)
			if got != tc.want {
				t.Errorf("got = `%s`, want `%s`", got, tc.want)
			}
			if ok != tc.wantOk {
				t.Errorf("ok = %t, want %t", ok, tc.wantOk)
			}
		})
	}
}

func TestRxNotFound(t *testing.T) {
	l := blush.NewRx(regexp.MustCompile("nooooo"), blush.NoColour)
	got, ok := l.Find("yessss")
	if got != "" {
		t.Errorf("got = %s, want `%s`", got, "")
	}
	if ok {
		t.Error("ok = true, want false")
	}
}

func TestRxFind(t *testing.T) {
	l := blush.NewRx(regexp.MustCompile("nooooo"), blush.NoColour)
	got, ok := l.Find("yessss")
	if got != "" {
		t.Errorf("got = %s, want `%s`", got, "")
	}
	if ok {
		t.Error("ok = true, want false")
	}

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
		t.Run(tc.name, func(t *testing.T) {
			l := blush.NewRx(regexp.MustCompile(tc.search), tc.colour)
			got, ok := l.Find(tc.input)
			if got != tc.want {
				t.Errorf("got = `%s`, want `%s`", got, tc.want)
			}
			if ok != tc.wantOk {
				t.Errorf("ok = %t, want %t", ok, tc.wantOk)
			}
			if l.Colour() != tc.colour {
				t.Errorf("l.Colour() = %v, want %v", l.Colour(), tc.colour)
			}
		})
	}
}

func TestIexactNotFound(t *testing.T) {
	l := blush.NewIexact("nooooo", blush.NoColour)
	got, ok := l.Find("yessss")
	if got != "" {
		t.Errorf("got = %s, want `%s`", got, "")
	}
	if ok {
		t.Error("ok = true, want false")
	}
}

func TestIexact(t *testing.T) {
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
		t.Run(tc.name, func(t *testing.T) {
			l := blush.NewIexact(tc.search, tc.colour)
			got, ok := l.Find(tc.input)
			if got != tc.want {
				t.Errorf("got = `%s`, want `%s`", got, tc.want)
			}
			if ok != tc.wantOk {
				t.Errorf("ok = %t, want %t", ok, tc.wantOk)
			}
			if l.Colour() != tc.colour {
				t.Errorf("l.Colour() = %v, want %v", l.Colour(), tc.colour)
			}
		})
	}
}

func TestRxInsensitiveFind(t *testing.T) {
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
		t.Run(tc.name, func(t *testing.T) {
			l := blush.NewLocator(tc.colour, tc.search, true)
			if _, ok := l.(blush.Rx); !ok {
				t.Fatalf("l = %T, want blush.Rx", l)
			}
			got, ok := l.Find(tc.input)
			if got != tc.want {
				t.Errorf("got = `%s`, want `%s`", got, tc.want)
			}
			if ok != tc.wantOk {
				t.Errorf("ok = %t, want %t", ok, tc.wantOk)
			}
		})
	}

	rx := blush.NewLocator("b", "A{3}", true)
	want := "this " + blush.Colourise("aaa", blush.Blue) + "meeting"
	got, ok := rx.Find("this aaameeting")
	if got != want {
		t.Errorf("got = `%s`, want `%s`", got, want)
	}
	if !ok {
		t.Error("ok = false, want true")
	}
}

func TestColourGroup(t *testing.T) {
	for i := 0; i < 100; i++ {
		l := blush.NewLocator(fmt.Sprintf("b%d", i), "aaa", false)
		e, ok := l.(blush.Exact)
		if !ok {
			t.Fatalf("fact check: l = %T, want blush.Exact", l)
		}
		if e.Colour().Foreground != blush.FgBlue {
			t.Errorf("e.Colour().Foreground = %v, want %v", e.Colour().Foreground, blush.FgBlue)
		}
		if e.Colour().Background == blush.NoRGB {
			t.Errorf("e.Colour().Background = %v, want a different RGB", e.Colour().Background)
		}
		if e.Colour().Foreground == e.Colour().Background {
			t.Errorf("e.Colour().Foreground = %v, e.Colour().Background = %v: want a different colour", e.Colour().Foreground, e.Colour().Background)
		}
	}
}

func TestColourNewGroup(t *testing.T) {
	var (
		e  blush.Exact
		ok bool
	)
	l := blush.NewLocator("b1", "aaa", false)
	if e, ok = l.(blush.Exact); !ok {
		t.Fatalf("fact check: l = %T, want blush.Exact", l)
	}

	c1 := e.Colour()
	l = blush.NewLocator("b1", "aaa", false)
	e = l.(blush.Exact)
	if e.Colour() != c1 {
		t.Errorf("e.Colour() = %v, want %v", e.Colour(), c1)
	}

	l = blush.NewLocator("b2", "aaa", false)
	e = l.(blush.Exact)
	if e.Colour().Foreground != blush.FgBlue {
		t.Errorf("e.Colour().Foreground = %v, want %v", e.Colour().Foreground, blush.FgBlue)
	}
	if e.Colour() == c1 {
		t.Errorf("e.Colour() = %v, want a different Colour", e.Colour())
	}
	if e.Colour().Background == c1.Background {
		t.Errorf("e.Colour().Background = %v, c1.Background = %v, want a different Colour", e.Colour().Background, c1.Background)
	}
}
