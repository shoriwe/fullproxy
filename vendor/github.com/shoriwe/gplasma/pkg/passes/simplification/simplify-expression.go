package simplification

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
	"reflect"
)

func (simplify *simplifyPass) Expression(expr ast.Expression) ast2.Expression {
	if expr == nil {
		return &ast2.None{}
	}
	switch e := expr.(type) {
	case *ast.ArrayExpression:
		return simplify.Array(e)
	case *ast.TupleExpression:
		return simplify.Tuple(e)
	case *ast.HashExpression:
		return simplify.Hash(e)
	case *ast.Identifier:
		return simplify.Identifier(e)
	case *ast.BasicLiteralExpression:
		return simplify.Literal(e)
	case *ast.BinaryExpression:
		return simplify.Binary(e)
	case *ast.UnaryExpression:
		return simplify.Unary(e)
	case *ast.ParenthesesExpression:
		return simplify.Parentheses(e)
	case *ast.LambdaExpression:
		return simplify.Lambda(e)
	case *ast.GeneratorExpression:
		return simplify.GeneratorExpr(e)
	case *ast.SelectorExpression:
		return simplify.Selector(e)
	case *ast.MethodInvocationExpression:
		return simplify.Call(e)
	case *ast.IndexExpression:
		return simplify.Index(e)
	case *ast.IfOneLinerExpression:
		return simplify.IfOneLiner(e)
	case *ast.UnlessOneLinerExpression:
		return simplify.UnlessOneLiner(e)
	case *ast.SuperExpression:
		return simplify.Super(e)
	default:
		panic(fmt.Sprintf("unknown expression type %s", reflect.TypeOf(expr).String()))
	}
}
