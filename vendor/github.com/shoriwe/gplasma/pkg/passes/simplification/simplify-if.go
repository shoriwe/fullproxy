package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simplify *simplifyPass) If(if_ *ast.IfStatement) *ast2.If {
	body := make([]ast2.Node, 0, len(if_.Body))
	for _, node := range if_.Body {
		body = append(body, simplify.Node(node))
	}
	root := &ast2.If{
		Condition: simplify.Expression(if_.Condition),
		Body:      body,
		Else:      nil,
	}
	lastIf := root
	for _, elif := range if_.ElifBlocks {
		elifBody := make([]ast2.Node, 0, len(elif.Body))
		for _, node := range elif.Body {
			elifBody = append(elifBody, simplify.Node(node))
		}
		newLastIf := &ast2.If{
			Condition: simplify.Expression(elif.Condition),
			Body:      elifBody,
			Else:      nil,
		}
		lastIf.Else = []ast2.Node{newLastIf}
		lastIf = newLastIf
	}
	lastIf.Else = make([]ast2.Node, 0, len(if_.Else))
	for _, node := range if_.Else {
		lastIf.Else = append(lastIf.Else, simplify.Node(node))
	}
	return root
}
