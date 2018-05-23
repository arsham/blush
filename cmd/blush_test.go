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

type stdFile struct {
	f *os.File
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

func (s *stdFile) String() string {
	s.f.Seek(0, 0)
	buf := new(bytes.Buffer)
	buf.ReadFrom(s.f)
	return buf.String()
}

func setup(t *testing.T, args string) (stdout, stderr *stdFile, cleanup func()) {
	oldArgs := os.Args
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	oldFatalErr := cmd.FatalErr

	stdout, outCleanup := newStdFile(t, "stdout")
	os.Stdout = stdout.f

	stderr, errCleanup := newStdFile(t, "stderr")
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
		t.Errorf("didn't expect printing anything: %s", stdout.String())
	}
	if !strings.Contains(stderr.String(), cmd.ErrNoInput.Error()) {
		t.Errorf("stderr = `%s`, want `%s` in it", stderr.String(), cmd.ErrNoInput.Error())
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

func TestMainMatchRegexp(t *testing.T) {
	match := blush.Colourise("TOKEN", blush.FgBlue)
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	location := path.Join(pwd, "../blush/testdata")

	stdout, stderr, cleanup := setup(t, `-b TOKEN -r .* `+location)
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

func TestFindFiles(t *testing.T) {
	tcs := []struct {
		name  string
		input []string
		count int
	}{
		{"path", []string{"/"}, 1},
		{"path duplicate", []string{"/", "/"}, 1},
		{"path trailing spaces", []string{"/", "         "}, 1},
		{"path with prefix", []string{"something else", "/"}, 1},
		{"paths", []string{"/", "/dev"}, 2},
		{"file", []string{"/dev/null"}, 1},
		{"file duplicate", []string{"/dev/null", "/dev/null"}, 1},
		{"file trailing spaces", []string{"/dev/null", "     "}, 1},
		{"file with prefix", []string{"something else", "/dev/null"}, 1},
		{"files", []string{"/dev/null", "/dev/zero"}, 2},
		{"file and path", []string{"/dev/null", "/dev"}, 2},
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
				return
			}
			if len(b.Paths) != tc.count {
				t.Errorf("len(b.Paths) = %d, want %d", len(b.Paths), tc.count)
			}
		})
	}
}

func TestOtherArgs(t *testing.T) {
	tcs := []struct {
		input []string
		f     func(blush.Blush) bool
	}{
		{[]string{"/"},
			func(b blush.Blush) bool { return !b.Recursive && !b.Colouring && !b.Insensitive },
		},
		{[]string{"-i", "/"},
			func(b blush.Blush) bool { return !b.Recursive && !b.Colouring && b.Insensitive },
		},
		{[]string{"-i", "-i", "/"},
			func(b blush.Blush) bool { return !b.Recursive && !b.Colouring && b.Insensitive },
		},
		{[]string{"-i", "-R", "/"},
			func(b blush.Blush) bool { return b.Recursive && !b.Colouring && b.Insensitive },
		},
		{[]string{"-R", "-i", "/"},
			func(b blush.Blush) bool { return b.Recursive && !b.Colouring && b.Insensitive },
		},
		{[]string{"aaa", "-i", "/"},
			func(b blush.Blush) bool { return !b.Recursive && !b.Colouring && b.Insensitive },
		},
		{[]string{"-i", "aaa", "-i", "/"},
			func(b blush.Blush) bool { return !b.Recursive && !b.Colouring && b.Insensitive },
		},
		{[]string{"aaa", "-i", "-i", "/"},
			func(b blush.Blush) bool { return !b.Recursive && !b.Colouring && b.Insensitive },
		},
		{[]string{"aaa", "-C", "-i", "/"},
			func(b blush.Blush) bool { return !b.Recursive && b.Colouring && b.Insensitive },
		},
		{[]string{"aaa", "-i", "-C", "/"},
			func(b blush.Blush) bool { return !b.Recursive && b.Colouring && b.Insensitive },
		},
		{[]string{"aaa", "-i", "-C", "-R", "/"},
			func(b blush.Blush) bool { return b.Recursive && b.Colouring && b.Insensitive },
		},
	}
	for _, tc := range tcs {
		input := append([]string{"blush"}, tc.input...)
		b, err := cmd.GetBlush(input)
		if err != nil {
			t.Errorf("GetBlush(%s): err = %s, want nil", tc.input, err)
		}
		if b == nil {
			t.Errorf("GetBlush(%s): b = nil, want *Blush", tc.input)
			continue
		}
		if !tc.f(*b) {
			t.Errorf("failed on: %v", tc.input)
		}
	}
}

func TestColourArgs(t *testing.T) {
	aaa := blush.Exact("aaa")
	bbb := blush.Exact("bbb")
	tcs := []struct {
		name  string
		input []string
		want  []blush.Arg
	}{
		{"empty", []string{"/"}, []blush.Arg{}},
		{"1-no colour", []string{"aaa", "/"}, []blush.Arg{
			blush.Arg{Colour: blush.DefaultColour, Find: aaa},
		}},
		{"1-colour", []string{"-b", "aaa", "/"}, []blush.Arg{
			blush.Arg{Colour: blush.FgBlue, Find: aaa},
		}},
		{"1-colour long", []string{"--blue", "aaa", "/"}, []blush.Arg{
			blush.Arg{Colour: blush.FgBlue, Find: aaa},
		}},
		{"2-no colour", []string{"aaa", "bbb", "/"}, []blush.Arg{
			blush.Arg{Colour: blush.DefaultColour, Find: aaa},
			blush.Arg{Colour: blush.DefaultColour, Find: bbb},
		}},
		{"2-colour", []string{"-b", "aaa", "bbb", "/"}, []blush.Arg{
			blush.Arg{Colour: blush.FgBlue, Find: aaa},
			blush.Arg{Colour: blush.FgBlue, Find: bbb},
		}},
		{"2-two colours", []string{"-b", "aaa", "-g", "bbb", "/"}, []blush.Arg{
			blush.Arg{Colour: blush.FgBlue, Find: aaa},
			blush.Arg{Colour: blush.FgGreen, Find: bbb},
		}},
		{"red", []string{"-r", "aaa", "--red", "bbb", "/"}, []blush.Arg{
			blush.Arg{Colour: blush.FgRed, Find: aaa},
			blush.Arg{Colour: blush.FgRed, Find: bbb},
		}},
		{"green", []string{"-g", "aaa", "--green", "bbb", "/"}, []blush.Arg{
			blush.Arg{Colour: blush.FgGreen, Find: aaa},
			blush.Arg{Colour: blush.FgGreen, Find: bbb},
		}},
		{"blue", []string{"-b", "aaa", "--blue", "bbb", "/"}, []blush.Arg{
			blush.Arg{Colour: blush.FgBlue, Find: aaa},
			blush.Arg{Colour: blush.FgBlue, Find: bbb},
		}},
		{"white", []string{"-w", "aaa", "--white", "bbb", "/"}, []blush.Arg{
			blush.Arg{Colour: blush.FgWhite, Find: aaa},
			blush.Arg{Colour: blush.FgWhite, Find: bbb},
		}},
		{"black", []string{"-bl", "aaa", "--black", "bbb", "/"}, []blush.Arg{
			blush.Arg{Colour: blush.FgBlack, Find: aaa},
			blush.Arg{Colour: blush.FgBlack, Find: bbb},
		}},
		{"cyan", []string{"-cy", "aaa", "--cyan", "bbb", "/"}, []blush.Arg{
			blush.Arg{Colour: blush.FgCyan, Find: aaa},
			blush.Arg{Colour: blush.FgCyan, Find: bbb},
		}},
		{"magenta", []string{"-mg", "aaa", "--magenta", "bbb", "/"}, []blush.Arg{
			blush.Arg{Colour: blush.FgMagenta, Find: aaa},
			blush.Arg{Colour: blush.FgMagenta, Find: bbb},
		}},
		{"yellow", []string{"-yl", "aaa", "--yellow", "bbb", "/"}, []blush.Arg{
			blush.Arg{Colour: blush.FgYellow, Find: aaa},
			blush.Arg{Colour: blush.FgYellow, Find: bbb},
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
			if !argsEqual(b.Args, tc.want) {
				t.Errorf("(%s): b.Args = %v, want %v", tc.input, b.Args, tc.want)
			}
		})
	}
}

func argsEqual(a, b []blush.Arg) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	isIn := func(a blush.Arg, haystack []blush.Arg) bool {
		for _, b := range haystack {
			af, bf := a.Find.(blush.Exact), b.Find.(blush.Exact)
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
