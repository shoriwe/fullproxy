package simplification

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simplify *simplifyPass) nextAnonIdentifier() *ast2.Identifier {
	ident := simplify.currentAnonIdent
	simplify.currentAnonIdent++
	return &ast2.Identifier{
		Symbol: fmt.Sprintf("____simplify_%d", ident),
	}
}
