package transformations_1

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
	"reflect"
)

func (transform *transformPass) Statement(stmt ast2.Statement) []ast3.Node {
	switch s := stmt.(type) {
	case *ast2.Assignment:
		return transform.Assignment(s)
	case *ast2.DoWhile:
		return transform.DoWhile(s)
	case *ast2.While:
		return transform.While(s)
	case *ast2.If:
		return transform.If(s)
	case *ast2.Module:
		return transform.Module(s)
	case *ast2.FunctionDefinition:
		return transform.Function(s)
	case *ast2.GeneratorDefinition:
		return transform.GeneratorDef(s)
	case *ast2.Class:
		return transform.Class(s)
	case *ast2.Return:
		return transform.Return(s)
	case *ast2.Yield:
		return transform.Yield(s)
	case *ast2.Continue:
		return transform.Continue(s)
	case *ast2.Break:
		return transform.Break(s)
	case *ast2.Pass:
		return transform.Pass(s)
	case *ast2.Delete:
		return transform.Delete(s)
	case *ast2.Defer:
		return transform.Defer(s)
	default:
		panic(fmt.Sprintf("unknown statement type %s", reflect.TypeOf(s).String()))
	}
}
