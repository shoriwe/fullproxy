package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simplify *simplifyPass) IfOneLiner(if_ *ast.IfOneLinerExpression) *ast2.IfOneLiner {
	return &ast2.IfOneLiner{
		Condition: simplify.Expression(if_.Condition),
		Result:    simplify.Expression(if_.Result),
		Else:      simplify.Expression(if_.ElseResult),
	}
}

func (simplify *simplifyPass) UnlessOneLiner(unless *ast.UnlessOneLinerExpression) *ast2.IfOneLiner {
	return &ast2.IfOneLiner{
		Condition: &ast2.Unary{
			Operator: ast2.Not,
			X:        simplify.Expression(unless.Condition),
		},
		Result: simplify.Expression(unless.Result),
		Else:   simplify.Expression(unless.ElseResult),
	}
}
