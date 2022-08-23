package transformations_1

import (
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
)

func (transform *transformPass) If(if_ *ast2.If) []ast3.Node {
	elseLabel := transform.nextLabel()
	endLabel := transform.nextLabel()
	condition := &ast3.IfJump{
		Condition: transform.Expression(&ast2.Unary{
			Operator: ast2.Not,
			X:        if_.Condition,
		}),
		Target: elseLabel,
	}
	// Switch setup
	var switchSetup []ast3.Node
	if if_.SwitchSetup != nil {
		switchSetup = transform.Assignment(if_.SwitchSetup)
	}
	// If
	ifBody := make([]ast3.Node, 0, 1+len(if_.Body))
	for _, node := range if_.Body {
		ifBody = append(ifBody, transform.Node(node)...)
	}
	ifBody = append(ifBody, &ast3.Jump{Target: endLabel})
	// Else
	elseBody := make([]ast3.Node, 0, 2+len(if_.Else))
	elseBody = append(elseBody, elseLabel)
	for _, node := range if_.Else {
		elseBody = append(elseBody, transform.Node(node)...)
	}
	elseBody = append(elseBody, &ast3.Jump{Target: endLabel})
	//
	result := make([]ast3.Node, 0, 3+len(ifBody)+len(elseBody))
	if switchSetup != nil {
		result = append(result, switchSetup...)
	}
	result = append(result, condition)
	result = append(result, ifBody...)
	result = append(result, elseBody...)
	result = append(result, endLabel)
	return result
}
