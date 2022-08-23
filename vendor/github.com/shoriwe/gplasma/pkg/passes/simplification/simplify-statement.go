package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simplify *simplifyPass) Statement(stmt ast.Statement) ast2.Statement {
	switch s := stmt.(type) {
	case *ast.AssignStatement:
		return simplify.Assign(s)
	case *ast.DoWhileStatement:
		return simplify.DoWhile(s)
	case *ast.WhileLoopStatement:
		return simplify.While(s)
	case *ast.UntilLoopStatement:
		return simplify.Until(s)
	case *ast.ForLoopStatement:
		return simplify.For(s)
	case *ast.IfStatement:
		return simplify.If(s)
	case *ast.UnlessStatement:
		return simplify.Unless(s)
	case *ast.SwitchStatement:
		return simplify.Switch(s)
	case *ast.ModuleStatement:
		return simplify.Module(s)
	case *ast.FunctionDefinitionStatement:
		return simplify.Function(s)
	case *ast.GeneratorDefinitionStatement:
		return simplify.GeneratorDef(s)
	case *ast.InterfaceStatement:
		return simplify.Interface(s)
	case *ast.ClassStatement:
		return simplify.Class(s)
	case *ast.ReturnStatement:
		return simplify.Return(s)
	case *ast.YieldStatement:
		return simplify.Yield(s)
	case *ast.ContinueStatement:
		return simplify.Continue(s)
	case *ast.BreakStatement:
		return simplify.Break(s)
	case *ast.PassStatement:
		return simplify.Pass(s)
	case *ast.DeleteStatement:
		return simplify.Delete(s)
	case *ast.DeferStatement:
		return simplify.Defer(s)
	default:
		panic("unknown statement type")
	}
}
