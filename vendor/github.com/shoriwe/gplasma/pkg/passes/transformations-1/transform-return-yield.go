package transformations_1

import (
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
)

func (transform *transformPass) Yield(yield *ast2.Yield) []ast3.Node {
	return []ast3.Node{&ast3.Yield{
		Result: transform.Expression(yield.Result),
	}}
}

func (transform *transformPass) Return(ret *ast2.Return) []ast3.Node {
	return []ast3.Node{&ast3.Return{
		Result: transform.Expression(ret.Result),
	}}
}
