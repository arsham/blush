package blush_test

import (
	"regexp"
	"testing"

	"github.com/arsham/blush/blush"
)

func TestNewLocatorExact(t *testing.T) {
	l := blush.NewLocator("aaa")
	if _, ok := l.(blush.Exact); !ok {
		t.Errorf("l = %T, want *blush.Exact", l)
	}
	l = blush.NewLocator("*aaa")
	if _, ok := l.(blush.Exact); !ok {
		t.Errorf("l = %T, want *blush.Exact", l)
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
			l := blush.NewLocator(tc.input)
			if _, ok := l.(blush.Rx); !ok {
				t.Errorf("l = %T, want *blush.Rx", l)
			}
		})
	}
}

func TestExactFind(t *testing.T) {
	l := blush.Exact("nooooo")
	got, ok := l.Find("yessss", blush.NoColour)
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
	}{
		{"exact no colour", "aaa", blush.NoColour, "aaa", "aaa"},
		{"some parts no colour", "aaa", blush.NoColour, "bb aaa bb", "bb aaa bb"},
		{"exact blue", "aaa", blush.FgBlue, "aaa", blush.Colourise("aaa", blush.FgBlue)},
		{"some parts blue", "aaa", blush.FgBlue, "bb aaa bb", "bb " + blush.Colourise("aaa", blush.FgBlue) + " bb"},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			l := blush.Exact(tc.search)
			got, ok := l.Find(tc.input, tc.colour)
			if got != tc.want {
				t.Errorf("got = `%s`, want `%s`", got, tc.want)
			}
			if !ok {
				t.Error("ok = false, want true")
			}
		})
	}
}

func TestRxFind(t *testing.T) {
	l := blush.Rx{regexp.MustCompile("nooooo")}
	got, ok := l.Find("yessss", blush.NoColour)
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
	}{
		{"exact no colour", "(^aaa$)", blush.NoColour, "aaa", "aaa"},
		{"some parts no colour", "(aaa)", blush.NoColour, "bb aaa bb", "bb aaa bb"},
		{"exact blue", "(aaa)", blush.FgBlue, "aaa", blush.Colourise("aaa", blush.FgBlue)},
		{"some parts blue", "(aaa)", blush.FgBlue, "bb aaa bb", "bb " + blush.Colourise("aaa", blush.FgBlue) + " bb"},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			l := blush.Rx{regexp.MustCompile(tc.search)}
			got, ok := l.Find(tc.input, tc.colour)
			if got != tc.want {
				t.Errorf("got = `%s`, want `%s`", got, tc.want)
			}
			if !ok {
				t.Error("ok = false, want true")
			}
		})
	}

	rx := blush.NewLocator("a{3}")
	want := "this " + blush.Colourise("aaa", blush.FgBlue) + "meeting"
	got, ok = rx.Find("this aaameeting", blush.FgBlue)
	if got != want {
		t.Errorf("got = `%s`, want `%s`", got, want)
	}
	if !ok {
		t.Error("ok = false, want true")
	}
}
