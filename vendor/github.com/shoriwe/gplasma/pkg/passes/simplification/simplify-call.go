package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simplify *simplifyPass) Call(call *ast.MethodInvocationExpression) *ast2.FunctionCall {
	arguments := make([]ast2.Expression, 0, len(call.Arguments))
	for _, argument := range call.Arguments {
		arguments = append(arguments, simplify.Expression(argument))
	}
	return &ast2.FunctionCall{
		Function:  simplify.Expression(call.Function),
		Arguments: arguments,
	}
}
