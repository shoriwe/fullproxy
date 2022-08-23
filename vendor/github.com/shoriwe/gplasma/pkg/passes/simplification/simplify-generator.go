package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simplify *simplifyPass) GeneratorExpr(generator *ast.GeneratorExpression) *ast2.Generator {
	receivers := make([]*ast2.Identifier, 0, len(generator.Receivers))
	for _, receiver := range generator.Receivers {
		receivers = append(receivers, simplify.Identifier(receiver))
	}
	return &ast2.Generator{
		Operation: simplify.Expression(generator.Operation),
		Receivers: receivers,
		Source:    simplify.Expression(generator.Source),
	}
}

func (simplify *simplifyPass) GeneratorDef(generator *ast.GeneratorDefinitionStatement) *ast2.GeneratorDefinition {
	arguments := make([]*ast2.Identifier, 0, len(generator.Arguments))
	for _, argument := range generator.Arguments {
		arguments = append(arguments, simplify.Identifier(argument))
	}
	body := make([]ast2.Node, 0, len(generator.Body))
	for _, node := range generator.Body {
		body = append(body, simplify.Node(node))
	}
	return &ast2.GeneratorDefinition{
		Name:      simplify.Identifier(generator.Name),
		Arguments: arguments,
		Body:      body,
	}
}
