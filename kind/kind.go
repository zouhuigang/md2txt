/*
Type definitions for markdown elements.
*/
package kind

//go:generate stringer -type=Kind
type Kind int

//go:generate stringer -type=ElementType
type ElementType int

// specific types
const (
	// block types
	Head Kind = iota
	Paragraph
	List
	QuoteBlock
	CodeBlock
	Rule
	// inline types
	Emphasis
	Strong
	Link
	Code
	Image
)

// element types
const (
	Block ElementType = iota
	Inline
)
