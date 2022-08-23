package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simplify *simplifyPass) Super(super *ast.SuperExpression) *ast2.Super {
	return &ast2.Super{
		X: simplify.Expression(super.X),
	}
}
