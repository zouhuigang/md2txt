package md2txt

import (
	"regexp"
	"unicode/utf8"

	"github.com/ggaaooppeenngg/md2txt/kind"
)

const (
	BASIC = iota // Basic Markdown based on http://daringfireball.net/projects/markdown/syntax
	GFM          // Github Flavored Markdown
)

type Element interface {
	Type() kind.Kind
	Content() []byte
}

// Head represents element beginning with '#'
type Head struct {
	level   int
	content []byte
}

func (h *Head) Content() []byte {
	return h.content
}
func (h *Head) Type() kind.Kind {
	return kind.Head
}

// Paragraph represents paragraph.
type Paragraph struct {
	Content string
}

// BlockQuote represents element beginning with '>'
type BlockQuote struct {
	Level   int // level of recursive layer
	Content string
}

// List represents element beginning with '*'|'+'|'-'|digit
type List struct {
	Level   int // recursive level
	Content string
}

// Code represents element beginning with one tab or at least a 4 spaces.
type Code struct {
	Level   int // recursive level
	Content string
}

// Rule represents horizontal rules
type Rule struct {
}

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

type Position struct {
	Row    int
	Colunm int
}

type stateFn func(p *parser) stateFn

type parser struct {
	src []byte
	pos Position

	start  int // start index
	cur    int // current index
	length int // length of scanned content

	state       stateFn
	elementChan chan Element
}

func newParser(src []byte) *parser {
	p := &parser{
		src:         src,
		elementChan: make(chan Element),
	}
	go p.run()
	return p
}
func Parse() {

}

const eof = -1

func (p *parser) element() Element {
	e := <-p.elementChan
	return e
}

// consume repeatly consume rune equal to r.
func (p *parser) consume(r rune) int {
	var count int
	for pr := p.peek(); pr == r; {
		pr = p.next()
		count++
	}
	return count
}

// next returns next rune,if reach end of file returns EOF.
func (p *parser) next() rune {
	if p.cur >= len(p.src) {
		p.length = 0
		return eof
	}
	r, w := utf8.DecodeRune(p.src[p.cur:])
	p.length = w
	p.cur += p.length
	p.pos.Colunm += p.length
	return r
}

// peek peek a rune from the src,
func (p *parser) peek() rune {
	if p.cur >= len(p.src) {
		return eof
	}
	r, _ := utf8.DecodeRune(p.src[p.cur:])
	return r
}

// backup backup a rune to the src.
func (p *parser) backup() {
	p.pos.Colunm -= p.length
	p.cur -= p.length
}

// emit send element
func (p *parser) emit(e Element) {
	p.elementChan <- e
	if e.Type() == kind.EOF {
		close(p.elementChan)
	}
	p.start = p.cur
}

// lines returns line number of the input
func numberOfLines(input []byte) int {
	var count int
	for _, v := range input {
		if v == '\n' {
			count++
		}
	}
	return count
}

func parseBegin(p *parser) stateFn {
	switch r := p.peek(); {
	case r == '#':
		return parseHead
	case r == eof:
		return nil
	default:
		return nil
	}

}

func parseHead(p *parser) stateFn {
	level := p.consume('#')
	for r := p.next(); r != '\n' && r != eof; {
		r = p.next()
	}
	content := p.src[p.start:p.cur]
	content = regexp.MustCompile("#*\n").ReplaceAll(content, []byte{})
	content = regexp.MustCompile("^#+").ReplaceAll(content, []byte{})
	head := &Head{level, content}
	p.emit(head)
	return parseBegin
}

func (p *parser) run() {
	for p.state = parseBegin; p.state != nil; {
		p.state = p.state(p)
	}
}
