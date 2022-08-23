package transformations_1

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
	"reflect"
)

func (transform *transformPass) Expression(expr ast2.Expression) ast3.Expression {
	switch e := expr.(type) {
	case *ast2.Binary:
		return transform.Binary(e)
	case *ast2.Unary:
		return transform.Unary(e)
	case *ast2.IfOneLiner:
		return transform.IfOneLiner(e)
	case *ast2.Array:
		return transform.Array(e)
	case *ast2.Tuple:
		return transform.Tuple(e)
	case *ast2.Hash:
		return transform.Hash(e)
	case *ast2.Identifier:
		return transform.Identifier(e)
	case *ast2.Integer:
		return transform.Integer(e)
	case *ast2.Float:
		return transform.Float(e)
	case *ast2.String:
		return transform.String(e)
	case *ast2.Bytes:
		return transform.Bytes(e)
	case *ast2.True:
		return transform.True(e)
	case *ast2.False:
		return transform.False(e)
	case *ast2.None:
		return transform.None(e)
	case *ast2.Lambda:
		return transform.Lambda(e)
	case *ast2.Generator:
		return transform.GeneratorExpr(e)
	case *ast2.Selector:
		return transform.Selector(e)
	case *ast2.FunctionCall:
		return transform.Call(e)
	case *ast2.Index:
		return transform.Index(e)
	case *ast2.Super:
		return transform.Super(e)
	default:
		panic(fmt.Sprintf("unknown expression type %s", reflect.TypeOf(e).String()))
	}
}
