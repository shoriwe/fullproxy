package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simplify *simplifyPass) Defer(d *ast.DeferStatement) *ast2.Defer {
	return &ast2.Defer{
		X: simplify.Expression(d.X),
	}
}
