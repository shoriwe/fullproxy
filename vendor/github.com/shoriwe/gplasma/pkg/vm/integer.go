package vm

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/tools"
	"math"
)

func (p *Plasma) NewInteger(context *Context, isBuiltIn bool, value int64) *Value {
	integer := p.NewValue(context, isBuiltIn, IntegerName, nil, context.PeekSymbolTable())
	integer.BuiltInTypeId = IntegerId
	integer.Integer = value
	p.IntegerInitialize(isBuiltIn)(context, integer)
	integer.SetOnDemandSymbol(Self,
		func() *Value {
			return integer
		},
	)
	return integer
}

func (p *Plasma) IntegerInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(context *Context, object *Value) *Value {
		object.SetOnDemandSymbol(NegateBits,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							return p.NewInteger(context, false,
								^self.Integer,
							), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Negative,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							return p.NewInteger(context, false,
								-self.Integer,
							), true
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
							switch right.BuiltInTypeId {
							case IntegerId:
								return p.NewInteger(context, false,
									self.Integer+right.Integer,
								), true
							case FloatId:
								return p.NewFloat(context, false,
									float64(self.Integer)+right.Float,
								), true
							default:
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName), false
							}
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
							switch left.BuiltInTypeId {
							case IntegerId:
								return p.NewInteger(context, false,
									left.Integer+self.Integer,
								), true
							case FloatId:
								return p.NewFloat(context, false,
									left.Float+float64(self.Integer),
								), true
							default:
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Sub,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							switch right.BuiltInTypeId {
							case IntegerId:
								return p.NewInteger(context, false,
									self.Integer-right.Integer,
								), true
							case FloatId:
								return p.NewFloat(context, false,
									float64(self.Integer)-right.Float,
								), true
							default:
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightSub,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							switch left.BuiltInTypeId {
							case IntegerId:
								return p.NewInteger(context, false,
									left.Integer-self.Integer,
								), true
							case FloatId:
								return p.NewFloat(context, false,
									left.Float-float64(self.Integer),
								), true
							default:
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName), false
							}
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
								return p.NewInteger(context, false,
									self.Integer*right.Integer,
								), true
							case FloatId:
								return p.NewFloat(context, false,
									float64(self.Integer)*right.Float,
								), true
							case StringId:
								return p.NewString(context, false, tools.Repeat(right.String, self.Integer)), true
							case BytesId:
								content, repetitionError := p.Repeat(context, right.Content, self.Integer)
								if repetitionError != nil {
									return repetitionError, false
								}
								return p.NewTuple(context, false, content), true
							case ArrayId:
								content, repetitionError := p.Repeat(context, right.Content, self.Integer)
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
								return p.NewInteger(context, false,
									left.Integer*self.Integer,
								), true
							case FloatId:
								return p.NewFloat(context, false,
									left.Float*float64(self.Integer),
								), true
							case StringId:
								return p.NewString(context, false, tools.Repeat(left.String, self.Integer)), true
							case BytesId:
								content, repetitionError := p.Repeat(context, left.Content, self.Integer)
								if repetitionError != nil {
									return repetitionError, false
								}
								return p.NewTuple(context, false, content), true
							case ArrayId:
								content, repetitionError := p.Repeat(context, left.Content, self.Integer)
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
		object.SetOnDemandSymbol(Div,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							switch right.BuiltInTypeId {
							case IntegerId:
								return p.NewFloat(context, false,
									float64(self.Integer)/float64(right.Integer),
								), true
							case FloatId:
								return p.NewFloat(context, false,
									float64(self.Integer)/right.Float,
								), true
							default:
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightDiv,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							switch left.BuiltInTypeId {
							case IntegerId:
								return p.NewFloat(context, false,
									float64(left.Integer)/float64(self.Integer),
								), true
							case FloatId:
								return p.NewFloat(context, false,
									left.Float/float64(self.Integer),
								), true
							default:
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(FloorDiv,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							switch right.BuiltInTypeId {
							case IntegerId:
								return p.NewInteger(context, false,
									self.Integer/right.Integer,
								), true
							case FloatId:
								return p.NewInteger(context, false,
									self.Integer/int64(right.Float),
								), true
							default:
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightFloorDiv,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							switch left.BuiltInTypeId {
							case IntegerId:
								return p.NewInteger(context, false,
									left.Integer/self.Integer,
								), true
							case FloatId:
								return p.NewInteger(context, false,
									int64(left.Float)/self.Integer,
								), true
							default:
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Mod,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							switch right.BuiltInTypeId {
							case IntegerId:
								return p.NewFloat(context, false,
									math.Mod(float64(self.Integer), float64(right.Integer)),
								), true
							case FloatId:
								return p.NewFloat(context, false,
									math.Mod(float64(self.Integer), right.Float),
								), true
							default:
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightMod,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							switch left.BuiltInTypeId {
							case IntegerId:
								return p.NewInteger(context, false,
									left.Integer%self.Integer,
								), true
							case FloatId:
								return p.NewInteger(context, false,
									int64(left.Float)%self.Integer,
								), true
							default:
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Pow,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							switch right.BuiltInTypeId {
							case IntegerId:
								return p.NewFloat(context, false,
									math.Pow(float64(self.Integer), float64(right.Integer)),
								), true
							case FloatId:
								return p.NewFloat(context, false,
									math.Pow(float64(self.Integer), right.Float),
								), true
							default:
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightPow,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							switch left.BuiltInTypeId {
							case IntegerId:
								return p.NewFloat(context, false,
									math.Pow(
										float64(left.Integer),
										float64(self.Integer),
									),
								), true
							case FloatId:
								return p.NewFloat(context, false,
									math.Pow(
										left.Float,
										float64(self.Integer),
									),
								), true
							default:
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(BitXor,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							if !right.IsTypeById(BytesId) {
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName), false
							}
							return p.NewInteger(context, false,
								self.Integer^right.Integer,
							), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightBitXor,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							if !left.IsTypeById(IntegerId) {
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName), false
							}
							return p.NewInteger(context, false,
								left.Integer^self.Integer,
							), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(BitAnd,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							if !right.IsTypeById(BytesId) {
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName), false
							}
							return p.NewInteger(context, false,
								self.Integer&right.Integer,
							), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightBitAnd,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							if !left.IsTypeById(IntegerId) {
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName), false
							}
							return p.NewInteger(context, false,
								left.Integer&self.Integer,
							), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(BitOr,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							if !right.IsTypeById(BytesId) {
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName), false
							}
							return p.NewInteger(context, false,
								self.Integer|right.Integer,
							), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightBitOr,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							if !left.IsTypeById(IntegerId) {
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName), false
							}
							return p.NewInteger(context, false,
								left.Integer|self.Integer,
							), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(BitLeft,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							if !right.IsTypeById(BytesId) {
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName), false
							}
							return p.NewInteger(context, false,
								self.Integer<<uint(right.Integer),
							), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightBitLeft,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							if !left.IsTypeById(IntegerId) {
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName), false
							}
							return p.NewInteger(context, false,
								left.Integer<<uint(self.Integer),
							), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(BitRight,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							if !right.IsTypeById(BytesId) {
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName), false
							}
							return p.NewInteger(context, false,
								self.Integer>>uint(right.Integer),
							), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightBitRight,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							if !left.IsTypeById(IntegerId) {
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName), false
							}
							return p.NewInteger(context, false,
								left.Integer>>uint(self.Integer),
							), true
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
							switch right.BuiltInTypeId {
							case IntegerId:
								return p.InterpretAsBool(self.Integer == right.Integer), true
							case FloatId:
								return p.InterpretAsBool(float64(self.Integer) == right.Float), true
							default:
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName), false
							}
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
							switch left.BuiltInTypeId {
							case IntegerId:
								return p.InterpretAsBool(left.Integer == self.Integer), true
							case FloatId:
								return p.InterpretAsBool(left.Float == float64(self.Integer)), true
							default:
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName), false
							}
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
							switch right.BuiltInTypeId {
							case IntegerId:
								return p.InterpretAsBool(self.Integer != right.Integer), true
							case FloatId:
								return p.InterpretAsBool(float64(self.Integer) != (right.Float)), true
							default:
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName), false
							}
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
							switch left.BuiltInTypeId {
							case IntegerId:
								return p.InterpretAsBool(left.Integer != self.Integer), true
							case FloatId:
								return p.InterpretAsBool(left.Float != float64(self.Integer)), true
							default:
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(GreaterThan,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							switch right.BuiltInTypeId {
							case IntegerId:
								return p.InterpretAsBool(self.Integer > right.Integer), true
							case FloatId:
								return p.InterpretAsBool(float64(self.Integer) > right.Float), true
							default:
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightGreaterThan,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							switch left.BuiltInTypeId {
							case IntegerId:
								return p.InterpretAsBool(left.Integer > self.Integer), true
							case FloatId:
								return p.InterpretAsBool(left.Float > float64(self.Integer)), true
							default:
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(LessThan,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							switch right.BuiltInTypeId {
							case IntegerId:
								return p.InterpretAsBool(self.Integer < right.Integer), true
							case FloatId:
								return p.InterpretAsBool(float64(self.Integer) < right.Float), true
							default:
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightLessThan,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							switch left.BuiltInTypeId {
							case IntegerId:
								return p.InterpretAsBool(left.Integer < self.Integer), true
							case FloatId:
								return p.InterpretAsBool(left.Float < float64(self.Integer)), true
							default:
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(GreaterThanOrEqual,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							switch right.BuiltInTypeId {
							case IntegerId:
								return p.InterpretAsBool(self.Integer >= right.Integer), true
							case FloatId:
								return p.InterpretAsBool(float64(self.Integer) >= right.Float), true
							default:
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightGreaterThanOrEqual,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							switch left.BuiltInTypeId {
							case IntegerId:
								return p.InterpretAsBool(left.Integer >= self.Integer), true
							case FloatId:
								return p.InterpretAsBool(left.Float >= float64(self.Integer)), true
							default:
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(LessThanOrEqual,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							right := arguments[0]
							switch right.BuiltInTypeId {
							case IntegerId:
								return p.InterpretAsBool(self.Integer <= right.Integer), true
							case FloatId:
								return p.InterpretAsBool(float64(self.Integer) <= right.Float), true
							default:
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName, FloatName), false
							}
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightLessThanOrEqual,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							left := arguments[0]
							switch left.BuiltInTypeId {
							case IntegerId:
								return p.InterpretAsBool(left.Integer <= self.Integer), true
							case FloatId:
								return p.InterpretAsBool(left.Float <= float64(self.Integer)), true
							default:
								return p.NewInvalidTypeError(context, left.TypeName(), IntegerName, FloatName), false
							}
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
							return p.NewInteger(context, false, self.Integer), true
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
							return p.NewInteger(context, false, self.Integer), true
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
							return p.NewInteger(context, false, self.Integer), true
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
							return p.NewFloat(context, false,
								float64(self.Integer),
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
							return p.NewString(context, false, fmt.Sprint(self.Integer)), true
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
							return p.InterpretAsBool(self.Integer != 0), true
						},
					),
				)
			},
		)
		return nil
	}
}
