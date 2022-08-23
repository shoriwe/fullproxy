package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simplify *simplifyPass) Delete(del *ast.DeleteStatement) *ast2.Delete {
	var x ast2.Assignable
	switch dx := del.X.(type) {
	case *ast.Identifier:
		x = simplify.Identifier(dx)
	case *ast.IndexExpression:
		x = simplify.Index(dx)
	case *ast.SelectorExpression:
		x = simplify.Selector(dx)
	default:
		panic("unknown selector type")
	}
	return &ast2.Delete{
		X: x,
	}
}
