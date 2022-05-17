package vm

import (
	"bytes"
	"encoding/binary"
)

func (p *Plasma) NewBytes(context *Context, isBuiltIn bool, content []uint8) *Value {
	bytes_ := p.NewValue(context, isBuiltIn, BytesName, nil, context.PeekSymbolTable())
	bytes_.BuiltInTypeId = BytesId
	bytes_.Bytes = content
	p.BytesInitialize(isBuiltIn)(context, bytes_)
	bytes_.SetOnDemandSymbol(Self,
		func() *Value {
			return bytes_
		},
	)
	return bytes_
}

func (p *Plasma) BytesInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(context *Context, object *Value) *Value {
		object.SetOnDemandSymbol(Length,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							return p.NewInteger(context, false, int64(len(self.Bytes))), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Add,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							if !right.IsTypeById(BytesId) {
								return p.NewInvalidTypeError(context, right.TypeName(), BytesName), false
							}
							var newContent []uint8
							copy(newContent, self.Bytes)
							newContent = append(newContent, right.Bytes...)
							return p.NewBytes(context, false, newContent), true
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
							left := arguments[0]
							if !left.IsTypeById(BytesId) {
								return p.NewInvalidTypeError(context, left.TypeName(), BytesName), false
							}
							var newContent []uint8
							copy(newContent, left.Bytes)
							newContent = append(newContent, self.Bytes...)
							return p.NewBytes(context, false, newContent), true
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
							if !right.IsTypeById(BytesId) {
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName), false
							}
							return p.NewBytes(context, false, bytes.Repeat(self.Bytes, int(right.Integer))), true
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
							if !left.IsTypeById(IntegerId) {
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName), false
							}
							return p.NewBytes(context, false, bytes.Repeat(left.Bytes, int(self.Integer))), true
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
							right := arguments[0]
							if !right.IsTypeById(BytesId) {
								return p.GetFalse(), true
							}
							if len(self.Bytes) != len(right.Bytes) {
								return p.GetFalse(), true
							}
							return p.InterpretAsBool(bytes.Compare(self.Bytes, right.Bytes) == 0), true
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
							if !left.IsTypeById(BytesId) {
								return p.GetFalse(), true
							}
							if len(left.Bytes) != len(self.Bytes) {
								return p.GetFalse(), true
							}
							return p.InterpretAsBool(bytes.Compare(left.Bytes, self.Bytes) == 0), true
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
							if !right.IsTypeById(BytesId) {
								return p.GetFalse(), true
							}
							if len(self.Bytes) != len(right.Bytes) {
								return p.GetFalse(), true
							}
							return p.InterpretAsBool(bytes.Compare(self.Bytes, right.Bytes) != 0), true
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
							if !left.IsTypeById(BytesId) {
								return p.GetFalse(), true
							}
							if len(left.Bytes) != len(self.Bytes) {
								return p.GetFalse(), true
							}
							return p.InterpretAsBool(bytes.Compare(left.Bytes, self.Bytes) != 0), true
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
							selfHash := p.HashBytes(append(self.Bytes, []byte("Bytes")...))
							return p.NewInteger(context, false, selfHash), true
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
							var newBytes []uint8
							copy(newBytes, self.Bytes)
							return p.NewBytes(context, false, newBytes), true
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
							return p.BytesIndex(context, self, arguments[0])
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
							return p.BytesIterator(context, self)
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
							return p.NewInteger(context, false,
								int64(binary.BigEndian.Uint32(self.Bytes)),
							), true
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
							return p.NewString(context, false, string(self.Bytes)), true
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
							return p.InterpretAsBool(len(self.Bytes) != 0), true
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
							return p.BytesToContent(context, self, ArrayId)
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
							return p.BytesToContent(context, self, TupleId)
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Replace,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							if !arguments[0].IsTypeById(BytesId) {
								return p.NewInvalidTypeError(context, arguments[0].GetClass(p).Name, BytesName), false
							}
							var subBytes []*Value
							for _, s := range bytes.Split(self.Bytes, arguments[0].Bytes) {
								subBytes = append(subBytes, p.NewBytes(context, false, s))
							}
							return p.NewArray(context, false, subBytes), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Replace,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 2,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							if !arguments[0].IsTypeById(BytesId) {
								return p.NewInvalidTypeError(context, arguments[0].GetClass(p).Name, BytesName), false
							}
							if !arguments[1].IsTypeById(BytesId) {
								return p.NewInvalidTypeError(context, arguments[1].GetClass(p).Name, BytesName), false
							}
							return p.NewBytes(context, false, bytes.ReplaceAll(self.Bytes, arguments[0].Bytes, arguments[1].Bytes)), true
						},
					),
				)
			},
		)
		return nil
	}
}
