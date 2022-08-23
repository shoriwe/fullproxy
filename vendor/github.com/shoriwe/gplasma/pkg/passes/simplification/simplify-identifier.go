package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simplify *simplifyPass) Identifier(ident *ast.Identifier) *ast2.Identifier {
	return &ast2.Identifier{
		Symbol: ident.Token.String(),
	}
}
