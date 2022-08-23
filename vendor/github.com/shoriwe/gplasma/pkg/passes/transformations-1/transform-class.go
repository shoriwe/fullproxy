package transformations_1

import (
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
)

func (transform *transformPass) Class(class *ast2.Class) []ast3.Node {
	body := make([]ast3.Node, len(class.Body))
	for _, node := range class.Body {
		body = append(body, transform.Node(node)...)
	}
	bases := make([]ast3.Expression, 0, len(class.Bases))
	for _, base := range class.Bases {
		bases = append(bases, transform.Expression(base))
	}
	return []ast3.Node{
		&ast3.Assignment{
			Left: transform.Identifier(class.Name),
			Right: &ast3.Class{
				Bases: bases,
				Body:  body,
			},
		},
	}
}
