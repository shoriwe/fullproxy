package transformations_1

import (
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
)

func (transform *transformPass) Lambda(lambda *ast2.Lambda) *ast3.Function {
	arguments := make([]*ast3.Identifier, 0, len(lambda.Arguments))
	for _, argument := range lambda.Arguments {
		arguments = append(arguments, transform.Identifier(argument))
	}
	return &ast3.Function{
		Expression: nil,
		Arguments:  arguments,
		Body: []ast3.Node{
			&ast3.Return{
				Statement: nil,
				Result:    transform.Expression(lambda.Result),
			},
		},
	}
}
