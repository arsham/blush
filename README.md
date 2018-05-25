# Blush

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![GoDoc](https://godoc.org/github.com/arsham/blush?status.svg)](http://godoc.org/github.com/arsham/blush)
[![Build Status](https://travis-ci.org/arsham/blush.svg?branch=master)](https://travis-ci.org/arsham/blush)
[![Coverage Status](https://codecov.io/gh/arsham/blush/branch/master/graph/badge.svg)](https://codecov.io/gh/arsham/blush)
[![Go Report Card](https://goreportcard.com/badge/github.com/arsham/blush)](https://goreportcard.com/report/github.com/arsham/blush)

With Blush, you can grep with colours, many colours!

## Usage

### Grep Method

This method greps the line that matches the input:
```bash
$ blush -b "first search" -g "second one" -g "and another one" files/paths
```

Any occurrence of `first search` will be in blue, `second one` and `and another one`
are in green.

### Colouring Method

With this method all texts are shown, but the matching words are coloured. You
can activate this mode by providing `--colour` or `-C` argument.

### Piping

Blush can also read from a pipe:
```bash
$ cat FILENAME | blush -b "print in blue" -g "in green" -g "another green"
$ cat FILENAME | blush "some text"
```

## Arguments

```
+-------------+----------+------------------------------+
|   Argument  | Shortcut |            Notes             |
+-------------+----------+------------------------------+
| --colour    | -C       | Colour, don't drop anything. |
| N/A         | -i       | Case insensitive matching    |
| N/A         | -R       | Recursive                    |
| --no-colour | N/A      | Doesn't colourize matches.   |
| --no-color  | N/A      | Same as --no-colour          |
+-------------+----------+------------------------------+
```

File names or paths are matched from the end. Any argument that doesn't match
any files or paths are considered as regular expression. If regular expressions
are not followed by colouring arguments are coloured based on previously
provided colour:

```bash
$ blush -b match1 match3 FILENAME
```

### Notes

* If no colour is provided, blush will choose blue.
* If you only provide file/path, it will print them out without colouring.
* If the matcher contains only alphabets and numbers, a non-regular expression is applied to search.

### Colour Groups

You can provide a number for a colour argument to create a colour group:

```bash
$ blush -b1 match1 -b2 match2 -b1 match3 FILENAME
```

Both `match1` and `match3` will be shown with the same `random blue` colour,
while `match2` will be another random blue colour. This means the numbers will
create colour groups.

You also can provide a colour with a series of grep requests:

```bash
$ blush -b match1 match3 -g match2 FILENAME
```

### Colours

You can choose a pre-defined colour, or pass it your own colour with a hash:

```
+-----------+----------+
|  Argument | Shortcut |
+-----------+----------+
| --red     | -r       |
| --green   | -g       |
| --blue    | -b       |
| --white   | -w       |
| --black   | -bl      |
| --yellow  | -yl      |
| --magenta | -mg      |
| --cyan    | -cy      |
| --#11bb22 | --#1b2   |
+-----------+----------+

```

### Complex Grep

You must put your complex grep into quotations:

```bash
$ blush -b "^age: [0-9]+" FILENAME
```

## Roadmap

* [ ] user defined colours.
* [ ] invert match (-v).
* [ ] implement all grep arguments.
* [ ] config files.
* [ ] internal pager and fuzzy search.
