package transformations_1

import (
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
)

func (transform *transformPass) Call(call *ast2.FunctionCall) *ast3.Call {
	arguments := make([]ast3.Expression, 0, len(call.Arguments))
	for _, argument := range call.Arguments {
		arguments = append(arguments, transform.Expression(argument))
	}
	return &ast3.Call{
		Function:  transform.Expression(call.Function),
		Arguments: arguments,
	}
}
