package md2txt

import (
	"testing"
)

// Test head
func TestHead(t *testing.T) {
	out := Md2TxtString("#h1")
	if out != "h1" {
		t.Fail()
	}
	out = Md2TxtString(`#h1
##h2`)
	if out != `h1
h2` {
		t.Log(out)
		t.Fail()
	}
}

// Test italic
func TestItalic(t *testing.T) {

}
