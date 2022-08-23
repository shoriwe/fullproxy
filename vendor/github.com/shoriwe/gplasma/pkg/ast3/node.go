package ast3

type (
	Node interface {
		N3()
	}
	Program []Node
)
