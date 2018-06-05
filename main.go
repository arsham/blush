package main

import (
	"github.com/arsham/blush/cmd"
)

func main() {
	// defer profile.Start(profile.MemProfile, profile.CPUProfile).Stop()
	// defer profile.Start( profile.TraceProfile).Stop()
	cmd.Main()
}
