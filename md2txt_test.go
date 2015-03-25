package md2txt

import (
	"testing"
)

func TestHead(t *testing.T) {
	p := newParser([]byte("#Head\n"))
	e := p.element()
	if string(e.Content()) != "Head" {
		t.Logf("%s", e.Content())
		t.Fail()
	}
}
