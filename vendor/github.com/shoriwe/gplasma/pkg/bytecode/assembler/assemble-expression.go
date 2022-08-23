package assembler

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast3"
	"reflect"
)

func (a *assembler) Expression(expr ast3.Expression) []byte {
	if expr == nil {
		return nil
	}
	switch e := expr.(type) {
	case *ast3.Function:
		return a.Function(e)
	case *ast3.Class:
		return a.Class(e)
	case *ast3.Call:
		return a.Call(e)
	case *ast3.Array:
		return a.Array(e)
	case *ast3.Tuple:
		return a.Tuple(e)
	case *ast3.Hash:
		return a.Hash(e)
	case *ast3.Identifier:
		return a.Identifier(e)
	case *ast3.Integer:
		return a.Integer(e)
	case *ast3.Float:
		return a.Float(e)
	case *ast3.String:
		return a.String(e)
	case *ast3.Bytes:
		return a.Bytes(e)
	case *ast3.True:
		return a.True(e)
	case *ast3.False:
		return a.False(e)
	case *ast3.None:
		return a.None(e)
	case *ast3.Selector:
		return a.Selector(e)
	case *ast3.Index:
		return a.Index(e)
	case *ast3.Super:
		return a.Super(e)
	default:
		panic(fmt.Sprintf("unknown expression type %s", reflect.TypeOf(e).String()))
	}
}
