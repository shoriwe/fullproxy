package transformations_1

import (
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
)

func (transform *transformPass) While(while *ast2.While) []ast3.Node {
	setup := make([]ast3.Node, 0, len(while.Setup))
	for _, node := range while.Setup {
		setup = append(setup, transform.Node(node)...)
	}
	startLabel := transform.nextLabel()
	endLabel := transform.nextLabel()
	condition := &ast3.IfJump{
		Condition: transform.Expression(&ast2.Unary{
			Operator: ast2.Not,
			X:        while.Condition,
		}),
		Target: endLabel,
	}
	body := make([]ast3.Node, 0, len(while.Body))
	for _, node := range while.Body {
		body = append(body, transform.Node(node)...)
	}
	for _, node := range body {
		switch n := node.(type) {
		case *ast3.ContinueJump:
			if n.Target == nil {
				n.Target = startLabel
			}
		case *ast3.BreakJump:
			if n.Target == nil {
				n.Target = endLabel
			}
		}
	}
	result := make([]ast3.Node, 0, 3+len(setup)+len(body))
	result = append(result, setup...)
	result = append(result, startLabel, condition)
	result = append(result, body...)
	result = append(result, &ast3.Jump{
		Target: startLabel,
	})
	result = append(result, endLabel)
	return result
}
