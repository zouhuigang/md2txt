package kind

type Kind int
type ElementType int

// specific types
const (
	Head Kind = iota
	Paragraph
	List
	CodeBlock
	Code
	Rule

	Emphasis
	Strong
)

// element types
const (
	Block = iota
	Inline
)
