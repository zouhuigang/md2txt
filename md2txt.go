package md2txt

import (
	"bytes"
	"regexp"
	"unicode/utf8"

	"github.com/ggaaooppeenngg/md2txt/kind"
)

const (
	BASIC = iota // Basic Markdown based on http://daringfireball.net/projects/markdown/syntax
	GFM          // Github Flavored Markdown
)
const (
	tab    = "\t"
	sapce4 = "    s"
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

func (h *Head) Content() []byte { return h.content }
func (h *Head) Type() kind.Kind { return kind.Head }

// Paragraph represents paragraph.
type Paragraph struct {
	content []byte
}

func (p *Paragraph) Content() []byte { return p.content }
func (p *Paragraph) Type() kind.Kind { return kind.Paragraph }

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

func (l *List) Type() kind.Kind { return kind.List }

// TODO:handle sub elements
func (l *List) Content() []byte {
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

// peek peek nth rune from the p.cur,
// default is 1.
func (p *parser) peek(i ...int) rune {
	if len(i) == 0 {
		if p.cur >= len(p.src) {
			return eof
		}
		r, _ := utf8.DecodeRune(p.src[p.cur:])
		return r
	} else {
		if len(i) > 1 {
			panic("no more than one arguments allowed")
		}
		nth := i[0]
		pos := p.cur
		var (
			r rune
			w int
		)
		for k := 0; k < nth; k++ {
			r, w = utf8.DecodeRune(p.src[pos:])
			pos += w
		}
		return r
	}
}

// backup backup a rune to the src.
func (p *parser) backup() {
	p.pos.Colunm -= p.length
	p.cur -= p.length
}

// emit send element
func (p *parser) emit(e Element) {
	p.elementChan <- e
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
	case r == '-' || r == '*' || r == '+':
		r1 := p.peek(2)
		if r1 == ' ' {
			return parseList
		}
		fallthrough
	default:
		return parseParagraph
	case r == eof:
		return nil
	}

}

func parseHead(p *parser) stateFn {
	level := p.consume('#')
	for r := p.next(); r != '\n' && r != eof; {
		r = p.next()
	}
	content := p.src[p.start:p.cur]
	// deliminate suffix and prefix '#'
	content = regexp.MustCompile("#*\n$").ReplaceAll(content, []byte{})
	content = regexp.MustCompile("^#+").ReplaceAll(content, []byte{})
	head := &Head{level, content}
	p.emit(head)
	return parseBegin
}

// parse text with no prefix.
// NOTICE:if followed by '---'|'====',
// emitted as Head Type else Paragraph Type.
func parseParagraph(p *parser) stateFn {
	for r := p.next(); r != '\n' && r != eof; {
		r = p.next()
	}
	content := p.src[p.start:p.cur]
	r := p.peek()
	// Head type has tailling ----- (H2) or ====== (H1)
	if r == '-' || r == '=' {
		p.consume(r)
		r1 := p.peek(2)
		if r1 == '\n' {
			p.next()
		}
		content := p.src[p.start:p.cur]
		content = regexp.MustCompile("\n?"+string(r)+"*\n?").ReplaceAll(content, []byte{})

		var level int
		if r == '-' {
			level = 2
		}
		if r == '=' {
			level = 1
		}

		head := &Head{level, content}
		p.emit(head)
		return parseBegin
	}
	content = regexp.MustCompile("\n{0,2}$").ReplaceAll(content, []byte{})
	paragraph := &Paragraph{content}
	p.emit(paragraph)
	return parseBegin

}

// parseList parse lists with embedded sub elements.
func parseList(p *parser) stateFn {
	marker := p.peek()
	list := &List{}
	start := p.start
	for {
		for r := p.next(); r != '\n' && r != eof; {
			r = p.next()
		}
		content := p.src[start:p.cur]
		var escape = ""
		if marker == '+' || marker == '*' {
			escape = "\\"
		}
		content = regexp.MustCompile("^"+escape+string(marker)+" ").ReplaceAll(content, []byte{})
		content = regexp.MustCompile("\n$").ReplaceAll(content, []byte{})
		item := &Item{content: content}
		list.items = append(list.items, item)
		// if forsee Sprinf("%s ",marker),
		// parse another list.
		r := p.peek(1)
		r1 := p.peek(2)
		if r != marker && r1 != ' ' {
			p.emit(list)
			return parseBegin
		}
		start = p.cur
	}
}

func (p *parser) run() {
	for p.state = parseBegin; p.state != nil; {
		p.state = p.state(p)
	}
	close(p.elementChan)
}
