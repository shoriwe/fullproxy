package transformations_1

import (
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
)

func (transform *transformPass) IfOneLiner(iol *ast2.IfOneLiner) *ast3.Call {
	end := transform.nextLabel()
	elseLabel := transform.nextLabel()
	condition := &ast3.IfJump{
		Condition: transform.Expression(&ast2.Unary{
			Operator: ast2.Not,
			X:        iol.Condition,
		}),
		Target: elseLabel,
	}
	result := transform.Expression(iol.Result)
	endJump := &ast3.Jump{
		Target: end,
	}
	else_ := transform.Expression(iol.Else)
	var body []ast3.Node
	body = append(body, condition)
	body = append(body, result)
	body = append(body, endJump)
	body = append(body, elseLabel)
	body = append(body, else_)
	body = append(body, end)
	return &ast3.Call{
		Function: &ast3.Function{
			Body: body,
		},
	}
}
