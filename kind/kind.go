package kind

type Kind int
type ElementType int

const (
	Head Kind = iota
	Paragraph
	List
	CodeBlock
	Code
	Rule
	EOF
)
const (
	Block = iota
	Inline
)
