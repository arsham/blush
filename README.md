# Blush

[![PkgGoDev](https://pkg.go.dev/badge/github.com/arsham/dbtools)](https://pkg.go.dev/github.com/arsham/dbtools)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/arsham/dbtools)
[![Build Status](https://github.com/arsham/dbtools/actions/workflows/go.yml/badge.svg)](https://github.com/arsham/dbtools/actions/workflows/go.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Coverage Status](https://codecov.io/gh/arsham/blush/branch/master/graph/badge.svg)](https://codecov.io/gh/arsham/blush)
[![Go Report Card](https://goreportcard.com/badge/github.com/arsham/blush)](https://goreportcard.com/report/github.com/arsham/blush)

With Blush, you can highlight matches with any colours of your choice.

![1](https://user-images.githubusercontent.com/428611/164768864-e9713ac3-0097-4435-8bcb-577dbf7b9931.png)

1. [Install](#install)
2. [Usage](#usage)
   - [Note](#note)
   - [Normal Mode](#normal-mode)
   - [Dropping Unmatched](#dropping-unmatched)
   - [Piping](#piping)
3. [Arguments](#arguments)
   - [Notes](#notes)
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
$ go install github.com/arsham/blush@latest
```

Make sure you have `go>=1.18` installed.

## Usage

Blush can read from a file or a pipe:

```bash
$ cat FILENAME | blush -b "print in blue" -g "in green" -g "another green"
$ cat FILENAME | blush "some text"
$ blush -b "print in blue" -g "in green" -g "another green" FILENAME
$ blush "some text" FILENAME
```

### Note

Although this program has a good performance, but performance is not the main
concern. There are other tools you should use if you are searching in large
files. Two examples:

- [Ripgrep](https://github.com/BurntSushi/ripgrep)
- [The Silver Searcher](https://github.com/ggreer/the_silver_searcher)

### Normal Mode

This method shows matches with the given input:

```bash
$ blush -b "first search" -g "second one" -g "and another one" files/paths
```

Any occurrence of `first search` will be in blue, `second one` and `and another one`
are in green.

![2](https://user-images.githubusercontent.com/428611/164768874-bf687313-c103-449b-bb57-6fdcea51fc5d.png)

### Dropping Unmatched

By default, unmatched lines are not dropped. But you can use the `-d` flag to
drop them:

![3](https://user-images.githubusercontent.com/428611/164768875-c9aa3e47-7db0-454f-8a55-1e2bff332c69.png)

## Arguments

| Argument      | Shortcut | Notes                                           |
| :------------ | :------- | :---------------------------------------------- |
| N/A           | -i       | Case insensitive matching.                      |
| N/A           | -R       | Recursive matching.                             |
| --no-filename | -h       | Suppress the prefixing of file names on output. |
| --drop        | -d       | Drop unmatched lines                            |

File names or paths are matched from the end. Any argument that doesn't match
any files or paths are considered as regular expression. If regular expressions
are not followed by colouring arguments are coloured based on previously
provided colour:

```bash
$ blush -b match1 match2 FILENAME
```

![4](https://user-images.githubusercontent.com/428611/164768879-f9b73b2c-b6bb-4cf5-a98a-e51535fa554a.png)

### Notes

- If no colour is provided, blush will choose blue.
- If you only provide file/path, it will print them out without colouring.
- If the matcher contains only alphabets and numbers, a non-regular expression is applied to search.

## Colour Groups

You can provide a number for a colour argument to create a colour group:

```bash
$ blush -r1 match1 -r2 match2 -r1 match3 FILENAME
```

![5](https://user-images.githubusercontent.com/428611/164768882-5ce57477-e9d5-4170-ac10-731e9391cbee.png)

All matches will be shown as blue. But `match1` and `match3` will have a
different background colour than `match2`. This means the numbers will create
colour groups.

You also can provide a colour with a series of match requests:

```bash
$ blush -r match1 match3 -g match2 FILENAME
```

## Colours

You can choose a pre-defined colour, or pass it your own colour with a hash:

| Argument  | Shortcut |
| :-------- | :------- |
| --red     | -r       |
| --green   | -g       |
| --blue    | -b       |
| --white   | -w       |
| --black   | -bl      |
| --yellow  | -yl      |
| --magenta | -mg      |
| --cyan    | -cy      |

You can also pass an RGB colour. It can be in short form (--#1b2, -#1b2), or
long format (--#11bb22, -#11bb22).

![6](https://user-images.githubusercontent.com/428611/164768883-154b4fd9-946f-43eb-b3f5-ede6027c3eda.png)

## Complex Grep

You must put your complex grep into quotations:

```bash
$ blush -b "^age: [0-9]+" FILENAME
```

![7](https://user-images.githubusercontent.com/428611/164768886-5b94b8fa-77e2-4617-80f2-040edce18660.png)

## Suggestions

This tool is made to make your experience in terminal a more pleasant. Please
feel free to make any suggestions or request features by creating an issue.

## License

Use of this source code is governed by the MIT License. License file can be
found in the [LICENSE](./LICENSE) file.
