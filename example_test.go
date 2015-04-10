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
	// Output:
	// quote
}

func ExampleList() {
	ret := Parse([]byte(`*   Lorem ipsum dolor sit amet, consectetuer adipiscing elit.
    Aliquam hendrerit mi posuere lectus. Vestibulum enim wisi,
    viverra nec, fringilla in, laoreet vitae, risus.
*   Donec sit amet nisl. Aliquam semper ipsum sit amet velit.
    Suspendisse id sem consectetuer libero luctus adipiscing.`), BASIC)
	fmt.Printf("%s", ret)
	// Output:
	// Lorem ipsum dolor sit amet, consectetuer adipiscing elit.
	// Aliquam hendrerit mi posuere lectus. Vestibulum enim wisi,
	// viverra nec, fringilla in, laoreet vitae, risus.
	// Donec sit amet nisl. Aliquam semper ipsum sit amet velit.
	// Suspendisse id sem consectetuer libero luctus adipiscing.
}
