package vm

func (p *Plasma) QuickGetBool(context *Context, value *Value) (bool, *Value) {
	if value.BuiltInTypeId == BoolId {
		return value.Bool, nil
	}
	valueToBool, getError := value.Get(p, context, ToBool)
	if getError != nil {
		return false, getError
	}
	valueBool, success := p.CallFunction(context, valueToBool)
	if !success {
		return false, valueBool
	}
	if !valueBool.IsTypeById(BoolId) {
		return false, p.NewInvalidTypeError(context, value.TypeName(), BoolName)
	}
	return valueBool.Bool, nil
}

func (p *Plasma) NewBool(context *Context, isBuiltIn bool, value bool) *Value {
	bool_ := p.NewValue(context, isBuiltIn, BoolName, nil, context.PeekSymbolTable())
	bool_.BuiltInTypeId = BoolId
	bool_.Bool = value
	p.BoolInitialize(isBuiltIn)(context, bool_)
	bool_.SetOnDemandSymbol(Self,
		func() *Value {
			return bool_
		},
	)
	return bool_
}

func (p *Plasma) GetFalse() *Value {
	return p.ForceMasterGetAny(FalseName)
}

func (p *Plasma) GetTrue() *Value {
	return p.ForceMasterGetAny(TrueName)
}

func (p *Plasma) BoolInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(context *Context, object *Value) *Value {
		object.SetOnDemandSymbol(Equals,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							if !right.IsTypeById(BoolId) {
								return p.GetFalse(), true
							}
							return p.InterpretAsBool(self.Bool == right.Bool), true
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
							left := arguments[0]
							if !left.IsTypeById(BoolId) {
								return p.GetFalse(), true
							}
							return p.InterpretAsBool(left.Bool == self.Bool), true
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
							right := arguments[0]
							if right.BuiltInTypeId != BoolId {
								return p.GetFalse(), true
							}
							return p.InterpretAsBool(self.Bool != right.Bool), true
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
							left := arguments[0]
							if left.BuiltInTypeId != BoolId {
								return p.GetFalse(), true
							}
							return p.InterpretAsBool(left.Bool != self.Bool), true
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
							return p.InterpretAsBool(self.Bool), true
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
							if self.Bool {
								return p.NewInteger(context, false, 1), true
							}
							return p.NewInteger(context, false, 0), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToInteger,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							if self.Bool {
								return p.NewInteger(context, false, 1), true
							}
							return p.NewInteger(context, false, 0), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToFloat,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							if self.Bool {
								return p.NewFloat(context, false, 1), true
							}
							return p.NewFloat(context, false, 0), true
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
							if self.Bool {
								return p.NewString(context, false, TrueName), true
							}
							return p.NewString(context, false, FalseName), true
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
							return p.InterpretAsBool(self.Bool), true
						},
					),
				)
			},
		)
		return nil
	}
}
