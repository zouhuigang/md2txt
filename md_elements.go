package md2txt

import (
	"bytes"

	"github.com/ggaaooppeenngg/md2txt/kind"
)

type Element interface {
	Type() kind.Kind
	Content() []byte
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
	content []byte
	subEles []Element
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

type Link struct {
	Text  string
	Title string
	URL   string
}

type Emphasis struct {
	Content string
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
