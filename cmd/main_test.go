package cmd_test

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/arsham/blush/blush"
	"github.com/arsham/blush/cmd"
)

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
	file, cl := getPipe(t)
	defer cl()
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

func TestMainMatch(t *testing.T) {
	match := blush.Colourise("TOKEN", blush.FgBlue)
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	location := path.Join(pwd, "../blush/testdata")

	tcs := []struct {
		name  string
		input string
	}{
		{"exact sensitive", "-b TOKEN"},
		{"exact insensitive", "-i -b TOKEN"},
		{"regexp sensitive", "-b TOK[EN]{2}"},
		{"regexp insensitive", "-i -b tok[en]{2}"},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, cleanup := setup(t, fmt.Sprintf("%s %s", tc.input, location))
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
		})
	}
}

func TestMainMatchNoCut(t *testing.T) {
	matches := []string{"TOKEN", "ONE", "TWO", "THREE", "FOUR"}
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	location := path.Join(pwd, "../blush/testdata")

	tcs := []struct {
		name, input string
	}{
		{"short", "-C"},
		{"long", "--colour"},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, cleanup := setup(t, fmt.Sprintf("%s -b %s %s", tc.input, leaveMeHere, location))
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
		})
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
	aaa := "aaa"
	bbb := "bbb"
	tcs := []struct {
		name  string
		input []string
		want  []blush.Finder
	}{
		{"empty", []string{"/"}, []blush.Finder{}},
		{"1-default colour", []string{"aaa", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.DefaultColour),
		}},
		{"1-no colour", []string{"--no-colour", "aaa", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.NoColour),
		}},
		{"1-no colour american", []string{"--no-color", "aaa", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.NoColour),
		}},
		{"1-colour", []string{"-b", "aaa", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.FgBlue),
		}},
		{"1-colour long", []string{"--blue", "aaa", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.FgBlue),
		}},
		{"2-default colour", []string{"aaa", "bbb", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.DefaultColour),
			blush.NewExact(bbb, blush.DefaultColour),
		}},
		{"2-no colour", []string{"--no-colour", "aaa", "bbb", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.NoColour),
			blush.NewExact(bbb, blush.NoColour),
		}},
		{"2-no colour american", []string{"--no-color", "aaa", "bbb", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.NoColour),
			blush.NewExact(bbb, blush.NoColour),
		}},
		{"2-colour", []string{"-b", "aaa", "bbb", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.FgBlue),
			blush.NewExact(bbb, blush.FgBlue),
		}},
		{"2-two colours", []string{"-b", "aaa", "-g", "bbb", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.FgBlue),
			blush.NewExact(bbb, blush.FgGreen),
		}},
		{"red", []string{"-r", "aaa", "--red", "bbb", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.FgRed),
			blush.NewExact(bbb, blush.FgRed),
		}},
		{"green", []string{"-g", "aaa", "--green", "bbb", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.FgGreen),
			blush.NewExact(bbb, blush.FgGreen),
		}},
		{"blue", []string{"-b", "aaa", "--blue", "bbb", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.FgBlue),
			blush.NewExact(bbb, blush.FgBlue),
		}},
		{"white", []string{"-w", "aaa", "--white", "bbb", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.FgWhite),
			blush.NewExact(bbb, blush.FgWhite),
		}},
		{"black", []string{"-bl", "aaa", "--black", "bbb", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.FgBlack),
			blush.NewExact(bbb, blush.FgBlack),
		}},
		{"cyan", []string{"-cy", "aaa", "--cyan", "bbb", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.FgCyan),
			blush.NewExact(bbb, blush.FgCyan),
		}},
		{"magenta", []string{"-mg", "aaa", "--magenta", "bbb", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.FgMagenta),
			blush.NewExact(bbb, blush.FgMagenta),
		}},
		{"yellow", []string{"-yl", "aaa", "--yellow", "bbb", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.FgYellow),
			blush.NewExact(bbb, blush.FgYellow),
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
			if !argsEqual(b.Finders, tc.want) {
				t.Errorf("(%s): b.Args = %v, want %v", tc.input, b.Finders, tc.want)
			}
		})
	}
}

func TestWithFilename(t *testing.T) {
	tcs := []struct {
		name  string
		input []string
		want  bool
	}{
		{"with filename", []string{"blush", "/"}, true},
		{"no filename", []string{"blush", "-h", "aaa", "/"}, false},
		{"no filename long", []string{"blush", "--no-filename", "aaa", "/"}, false},
		{"no filename both", []string{"blush", "-h", "--no-filename", "aaa", "/"}, false},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			b, err := cmd.GetBlush(tc.input)
			if err != nil {
				t.Errorf("GetBlush(): err = %s, want nil", err)
			}
			if b == nil {
				t.Fatal("GetBlush(): b = nil, want *Blush")
			}
			if b.WithFileName != tc.want {
				t.Errorf("b.WithFileName = %t, want %t", b.WithFileName, tc.want)
			}
		})
	}
}
