package cmd_test

import (
	"fmt"

	"github.com/arsham/blush/blush"
	"github.com/arsham/blush/cmd"
)

type colourer interface {
	Colour() blush.Colour
}

func ExampleGetBlush_red() {
	input := []string{"blush", "--red", "term", "/"}
	b, err := cmd.GetBlush(input)
	fmt.Println("err == nil:", err == nil)
	fmt.Println("Finders count:", len(b.Finders))
	c := b.Finders[0].(colourer)
	fmt.Println("Is red:", c.Colour() == blush.FgRed)

	// Output:
	// err == nil: true
	// Finders count: 1
	// Is red: true
}

func ExampleGetBlush_multiColour() {
	input := []string{"-b", "term1", "-g", "term2", "/"}
	b, err := cmd.GetBlush(input)
	c1 := b.Finders[0].(colourer)
	c2 := b.Finders[1].(colourer)
	fmt.Println("err == nil:", err == nil)
	fmt.Println("Finders count:", len(b.Finders))
	fmt.Println("Is blue:", c1.Colour() == blush.FgBlue)
	fmt.Println("Is green:", c2.Colour() == blush.FgGreen)

	// Output:
	// err == nil: true
	// Finders count: 2
	// Is blue: true
	// Is green: true
}
