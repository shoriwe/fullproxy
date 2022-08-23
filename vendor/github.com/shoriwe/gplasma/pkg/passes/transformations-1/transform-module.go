package transformations_1

import (
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
	magic_functions "github.com/shoriwe/gplasma/pkg/common/magic-functions"
)

func (transform *transformPass) Module(module *ast2.Module) []ast3.Node {
	body := make([]ast3.Node, 0, len(module.Body))
	for _, node := range module.Body {
		body = append(body, transform.Node(node)...)
	}
	body = append(body, &ast3.Assignment{
		Left: &ast3.Identifier{
			Symbol: magic_functions.Init,
		},
		Right: &ast3.Function{},
	})
	moduleAssignment := &ast3.Assignment{
		Left: transform.Identifier(module.Name),
		Right: &ast3.Call{
			Function: &ast3.Class{
				Body: body,
			},
		},
	}
	return []ast3.Node{moduleAssignment}
}
