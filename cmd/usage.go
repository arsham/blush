package cmd

import "errors"

// These variables are used for showing help messages on command line.
var (
	errShowHelp = errors.New("show errors")

	Help = "Usage: blush [OPTION]... PATTERN [FILE]...\nTry 'blush --help' for more information."
	// nolint:misspell // it's ok.
	Usage = `Usage: blush [OPTION]... PATTERN [FILE]...
Colours:
    -r, --red       Match decorated with red colour. See Stock Colours section.
    -r[G], --red[G] Matches are grouped with the group number.
                    Example: blush -b1 match filename
    -#RGB, --#RGB   Use user defined colour schemas.
                    Example: blush -#1eF match filename
    -#RRGGBB, --#RRGGBB Same as -#RGB/--#RGB.

Pattern:
    You can use simple pattern or regexp. If your pattern expands between
    multiple words or has space in between, you should put them in quotations.

Stock Colours:
    -r, --red
    -g, --green
    -b, --blue
    -w, --white
    -bl, --black
    -yl, --yellow
    -mg, --magenta
    -cy, --cyan

Control arguments:
    -C, --colour            Don't drop unmatched lines.
    -i                      Case insensitive match.
    --no-colour, --no-color Don't colourise the output.
    -h, --no-filename       Suppress the prefixing of file names on output.

Multi match colouring:
    blush -b match1 [match2]...: will colourise all matches with the same colour.

Using pipes:
    cat FILE | blush -b match [-g match]...
`
)
