package cmd_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/arsham/blush/blush"
	"github.com/arsham/blush/cmd"
)

var leaveMeHere = "LEAVEMEHERE"

type stdFile struct {
	f *os.File
}

func (s *stdFile) String() string {
	s.f.Seek(0, 0)
	buf := new(bytes.Buffer)
	buf.ReadFrom(s.f)
	return buf.String()
}

func newStdFile(t *testing.T, name string) (*stdFile, func()) {
	f, err := ioutil.TempFile("", name)
	if err != nil {
		t.Fatal(err)
	}
	sf := &stdFile{f}
	return sf, func() {
		f.Close()
		os.Remove(f.Name())
	}
}

func setup(t *testing.T, args string) (stdout, stderr *stdFile, cleanup func()) {
	oldArgs := os.Args
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	oldFatalErr := cmd.FatalErr

	stdout, outCleanup := newStdFile(t, "stdout")
	stderr, errCleanup := newStdFile(t, "stderr")
	os.Stdout = stdout.f
	os.Stderr = stderr.f

	os.Args = []string{"blush"}
	if len(args) > 1 {
		os.Args = append(os.Args, strings.Split(args, " ")...)
	}
	cmd.FatalErr = func(s string) {
		fmt.Fprintf(os.Stderr, "%s\n", s)
	}

	cleanup = func() {
		outCleanup()
		errCleanup()
		os.Args = oldArgs
		os.Stdout = oldStdout
		os.Stderr = oldStderr
		cmd.FatalErr = oldFatalErr
	}
	return stdout, stderr, cleanup
}

func TestMainNoArgs(t *testing.T) {
	stdout, stderr, cleanup := setup(t, "")
	defer cleanup()
	cmd.Main()
	if len(stdout.String()) > 0 {
		t.Errorf("didn't expect any stdout, got: %s", stdout.String())
	}
	if !strings.Contains(stderr.String(), cmd.ErrNoInput.Error()) {
		t.Errorf("stderr = `%s`, want `%s` in it", stderr.String(), cmd.ErrNoInput.Error())
	}
}

func TestPipeInput(t *testing.T) {
	oldStdin := os.Stdin
	stdout, stderr, cleanup := setup(t, "findme")
	defer func() {
		cleanup()
		os.Stdin = oldStdin
	}()
	file, err := ioutil.TempFile("", "blush_pipe")
	if err != nil {
		t.Fatal(err)
	}
	name := file.Name()
	rmFile := func() {
		if err := os.Remove(name); err != nil {
			t.Error(err)
		}
	}
	defer rmFile()
	file.Close()
	rmFile()
	file, err = os.OpenFile(name, os.O_CREATE|os.O_RDWR, os.ModeCharDevice|os.ModeDevice)
	if err != nil {
		t.Fatal(err)
	}
	file.WriteString("you can findme here")
	os.Stdin = file
	file.Seek(0, 0)
	cmd.Main()
	if len(stderr.String()) > 0 {
		t.Errorf("didn't expect printing anything: %s", stderr.String())
	}
	if !strings.Contains(stdout.String(), "findme") {
		t.Errorf("stdout = `%s`, want `%s` in it", stdout.String(), "findme")
	}
}

func TestMainMatchExact(t *testing.T) {
	match := blush.Colourise("TOKEN", blush.FgBlue)
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	location := path.Join(pwd, "../blush/testdata")

	stdout, stderr, cleanup := setup(t, "-b TOKEN "+location)
	defer cleanup()
	cmd.Main()

	if len(stderr.String()) > 0 {
		t.Errorf("len(stderr.String()) = %d, want 0: `%s`", len(stderr.String()), stderr.String())
	}
	if len(stdout.String()) == 0 {
		t.Errorf("len(stdout.String()) = %d, want > 0", len(stdout.String()))
	}
	if !strings.Contains(stdout.String(), match) {
		t.Errorf("want `%s` in `%s`", match, stdout.String())
	}
}

func TestMainMatchIExact(t *testing.T) {
	match := blush.Colourise("TOKEN", blush.FgBlue)
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	location := path.Join(pwd, "../blush/testdata")

	stdout, stderr, cleanup := setup(t, "-i -b token "+location)
	defer cleanup()
	cmd.Main()

	if len(stderr.String()) > 0 {
		t.Errorf("len(stderr.String()) = %d, want 0: `%s`", len(stderr.String()), stderr.String())
	}
	if len(stdout.String()) == 0 {
		t.Errorf("len(stdout.String()) = %d, want > 0", len(stdout.String()))
	}
	if !strings.Contains(stdout.String(), match) {
		t.Errorf("want `%s` in `%s`", match, stdout.String())
	}
}

func TestMainMatchRegexp(t *testing.T) {
	match := blush.Colourise("TOKEN", blush.FgBlue)
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	location := path.Join(pwd, "../blush/testdata")

	stdout, stderr, cleanup := setup(t, `-b TOK[EN]{2} `+location)
	defer cleanup()
	cmd.Main()

	if len(stderr.String()) > 0 {
		t.Errorf("len(stderr.String()) = %d, want 0: `%s`", len(stderr.String()), stderr.String())
	}
	if len(stdout.String()) == 0 {
		t.Errorf("len(stdout.String()) = %d, want > 0", len(stdout.String()))
	}
	if !strings.Contains(stdout.String(), match) {
		t.Errorf("want `%s` in `%s`", match, stdout.String())
	}
}

func TestMainMatchRegexpInsensitive(t *testing.T) {
	match := blush.Colourise("TOKEN", blush.FgBlue)
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	location := path.Join(pwd, "../blush/testdata")

	stdout, stderr, cleanup := setup(t, `-i -b tok[en]{2} `+location)
	defer cleanup()
	cmd.Main()

	if len(stderr.String()) > 0 {
		t.Errorf("len(stderr.String()) = %d, want 0: `%s`", len(stderr.String()), stderr.String())
	}
	if len(stdout.String()) == 0 {
		t.Errorf("len(stdout.String()) = %d, want > 0", len(stdout.String()))
	}
	if !strings.Contains(stdout.String(), match) {
		t.Errorf("want `%s` in `%s`", match, stdout.String())
	}
}

func TestMainMatchNoCut(t *testing.T) {
	matches := []string{"TOKEN", "ONE", "TWO", "THREE", "FOUR"}
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	location := path.Join(pwd, "../blush/testdata")

	stdout, stderr, cleanup := setup(t, fmt.Sprintf("-C -b %s %s", leaveMeHere, location))
	defer cleanup()
	cmd.Main()

	if len(stderr.String()) > 0 {
		t.Errorf("len(stderr.String()) = %d, want 0: `%s`", len(stderr.String()), stderr.String())
	}
	if len(stdout.String()) == 0 {
		t.Errorf("len(stdout.String()) = %d, want > 0", len(stdout.String()))
	}
	for _, s := range matches {
		if !strings.Contains(stdout.String(), s) {
			t.Errorf("want `%s` in `%s`", s, stdout.String())
		}
	}
}

func TestNoFiles(t *testing.T) {
	fileName := "test"
	b, err := cmd.GetBlush([]string{fileName})
	if err == nil {
		t.Error("GetBlush(): err = nil, want error")
	}
	if b != nil {
		t.Errorf("GetBlush(): b = %v, want nil", b)
	}
}

func TestColourArgs(t *testing.T) {
	aaa := blush.Exact("aaa")
	bbb := blush.Exact("bbb")
	tcs := []struct {
		name  string
		input []string
		want  []blush.ColourLocator
	}{
		{"empty", []string{"/"}, []blush.ColourLocator{}},
		{"1-no colour", []string{"aaa", "/"}, []blush.ColourLocator{
			blush.ColourLocator{Colour: blush.DefaultColour, Locator: aaa},
		}},
		{"1-colour", []string{"-b", "aaa", "/"}, []blush.ColourLocator{
			blush.ColourLocator{Colour: blush.FgBlue, Locator: aaa},
		}},
		{"1-colour long", []string{"--blue", "aaa", "/"}, []blush.ColourLocator{
			blush.ColourLocator{Colour: blush.FgBlue, Locator: aaa},
		}},
		{"2-no colour", []string{"aaa", "bbb", "/"}, []blush.ColourLocator{
			blush.ColourLocator{Colour: blush.DefaultColour, Locator: aaa},
			blush.ColourLocator{Colour: blush.DefaultColour, Locator: bbb},
		}},
		{"2-colour", []string{"-b", "aaa", "bbb", "/"}, []blush.ColourLocator{
			blush.ColourLocator{Colour: blush.FgBlue, Locator: aaa},
			blush.ColourLocator{Colour: blush.FgBlue, Locator: bbb},
		}},
		{"2-two colours", []string{"-b", "aaa", "-g", "bbb", "/"}, []blush.ColourLocator{
			blush.ColourLocator{Colour: blush.FgBlue, Locator: aaa},
			blush.ColourLocator{Colour: blush.FgGreen, Locator: bbb},
		}},
		{"red", []string{"-r", "aaa", "--red", "bbb", "/"}, []blush.ColourLocator{
			blush.ColourLocator{Colour: blush.FgRed, Locator: aaa},
			blush.ColourLocator{Colour: blush.FgRed, Locator: bbb},
		}},
		{"green", []string{"-g", "aaa", "--green", "bbb", "/"}, []blush.ColourLocator{
			blush.ColourLocator{Colour: blush.FgGreen, Locator: aaa},
			blush.ColourLocator{Colour: blush.FgGreen, Locator: bbb},
		}},
		{"blue", []string{"-b", "aaa", "--blue", "bbb", "/"}, []blush.ColourLocator{
			blush.ColourLocator{Colour: blush.FgBlue, Locator: aaa},
			blush.ColourLocator{Colour: blush.FgBlue, Locator: bbb},
		}},
		{"white", []string{"-w", "aaa", "--white", "bbb", "/"}, []blush.ColourLocator{
			blush.ColourLocator{Colour: blush.FgWhite, Locator: aaa},
			blush.ColourLocator{Colour: blush.FgWhite, Locator: bbb},
		}},
		{"black", []string{"-bl", "aaa", "--black", "bbb", "/"}, []blush.ColourLocator{
			blush.ColourLocator{Colour: blush.FgBlack, Locator: aaa},
			blush.ColourLocator{Colour: blush.FgBlack, Locator: bbb},
		}},
		{"cyan", []string{"-cy", "aaa", "--cyan", "bbb", "/"}, []blush.ColourLocator{
			blush.ColourLocator{Colour: blush.FgCyan, Locator: aaa},
			blush.ColourLocator{Colour: blush.FgCyan, Locator: bbb},
		}},
		{"magenta", []string{"-mg", "aaa", "--magenta", "bbb", "/"}, []blush.ColourLocator{
			blush.ColourLocator{Colour: blush.FgMagenta, Locator: aaa},
			blush.ColourLocator{Colour: blush.FgMagenta, Locator: bbb},
		}},
		{"yellow", []string{"-yl", "aaa", "--yellow", "bbb", "/"}, []blush.ColourLocator{
			blush.ColourLocator{Colour: blush.FgYellow, Locator: aaa},
			blush.ColourLocator{Colour: blush.FgYellow, Locator: bbb},
		}},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			input := append([]string{"blush"}, tc.input...)
			b, err := cmd.GetBlush(input)
			if err != nil {
				t.Errorf("GetBlush(): err = %s, want nil", err)
			}
			if b == nil {
				t.Error("GetBlush(): b = nil, want *Blush")
			}
			if !argsEqual(b.Locator, tc.want) {
				t.Errorf("(%s): b.Args = %v, want %v", tc.input, b.Locator, tc.want)
			}
		})
	}
}

func argsEqual(a, b []blush.ColourLocator) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	isIn := func(a blush.ColourLocator, haystack []blush.ColourLocator) bool {
		for _, b := range haystack {
			af, bf := a.Locator.(blush.Exact), b.Locator.(blush.Exact)
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
