package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simplify *simplifyPass) Array(array *ast.ArrayExpression) *ast2.Array {
	values := make([]ast2.Expression, 0, len(array.Values))
	for _, value := range array.Values {
		values = append(values, simplify.Expression(value))
	}
	return &ast2.Array{
		Values: values,
	}
}
