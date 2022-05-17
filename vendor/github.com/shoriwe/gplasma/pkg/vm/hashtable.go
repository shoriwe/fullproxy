package vm

type KeyValue struct {
	Key   *Value
	Value *Value
}

func (p *Plasma) NewHashTable(context *Context, isBuiltIn bool) *Value {
	hashTable := p.NewValue(context, isBuiltIn, HashName, nil, context.PeekSymbolTable())
	hashTable.BuiltInTypeId = HashTableId
	p.HashTableInitialize(isBuiltIn)(context, hashTable)
	hashTable.SetOnDemandSymbol(Self,
		func() *Value {
			return hashTable
		},
	)
	return hashTable
}

func (p *Plasma) HashTableInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(context *Context, object *Value) *Value {
		object.SetOnDemandSymbol(Length,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							return p.NewInteger(context, false, int64(len(self.KeyValues))), true
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
							return p.HashEquals(context, self, arguments[0])
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
							return p.HashEquals(context, arguments[0], self)
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
							return p.HashNotEquals(context, self, arguments[0])
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
							return p.HashNotEquals(context, arguments[0], self)
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
							return p.HashContains(context, self, arguments[0])
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
							return p.NewUnhashableTypeError(context, p.ForceMasterGetAny(HashName)), false
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol(Copy,
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
		object.SetOnDemandSymbol(Index,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 1,
						func(self *Value, arguments ...*Value) (*Value, bool) {
							return p.HashIndex(context, self, arguments[0])
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
							return p.HashIndexAssign(context, self, arguments[0], arguments[1])
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
							return p.HashIterator(context, self)
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
							result := "{"
							var (
								valueToString *Value
								valueString   *Value
							)
							for _, keyValues := range self.KeyValues {
								for _, keyValue := range keyValues {
									keyToString, getError := keyValue.Key.Get(p, context, ToString)
									if getError != nil {
										return getError, false
									}
									keyString, success := p.CallFunction(context, keyToString)
									if !success {
										return keyString, false
									}
									result += keyString.String
									valueToString, getError = keyValue.Value.Get(p, context, ToString)
									if getError != nil {
										return getError, false
									}
									valueString, success = p.CallFunction(context, valueToString)
									if !success {
										return valueString, false
									}
									result += ": " + valueString.String + ", "
								}
							}
							if len(result) > 1 {
								result = result[:len(result)-2]
							}
							return p.NewString(context, false, result+"}"), true
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
							return p.InterpretAsBool(len(self.KeyValues) > 0), true
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
							return p.HashToContent(context, self, ArrayId)
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
							return p.HashToContent(context, self, TupleId)
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
							keyHash, success := p.Hash(context, arguments[0])
							if !success {
								return keyHash, false
							}
							if _, found := self.KeyValues[keyHash.Integer]; !found {
								return p.NewKeyNotFoundError(context, arguments[0]), false
							}
							for index, keyValue := range self.KeyValues[keyHash.Integer] {
								doesEquals, equalsError := p.Equals(context, keyValue.Key, arguments[0])
								if equalsError != nil {
									return equalsError, false
								}
								if doesEquals {
									self.KeyValues[keyHash.Integer] = append(self.KeyValues[keyHash.Integer][:index], self.KeyValues[keyHash.Integer][index+1:]...)
									return p.GetNone(), true
								}
							}
							return p.NewKeyNotFoundError(context, arguments[0]), false
						},
					),
				)
			},
		)
		return nil
	}
}
