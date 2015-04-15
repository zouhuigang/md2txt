package md2txt

import (
	"bytes"

	"github.com/ggaaooppeenngg/md2txt/kind"
)

// Block is the block element interface.
type Block interface {
	Type() kind.Kind // kind of block.
	Content() []byte // pure text including inline.
}

// TODO:support block html
type BlockHtml struct {
}

// Head represents element beginning with '#'
type Head struct {
	level   int // head type h1,h2,...h6
	content []byte
}

func (h Head) Content() []byte { return h.content }
func (h Head) Type() kind.Kind { return kind.Head }

// Paragraph represents paragraph.
type Paragraph struct {
	content []byte
}

func (p Paragraph) Content() []byte {
	var (
		sp     = newSpanParser(p.content)
		spans  []Span
		length int
	)
	for s := sp.element(); s != nil; s = sp.element() {
		spans = append(spans, s)
	}
	p.content = sp.src

	for _, v := range spans {
		p.content = append(p.content[:length+v.StartPos()], append(v.Content(), p.content[length+v.StartPos():]...)...)
		length += len(v.Content())
	}
	return p.content
}

func (p Paragraph) Type() kind.Kind { return kind.Paragraph }

// BlockQuote represents element beginning with '>'
type QuoteBlock struct {
	level     int // level of recursive layer
	content   []byte
	subBlocks []Block
}

func (q QuoteBlock) Content() []byte {
	var contents [][]byte
	for _, v := range q.subBlocks {
		contents = append(contents, v.Content())
	}
	return bytes.Join(contents, []byte("\n"))
}

func (q QuoteBlock) Type() kind.Kind { return kind.QuoteBlock }

// List represents element beginning with '*'|'+'|'-'|digit
type List struct {
	level int // recursive level
	items []*Item
}

// list has no inline but subitems have inline elements.
func (l List) Type() kind.Kind { return kind.List }

// TODO:handle sub elements
func (l List) Content() []byte {
	var output [][]byte
	for _, v := range l.items {
		output = append(output, v.content)
		for _, b := range v.subBlocks {
			output = append(output, b.Content())
		}
	}
	return bytes.Join(output, []byte("\n"))
}

// list item.
type Item struct {
	content   []byte
	subBlocks []Block
}

// CodeBlock represents element beginning with one tab or at least a 4 spaces.
type CodeBlock struct {
	level   int // recursive level
	content []byte
}

func (c CodeBlock) Content() []byte { return c.content }
func (c CodeBlock) Type() kind.Kind { return kind.CodeBlock }

// Rule represents horizontal rules
type Rule struct {
}

func (r Rule) Content() []byte { return []byte{} }
func (r Rule) Type() kind.Kind { return kind.Rule }

// inline span elements.
type Span interface {
	StartPos() int
	Content() []byte
	Type() kind.Kind
}

type Emphasis struct {
	start   int
	content []byte
}

func (e Emphasis) Type() kind.Kind { return kind.Emphasis }
func (e Emphasis) Content() []byte { return e.content }
func (e Emphasis) StartPos() int   { return e.start }

type Strong struct {
	start   int
	content []byte
}

func (s Strong) Type() kind.Kind { return kind.Strong }
func (s Strong) Content() []byte { return s.content }
func (s Strong) StartPos() int   { return s.start }

type Code struct {
	start   int
	content []byte
}

func (c Code) Type() kind.Kind { return kind.Code }
func (c Code) Content() []byte { return c.content }
func (c Code) StartPos() int   { return c.start }

type Link struct {
	start int
	id    []byte
	text  []byte
	title []byte
	url   []byte
}

func (l Link) Type() kind.Kind { return kind.Link }
func (l Link) Content() []byte { return bytes.Join([][]byte{l.text, l.title, l.url}, []byte{}) }
func (l Link) StartPos() int   { return l.start }

type Image struct {
	start int
	id    []byte
	text  []byte
	title []byte
	link  []byte
}

func (i Image) Type() kind.Kind { return kind.Image }
func (i Image) Content() []byte { return bytes.Join([][]byte{i.text, i.title, i.link}, []byte{}) }
func (i Image) StartPos() int   { return i.start }

// TODO: support inline html
type InlineHTML struct {
}
