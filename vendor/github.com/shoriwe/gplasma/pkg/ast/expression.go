package ast

import (
	lexer2 "github.com/shoriwe/gplasma/pkg/lexer"
)

type (
	Expression interface {
		Node
		E()
	}

	ArrayExpression struct {
		Expression
		Values []Expression
	}

	TupleExpression struct {
		Expression
		Values []Expression
	}

	KeyValue struct {
		Key   Expression
		Value Expression
	}

	HashExpression struct {
		Expression
		Values []*KeyValue
	}

	Identifier struct {
		Expression
		Token *lexer2.Token
	}

	BasicLiteralExpression struct {
		Expression
		Token       *lexer2.Token
		Kind        lexer2.Kind
		DirectValue lexer2.DirectValue
	}

	BinaryExpression struct {
		Expression
		LeftHandSide  Expression
		Operator      *lexer2.Token
		RightHandSide Expression
	}

	UnaryExpression struct {
		Expression
		Operator *lexer2.Token
		X        Expression
	}

	ParenthesesExpression struct {
		Expression
		X Expression
	}

	LambdaExpression struct {
		Expression
		Arguments []*Identifier
		Code      Expression
	}

	GeneratorExpression struct {
		Expression
		Operation Expression
		Receivers []*Identifier
		Source    Expression
	}

	SelectorExpression struct {
		Expression
		X          Expression
		Identifier *Identifier
	}

	MethodInvocationExpression struct {
		Expression
		Function  Expression
		Arguments []Expression
	}

	IndexExpression struct {
		Expression
		Source Expression
		Index  Expression
	}

	IfOneLinerExpression struct {
		Expression
		Result     Expression
		Condition  Expression
		ElseResult Expression
	}

	UnlessOneLinerExpression struct {
		Expression
		Result     Expression
		Condition  Expression
		ElseResult Expression
	}

	SuperExpression struct {
		Expression
		X Expression
	}
)
