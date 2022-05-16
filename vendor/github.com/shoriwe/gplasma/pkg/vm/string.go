package vm

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/tools"
	"strconv"
	"strings"
)

func (p *Plasma) NewString(context *Context, isBuiltIn bool, value string) *Value {
	string_ := p.NewValue(context, isBuiltIn, StringName, nil, context.PeekSymbolTable())
	string_.BuiltInTypeId = StringId
	string_.String = value

	p.StringInitialize(isBuiltIn)(context, string_)
	string_.SetOnDemandSymbol(Self,
		func() *Value {
			return string_
		},
	)
	return string_
}

func (p *Plasma) StringInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(context *Context, object *Value) *Value {
		object.SetOnDemandSymbol(Length,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							return p.NewInteger(context, false, int64(len(self.String))), true
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
							if !right.IsTypeById(StringId) {
								return p.NewInvalidTypeError(context, right.TypeName(), StringName), false
							}
							return p.NewString(context, false,
								self.String+right.String,
							), true
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
							if !left.IsTypeById(StringId) {
								return p.NewInvalidTypeError(context, left.TypeName(), StringName), false
							}
							return p.NewString(context, false,
								left.String+self.String,
							), true
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
							if !right.IsTypeById(IntegerId) {
								return p.NewInvalidTypeError(context, right.TypeName(), IntegerName), false
							}
							return p.NewString(context, false,
								tools.Repeat(self.String, right.Integer),
							), true
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
							return p.NewString(context, false,
								tools.Repeat(self.String, left.Integer),
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
							if !right.IsTypeById(StringId) {
								return p.GetFalse(), true
							}
							return p.InterpretAsBool(self.String == right.String), true
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
							if !left.IsTypeById(StringId) {
								return p.GetFalse(), true
							}
							return p.InterpretAsBool(left.String == self.String), true
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
							if !right.IsTypeById(StringId) {
								return p.GetTrue(), true
							}
							return p.InterpretAsBool(self.String != right.String), true
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
							if !left.IsTypeById(StringId) {
								return p.GetTrue(), true
							}
							return p.InterpretAsBool(left.String != self.String), true
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
							if self.GetHash() == 0 {
								stringHash := p.HashString(fmt.Sprintf("%s-%s", self.String, StringName))
								self.SetHash(stringHash)
							}
							return p.NewInteger(context, false, self.GetHash()), true
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
							return p.NewString(context, false,
								self.String,
							), true
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
							return p.StringIndex(context, self, arguments[0])
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
							return p.StringIterator(context, self)
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
							number, parsingError := strconv.ParseInt(strings.ReplaceAll(self.String, "_", ""), 10, 64)
							if parsingError != nil {
								return p.NewIntegerParsingError(context), false
							}
							return p.NewInteger(context, false, number), true
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
							number, parsingError := strconv.ParseFloat(strings.ReplaceAll(self.String, "_", ""), 64)
							if parsingError != nil {
								return p.NewFloatParsingError(context), false
							}
							return p.NewFloat(context, false, number), true
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
							return p.NewString(context, false,
								self.String,
							), true
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
							return p.InterpretAsBool(len(self.String) > 0), true
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
							return p.StringToContent(context, self, ArrayId)
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
							return p.StringToContent(context, self, TupleId)
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(ToBytes,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							return p.NewBytes(context, false, []byte(self.String)), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Lower,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							return p.NewString(context, false, strings.ToLower(self.String)), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Upper,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							return p.NewString(context, false, strings.ToUpper(self.String)), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Split,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							if !arguments[0].IsTypeById(StringId) {
								return p.NewInvalidTypeError(context, arguments[0].GetClass(p).Name, StringName), false
							}
							var subStrings []*Value
							for _, s := range strings.Split(self.String, arguments[0].String) {
								subStrings = append(subStrings, p.NewString(context, false, s))
							}
							return p.NewArray(context, false, subStrings), true
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
							if !arguments[0].IsTypeById(StringId) {
								return p.NewInvalidTypeError(context, arguments[0].GetClass(p).Name, StringName), false
							}
							if !arguments[1].IsTypeById(StringId) {
								return p.NewInvalidTypeError(context, arguments[1].GetClass(p).Name, StringName), false
							}
							return p.NewString(context, false, strings.ReplaceAll(self.String, arguments[0].String, arguments[1].String)), true
						},
					),
				)
			},
		)
		return nil
	}
}
