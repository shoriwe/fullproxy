package assembler

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast3"
	"reflect"
)

func (a *assembler) Statement(stmt ast3.Statement) []byte {
	switch s := stmt.(type) {
	case *ast3.Assignment:
		return a.Assignment(s)
	case *ast3.Label:
		return a.Label(s)
	case *ast3.Jump:
		return a.Jump(s)
	case *ast3.ContinueJump:
		return a.ContinueJump(s)
	case *ast3.BreakJump:
		return a.BreakJump(s)
	case *ast3.IfJump:
		return a.IfJump(s)
	case *ast3.Return:
		return a.Return(s)
	case *ast3.Yield:
		return a.Yield(s)
	case *ast3.Delete:
		return a.Delete(s)
	case *ast3.Defer:
		return a.Defer(s)
	default:
		panic(fmt.Sprintf("unknown type of statement %s", reflect.TypeOf(s).String()))
	}
}
