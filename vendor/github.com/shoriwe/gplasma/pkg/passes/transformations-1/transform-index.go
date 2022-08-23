package transformations_1

import (
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
)

func (transform *transformPass) Index(index *ast2.Index) *ast3.Index {
	return &ast3.Index{
		Source: transform.Expression(index.Source),
		Index:  transform.Expression(index.Index),
	}
}
