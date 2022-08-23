package transformations_1

import (
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
)

func (transform *transformPass) Delete(del *ast2.Delete) []ast3.Node {
	return []ast3.Node{
		&ast3.Delete{
			X: transform.Expression(del.X).(ast3.Assignable),
		},
	}
}
