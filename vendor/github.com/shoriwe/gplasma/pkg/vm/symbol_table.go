package vm

import (
	"github.com/shoriwe/gplasma/pkg/errors"
)

type SymbolTable struct {
	Parent  *SymbolTable
	Symbols map[string]*Value
}

func (symbolTable *SymbolTable) Set(s string, object *Value) {
	symbolTable.Symbols[s] = object
}

func (symbolTable *SymbolTable) GetSelf(symbol string) (*Value, *errors.Error) {
	result, found := symbolTable.Symbols[symbol]
	if !found {
		return nil, errors.NewNameNotFoundError()
	}
	return result, nil
}

func (symbolTable *SymbolTable) GetAny(symbol string) (*Value, *errors.Error) {
	for source := symbolTable; source != nil; source = source.Parent {
		result, found := source.Symbols[symbol]
		if found {
			return result, nil
		}
	}
	return nil, errors.NewNameNotFoundError()
}

func NewSymbolTable(parentSymbols *SymbolTable) *SymbolTable {
	return &SymbolTable{
		Parent:  parentSymbols,
		Symbols: map[string]*Value{},
	}
}
