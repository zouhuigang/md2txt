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

func TestQuote(t *testing.T) {
	p := newParser([]byte(">quote\n"))
	b := p.element().(*QuoteBlock)
	v := b.subBlocks[0]
	if string(v.Content()) != `quote` {
		t.Fail()
	}
}
func TestRecursiveQuote(t *testing.T) {
	p := newParser([]byte(`>quote1
>
>>quote2
>
>quote1`))
	b := p.element().(*QuoteBlock)
	if string(b.Content()) != `quote1
quote2
quote1` {
		t.Logf("%s", b.Content())
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
	sp := newSpanParser([]byte("*emphasis*"))
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

	sp1 := newSpanParser([]byte("__strong__"))
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

	sp2 := newSpanParser([]byte("un*frigging*believable"))
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
		t.Fail()
	}
}
func TestCode(t *testing.T) {
	sp := newSpanParser([]byte("It is 'code'"))
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
	sp := newSpanParser([]byte("It is [link](ref \"title\")"))
	s := sp.element()
	if s.Type() != kind.Link {
		t.Fail()
	}
	if string(s.Content()) != "linktitleref" {
		t.Logf("%s", s.Content())
		t.Fail()
	}
	if string(sp.src) != "It is " {
		t.Logf("%s", sp.src)
		t.Fail()
	}
}

func TestImage(t *testing.T) {
	sp := newSpanParser([]byte("It is ![image](ref \"title\")"))
	s := sp.element()
	if s.Type() != kind.Image {
		t.Fail()
	}
	if string(s.Content()) != "imagetitleref" {
		t.Logf("%s", s.Content())
		t.Fail()
	}
	if string(sp.src) != "It is " {
		t.Logf("%s", sp.src)
		t.Fail()
	}

}

func TestReference(t *testing.T) {
	sp := newSpanParser([]byte("[id]: link \"title\"\n"))
	sp.element()
	ref := sp.ref["id"]
	if string(ref.link) != "link" {
		t.Fail()
	}
	if string(ref.title) != "title" {
		t.Fail()
	}
}

func TestExample1(t *testing.T) {
	s := newParser([]byte(`This is a *paragraph*,this is a
[link](http://www.baidu.com "百度") and this is a ![image](http://www.baidu.com "百度图片")`))
	p := s.element()
	if string(p.Content()) != `This is a paragraph,this is a
link百度http://www.baidu.com and this is a image百度图片http://www.baidu.com` {
		t.Fail()
	}
}
