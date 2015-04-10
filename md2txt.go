/*
Package md2txt implements a tool to convert markdown to pure text.
It uses no regexp,and gains more efficiency.
*/
package md2txt

import (
	"bytes"
	"math"
	"regexp"
	"unicode"
	"unicode/utf8"

	"github.com/ggaaooppeenngg/md2txt/kind"
)

type EXT int

const (
	BASIC EXT = iota // Basic Markdown based on http://daringfireball.net/projects/markdown/syntax
	GFM              // Github Flavored Markdown
)

const (
	tab    = "\t"
	sapce4 = "    "
)

// state for block parser.
type stateFn func(p *blockParser) stateFn

// state for span parser.
type spanStateFn func(p *spanParser) spanStateFn

// parser is a main part for a parsing procedure,
// it provides convinient methods for parsing.
type parser struct {
	src []byte // source bytes slice.

	start  int // start index.
	cur    int // current index.
	length int // length of scanned content.
}

// reference is used in link or image,
// it's format is [id]: url "title".
type reference struct {
	link  []byte
	title []byte
}

// span parser aims at span elements parsing.
type spanParser struct {
	*parser
	ref      map[string]*reference
	state    spanStateFn
	spanChan chan Span
}

// element gets a span from the channel,
// return nil if no more span elements.
func (p *spanParser) element() Span { return <-p.spanChan }

// emit emits a span element to the channel.
func (p *spanParser) emit(s Span) {
	p.spanChan <- s
	p.start = p.cur
}

// run is the main procedure of the state machine for span elements parsing,
// when it runs into the end close the channel.
func (p *spanParser) run() {
	for p.state = parseSpan; p.state != nil; {
		p.state = p.state(p)
	}
	close(p.spanChan)
}

type blockParser struct {
	*parser
	state     stateFn
	blockChan chan Block
}

// element gets a block from the channel,
// return nil if no more span elements.
func (p *blockParser) element() Block { return <-p.blockChan }

// emit emits a block element to the channel.
func (p *blockParser) emit(b Block) {
	p.blockChan <- b
	p.start = p.cur
}

// run is the main procedure of the state machine for block elements parsing,
// when it runs into the end close the channel.
func (p *blockParser) run() {
	for p.state = parseBegin; p.state != nil; {
		p.state = p.state(p)
	}
	close(p.blockChan)
}

// newParser returns a blockParser for parsing src.
func newParser(src []byte) *blockParser {
	p := &parser{
		src: src,
	}
	bp := &blockParser{parser: p, blockChan: make(chan Block)}
	go bp.run()
	return bp
}

// newSpanParser returns a spanParser for parsing src.
func newSpanParser(src []byte) *spanParser {
	p := &parser{
		src: src,
	}
	sp := &spanParser{parser: p, ref: make(map[string]*reference), spanChan: make(chan Span)}
	go sp.run()
	return sp
}

const eof = -1

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

// ignore current rune
func (p *parser) ignore() {
	p.start = p.cur
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
	//p.pos.Colunm += p.length
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
	//p.pos.Colunm -= p.length
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
func parseHead(p *blockParser) stateFn {
	level := p.consume('#')
	for r := p.next(); r != '\n' && r != eof; {
		r = p.next()
	}
	content := p.src[p.start:p.cur]
	// deliminate headiing and tailing shaps or spaces.
	content = bytes.TrimFunc(content, func(r rune) bool {
		if r == '#' || r == ' ' {
			return true
		}
		return false
	})
	content = bytes.TrimSpace(content)
	head := &Head{level, content}
	p.emit(head)
	return parseBegin
}

// parse text with no prefix.
// NOTICE:if followed by '---'|'====',
// emitted as Head Type else Paragraph Type.
func parseParagraph(p *blockParser) stateFn {
	for r := p.next(); r != eof; r = p.next() {
		if r == '\n' {
			r = p.peek()
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
			if r == '\n' {
				p.next()
				goto emit
			}
		}
	}
emit:
	content := p.src[p.start:p.cur]
	content = regexp.MustCompile("\n{0,2}$").ReplaceAll(content, []byte{})
	paragraph := &Paragraph{content: content}
	p.emit(paragraph)
	return parseBegin

}

// parseOrderlist parses order lists with embedded sub elements.
func parseOrderList(p *blockParser) stateFn {
	list := &List{}
	start := p.start
	for {
		for r := p.next(); r != eof; {
			if r == '\n' {
				r1 := p.peek()
				if r1 == '\n' ||
					r1 == eof ||
					regexp.MustCompile(`^\d+\.  `).Match(p.src[p.cur:]) {
					break
				}
				if p.forsee(' ', ' ', ' ', ' ') {
					p.src = append(p.src[:p.cur], p.src[p.cur+4:]...)
					continue
				}
				if p.forsee('\t') {
					p.src = append(p.src[:p.cur], p.src[p.cur+1:]...)
					continue
				}

			}
			r = p.next()
		}
		content := p.src[start:p.cur]
		content = regexp.MustCompile(`^\d+\.  `).ReplaceAllLiteral(content, []byte{})
		content = bytes.TrimRightFunc(content, func(r rune) bool {
			if r == '\n' {
				return true
			}
			return false
		})
		content = regexp.MustCompile("\n$").ReplaceAll(content, []byte{})
		item := &Item{content: content}
		list.items = append(list.items, item)

		if !regexp.MustCompile(`^\d+\.  `).Match(p.src[p.cur:]) {
			p.emit(list)
			return parseBegin
		}
		start = p.cur
	}
}

// parseUnorderList parses unorder lists with embedded sub elements.
func parseUnorderList(p *blockParser) stateFn {
	marker := p.peek()
	escape := ""
	if marker == '+' || marker == '*' {
		escape = "\\"
	}
	list := &List{}
	start := p.start
	for {
		for r := p.next(); r != eof; {
			if r == '\n' {
				r1 := p.peek()
				if r1 == '\n' ||
					r1 == eof ||
					regexp.MustCompile("^"+escape+string(marker)+"   ").Match(p.src[p.cur:]) {
					break
				}
				if p.forsee(' ', ' ', ' ', ' ') {
					p.src = append(p.src[:p.cur], p.src[p.cur+4:]...)
					continue
				}
				if p.forsee('\t') {
					p.src = append(p.src[:p.cur], p.src[p.cur+1:]...)
					continue
				}
			}
			r = p.next()
		}
		content := p.src[start:p.cur]
		content = regexp.MustCompile("^"+escape+string(marker)+"   ").ReplaceAll(content, []byte{})
		content = bytes.TrimRightFunc(content, func(r rune) bool {
			if r == '\n' {
				return true
			}
			return false
		})

		item := &Item{content: content}
		list.items = append(list.items, item)
		// if forsee Sprinf("%s ",marker),
		// parse another list item,else emit list.
		if !p.forsee(marker, ' ', ' ', ' ') {
			p.emit(list)
			return parseBegin
		}
		start = p.cur
	}
}

// parseCode parses code beginning with 4 sapces or 1 tab.
func parseCodeBlock(p *blockParser) stateFn {
	codeBlock := &CodeBlock{}
	start := p.start
	var marker string
	r := p.peek()
	marker = "\t"
	if r == ' ' {
		marker = "    " // 4 spaces
	}
	for {
		for r := p.next(); r != '\n' && r != eof; {
			r = p.next()
		}
		content := p.src[start:p.cur]
		content = regexp.MustCompile("^"+marker).ReplaceAll(content, []byte{})
		content = bytes.TrimRightFunc(content, func(r rune) bool {
			if r == '\n' {
				return true
			}
			return false
		})
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
func parseRule(p *blockParser) stateFn {
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

// parseQuote is the parser for state of quote.
func parseQuote(p *blockParser) stateFn {
	var (
		r       rune
		content []byte
	)
	for {
		for r = p.next(); r != '\n' && r != eof; r = p.next() {
		}

		r1 := p.peek()
		if r == '\n' && (r1 == '\n' || r1 == eof) {
			break
		}
		if r == eof {
			break
		}
	}
	if r == '\n' {
		content = p.src[p.start : p.cur-1]
		p.next()
	}
	if r == eof {
		content = p.src[p.start:p.cur]
	}
	lines := bytes.Split(content, []byte{'\n'})
	for i := 0; i < len(lines); i++ {
		// remove heading '>' and space.
		if len(lines[i]) > 0 && lines[i][0] == '>' {
			lines[i] = lines[i][1:]
			if len(lines[i]) > 0 && lines[i][0] == ' ' {
				lines[i] = lines[i][1:]
			}
		}
	}
	content = bytes.Join(lines, []byte{'\n'})
	np := newParser(content)
	quote := &QuoteBlock{}
	for b := np.element(); b != nil; b = np.element() {
		quote.subBlocks = append(quote.subBlocks, b)
	}
	p.emit(quote)
	return parseBegin
}

// parseError is error handler when account for errors.
func parseError(p *blockParser) stateFn {
	return nil
}

// block main parsing.
func parseBegin(p *blockParser) stateFn {
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
			if p.forsee(r, ' ', ' ', ' ') {
				return parseUnorderList
			}
		}
		if r1 == r && p.forsee(r, r, r) {
			return parseRule
		}
		return parseParagraph
	case r == '_':
		r1 := p.peek(2)
		if r1 == r && (p.forsee(r, r, r) || p.forsee(r, ' ', r, ' ', r)) {
			return parseRule
		}
		return parseParagraph
	case unicode.IsDigit(r):
		if regexp.MustCompile(`\d+\.  `).Match(p.src[p.cur:]) {
			return parseOrderList
		}
		return parseParagraph
	case r == '>':
		return parseQuote
	case r == '\t' || (r == ' ' && p.forsee(' ', ' ', ' ')):
		return parseCodeBlock
	case r == '\n':
		p.next()
		p.ignore()
		return parseBegin
	case r == eof:
		return nil
	default:
		return parseParagraph
	}
}

// -----------span parsing----------

// parse emphasis or strong span
func parseEmphasis(p *spanParser) spanStateFn {
	start := p.cur
	marker := p.peek()
	n := p.consume(marker, 2)
	t := kind.Strong
	if n == 1 {
		t = kind.Emphasis
	}

	for r := p.next(); r != marker && r != '\n' && r != eof; {
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

// cut content with left and right wrapped from the src.
func cut(left, right byte, begin int, src []byte) (remain []byte, cut []byte) {
	start := bytes.IndexByte(src[begin:], left)
	end := start + 1 + bytes.IndexByte(src[begin+start+1:], right)
	cut = make([]byte, end-start-1)
	copy(cut, src[begin+start+1:begin+end])
	remain = append(src[:begin+start], src[begin+end+1:]...)
	return
}

// parseRef parses links and images including references referring to the previous links or images.
func parseRef(p *spanParser) spanStateFn {
	var (
		text  []byte
		title []byte
		ref   []byte
		id    []byte
		k     kind.Kind
	)

	r := p.peek()
	if r == '[' {
		k = kind.Link
		i1 := bytes.IndexByte(p.src[p.cur:], ']')
		i2 := bytes.IndexByte(p.src[p.cur:], ':')
		// reference.
		if i1+1 == i2 {
			p.src, id = cut('[', ']', p.cur, p.src)
			p.src, ref = cut(':', '\n', p.cur, p.src)
			ref, title = cut('"', '"', 0, ref)
			ref = bytes.TrimSpace(ref)
			p.ref[string(id)] = &reference{ref, title}
			return parseSpan
		}
	}
	if r == '!' && p.peek(2) == '[' {
		p.src = append(p.src[:p.cur], p.src[p.cur+1:]...)
		k = kind.Image
	}

	p.src, text = cut('[', ']', p.cur, p.src)
	r = p.peek()

	if r == '[' {
		p.src, id = cut('[', ']', p.cur, p.src)
	}

	if r == '(' {
		p.src, ref = cut('(', ')', p.cur, p.src)
		ref, title = cut('"', '"', 0, ref)
		ref = bytes.TrimSpace(ref)
	}
	if k == kind.Image {
		p.emit(&Image{p.cur, id, text, title, ref})
	}

	if k == kind.Link {
		p.emit(&Link{p.cur, id, text, title, ref})
	}

	return parseSpan
}

// parseCode parses span code.
func parseCode(p *spanParser) spanStateFn {
	indexOfNewLine := bytes.IndexByte(p.src[p.cur:], '\n')
	if indexOfNewLine == -1 {
		indexOfNewLine = len(p.src[p.cur:])
	}
	indexOfBacktick := bytes.LastIndex(p.src[p.cur:p.cur+indexOfNewLine], []byte{'`'})

	var content = make([]byte, indexOfBacktick+1)
	copy(content, p.src[p.cur:p.cur+indexOfBacktick+1])
	content = content[1 : len(content)-1]
	p.src = append(p.src[:p.cur], p.src[indexOfBacktick+1:]...)
	p.emit(&Code{p.cur, content})
	return parseSpan
}

var escapeRunes = "\\'*_{}[]()#+-.!"

// isEscapeRune returns true if r needs escaping.
func isEscapeRune(r rune) bool {
	for _, v := range escapeRunes {
		if rune(v) == r {
			return true
		}
	}
	return false
}

func findRune(src []byte, r byte) int {
	i := bytes.IndexByte(src, r)
	if i == -1 {
		i = 1 << 32 // big enough as eof
	}
	return i

}

// span main parsing.
func parseSpan(p *spanParser) spanStateFn {
	for {
		switch r := p.peek(); {
		case r == '\\':
			r1 := p.peek(2)
			// escape runes
			if isEscapeRune(r1) {
				p.next()
				p.merge()
				p.next()
			}
			return parseSpan
		case r == '`':
			return parseCode
		case r == '!' || r == '[':
			return parseRef
		case r == '*' || r == '_':
			if findRune(p.src[p.cur+1:], byte(r)) < findRune(p.src[p.cur+1:], '\n') {
				return parseEmphasis
			}
			p.next()
			p.ignore()
		case r == eof:
			return nil
		default:
			p.next()
			p.ignore()
		}
	}
}

// Parse parses src with ext as extension,and returns pure text content.
func Parse(src []byte, ext EXT) []byte {
	if ext == BASIC {
	}
	p := newParser(src)
	var contents [][]byte
	for block := p.element(); block != nil; block = p.element() {
		contents = append(contents, block.Content())
	}
	return bytes.Join(contents, []byte("\n"))
}
