package transformations_1

import (
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
)

func (transform *transformPass) Hash(hash *ast2.Hash) *ast3.Hash {
	values := make([]*ast3.KeyValue, 0, len(hash.Values))
	for _, keyValue := range hash.Values {
		values = append(values, &ast3.KeyValue{
			Key:   transform.Expression(keyValue.Key),
			Value: transform.Expression(keyValue.Value),
		})
	}
	return &ast3.Hash{
		Values: values,
	}
}
