package vm

func (p *Plasma) NewModule(context *Context, isBuiltIn bool) *Value {
	module := p.NewValue(context, isBuiltIn, ModuleName, nil, context.PeekSymbolTable())
	module.BuiltInTypeId = ModuleId
	module.SetOnDemandSymbol(Self,
		func() *Value {
			return module
		},
	)
	return module
}
