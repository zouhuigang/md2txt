package md2txt

import (
	"testing"

	"github.com/ggaaooppeenngg/md2txt/kind"
)

func TestHead(t *testing.T) {
	p := newParser([]byte("#头部\n"))
	e := p.element()
	if string(e.Content()) != "头部" {
		t.Logf("%s", e.Content())
		t.Fail()
	}
	if e.Type() != kind.Head {
		t.Fail()
	}
}

func TestParagraphHead(t *testing.T) {
	p := newParser([]byte("一级头部\n======\n"))
	e := p.element()
	if string(e.Content()) != "一级头部" {
		t.Logf("%s", e.Content())
		t.Fail()
	}
	if e.Type() != kind.Head {
		t.Fail()
	}
}

func TestParagraph(t *testing.T) {
	p := newParser([]byte("一级头部\n"))
	e := p.element()
	if string(e.Content()) != "一级头部" {
		t.Logf("%s", e.Content())
		t.Fail()
	}
	if e.Type() != kind.Paragraph {
		t.Fail()
	}

}
