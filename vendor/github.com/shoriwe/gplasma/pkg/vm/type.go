package vm

func (p *Plasma) NewType(
	context *Context,
	isBuiltIn bool,
	typeName string,
	parent *SymbolTable,
	subclasses []*Value,
	constructor Constructor,
) *Value {
	result := p.NewValue(context, isBuiltIn, TypeName, subclasses, parent)
	result.BuiltInTypeId = TypeId
	result.Constructor = constructor
	result.Name = typeName
	result.SetOnDemandSymbol(ToString,
		func() *Value {
			return p.NewFunction(context, isBuiltIn, result.symbols,
				NewBuiltInClassFunction(result, 0,
					func(_ *Value, _ ...*Value) (*Value, bool) {
						return p.NewString(context, false, "Type@"+typeName), true
					},
				),
			)
		},
	)
	result.SetOnDemandSymbol(Self,
		func() *Value {
			return result
		},
	)
	return result
}
