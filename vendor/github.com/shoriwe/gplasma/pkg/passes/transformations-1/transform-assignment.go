package transformations_1

import (
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
)

func (transform *transformPass) Assignment(assignment *ast2.Assignment) []ast3.Node {
	return []ast3.Node{&ast3.Assignment{
		Statement: nil,
		Left:      transform.Expression(assignment.Left).(ast3.Assignable),
		Right:     transform.Expression(assignment.Right),
	}}
}
