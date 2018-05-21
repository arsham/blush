# Blush
Grep with colours.

## Usage

### Grep Method

This method greps the line that matches the input:
```bash
$ blush -b "print in blue" -g "in green" -g "another green" files/paths
```

### Match All Method

With this method all texts are shown, but the matching words are coloured. You
can activate this mode by providing `--color` argument.

### Piping

You can pipe your input as well:
```bash
$ cat FILENAME | blush -b "print in blue" -g "in green" -g "another green"
$ cat FILENAME | blush "some text"
```

## Arguments

```
+----------+----------+-------------------------------+
| Argument | Shortcut |             Notes             |
+----------+----------+-------------------------------+
| --colour | -c       | Don't drop non-matched lines. |
| --rand   | N/A      | Chooses a random colour.      |
| N/A      | -i       | Case insensitive matching     |
| N/A      | -R       | Recursive                     |
+----------+----------+-------------------------------+
```

File names or paths are matched from the end. Any argument that doesn't match
any files or paths are considered as regular expression. If regular expressions
are not followed by colouring arguments are coloured based on previously
provided colour:

```bash
$ blush -b match1 match3 FILENAME
```

### Notes

* If no colour is provided, blush will choose a different colour for each regexp.
* If you only provide file/path, it will print them out without colouring.
* If the matcher contains only alphabets and numbers, a non-regular expression is applied to search.

### Colour Groups

You can provide a number as the colour arguments to create a colour group. For
example, with:

```bash
$ blush -1 match1 -2 match2 -1 match3 FILENAME
```

Both `match1` and `match3` will be shown with the same `random` colour, while
`match2` will be another random colour. This means the numbers will create
colour groups.

You also can provide a colour with a series of grep requests:

```bash
$ blush -b match1 match3 -g match2 FILENAME
```

### Colours

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

### Complex Grep

You must put your complex grep into quotations:

```bash
$ blush -b "^age: [0-9]+" FILENAME
```
