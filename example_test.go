package md2txt

import (
	"fmt"
)

func ExampleParse() {
	ret := Parse([]byte(`This is a list:

*   item
*   item
*   item

and this is a new paragraph.`), BASIC)

	fmt.Printf("%s", ret)
	// Output:
	// This is a list:
	// item
	// item
	// item
	// and this is a new paragraph.
}

func ExampleHead_H1() {
	ret := Parse([]byte(`This is an H1
=============`), BASIC)
	fmt.Printf("%s", ret)
	// Output:
	// This is an H1
}

func ExampleHead_H2() {
	ret := Parse([]byte(`## This is an H2`), BASIC)
	fmt.Printf("%s", ret)
	// Output:
	// This is an H2
}

func ExampleQuote() {
	ret := Parse([]byte(`> quote`), BASIC)
	fmt.Printf("%s", ret)
	// Output
	// quote
}
