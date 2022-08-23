package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simplify *simplifyPass) Until(until *ast.UntilLoopStatement) *ast2.While {
	var (
		body                      = make([]ast2.Node, 0, len(until.Body))
		condition ast2.Expression = &ast2.Unary{
			Operator: ast2.Not,
			X:        simplify.Expression(until.Condition),
		}
	)
	for _, node := range until.Body {
		body = append(body, simplify.Node(node))
	}
	return &ast2.While{
		Body:      body,
		Condition: condition,
	}
}
