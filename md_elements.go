package md2txt

import (
	"bytes"

	"github.com/ggaaooppeenngg/md2txt/kind"
)

type Block interface {
	Spans() []Inline // span elements including pure text.
	Type() kind.Kind // kind of block.
	Content() []byte // pure text including inline.
}

// Head represents element beginning with '#'
type Head struct {
	level   int // head type h1,h2,...h6
	content []byte
}

// Head has no spans.
func (h Head) Spans() []Inline { return []Inline{} }
func (h Head) Content() []byte { return h.content }
func (h Head) Type() kind.Kind { return kind.Head }

// Paragraph represents paragraph.
type Paragraph struct {
	content []byte
}

func (p Paragraph) Spans() []Inline { return []Inline{} }
func (p Paragraph) Content() []byte { return p.content }
func (p Paragraph) Type() kind.Kind { return kind.Paragraph }

// BlockQuote represents element beginning with '>'
type BlockQuote struct {
	Level   int // level of recursive layer
	Content string
}

// List represents element beginning with '*'|'+'|'-'|digit
type List struct {
	level int // recursive level
	items []*Item
}

// list has no inline but subitems have inline elements.
func (l List) Spans() []Inline { return []Inline{} }
func (l List) Type() kind.Kind { return kind.List }

// TODO:handle sub elements
func (l List) Content() []byte {
	var output [][]byte
	for _, v := range l.items {
		output = append(output, v.content)
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

func (c CodeBlock) Spans() []Inline { return []Inline{} }
func (c CodeBlock) Content() []byte { return c.content }
func (c CodeBlock) Type() kind.Kind { return kind.CodeBlock }

// Rule represents horizontal rules
type Rule struct {
}

func (r Rule) Spans() []Inline { return []Inline{} }
func (r Rule) Content() []byte { return []byte{} }
func (r Rule) Type() kind.Kind { return kind.Rule }

// inline span elements.
type Inline interface {
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
	Text  string
	Title string
	URL   string
}

type Image struct {
	Text  string
	Title string
	Link  string
}

type InlineHTML struct {
}

type EOF struct {
}
