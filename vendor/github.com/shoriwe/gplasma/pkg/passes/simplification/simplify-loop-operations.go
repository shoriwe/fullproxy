package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simplify *simplifyPass) Continue(c *ast.ContinueStatement) *ast2.Continue {
	return &ast2.Continue{}
}

func (simplify *simplifyPass) Break(c *ast.BreakStatement) *ast2.Break {
	return &ast2.Break{}
}
