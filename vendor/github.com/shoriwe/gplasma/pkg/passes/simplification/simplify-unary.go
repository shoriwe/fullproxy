package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/lexer"
)

func (simplify *simplifyPass) Unary(unary *ast.UnaryExpression) *ast2.Unary {
	var operator ast2.UnaryOperator
	switch unary.Operator.DirectValue {
	case lexer.Not, lexer.SignNot:
		operator = ast2.Not
	case lexer.NegateBits:
		operator = ast2.NegateBits
	case lexer.Add:
		operator = ast2.Positive
	case lexer.Sub:
		operator = ast2.Negative
	}
	return &ast2.Unary{
		Operator: operator,
		X:        simplify.Expression(unary.X),
	}
}
