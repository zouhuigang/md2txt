package kind

type Kind int

const (
	Head Kind = iota
	Paragraph
	List
	CodeBlock
	Code
	EOF
)
