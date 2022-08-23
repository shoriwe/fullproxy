package ast

import (
	"github.com/shoriwe/gplasma/pkg/lexer"
)

type (
	Statement interface {
		S()
		Node
	}

	AssignStatement struct {
		Statement
		LeftHandSide   Expression // Identifiers or Selectors
		AssignOperator *lexer.Token
		RightHandSide  Expression
	}

	DoWhileStatement struct {
		Statement
		Condition Expression
		Body      []Node
	}

	WhileLoopStatement struct {
		Statement
		Condition Expression
		Body      []Node
	}

	UntilLoopStatement struct {
		Statement
		Condition Expression
		Body      []Node
	}

	ForLoopStatement struct {
		Statement
		Receivers []*Identifier
		Source    Expression
		Body      []Node
	}

	ElifBlock struct {
		Condition Expression
		Body      []Node
	}

	IfStatement struct {
		Statement
		Condition  Expression
		Body       []Node
		ElifBlocks []ElifBlock
		Else       []Node
	}

	UnlessStatement struct {
		Statement
		Condition  Expression
		Body       []Node
		ElifBlocks []ElifBlock
		Else       []Node
	}

	CaseBlock struct {
		Cases []Expression
		Body  []Node
	}

	SwitchStatement struct {
		Statement
		Target     Expression
		CaseBlocks []*CaseBlock
		Default    []Node
	}

	ModuleStatement struct {
		Statement
		Name *Identifier
		Body []Node
	}

	FunctionDefinitionStatement struct {
		Statement
		Name      *Identifier
		Arguments []*Identifier
		Body      []Node
	}

	GeneratorDefinitionStatement struct {
		Statement
		Name      *Identifier
		Arguments []*Identifier
		Body      []Node
	}

	InterfaceStatement struct {
		Statement
		Name              *Identifier
		Bases             []Expression
		MethodDefinitions []*FunctionDefinitionStatement
	}

	ClassStatement struct {
		Statement
		Name  *Identifier
		Bases []Expression // Identifiers and selectors
		Body  []Node
	}

	BeginStatement struct {
		Statement
		Body []Node
	}

	EndStatement struct {
		Statement
		Body []Node
	}

	ReturnStatement struct {
		Statement
		Results []Expression
	}

	YieldStatement struct {
		Statement
		Results []Expression
	}

	ContinueStatement struct {
		Statement
	}

	BreakStatement struct {
		Statement
	}

	PassStatement struct {
		Statement
	}

	DeleteStatement struct {
		Statement
		X Expression
	}

	DeferStatement struct {
		Statement
		X *MethodInvocationExpression
	}
)
