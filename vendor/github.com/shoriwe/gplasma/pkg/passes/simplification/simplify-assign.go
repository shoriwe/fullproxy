package simplification

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/lexer"
)

func (simplify *simplifyPass) Assign(assign *ast.AssignStatement) *ast2.Assignment {
	var (
		left  ast2.Assignable
		right = simplify.Expression(assign.RightHandSide)
	)
	switch l := assign.LeftHandSide.(type) {
	case *ast.Identifier:
		left = simplify.Identifier(l)
	case *ast.SelectorExpression:
		left = simplify.Selector(l)
	case *ast.IndexExpression:
		left = simplify.Index(l)
	default:
		panic("invalid identifier left hand side type")
	}
	switch assign.AssignOperator.DirectValue {
	case lexer.Assign:
		break
	case lexer.BitwiseOrAssign:
		right = &ast2.Binary{
			Left:     left,
			Right:    right,
			Operator: ast2.BitwiseOr,
		}
	case lexer.BitwiseXorAssign:
		right = &ast2.Binary{
			Left:     left,
			Right:    right,
			Operator: ast2.BitwiseXor,
		}
	case lexer.BitwiseAndAssign:
		right = &ast2.Binary{
			Left:     left,
			Right:    right,
			Operator: ast2.BitwiseAnd,
		}
	case lexer.BitwiseLeftAssign:
		right = &ast2.Binary{
			Left:     left,
			Right:    right,
			Operator: ast2.BitwiseLeft,
		}
	case lexer.BitwiseRightAssign:
		right = &ast2.Binary{
			Left:     left,
			Right:    right,
			Operator: ast2.BitwiseRight,
		}
	case lexer.AddAssign:
		right = &ast2.Binary{
			Left:     left,
			Right:    right,
			Operator: ast2.Add,
		}
	case lexer.SubAssign:
		right = &ast2.Binary{
			Left:     left,
			Right:    right,
			Operator: ast2.Sub,
		}
	case lexer.StarAssign:
		right = &ast2.Binary{
			Left:     left,
			Right:    right,
			Operator: ast2.Mul,
		}
	case lexer.DivAssign:
		right = &ast2.Binary{
			Left:     left,
			Right:    right,
			Operator: ast2.Div,
		}
	case lexer.FloorDivAssign:
		right = &ast2.Binary{
			Left:     left,
			Right:    right,
			Operator: ast2.FloorDiv,
		}
	case lexer.ModulusAssign:
		right = &ast2.Binary{
			Left:     left,
			Right:    right,
			Operator: ast2.Modulus,
		}
	case lexer.PowerOfAssign:
		right = &ast2.Binary{
			Left:     left,
			Right:    right,
			Operator: ast2.PowerOf,
		}
	default:
		panic(fmt.Sprintf("unknown binary operator for assignment %d", assign.AssignOperator.DirectValue))
	}
	return &ast2.Assignment{
		Left:  left,
		Right: right,
	}
}
