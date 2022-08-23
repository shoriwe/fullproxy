package ast

type (
	Node interface {
		N()
	}

	Program struct {
		Node
		Begin *BeginStatement
		End   *EndStatement
		Body  []Node
	}
)
