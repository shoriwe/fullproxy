package vm

import "github.com/shoriwe/gplasma/pkg/tools"

func (p *Plasma) NewArray(context *Context, isBuiltIn bool, content []*Value) *Value {
	array := p.NewValue(context, isBuiltIn, ArrayName, nil, context.PeekSymbolTable())
	array.BuiltInTypeId = ArrayId
	array.Content = content
	p.ArrayInitialize(isBuiltIn)(context, array)
	array.SetOnDemandSymbol(Self,
		func() *Value {
			return array
		},
	)
	return array
}

func (p *Plasma) ArrayInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(context *Context, object *Value) *Value {
		object.SetOnDemandSymbol(Add,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							if !arguments[0].IsTypeById(ArrayId) {
								return p.NewInvalidTypeError(context, arguments[0].GetClass(p).Name, ArrayName), false
							}
							return p.NewArray(context, false, append(self.Content, arguments[0].Content...)), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightAdd,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							if !arguments[0].IsTypeById(ArrayId) {
								return p.NewInvalidTypeError(context, arguments[0].GetClass(p).Name, ArrayName), false
							}
							return p.NewArray(context, false, append(arguments[0].Content, self.Content...)), true
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
								content, repetitionError := p.Repeat(context, self.Content,
									right.Integer)
								if repetitionError != nil {
									return repetitionError, false
								}
								return p.NewArray(context, false, content), true
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
								content, repetitionError := p.Repeat(context, self.Content,
									left.Integer,
								)
								if repetitionError != nil {
									return repetitionError, false
								}
								return p.NewArray(context, false, content), true
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
						func(_ *Value, _ ...*Value) (*Value, bool) {
							return p.NewUnhashableTypeError(context, object.GetClass(p)), false
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
		object.SetOnDemandSymbol(Assign,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 2,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							return p.ContentAssign(context, self, arguments[0], arguments[1])
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
							return p.NewTuple(context, false, append([]*Value{}, self.Content...)), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Append,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							self.Content = append(self.Content, arguments[0])
							return p.GetNone(), true
						},
					),
				)
			},
		)
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
		object.SetOnDemandSymbol(Pop,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							length := len(self.Content)
							if length > 0 {
								result := self.Content[length-1]
								self.Content = self.Content[:length-1]
								return result, true
							}
							return p.NewIndexOutOfRange(context, 0, 0), false
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Delete,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							if !arguments[0].IsTypeById(IntegerId) {
								return p.NewInvalidTypeError(context, arguments[0].GetClass(p).Name, IntegerName), false
							}
							contentLength := len(self.Content)
							realIndex, indexCalculationError := tools.CalcIndex(arguments[0].Integer, contentLength)
							if indexCalculationError != nil {
								return p.NewIndexOutOfRange(context, contentLength, arguments[0].Integer), false
							}
							self.Content = append(self.Content[:realIndex], self.Content[realIndex+1:]...)
							return p.GetNone(), true
						},
					),
				)
			},
		)
		return nil
	}
}
