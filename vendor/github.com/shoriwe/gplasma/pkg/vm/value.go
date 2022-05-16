package vm

import (
	"fmt"
)

type OnDemandLoader func() *Value

const (
	ValueId = iota
	HashTableId
	IteratorId
	BoolId
	FunctionId
	IntegerId
	FloatId
	StringId
	BytesId
	ArrayId
	TupleId
	ModuleId
	TypeId
	NoneId
)

type Value struct {
	IsBuiltIn       bool
	id              int64
	typeName        string
	BuiltInTypeId   uint16
	class           *Value
	subClasses      []*Value
	symbols         *SymbolTable
	Callable        Callable
	Constructor     Constructor
	Name            string
	hash            int64
	Bool            bool
	String          string
	Bytes           []uint8
	Integer         int64
	Float           float64
	Content         []*Value
	KeyValues       map[int64][]*KeyValue
	onDemandSymbols map[string]OnDemandLoader
}

func (o *Value) AddKeyValue(hash int64, keyValue *KeyValue) {
	o.KeyValues[hash] = append(o.KeyValues[hash], keyValue)
}

func (o *Value) Id() int64 {
	return o.id
}

func (o *Value) SubClasses() []*Value {
	return o.subClasses
}

func (o *Value) Get(p *Plasma, context *Context, symbol string) (*Value, *Value) {
	result, getError := o.symbols.GetSelf(symbol)
	if getError != nil {
		loader, found := o.onDemandSymbols[symbol]
		if !found {
			return nil, p.NewObjectWithNameNotFoundError(context, o.GetClass(p), symbol)
		}
		result = loader()
		o.Set(p, context, symbol, result)
	}
	return result, nil
}

func (o *Value) SetOnDemandSymbol(symbol string, loader OnDemandLoader) {
	o.onDemandSymbols[symbol] = loader
}

func (o *Value) GetOnDemandSymbolLoader(symbol string) OnDemandLoader {
	return o.onDemandSymbols[symbol]
}
func (o *Value) GetOnDemandSymbols() map[string]OnDemandLoader {
	return o.onDemandSymbols
}

func (o *Value) Dir() map[string]byte {
	result := map[string]byte{}
	for symbol := range o.symbols.Symbols {
		result[symbol] = 0
	}
	for symbol := range o.onDemandSymbols {
		result[symbol] = 0
	}
	return result
}

func (o *Value) Set(p *Plasma, context *Context, symbol string, object *Value) *Value {
	if o.IsBuiltIn {
		return p.NewBuiltInSymbolProtectionError(context, symbol)
	}
	o.symbols.Set(symbol, object)
	return nil
}

func (o *Value) TypeName() string {
	return o.typeName
}

func (o *Value) SymbolTable() *SymbolTable {
	return o.symbols
}

func (o *Value) GetHash() int64 {
	return o.hash
}

func (o *Value) SetHash(newHash int64) {
	o.hash = newHash
}

func (o *Value) GetClass(p *Plasma) *Value {
	if o.class == nil { // This should only happen with built-ins
		o.class = p.ForceMasterGetAny(o.typeName)
	}
	return o.class
}

func (o *Value) SetClass(class *Value) {
	o.class = class
}

func (o *Value) Implements(class *Value) bool {
	if o.IsTypeById(TypeId) {
		if o == class {
			return true
		}
		for _, subClass := range o.subClasses {
			if subClass.Implements(class) {
				return true
			}
		}
		return false
	}
	if o.class == class {
		return true
	}
	for _, subClass := range o.subClasses {
		if subClass.Implements(class) {
			return true
		}
	}
	return false
}

func (o *Value) IsTypeById(id uint16) bool {
	return o.BuiltInTypeId == id
}

func (p *Plasma) NewValue(
	context *Context,
	isBuiltIn bool,
	typeName string,
	subClasses []*Value,
	parentSymbols *SymbolTable,
) *Value {
	result := &Value{
		id:              p.NextId(),
		typeName:        typeName,
		subClasses:      subClasses,
		symbols:         NewSymbolTable(parentSymbols),
		IsBuiltIn:       isBuiltIn,
		onDemandSymbols: map[string]OnDemandLoader{},
		BuiltInTypeId:   ValueId,
	}
	result.BuiltInTypeId = ValueId
	result.Bool = true
	result.String = ""
	result.Integer = 0
	result.Float = 0
	result.Content = []*Value{}
	result.KeyValues = map[int64][]*KeyValue{}
	result.Bytes = []uint8{}
	result.SetOnDemandSymbol(Self,
		func() *Value {
			return result
		},
	)
	p.ObjectInitialize(isBuiltIn)(context, result)
	return result
}

func (p *Plasma) ObjectInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(context *Context, object *Value) *Value {
		object.SetOnDemandSymbol(Initialize,
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
		object.SetOnDemandSymbol(Negate,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							selfBool, callError := p.QuickGetBool(context, self)
							if callError != nil {
								return callError, false
							}
							return p.InterpretAsBool(!selfBool), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(And,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							leftBool, callError := p.QuickGetBool(context, self)
							if callError != nil {
								return callError, false
							}
							var rightBool bool
							rightBool, callError = p.QuickGetBool(context, arguments[0])
							if callError != nil {
								return callError, false
							}
							return p.InterpretAsBool(leftBool && rightBool), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightAnd,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							leftBool, callError := p.QuickGetBool(context, arguments[0])
							if callError != nil {
								return callError, false
							}
							var rightBool bool
							rightBool, callError = p.QuickGetBool(context, self)
							if callError != nil {
								return callError, false
							}
							return p.InterpretAsBool(leftBool && rightBool), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Or,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							leftBool, callError := p.QuickGetBool(context, self)
							if callError != nil {
								return callError, false
							}
							var rightBool bool
							rightBool, callError = p.QuickGetBool(context, arguments[0])
							if callError != nil {
								return callError, false
							}
							return p.InterpretAsBool(leftBool || rightBool), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightOr,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							leftBool, callError := p.QuickGetBool(context, arguments[0])
							if callError != nil {
								return callError, false
							}
							var rightBool bool
							rightBool, callError = p.QuickGetBool(context, self)
							if callError != nil {
								return callError, false
							}
							return p.InterpretAsBool(leftBool || rightBool), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Xor,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							leftBool, callError := p.QuickGetBool(context, self)
							if callError != nil {
								return callError, false
							}
							var rightBool bool
							rightBool, callError = p.QuickGetBool(context, arguments[0])
							if callError != nil {
								return callError, false
							}
							return p.InterpretAsBool(leftBool != rightBool), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(RightXor,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							leftBool, callError := p.QuickGetBool(context, arguments[0])
							if callError != nil {
								return callError, false
							}
							var rightBool bool
							rightBool, callError = p.QuickGetBool(context, self)
							if callError != nil {
								return callError, false
							}
							return p.InterpretAsBool(leftBool != rightBool), true
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
							return p.InterpretAsBool(self.Id() == right.Id()), true
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
							return p.InterpretAsBool(self.Id() == right.Id()), true
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
							return p.InterpretAsBool(left.Id() == self.Id()), true
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
							return p.InterpretAsBool(self.Id() != right.Id()), true
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
							return p.InterpretAsBool(left.Id() != self.Id()), true
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
								objectHash := p.HashString(fmt.Sprintf("%v-%s-%d", self, self.TypeName(), self.Id()))
								self.SetHash(objectHash)
							}
							return p.NewInteger(context, false, self.GetHash()), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Class,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							return self.GetClass(p), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(SubClasses,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							var subClassesCopy []*Value
							for _, class := range self.SubClasses() {
								subClassesCopy = append(subClassesCopy, class)
							}
							return p.NewTuple(context, false, subClassesCopy), true
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
								fmt.Sprintf("%s{%s}-%X", ValueName, self.TypeName(), self.Id())), true
						},
					),
				)
			})
		object.SetOnDemandSymbol(ToBool,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(_ *Value, _ ...*Value) (*Value, bool) {
							return p.GetTrue(), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(GetInteger,
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
		object.SetOnDemandSymbol(GetBool,
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
		object.SetOnDemandSymbol(GetBytes,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							return p.NewBytes(context, false, self.Bytes), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(GetString,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							return p.NewString(context, false, self.String), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(GetFloat,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							return p.NewFloat(context, false, self.Float), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(GetContent,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							return p.NewArray(context, false, self.Content), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(GetKeyValues,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, _ ...*Value) (*Value, bool) {
							result := p.NewHashTable(context, false)
							result.KeyValues = self.KeyValues
							return result, true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(SetBool,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							self.Bool = arguments[0].Bool
							return p.GetNone(), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(SetBytes,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							self.Bytes = arguments[0].Bytes
							return p.GetNone(), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(SetString,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							self.String = arguments[0].String
							return p.GetNone(), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(SetInteger,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							self.Integer = arguments[0].Integer
							return p.GetNone(), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(SetFloat,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							self.Float = arguments[0].Float
							return p.GetNone(), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(SetContent,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							self.Content = arguments[0].Content
							return p.GetNone(), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(SetKeyValues,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							self.KeyValues = (arguments[0].KeyValues)
							return p.GetNone(), true
						},
					),
				)
			},
		)
		return nil
	}
}
