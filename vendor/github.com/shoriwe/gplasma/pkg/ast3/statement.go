package ast3

type (
	Statement interface {
		Node
		S3()
	}
	Assignment struct {
		Statement
		Left  Assignable
		Right Expression
	}
	Label struct {
		Statement
		Code int
	}
	Jump struct {
		Statement
		Target *Label
	}
	ContinueJump Jump
	BreakJump    Jump
	IfJump       struct {
		Statement
		Condition Expression
		Target    *Label
	}
	Return struct {
		Statement
		Result Expression
	}
	Yield Return

	Delete struct {
		Statement
		X Assignable
	}
	Defer struct {
		Statement
		X Expression
	}
)
