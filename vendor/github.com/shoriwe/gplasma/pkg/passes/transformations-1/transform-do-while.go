package transformations_1

import (
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
)

func (transform *transformPass) DoWhile(doWhile *ast2.DoWhile) []ast3.Node {
	startLabel := transform.nextLabel()
	endLabel := transform.nextLabel()
	condition := &ast3.IfJump{
		Condition: transform.Expression(doWhile.Condition),
		Target:    startLabel,
	}
	body := make([]ast3.Node, 0, len(doWhile.Body))
	for _, node := range doWhile.Body {
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
	result := make([]ast3.Node, 0, 3+len(body))
	result = append(result, startLabel)
	result = append(result, body...)
	result = append(result, condition, endLabel)
	return result
}
