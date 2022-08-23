package transformations_1

import (
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
)

func (transform *transformPass) Function(function *ast2.FunctionDefinition) []ast3.Node {
	arguments := make([]*ast3.Identifier, 0, len(function.Arguments))
	for _, argument := range function.Arguments {
		arguments = append(arguments, transform.Identifier(argument))
	}
	body := make([]ast3.Node, 0, len(function.Body))
	for _, node := range function.Body {
		body = append(body, transform.Node(node)...)
	}
	return []ast3.Node{&ast3.Assignment{
		Left: transform.Identifier(function.Name),
		Right: &ast3.Function{
			Arguments: arguments,
			Body:      body,
		},
	}}
}
