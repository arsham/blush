# Changelog

## v0.5.2
- Removed printing of Stdin when the results are piped in.

## v0.5.1
### Go 1.7 Support
- Added support for Go version 1.7

## v0.5.0
### Performance Enhancements
- Refactored Read and WriteTo methods to make Blush read the lines one by one.
- Sends read lines through a channel to be read later on with Read and WriteTo.
- Prevents reading the whole stream in one go.


## v0.4
### First Usable Release
- Reads from files and pipes.
- Uses colouring for matches.
- Groups matches from arguments.
- Accepts #RGB colours.
- Implements io.Reader and io.WriterTo for Blush struct.
