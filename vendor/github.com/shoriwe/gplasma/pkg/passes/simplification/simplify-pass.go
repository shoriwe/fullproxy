package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simplify *simplifyPass) Pass(pass *ast.PassStatement) *ast2.Pass {
	return &ast2.Pass{}
}
