package transformations_1

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
	magic_functions "github.com/shoriwe/gplasma/pkg/common/magic-functions"
)

func (transform *transformPass) Unary(unary *ast2.Unary) *ast3.Call {
	var (
		function string
		x        ast3.Expression
	)
	switch unary.Operator {
	case ast2.Not:
		function = magic_functions.Not
	case ast2.Positive:
		function = magic_functions.Positive
	case ast2.Negative:
		function = magic_functions.Negative
	case ast2.NegateBits:
		function = magic_functions.NegateBits
	default:
		panic(fmt.Sprintf("unknown binary operator %d", unary.Operator))
	}
	x = transform.Expression(unary.X)
	return &ast3.Call{
		Function: &ast3.Selector{
			X: x,
			Identifier: &ast3.Identifier{
				Symbol: function,
			},
		},
	}
}
