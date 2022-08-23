package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

type simplifyPass struct {
	currentAnonIdent uint
}

func (simplify *simplifyPass) Node(node ast.Node) ast2.Node {
	switch n := node.(type) {
	case ast.Statement:
		return simplify.Statement(n)
	case ast.Expression:
		return simplify.Expression(n)
	default:
		panic("unknown node type")
	}
}

func Simplify(program *ast.Program) (ast2.Program, error) {
	resultChan := make(chan ast2.Program, 1)
	errorChan := make(chan error, 1)
	go func(rChan chan ast2.Program, eChan chan error) {
		defer func() {
			err := recover()
			if err != nil {
				rChan <- nil
				eChan <- err.(error)
			}
		}()
		var (
			begin []ast2.Node
			body  []ast2.Node
			end   []ast2.Node
		)
		simp := simplifyPass{currentAnonIdent: 1}
		if program.Begin != nil {
			begin = make([]ast2.Node, 0, len(program.Begin.Body))
			for _, node := range program.Begin.Body {
				begin = append(begin, simp.Node(node))
			}
		}
		body = make([]ast2.Node, 0, len(program.Body))
		for _, node := range program.Body {
			body = append(body, simp.Node(node))
		}
		if program.End != nil {
			end = make([]ast2.Node, 0, len(program.End.Body))
			for _, node := range program.End.Body {
				end = append(end, simp.Node(node))
			}
		}
		result := make(ast2.Program, 0, len(begin)+len(body)+len(end))
		result = append(result, begin...)
		result = append(result, body...)
		result = append(result, end...)
		resultChan <- result
		errorChan <- nil
	}(resultChan, errorChan)
	return <-resultChan, <-errorChan
}
