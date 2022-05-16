package vm

func (p *Plasma) GetNone() *Value {
	return p.ForceMasterGetAny(None)
}

func (p *Plasma) NewNone(context *Context, isBuiltIn bool, parent *SymbolTable) *Value {
	result := p.NewValue(context, isBuiltIn, NoneName, nil, parent)
	result.BuiltInTypeId = NoneId
	p.NoneInitialize(isBuiltIn)(context, result)
	result.SetOnDemandSymbol(Self,
		func() *Value {
			return result
		},
	)
	return result
}

func (p *Plasma) NoneInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(context *Context, object *Value) *Value {
		object.SetOnDemandSymbol(Equals,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(_ *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							return p.InterpretAsBool(right.IsTypeById(NoneId)), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(NotEquals,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(_ *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							return p.InterpretAsBool(left.IsTypeById(NoneId)), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToString,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(_ *Value, _ ...*Value) (*Value, bool) {
							return p.NewString(context, false, "None"), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToBool,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(_ *Value, _ ...*Value) (*Value, bool) {
							return p.GetFalse(), true
						},
					),
				)
			},
		)
		return nil
	}
}
