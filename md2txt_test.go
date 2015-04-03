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

func TestList(t *testing.T) {
	p := newParser([]byte(`* item1
* item2
* item3`))
	e := p.element()
	if e == nil {
		t.Fail()
	}
	if string(e.Content()) != `item1
item2
item3` {
		t.Logf("%s", e.Content())
		t.Fail()
	}
	if e.Type() != kind.List {
		t.Fail()
	}
	p1 := newParser([]byte(`+ item1
+ item2
+ item3`))
	e1 := p1.element()
	if e1 == nil {
		t.Fail()
	}
	if string(e1.Content()) != `item1
item2
item3` {
		t.Logf("%s", e1.Content())
		t.Fail()
	}
	if e1.Type() != kind.List {
		t.Fail()
	}

}

func TestCodeBlock(t *testing.T) {
	p := newParser([]byte(`	codeblock1
	codeblock2`))
	e := p.element()
	if string(e.Content()) != `codeblock1
codeblock2` {
		t.Logf("'%s'", e.Content())
		t.Fail()
	}
}

func TestHorizontalRules(t *testing.T) {
	p := newParser([]byte(`***`))
	e := p.element()
	if e.Type() != kind.Rule {
		t.Fail()
	}
	p1 := newParser([]byte(`* * *`))
	e1 := p1.element()
	if e1.Type() != kind.Rule {
		t.Fail()
	}
}

func TestEmpahsis(t *testing.T) {
	p := newParser([]byte("*emphasis*"))
	e := p.element()
	sp := newSpanParser(e.Content())
	s := sp.element()
	if s.Type() != kind.Emphasis {
		t.Fail()
	}
	if string(s.Content()) != "emphasis" {
		t.Fail()
	}
	if s.StartPos() != 0 {
		t.Fail()
	}

	p1 := newParser([]byte("__strong__"))
	e1 := p1.element()
	sp1 := newSpanParser(e1.Content())
	s1 := sp1.element()
	if s1.Type() != kind.Strong {
		t.Fail()
	}
	if string(s1.Content()) != "strong" {
		t.Fail()
	}
	if s1.StartPos() != 0 {
		t.Fail()
	}

	p2 := newParser([]byte("un*frigging*believable"))
	e2 := p2.element()
	sp2 := newSpanParser(e2.Content())
	s2 := sp2.element()

	if s2.Type() != kind.Emphasis {
		t.Fail()
	}
	if string(s2.Content()) != "frigging" {
		t.Fail()
	}
	if s2.StartPos() != 2 {
		t.Fail()
	}
	if string(sp2.src) != "unbelievable" {
		println("!!!!")
		t.Fail()
	}
}
func TestCode(t *testing.T) {
	p := newParser([]byte("It is 'code'"))
	e := p.element()
	sp := newSpanParser(e.Content())
	s := sp.element()
	if s.Type() != kind.Code {
		t.Fail()
	}
	if string(s.Content()) != "code" {
		t.Logf("%s", s.Content())
		t.Fail()
	}
}
func TestLink(t *testing.T) {
	p := newParser([]byte("It is [link](ref \"title\")"))
	e := p.element()
	sp := newSpanParser(e.Content())
	s := sp.element()
	if s.Type() != kind.Link {
		t.Fail()
	}
	if string(s.Content()) != "linktitleref" {
		t.Logf("%s", s.Content())
		t.Fail()
	}
}
