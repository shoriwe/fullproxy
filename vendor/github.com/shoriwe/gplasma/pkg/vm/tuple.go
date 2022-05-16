package vm

func (p *Plasma) NewTuple(context *Context, isBuiltIn bool, content []*Value) *Value {
	tuple := p.NewValue(context, isBuiltIn, TupleName, nil, context.PeekSymbolTable())
	tuple.BuiltInTypeId = TupleId
	tuple.Content = content
	p.TupleInitialize(isBuiltIn)(context, tuple)
	tuple.SetOnDemandSymbol(Self,
		func() *Value {
			return tuple
		},
	)
	return tuple
}

func (p *Plasma) TupleInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(context *Context, object *Value) *Value {
		object.SetOnDemandSymbol(Length,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							return p.NewInteger(context, false, int64(len(self.Content))), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Mul,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							switch right.BuiltInTypeId {
							case IntegerId:
								content, repetitionError := p.Repeat(context, self.Content, right.Integer)
								if repetitionError != nil {
									return repetitionError, false
								}
								return p.NewTuple(context, false, content), true
							default:
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName, StringName, ArrayName, TupleName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightMul,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							switch left.BuiltInTypeId {
							case IntegerId:
								content, repetitionError := p.Repeat(context, self.Content, left.Integer)
								if repetitionError != nil {
									return repetitionError, false
								}
								return p.NewTuple(context, false, content), true
							default:
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName, StringName, ArrayName, TupleName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Equals,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							return p.ContentEquals(context, self, arguments[0])
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightEquals,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							return p.ContentEquals(context, arguments[0], self)
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(NotEquals,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							return p.ContentNotEquals(context, self, arguments[0])
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightNotEquals,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							return p.ContentNotEquals(context, arguments[0], self)
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Contains,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							return p.ContentContains(context, self, arguments[0])
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Hash,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							return p.HashContent(context, self)
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Copy,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							return p.ContentCopy(context, self)
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Index,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							return p.ContentIndex(context, self, arguments[0])
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
							return p.ContentIterator(context, self)
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToString,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							return p.ContentToString(context, self)
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToBool,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							return p.InterpretAsBool(len(self.Content) != 0), true
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
							return p.NewArray(context, false, append([]*Value{}, self.Content...)), true
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
							return p.NewTuple(context, false, append([]*Value{}, self.Content...)), true
						},
					),
				)
			},
		)
		return nil
	}
}
