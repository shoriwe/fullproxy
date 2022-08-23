package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simplify *simplifyPass) Selector(selector *ast.SelectorExpression) *ast2.Selector {
	return &ast2.Selector{
		X:          simplify.Expression(selector.X),
		Identifier: simplify.Identifier(selector.Identifier),
	}
}
