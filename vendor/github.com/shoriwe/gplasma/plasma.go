package gplasma

import (
	"github.com/shoriwe/gplasma/pkg/compiler/lexer"
	"github.com/shoriwe/gplasma/pkg/compiler/parser"
	"github.com/shoriwe/gplasma/pkg/reader"
	"github.com/shoriwe/gplasma/pkg/vm"
	"os"
)

type VirtualMachine struct {
	*vm.Plasma
}

func NewVirtualMachine() *VirtualMachine {
	return &VirtualMachine{Plasma: vm.NewPlasmaVM(os.Stdin, os.Stdout, os.Stderr)}
}

func (v *VirtualMachine) prepareCode(code string) ([]*vm.Code, *vm.Value) {
	root, parsingError := parser.NewParser(lexer.NewLexer(reader.NewStringReader(code))).Parse()
	if parsingError != nil {
		return nil, v.NewGoRuntimeError(v.BuiltInContext, parsingError.Error())
	}
	bytecode, compilationError := root.Compile()
	if compilationError != nil {
		return nil, v.NewGoRuntimeError(v.BuiltInContext, compilationError.Error())
	}
	return bytecode, nil
}

func (v *VirtualMachine) Execute(context *vm.Context, bytecode *vm.Bytecode) (*vm.Value, bool) {
	/*
		This only should be used when accessing developing a new module for the language
	*/
	context.PeekSymbolTable().Set(vm.IsMain, v.GetFalse())
	defer func() {
		_, found := context.PeekSymbolTable().Symbols[vm.IsMain]
		if found {
			delete(context.PeekSymbolTable().Symbols, vm.IsMain)
		}
	}()
	result, success := v.Plasma.Execute(context, bytecode)
	return result, success
}

func (v *VirtualMachine) ExecuteMain(mainScript string) (*vm.Value, bool) {
	bytecode, preparationError := v.prepareCode(mainScript)
	if preparationError != nil {
		return preparationError, false
	}
	context := v.NewContext()
	context.PeekSymbolTable().Set(vm.IsMain, v.GetTrue())
	defer func() {
		_, found := context.PeekSymbolTable().Symbols[vm.IsMain]
		if found {
			delete(context.PeekSymbolTable().Symbols, vm.IsMain)
		}
	}()
	return v.Plasma.Execute(context, vm.NewBytecodeFromArray(bytecode))
}
