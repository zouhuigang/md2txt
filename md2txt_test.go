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
