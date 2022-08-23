package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simplify *simplifyPass) Module(module *ast.ModuleStatement) *ast2.Module {
	body := make([]ast2.Node, 0, len(module.Body))
	for _, node := range module.Body {
		body = append(body, simplify.Node(node))
	}
	return &ast2.Module{
		Name: simplify.Identifier(module.Name),
		Body: body,
	}
}
