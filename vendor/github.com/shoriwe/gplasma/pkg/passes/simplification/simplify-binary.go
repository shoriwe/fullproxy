package simplification

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/lexer"
)

func (simplify *simplifyPass) Binary(binary *ast.BinaryExpression) *ast2.Binary {
	var operator ast2.BinaryOperator
	switch binary.Operator.DirectValue {
	case lexer.And:
		operator = ast2.And
	case lexer.Or:
		operator = ast2.Or
	case lexer.Xor:
		operator = ast2.Xor
	case lexer.In:
		operator = ast2.In
	case lexer.Is:
		operator = ast2.Is
	case lexer.Implements:
		operator = ast2.Implements
	case lexer.Equals:
		operator = ast2.Equals
	case lexer.NotEqual:
		operator = ast2.NotEqual
	case lexer.GreaterThan:
		operator = ast2.GreaterThan
	case lexer.GreaterOrEqualThan:
		operator = ast2.GreaterOrEqualThan
	case lexer.LessThan:
		operator = ast2.LessThan
	case lexer.LessOrEqualThan:
		operator = ast2.LessOrEqualThan
	case lexer.BitwiseOr:
		operator = ast2.BitwiseOr
	case lexer.BitwiseXor:
		operator = ast2.BitwiseXor
	case lexer.BitwiseAnd:
		operator = ast2.BitwiseAnd
	case lexer.BitwiseLeft:
		operator = ast2.BitwiseLeft
	case lexer.BitwiseRight:
		operator = ast2.BitwiseRight
	case lexer.Add:
		operator = ast2.Add
	case lexer.Sub:
		operator = ast2.Sub
	case lexer.Star:
		operator = ast2.Mul
	case lexer.Div:
		operator = ast2.Div
	case lexer.FloorDiv:
		operator = ast2.FloorDiv
	case lexer.Modulus:
		operator = ast2.Modulus
	case lexer.PowerOf:
		operator = ast2.PowerOf
	default:
		panic(fmt.Sprintf("unknown binary operator %d", binary.Operator.DirectValue))
	}
	return &ast2.Binary{
		Left:     simplify.Expression(binary.LeftHandSide),
		Right:    simplify.Expression(binary.RightHandSide),
		Operator: operator,
	}
}
