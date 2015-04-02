package md2txt

import (
	"math"
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
	go p.run(kind.Block)
	return p
}
func newSpanParser(src []byte) *parser {
	p := &parser{
		src:         src,
		elementChan: make(chan Element),
	}
	go p.run(kind.Inline)
	return p
}
func Parse() {

}

const eof = -1

//----------type dependent----------

// get element,type dependent
func (p *parser) element() Element {
	e := <-p.elementChan
	return e
}

// emit send element,type dependent
func (p *parser) emit(b Element) {
	p.elementChan <- b
	p.start = p.cur
}

//----------------------------------

// consume repeatly consume rune equal to r.
func (p *parser) consume(r rune, limit ...int64) int {
	if len(limit) > 1 {
		panic("only one argument is allowed")
	}
	var max int64 = math.MaxInt64
	if len(limit) == 1 {
		max = limit[0]
	}
	var count int
	for pr := p.peek(); pr == r && count < int(max); {
		pr = p.next()
		count++
	}
	return count
}

// merge escape runes(like \*,\_),to one rune.
func (p *parser) merge() {
	if p.cur+1 >= len(p.src) {
		return
	}
	p.src = append(p.src[:p.cur], p.src[p.cur+1:]...)
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

// forsee run of runes.
func (p *parser) forsee(rs ...rune) bool {
	pos := p.cur
	for k := 0; k < len(rs); k++ {
		r, w := utf8.DecodeRune(p.src[pos:])
		pos += w
		if r != rs[k] {
			return false
		}
	}
	return true
}

// runes return the runes consisting of the string.
func runes(s string) []rune {
	var runes []rune
	for {
		r, w := utf8.DecodeRuneInString(s)
		if w != 0 {
			runes = append(runes, r)
			s = s[w:]
		} else {
			break
		}
	}
	return runes
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

// parseHead parse head beginning with '#'
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
		// TODO: use forsee to detect next item.
		r := p.peek(1)
		r1 := p.peek(2)
		if r != marker && r1 != ' ' {
			p.emit(list)
			return parseBegin
		}
		start = p.cur
	}
}

// parseCode parses code beginning with 4 sapces or 1 tab.
func parseCodeBlock(p *parser) stateFn {
	codeBlock := &CodeBlock{}
	start := p.start
	var marker string
	r := p.peek()
	if r == ' ' {
		marker = "    " // 4 spaces
	}
	marker = "\t"
	for {
		for r := p.next(); r != '\n' && r != eof; {
			r = p.next()
		}
		content := p.src[start:p.cur]
		content = regexp.MustCompile("^"+marker).ReplaceAll(content, []byte{})
		codeBlock.content = append(codeBlock.content, content...)
		if !p.forsee(runes(marker)...) {
			break
		}
		start = p.cur
	}
	p.emit(codeBlock)
	return parseBegin

}

// parseRule parses rule begining,
// with more than three '*'|'+'|'-'
// (can be joined by one white sapce).
func parseRule(p *parser) stateFn {
	r := p.next()
	r1 := p.peek(2)
	if r1 == ' ' {
		for {
			r1 = p.next()
			if r1 != ' ' {
				break
			}
			r1 = p.next()
			if r1 != r {
				break
			}
		}
		p.emit(&Rule{})
		return parseBegin
	} else {
		p.consume(r)
		p.emit(&Rule{})
		return parseBegin
	}

}

// parseError is error handler when account for errors.
func parseError(p *parser) stateFn {
	return nil
}

// block main parsing.
func parseBegin(p *parser) stateFn {
	switch r := p.peek(); {
	case r == '#':
		return parseHead
	case r == '-' || r == '*' || r == '+':
		r1 := p.peek(2)
		if r1 == ' ' {
			// rule marker can be seperated by one white sapce.
			if r == '*' || r == '-' {
				r2 := p.peek(3)
				if r2 == r && p.forsee(r, ' ', r, ' ', r) {
					// fall to parseRule
					return parseRule
				}
			}
			return parseList
		}
		fallthrough
	case r == '_' || r == '*' || r == '-':
		r1 := p.peek(2)
		if r1 == ' ' {
			if p.forsee(r, ' ', r, ' ', r) {
				return parseRule
			}
		}
		if r1 == r && p.forsee(r, r, r) {
			return parseRule
		}
		fallthrough
	case r == '\t' || (r == ' ' && p.forsee(' ', ' ', ' ')):
		return parseCodeBlock
	default:
		return parseParagraph
	case r == eof:
		return nil
	}

}

//-----------span parsing----------

// parse emphasis or strong span
func parseEmphasis(p *parser) stateFn {
	start := p.cur
	marker := p.peek()
	n := p.consume(marker, 2)
	t := kind.Strong
	if n == 1 {
		t = kind.Emphasis
	}

	for r := p.next(); r != marker && r != '\n' && r != eof; {
		if r == '\\' {
			r1 := p.peek()
			if isEscapeRune(r1) {
				p.merge()
				p.next()
			}
		}
		r = p.next()
	}
	r := p.peek()
	f := string(marker)
	if f == "*" {
		f = "\\" + f
	}
	content := p.src[p.start:p.cur]
	content = regexp.MustCompile("^"+f+"+").ReplaceAll(content, []byte{})
	content = regexp.MustCompile(""+f+"+$").ReplaceAll(content, []byte{})
	p.src = append(p.src[:p.start], p.src[p.cur:]...)

	if r == marker && t == kind.Strong {
		p.next()
		p.emit(&Strong{start, content})
	} else {
		p.emit(&Emphasis{start, content})
	}
	return parseSpan
}

var escapeRunes = "\\'*_{}[]()#+-.!"

func isEscapeRune(r rune) bool {
	for _, v := range escapeRunes {
		if rune(v) == r {
			return true
		}
	}
	return false
}

// span main parsing.
func parseSpan(p *parser) stateFn {
	switch r := p.peek(); {
	case r == '\\':
		r1 := p.peek(2)
		// escape runes
		if isEscapeRune(r1) {
			p.next()
			p.merge()
			p.next()
		}
		fallthrough
	case r == '*' || r == '_':
		return parseEmphasis
	case r == eof:
		return nil
	default:
		return nil
	}
}

func (p *parser) run(t kind.ElementType) {
	if t == kind.Block {
		for p.state = parseBegin; p.state != nil; {
			p.state = p.state(p)
		}
	}
	if t == kind.Inline {
		for p.state = parseSpan; p.state != nil; {
			p.state = p.state(p)
		}
	}
	close(p.elementChan)
}
