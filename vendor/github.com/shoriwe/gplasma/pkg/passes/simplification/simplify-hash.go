package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simplify *simplifyPass) Hash(hash *ast.HashExpression) *ast2.Hash {
	values := make([]*ast2.KeyValue, 0, len(hash.Values))
	for _, value := range hash.Values {
		values = append(values, &ast2.KeyValue{
			Key:   simplify.Expression(value.Key),
			Value: simplify.Expression(value.Value),
		})
	}
	return &ast2.Hash{
		Values: values,
	}
}
