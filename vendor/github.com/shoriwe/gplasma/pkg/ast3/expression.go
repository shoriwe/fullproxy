package ast3

type (
	Expression interface {
		Node
		E3()
	}
	Assignable interface {
		Expression
		A2()
	}
	Function struct {
		Expression
		Arguments []*Identifier
		Body      []Node
	}
	Class struct {
		Expression
		Bases []Expression
		Body  []Node
	}
	Call struct {
		Expression
		Function  Expression
		Arguments []Expression
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

	Selector struct {
		Assignable
		X          Expression
		Identifier *Identifier
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
