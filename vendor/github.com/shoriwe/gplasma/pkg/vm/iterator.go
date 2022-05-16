package vm

func (p *Plasma) NewIterator(context *Context, isBuiltIn bool) *Value {
	iterator := p.NewValue(context, isBuiltIn, IteratorName, nil, context.PeekSymbolTable())
	iterator.BuiltInTypeId = IteratorId
	p.IteratorInitialize(isBuiltIn)(context, iterator)
	iterator.SetOnDemandSymbol(Self,
		func() *Value {
			return iterator
		},
	)
	return iterator
}

func (p *Plasma) IteratorInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(context *Context, object *Value) *Value {
		object.SetOnDemandSymbol(HasNext,
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
		object.SetOnDemandSymbol(Next,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(_ *Value, _ ...*Value) (*Value, bool) {
							return p.GetNone(), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Iter,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							return self, true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToTuple,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							return p.IterToContent(context, self, TupleId)
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToArray,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							return p.IterToContent(context, self, ArrayId)
						},
					),
				)
			},
		)
		return nil
	}
}
