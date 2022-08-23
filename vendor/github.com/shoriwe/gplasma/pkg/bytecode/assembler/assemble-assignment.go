package assembler

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast3"
	"github.com/shoriwe/gplasma/pkg/bytecode/opcodes"
	"github.com/shoriwe/gplasma/pkg/common"
	magic_functions "github.com/shoriwe/gplasma/pkg/common/magic-functions"
	"reflect"
)

func (a *assembler) Assignment(assign *ast3.Assignment) []byte {
	result := a.Expression(assign.Right)
	result = append(result, opcodes.Push)
	switch left := assign.Left.(type) {
	case *ast3.Identifier:
		result = append(result, opcodes.IdentifierAssign)
		result = append(result, common.IntToBytes(len(left.Symbol))...)
		result = append(result, []byte(left.Symbol)...)
	case *ast3.Selector:
		result = append(result, a.Expression(left.X)...)
		result = append(result, opcodes.Push)
		result = append(result, opcodes.SelectorAssign)
		result = append(result, common.IntToBytes(len(left.Identifier.Symbol))...)
		result = append(result, []byte(left.Identifier.Symbol)...)
	case *ast3.Index:
		return a.Call(&ast3.Call{
			Function: &ast3.Selector{
				X: left.Source,
				Identifier: &ast3.Identifier{
					Symbol: magic_functions.Set,
				},
			},
			Arguments: []ast3.Expression{left.Index, assign.Right},
		})
	default:
		panic(fmt.Sprintf("unknown left hand side type %s", reflect.TypeOf(left).String()))
	}
	return result
}
