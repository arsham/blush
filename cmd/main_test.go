package cmd_test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/arsham/blush/blush"
	"github.com/arsham/blush/cmd"
)

func TestMainNoArgs(t *testing.T) {
	stdout, stderr := setup(t, "")
	cmd.Main()
	assert.Empty(t, stdout.String())
	assert.Contains(t, stderr.String(), cmd.ErrNoInput.Error())
	assert.Contains(t, stderr.String(), cmd.Help)
}

func TestMainHelp(t *testing.T) {
	stdout, stderr := setup(t, "--help")
	cmd.Main()
	assert.Empty(t, stderr.String())
	assert.Contains(t, stdout.String(), cmd.Usage)
}

func TestPipeInput(t *testing.T) {
	oldStdin := os.Stdin
	stdout, stderr := setup(t, "findme")
	defer func() {
		os.Stdin = oldStdin
	}()
	file := getPipe(t)
	file.WriteString("you can findme here")
	os.Stdin = file
	file.Seek(0, 0)
	cmd.Main()
	assert.Empty(t, stderr.String())
	assert.Contains(t, stdout.String(), "findme")
}

func TestMainMatch(t *testing.T) {
	match := blush.Colourise("TOKEN", blush.Blue)
	pwd, err := os.Getwd()
	assert.NoError(t, err)
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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr := setup(t, fmt.Sprintf("%s %s", tc.input, location))
			cmd.Main()

			assert.Empty(t, stderr.String())
			assert.NotEmpty(t, stdout.String())
			assert.Contains(t, stdout.String(), match)
		})
	}
}

func TestMainMatchCut(t *testing.T) {
	matches := []string{"TOKEN", "ONE", "TWO", "THREE", "FOUR"}
	pwd, err := os.Getwd()
	assert.NoError(t, err)
	location := path.Join(pwd, "../blush/testdata")

	stdout, stderr := setup(t, fmt.Sprintf("-b %s %s", leaveMeHere, location))
	cmd.Main()
	assert.Empty(t, stderr.String())
	assert.NotEmpty(t, stdout.String())
	for _, s := range matches {
		assert.Contains(t, stdout.String(), s)
	}
}

func TestNoFiles(t *testing.T) {
	fileName := "test"
	b, err := cmd.GetBlush([]string{fileName})
	assert.Error(t, err)
	assert.Nil(t, b)
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
		{"1-no colour american", []string{"--no-colour", "aaa", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.NoColour),
		}},
		{"1-colour", []string{"-b", "aaa", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.Blue),
		}},
		{"1-colour long", []string{"--blue", "aaa", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.Blue),
		}},
		{"2-default colour", []string{"aaa", "bbb", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.DefaultColour),
			blush.NewExact(bbb, blush.DefaultColour),
		}},
		{"2-no colour", []string{"--no-colour", "aaa", "bbb", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.NoColour),
			blush.NewExact(bbb, blush.NoColour),
		}},
		{"2-no colour american", []string{"--no-colour", "aaa", "bbb", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.NoColour),
			blush.NewExact(bbb, blush.NoColour),
		}},
		{"2-colour", []string{"-b", "aaa", "bbb", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.Blue),
			blush.NewExact(bbb, blush.Blue),
		}},
		{"2-two colours", []string{"-b", "aaa", "-g", "bbb", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.Blue),
			blush.NewExact(bbb, blush.Green),
		}},
		{"red", []string{"-r", "aaa", "--red", "bbb", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.Red),
			blush.NewExact(bbb, blush.Red),
		}},
		{"green", []string{"-g", "aaa", "--green", "bbb", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.Green),
			blush.NewExact(bbb, blush.Green),
		}},
		{"blue", []string{"-b", "aaa", "--blue", "bbb", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.Blue),
			blush.NewExact(bbb, blush.Blue),
		}},
		{"white", []string{"-w", "aaa", "--white", "bbb", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.White),
			blush.NewExact(bbb, blush.White),
		}},
		{"black", []string{"-bl", "aaa", "--black", "bbb", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.Black),
			blush.NewExact(bbb, blush.Black),
		}},
		{"cyan", []string{"-cy", "aaa", "--cyan", "bbb", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.Cyan),
			blush.NewExact(bbb, blush.Cyan),
		}},
		{"magenta", []string{"-mg", "aaa", "--magenta", "bbb", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.Magenta),
			blush.NewExact(bbb, blush.Magenta),
		}},
		{"yellow", []string{"-yl", "aaa", "--yellow", "bbb", "/"}, []blush.Finder{
			blush.NewExact(aaa, blush.Yellow),
			blush.NewExact(bbb, blush.Yellow),
		}},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			input := append([]string{"blush"}, tc.input...)
			b, err := cmd.GetBlush(input)
			assert.NoError(t, err)
			assert.NotNil(t, b)
			assert.True(t, argsEqual(b.Finders, tc.want))
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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			b, err := cmd.GetBlush(tc.input)
			assert.NoError(t, err)
			assert.NotNil(t, b)
			assert.Equal(t, tc.want, b.WithFileName)
		})
	}
}
