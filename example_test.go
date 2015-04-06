package md2txt

import (
	"fmt"
)

func ExampleParse() {
	ret := Parse([]byte(`This is a list:

* item
* item
* item

and this is a new paragraph.`))

	fmt.Printf("%s", ret)
	// Output:
	// This is a list:item
	// item
	// item
	// and this is a new paragraph.
}
