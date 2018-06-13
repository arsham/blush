# Blush

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![GoDoc](https://godoc.org/github.com/arsham/blush?status.svg)](http://godoc.org/github.com/arsham/blush)
[![Build Status](https://travis-ci.org/arsham/blush.svg?branch=master)](https://travis-ci.org/arsham/blush)
[![Coverage Status](https://codecov.io/gh/arsham/blush/branch/master/graph/badge.svg)](https://codecov.io/gh/arsham/blush)
[![Go Report Card](https://goreportcard.com/badge/github.com/arsham/blush)](https://goreportcard.com/report/github.com/arsham/blush)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/4d4d4330fc2e44f18da6d8012d7432b9)](https://www.codacy.com/app/arsham/blush?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=arsham/blush&amp;utm_campaign=Badge_Grade)

With Blush, you can grep with any colours of your choice.

![Colored](http://i.imgur.com/RF19HYU.png)

1. [Install](#install)
    * [Update](#update)
2. [Usage](#usage)
    * [Match Method](#match-method)
    * [Colouring Method](#colouring-method)
    * [Piping](#piping)
3. [Arguments](#arguments)
    * [Notes](#notes)
4. [Colour Groups](#colour-groups)
5. [Colours](#colours)
6. [Complex Grep](#complex-grep)
7. [Suggestions](#suggestions)
8. [License](#license)

## Install

You can grab a binary from [releases](https://github.com/arsham/blush/releases)
page. If you prefer to install it manually you can get the code and install it
with the following command:

```bash
$ go get github.com/arsham/blush
```

Make sure you have `go>=1.7` installed.

### Update

In order to update the program:

```bash
$ cd $GOPATH/src/github.com/arsham/blush
$ make update
$ make install
```

## Usage

### Match Method

This method shows matches with the given input:
```bash
$ blush -b "first search" -g "second one" -g "and another one" files/paths
```

Any occurrence of `first search` will be in blue, `second one` and `and another one`
are in green.

![Colored](http://i.imgur.com/ghUTuva.png)

### Colouring Method

With this method all texts are shown, but the matching words are coloured. You
can activate this mode by providing `--colour` or `-C` argument.

![Colored](http://i.imgur.com/3CqzAUd.png)

### Piping

Blush can also read from a pipe:
```bash
$ cat FILENAME | blush -b "print in blue" -g "in green" -g "another green"
$ cat FILENAME | blush "some text"
```

## Arguments

```
+---------------+----------+------------------------------------------------+
|    Argument   | Shortcut |                     Notes                      |
+---------------+----------+------------------------------------------------+
| --colour      | -C       | Colour, don't drop anything.                   |
| N/A           | -i       | Case insensitive matching.                     |
| N/A           | -R       | Recursive matching.                            |
| --no-colour   | N/A      | Don't colourize matches.                       |
| --no-color    | N/A      | Same as --no-colour.                           |
| --no-filename | -h       | Suppress the prefixing of file names on output.|
+---------------+----------+------------------------------------------------+
```

File names or paths are matched from the end. Any argument that doesn't match
any files or paths are considered as regular expression. If regular expressions
are not followed by colouring arguments are coloured based on previously
provided colour:

```bash
$ blush -b match1 match2 FILENAME
```

![Colored](http://i.imgur.com/J6uZPQD.png)

### Notes

* If no colour is provided, blush will choose blue.
* If you only provide file/path, it will print them out without colouring.
* If the matcher contains only alphabets and numbers, a non-regular expression is applied to search.

## Colour Groups

You can provide a number for a colour argument to create a colour group:

```bash
$ blush -r1 match1 -r2 match2 -r1 match3 FILENAME
```

![Colored](http://i.imgur.com/cBnyrcy.png)

All matches will be shown as blue. But `match1` and `match3` will have a
different background colour than `match2`. This means the numbers will create
colour groups.

You also can provide a colour with a series of match requests:

```bash
$ blush -r match1 match3 -g match2 FILENAME
```

## Colours

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
+-----------+----------+
```

You can also pass an RGB colour. It can be in short form (--#1b2, -#1b2), or
long format (--#11bb22, -#11bb22).

![Colored](http://i.imgur.com/MkBIM9b.png)

## Complex Grep

You must put your complex grep into quotations:

```bash
$ blush -b "^age: [0-9]+" FILENAME
```
![Colored](http://i.imgur.com/hskdVhe.png)

## Suggestions

This tool is made to make your experience in terminal a more pleasant. Please
feel free to make any suggestions or request features by creating an issue.

Please see [changelog](./CHANGELOG.md) document for newest changes.

## License

Use of this source code is governed by the MIT License. License file can be
found in the [LICENSE](./LICENSE) file.
