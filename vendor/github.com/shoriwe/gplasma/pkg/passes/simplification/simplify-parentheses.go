package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simplify *simplifyPass) Parentheses(expression *ast.ParenthesesExpression) ast2.Expression {
	return simplify.Expression(expression.X)
}
