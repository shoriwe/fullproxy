package lexer

/*
	Token Kinds
*/

type Kind uint8

const (
	Unknown Kind = iota
	Comment
	Whitespace
	Literal
	IdentifierKind
	JunkKind
	Separator
	Punctuation
	Assignment
	Comparator
	Operator
	Keyword
	Boolean
	NoneType
	EOF
)
