package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simplify *simplifyPass) Function(f *ast.FunctionDefinitionStatement) *ast2.FunctionDefinition {
	arguments := make([]*ast2.Identifier, 0, len(f.Arguments))
	for _, argument := range f.Arguments {
		arguments = append(arguments, simplify.Identifier(argument))
	}
	body := make([]ast2.Node, 0, len(f.Body))
	for _, node := range f.Body {
		body = append(body, simplify.Node(node))
	}
	return &ast2.FunctionDefinition{
		Name:      simplify.Identifier(f.Name),
		Arguments: arguments,
		Body:      body,
	}
}
