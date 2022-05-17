package ast

import (
	"github.com/shoriwe/gplasma/pkg/errors"
	"github.com/shoriwe/gplasma/pkg/vm"
)

type Node interface {
	Compile() ([]*vm.Code, *errors.Error)
	CompilePush(push bool) ([]*vm.Code, *errors.Error)
}

func compileBody(body []Node) ([]*vm.Code, *errors.Error) {
	var result []*vm.Code
	for _, node := range body {
		nodeCode, compileError := node.Compile()
		if compileError != nil {
			return nil, compileError
		}
		result = append(result, nodeCode...)
	}
	return result, nil
}

type Program struct {
	Node
	Begin *BeginStatement
	End   *EndStatement
	Body  []Node
}

func (program *Program) Compile() ([]*vm.Code, *errors.Error) {
	var result []*vm.Code
	if program.Begin != nil {
		beginBody, beginCompilationError := program.Begin.Compile()
		if beginCompilationError != nil {
			return nil, beginCompilationError
		}
		result = append(result, beginBody...)
	}
	body, bodyCompilationError := compileBody(program.Body)
	if bodyCompilationError != nil {
		return nil, bodyCompilationError
	}
	result = append(result, body...)
	if program.End != nil {
		endBody, endCompilationError := program.End.Compile()
		if endCompilationError != nil {
			return nil, endCompilationError
		}
		result = append(result, endBody...)
	}
	return result, nil
}
