package transformations_1

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
	magic_functions "github.com/shoriwe/gplasma/pkg/common/magic-functions"
)

func (transform *transformPass) Binary(binary *ast2.Binary) *ast3.Call {
	var (
		function string
	)
	switch binary.Operator {
	case ast2.And:
		function = magic_functions.And
	case ast2.Or:
		function = magic_functions.Or
	case ast2.Xor:
		function = magic_functions.Xor
	case ast2.In:
		return transform.Call(
			&ast2.FunctionCall{
				Expression: nil,
				Function: &ast2.Selector{
					Assignable: nil,
					X:          binary.Right,
					Identifier: &ast2.Identifier{
						Assignable: nil,
						Symbol:     magic_functions.In,
					},
				},
				Arguments: []ast2.Expression{binary.Left},
			},
		)
	case ast2.Is:
		function = magic_functions.Is
	case ast2.Implements:
		function = magic_functions.Implements
	case ast2.Equals:
		function = magic_functions.Equal
	case ast2.NotEqual:
		function = magic_functions.NotEqual
	case ast2.GreaterThan:
		function = magic_functions.GreaterThan
	case ast2.GreaterOrEqualThan:
		function = magic_functions.GreaterOrEqualThan
	case ast2.LessThan:
		function = magic_functions.LessThan
	case ast2.LessOrEqualThan:
		function = magic_functions.LessOrEqualThan
	case ast2.BitwiseOr:
		function = magic_functions.BitwiseOr
	case ast2.BitwiseXor:
		function = magic_functions.BitwiseXor
	case ast2.BitwiseAnd:
		function = magic_functions.BitwiseAnd
	case ast2.BitwiseLeft:
		function = magic_functions.BitwiseLeft
	case ast2.BitwiseRight:
		function = magic_functions.BitwiseRight
	case ast2.Add:
		function = magic_functions.Add
	case ast2.Sub:
		function = magic_functions.Sub
	case ast2.Mul:
		function = magic_functions.Mul
	case ast2.Div:
		function = magic_functions.Div
	case ast2.FloorDiv:
		function = magic_functions.FloorDiv
	case ast2.Modulus:
		function = magic_functions.Modulus
	case ast2.PowerOf:
		function = magic_functions.PowerOf
	default:
		panic(fmt.Sprintf("unknown binary operator %d", binary.Operator))
	}
	return transform.Call(
		&ast2.FunctionCall{
			Expression: nil,
			Function: &ast2.Selector{
				Assignable: nil,
				X:          binary.Left,
				Identifier: &ast2.Identifier{
					Assignable: nil,
					Symbol:     function,
				},
			},
			Arguments: []ast2.Expression{binary.Right},
		},
	)
}
