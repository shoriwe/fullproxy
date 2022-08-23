package transformations_1

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
	"reflect"
)

type transformPass struct {
	currentLabel      int
	currentIdentifier int
}

func (transform *transformPass) Node(node ast2.Node) []ast3.Node {
	switch n := node.(type) {
	case ast2.Statement:
		return transform.Statement(n)
	case ast2.Expression:
		return []ast3.Node{transform.Expression(n)}
	default:
		panic(fmt.Sprintf("unknown node type %s", reflect.TypeOf(n).String()))
	}
}

func Transform(program ast2.Program) (ast3.Program, error) {
	resultChan := make(chan ast3.Program, 1)
	errorChan := make(chan error, 1)
	go func(rChan chan ast3.Program, eChan chan error) {
		defer func() {
			err := recover()
			if err != nil {
				rChan <- nil
				eChan <- err.(error)
			}
		}()
		result := make(ast3.Program, 0, len(program))
		transform := transformPass{
			currentLabel:      1,
			currentIdentifier: 1,
		}
		for _, node := range program {
			result = append(result, transform.Node(node)...)
		}
		rChan <- result
		eChan <- nil
	}(resultChan, errorChan)
	return <-resultChan, <-errorChan
}
