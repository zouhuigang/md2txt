package md2txt

import (
	"testing"

	"github.com/zouhuigang/md2txt/kind"
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

func TestQuoteOneLine(t *testing.T) {
	p := newParser([]byte(`> This is a blockquote with two paragraphs. Lorem ipsum dolor sit amet,
consectetuer adipiscing elit. Aliquam hendrerit mi posuere lectus.
Vestibulum enim wisi, viverra nec, fringilla in, laoreet vitae, risus.
`))
	for b := p.element(); b != nil; b = p.element() {
		if string(b.Content()) != `This is a blockquote with two paragraphs. Lorem ipsum dolor sit amet,
consectetuer adipiscing elit. Aliquam hendrerit mi posuere lectus.
Vestibulum enim wisi, viverra nec, fringilla in, laoreet vitae, risus.` {
			t.Fail()

		}
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

func TestItemSubBlocks(t *testing.T) {
	bs1, n := parseItemBlocks([]byte(`    subBlocks
in lazy mode
`))
	for _, b := range bs1 {
		if string(b.Content()) != `subBlocks
in lazy mode` {
			t.Logf("%s", b.Content())
			t.Fail()
		}
		if b.Type() != kind.Paragraph {
			t.Fail()
		}
		if n != 27 {
			t.Logf("%d", n)
			t.Fail()
		}
	}

	bs2, n := parseItemBlocks([]byte(`    > subBlocks
	> with heading indents

`))
	for _, b := range bs2 {
		if string(b.Content()) != `subBlocks
with heading indents` {
			t.Fail()
		}
		if b.Type() != kind.QuoteBlock {
			t.Fail()
		}
		if n != 41 {
			t.Logf("%d", n)
			t.Fail()
		}
	}
}

func TestQuoteContainingOtherBlocks(t *testing.T) {
	p := newParser([]byte(`> ## This is a header.
> 
> 1.  This is the first list item.
> 2.  This is the second list item.
> 
> Here's some example code:
> 
>     return shell_exec("echo $input | $markdown_script");`))
	e := p.element().(*QuoteBlock)
	var types = []kind.Kind{kind.Head, kind.List, kind.Paragraph, kind.CodeBlock}
	for i := 0; i < len(e.subBlocks); i++ {
		if e.subBlocks[i].Type() != types[i] {
			t.Fail()
		}
	}
}

func TestList(t *testing.T) {
	p := newParser([]byte(`*   Lorem ipsum dolor sit amet, consectetuer adipiscing elit.
    Aliquam hendrerit mi posuere lectus. Vestibulum enim wisi,
    viverra nec, fringilla in, laoreet vitae, risus.
*   Donec sit amet nisl. Aliquam semper ipsum sit amet velit.
    Suspendisse id sem consectetuer libero luctus adipiscing.`))
	e := p.element()
	if e == nil {
		t.Fail()
	}
	if string(e.Content()) != `Lorem ipsum dolor sit amet, consectetuer adipiscing elit.
Aliquam hendrerit mi posuere lectus. Vestibulum enim wisi,
viverra nec, fringilla in, laoreet vitae, risus.
Donec sit amet nisl. Aliquam semper ipsum sit amet velit.
Suspendisse id sem consectetuer libero luctus adipiscing.` {
		t.Logf("%s", e.Content())
		t.Fail()
	}
}

func TestOrderList(t *testing.T) {
	p := newParser([]byte(`1.  item1
2.  item2
3.  item3`))
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
}

func TestUnorderList(t *testing.T) {
	p := newParser([]byte(`*   item1
*   item2
*   item3`))
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

	p1 := newParser([]byte(`+   item1
+   item2
+   item3`))
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
	p2 := newParser([]byte(`* item
* item
* item`))

	e2 := p2.element()
	if e2 == nil {
		t.Fail()
	}

	if string(e2.Content()) != `* item
* item
* item` {
		t.Logf("%s", e2.Content())
		t.Fail()
	}
}

func TestListWithSubItems(t *testing.T) {
	p := newParser([]byte(`1.  This is a list item with two paragraphs. Lorem ipsum dolor
    sit amet, consectetuer adipiscing elit. Aliquam hendrerit
    mi posuere lectus.

    Vestibulum enim wisi, viverra nec, fringilla in, laoreet
    vitae, risus. Donec sit amet nisl. Aliquam semper ipsum
    sit amet velit.

2.  Suspendisse id sem consectetuer libero luctus adipiscing.`))
	e1 := p.element()
	if string(e1.Content()) != `This is a list item with two paragraphs. Lorem ipsum dolor
sit amet, consectetuer adipiscing elit. Aliquam hendrerit
mi posuere lectus.
Vestibulum enim wisi, viverra nec, fringilla in, laoreet
vitae, risus. Donec sit amet nisl. Aliquam semper ipsum
sit amet velit.
Suspendisse id sem consectetuer libero luctus adipiscing.` {
		t.Logf("%s", e1.Content())
		t.Fail()

	}

}
func TestCodeBlock(t *testing.T) {
	p := newParser([]byte(`	codeblock1`))
	e := p.element()
	if string(e.Content()) != `codeblock1` {
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

func TestEmphasis(t *testing.T) {
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
	sp := newSpanParser([]byte("It is `code`"))
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
