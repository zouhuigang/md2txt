/*
Type definitions for markdown elements.
*/
package kind

type Kind int
type ElementType int

// specific types
const (
	// block types
	Head Kind = iota
	Paragraph
	List
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
	Block = iota
	Inline
)
