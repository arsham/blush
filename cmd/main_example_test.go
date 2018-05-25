package cmd_test

import (
	"fmt"

	"github.com/arsham/blush/blush"
	"github.com/arsham/blush/cmd"
)

func ExampleGetBlush_red() {
	input := []string{"blush", "--red", "term", "/"}
	b, err := cmd.GetBlush(input)
	fmt.Println("err == nil:", err == nil)
	fmt.Println("Finders count:", len(b.Finders))
	fmt.Println("Is red:", b.Finders[0].Colour() == blush.FgRed)

	// Output:
	// err == nil: true
	// Finders count: 1
	// Is red: true
}

func ExampleGetBlush_multiColour() {
	input := []string{"-b", "term1", "-g", "term2", "/"}
	b, err := cmd.GetBlush(input)
	fmt.Println("err == nil:", err == nil)
	fmt.Println("Finders count:", len(b.Finders))
	fmt.Println("Is blue:", b.Finders[0].Colour() == blush.FgBlue)
	fmt.Println("Is green:", b.Finders[1].Colour() == blush.FgGreen)

	// Output:
	// err == nil: true
	// Finders count: 2
	// Is blue: true
	// Is green: true
}
