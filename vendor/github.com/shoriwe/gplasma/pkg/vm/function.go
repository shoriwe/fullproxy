package vm

func (p *Plasma) NewFunction(context *Context, isBuiltIn bool, parentSymbols *SymbolTable, callable Callable) *Value {
	function := p.NewValue(context, isBuiltIn, FunctionName, nil, parentSymbols)
	function.BuiltInTypeId = FunctionId
	function.Callable = callable
	function.SetOnDemandSymbol(Self,
		func() *Value {
			return function
		},
	)
	return function
}
