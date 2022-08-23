package ast2

const (
	Not UnaryOperator = iota
	Positive
	Negative
	NegateBits

	And BinaryOperator = iota
	Or
	Xor
	In
	Is
	Implements
	Equals
	NotEqual
	GreaterThan
	GreaterOrEqualThan
	LessThan
	LessOrEqualThan
	BitwiseOr
	BitwiseXor
	BitwiseAnd
	BitwiseLeft
	BitwiseRight
	Add
	Sub
	Mul
	Div
	FloorDiv
	Modulus
	PowerOf
)

type (
	BinaryOperator int
	UnaryOperator  int
	Expression     interface {
		Node
		E2()
	}

	Assignable interface {
		Expression
		A()
	}

	Binary struct {
		Expression
		Left, Right Expression
		Operator    BinaryOperator
	}

	Unary struct {
		Expression
		Operator UnaryOperator
		X        Expression
	}

	IfOneLiner struct {
		Expression
		Condition, Result, Else Expression
	}

	Array struct {
		Expression
		Values []Expression
	}

	Tuple struct {
		Expression
		Values []Expression
	}

	KeyValue struct {
		Key, Value Expression
	}

	Hash struct {
		Expression
		Values []*KeyValue
	}

	Identifier struct {
		Assignable
		Symbol string
	}

	Integer struct {
		Expression
		Value int64
	}

	Float struct {
		Expression
		Value float64
	}

	String struct {
		Expression
		Contents []byte
	}

	Bytes struct {
		Expression
		Contents []byte
	}

	True struct {
		Expression
	}

	False struct {
		Expression
	}

	None struct {
		Expression
	}

	Lambda struct {
		Expression
		Arguments []*Identifier
		Result    Expression
	}

	Generator struct {
		Expression
		Operation Expression
		Receivers []*Identifier
		Source    Expression
	}

	Selector struct {
		Assignable
		X          Expression
		Identifier *Identifier
	}

	FunctionCall struct {
		Expression
		Function  Expression
		Arguments []Expression
	}

	Index struct {
		Assignable
		Source Expression
		Index  Expression
	}

	Super struct {
		Expression
		X Expression
	}
)
