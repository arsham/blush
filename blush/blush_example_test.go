package blush_test

import (
	"bytes"
	"fmt"
	"io"

	"github.com/arsham/blush/blush"
)

func ExampleBlush() {
	f := blush.NewExact("sword", blush.Red)
	r := bytes.NewBufferString("He who lives by the sword, will surely also die")
	b := &blush.Blush{
		Finders: []blush.Finder{f},
		Reader:  io.NopCloser(r),
	}
	b.Close()
}

func ExampleBlush_Read() {
	var p []byte
	f := blush.NewExact("sin", blush.Red)
	r := bytes.NewBufferString("He who lives in sin, will surely live the lie")
	b := &blush.Blush{
		Finders: []blush.Finder{f},
		Reader:  io.NopCloser(r),
	}
	b.Read(p)
}

func ExampleBlush_Read_inDetails() {
	f := blush.NewExact("sin", blush.Red)
	r := bytes.NewBufferString("He who lives in sin, will surely live the lie")
	b := &blush.Blush{
		Finders: []blush.Finder{f},
		Reader:  io.NopCloser(r),
	}
	expect := fmt.Sprintf("He who lives in %s, will surely live the lie", f)

	// you should account for the additional characters for colour formatting.
	length := r.Len() - len("sin") + len(f.String())
	p := make([]byte, length)
	n, err := b.Read(p)
	fmt.Println("n == len(p):", n == len(p))
	fmt.Println("err:", err)
	fmt.Println("p == expect:", string(p) == expect)
	// by the way
	fmt.Println(`f == "sin":`, f.String() == "sin")

	// Output:
	// n == len(p): true
	// err: <nil>
	// p == expect: true
	// f == "sin": false
}

func ExampleBlush_WriteTo() {
	f := blush.NewExact("victor", blush.Red)
	r := bytes.NewBufferString("It is a shield of passion and strong will from this I am the victor instead of the kill\n")
	b := &blush.Blush{
		Finders: []blush.Finder{f},
		Reader:  io.NopCloser(r),
	}
	buf := new(bytes.Buffer)
	n, err := b.WriteTo(buf)

	expected := fmt.Sprintf("It is a shield of passion and strong will from this I am the %s instead of the kill\n", f)
	fmt.Println("err:", err)
	fmt.Println("n == len(expected):", int(n) == len(expected))
	fmt.Println("buf.String() == expected:", buf.String() == expected)

	// Output:
	// err: <nil>
	// n == len(expected): true
	// buf.String() == expected: true
}

func ExampleBlush_WriteTo_copy() {
	f := blush.NewExact("you feel", blush.Cyan)
	r := bytes.NewBufferString("Savour what you feel and what you see\n")
	b := &blush.Blush{
		Finders: []blush.Finder{f},
		Reader:  io.NopCloser(r),
	}
	buf := new(bytes.Buffer)
	n, err := io.Copy(buf, b)

	expected := fmt.Sprintf("Savour what %s and what you see\n", f)
	fmt.Println("err:", err)
	fmt.Println("n == len(expected):", int(n) == len(expected))
	fmt.Println("buf.String() == expected:", buf.String() == expected)

	// Output:
	// err: <nil>
	// n == len(expected): true
	// buf.String() == expected: true
}

func ExampleBlush_WriteTo_multiReader() {
	mg := blush.NewExact("truth", blush.Magenta)
	g := blush.NewExact("Life", blush.Green)
	r1 := bytes.NewBufferString("Life is like a mystery with many clues, but with few answers\n")
	r2 := bytes.NewBufferString("To tell us what it is that we can do to look for messages that keep us from the truth\n")
	mr := io.MultiReader(r1, r2)

	b := &blush.Blush{
		Finders: []blush.Finder{mg, g},
		Reader:  io.NopCloser(mr),
	}
	buf := new(bytes.Buffer)
	b.WriteTo(buf)
}

func ExampleBlush_WriteTo_multiReaderInDetails() {
	mg := blush.NewExact("truth", blush.Magenta)
	g := blush.NewExact("Life", blush.Green)
	r1 := bytes.NewBufferString("Life is like a mystery with many clues, but with few answers\n")
	r2 := bytes.NewBufferString("To tell us what it is that we can do to look for messages that keep us from the truth\n")
	mr := io.MultiReader(r1, r2)

	b := &blush.Blush{
		Finders: []blush.Finder{mg, g},
		Reader:  io.NopCloser(mr),
	}
	buf := new(bytes.Buffer)
	n, err := b.WriteTo(buf)

	line1 := fmt.Sprintf("%s is like a mystery with many clues, but with few answers\n", g)
	line2 := fmt.Sprintf("To tell us what it is that we can do to look for messages that keep us from the %s\n", mg)
	expected := line1 + line2
	fmt.Println("err:", err)
	fmt.Println("n == len(expected):", int(n) == len(expected))
	fmt.Println("buf.String() == expected:", buf.String() == expected)

	// Output:
	// err: <nil>
	// n == len(expected): true
	// buf.String() == expected: true
}
