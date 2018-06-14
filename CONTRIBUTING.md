# Contributing

1. [Dependencies](#dependencies)
2. [Testing](#testing)
3. [Benchmarking](#benchmarking)
4. [Pull Requests](#pull-requests)

## Dependencies

Dependency management is done with [glide](https://github.com/Masterminds/glide).
You can install it by running:
```bash
$ go get -u github.com/Masterminds/glide
```

Then make sure you run `glide install` to install the current version of
dependencies.

If you need to add a new dependency to the library, before you commit your
changes make sure you run:
```bash
$ glide get <DEPENDENCY>
```
as described in glide's documentations.

## Testing

Before you make a pull request, make sure all tests are passing. Here is a handy
snippet using [reflex](https://github.com/cespare/reflex) to run all tests every
time you change your codes:

```bash
$ reflex -d none -r "\.go$"  -- zsh -c "go test ./..."
```

If you need a separator between each run you can run:

```bash
$ reflex -d none -r "\.go$"  -- zsh -c "go test ./... ; repeat 100 printf '#'"
```

It's also a good idea to run tests with `-race` flag after the final iteration
to make sure you are not introducing any race conditions:

```bash
$ go test -race ./...
```

## Benchmarking

Benchmarking is done by running:

```bash
$ go test ./... -bench=.
```

## Pull Requests

Make sure each commit introduces one change at a time. This means if your
changes are changing a signature of a function and also adds a new feature, they
should be in two distinct commit. Make a new branch for your changes and make
the pull request based on that branch.

You can sign your commits with this command:
```bash
$ git commit -S
```

Please avoid the `Signed-off by ...` clause (-s switch).
