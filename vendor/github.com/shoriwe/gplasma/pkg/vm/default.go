package vm

import (
	"fmt"
)

/*
	Type         - (Done)
	Function     - (Done)
	 Value       - (Done)
	Bool         - (Done)
	Bytes        - (Done)
	String       - (Done)
	HashTable    - (Done)
	Integer      - (Done)
	Array        - (Done)
	Tuple        - (Done)
	Hash         - (Done)
	Expand       - (Done)
	Id           - (Done)
	Range        - (Done)
	Len          - (Done)
	DeleteFrom   - (Done)
	Dir          - (Done)
	Input        - (Done)
	ToString     - (Done)
	ToTuple      - (Done)
	ToArray      - (Done)
	ToInteger    - (Done)
	ToFloat      - (Done)
	ToBool       - (Done)
*/
func (p *Plasma) InitializeBuiltIn() {
	if !p.BuiltInContext.SymbolStack.HasNext() {

	}
	/*
		This is the master symbol table that is protected from writes
	*/

	// Types
	p.BuiltInContext.PeekSymbolTable().Set(TypeName, p.NewType(p.BuiltInContext, true, TypeName, nil, nil, NewBuiltInConstructor(p.ObjectInitialize(true))))
	//// Default Error Types
	exception := p.NewType(p.BuiltInContext, true, RuntimeError, p.BuiltInContext.PeekSymbolTable(), nil,
		NewBuiltInConstructor(p.RuntimeErrorInitialize),
	)
	p.BuiltInContext.PeekSymbolTable().Set(RuntimeError, exception)
	p.BuiltInContext.PeekSymbolTable().Set(InvalidTypeError,
		p.NewType(p.BuiltInContext, true, InvalidTypeError, p.BuiltInContext.PeekSymbolTable(), []*Value{exception},
			NewBuiltInConstructor(
				func(context *Context, object *Value) *Value {
					object.SetOnDemandSymbol(Initialize,
						func() *Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 2,
									func(self *Value, arguments ...*Value) (*Value, bool) {
										received := arguments[0]
										if !received.IsTypeById(StringId) {
											return p.NewInvalidTypeError(context, received.TypeName(), StringName), false
										}
										expecting := arguments[1]
										if !expecting.IsTypeById(StringId) {
											return p.NewInvalidTypeError(context, expecting.TypeName(), StringName), false
										}
										self.String = fmt.Sprintf("Expecting %s but received %s", expecting.String, received.String)
										return p.GetNone(), true
									},
								),
							)
						},
					)
					return nil
				},
			),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set(NotImplementedCallableError,
		p.NewType(p.BuiltInContext, true, NotImplementedCallableError, p.BuiltInContext.PeekSymbolTable(), []*Value{exception},
			NewBuiltInConstructor(
				func(context *Context, object *Value) *Value {
					object.SetOnDemandSymbol(Initialize,
						func() *Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 1,
									func(self *Value, arguments ...*Value) (*Value, bool) {
										methodNameObject := arguments[0]
										methodNameObjectToString, getError := methodNameObject.Get(p, context, ToString)
										if getError != nil {
											return getError, false
										}
										methodNameString, success := p.CallFunction(context, methodNameObjectToString)
										if !success {
											return methodNameString, false
										}
										if !methodNameString.IsTypeById(StringId) {
											return p.NewInvalidTypeError(context, methodNameString.TypeName(), StringName), false
										}
										self.String = fmt.Sprintf("Method %s not implemented", methodNameString.String)
										return p.GetNone(), true
									},
								),
							)
						},
					)
					return nil
				},
			),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set(ObjectConstructionError,
		p.NewType(p.BuiltInContext, true, ObjectConstructionError, p.BuiltInContext.PeekSymbolTable(), []*Value{exception},
			NewBuiltInConstructor(
				func(context *Context, object *Value) *Value {
					object.SetOnDemandSymbol(Initialize,
						func() *Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 2,
									func(self *Value, arguments ...*Value) (*Value, bool) {
										typeName := arguments[0]
										if !typeName.IsTypeById(StringId) {
											return p.NewInvalidTypeError(context, typeName.TypeName(), StringName), false
										}
										errorMessage := arguments[1]
										if !typeName.IsTypeById(StringId) {
											return p.NewInvalidTypeError(context, errorMessage.TypeName(), StringName), false
										}
										self.String = fmt.Sprintf("Could not construct object of Type: %s -> %s", typeName.String, errorMessage.String)
										return p.GetNone(), true
									},
								),
							)
						},
					)
					return nil
				},
			),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set(ObjectWithNameNotFoundError,
		p.NewType(p.BuiltInContext, true, ObjectWithNameNotFoundError, p.BuiltInContext.PeekSymbolTable(), []*Value{exception},
			NewBuiltInConstructor(
				func(context *Context, object *Value) *Value {
					object.SetOnDemandSymbol(Initialize,
						func() *Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 2,
									func(self *Value, arguments ...*Value) (*Value, bool) {
										objectType := arguments[0]
										if !objectType.IsTypeById(TypeId) {
											return p.NewInvalidTypeError(context, objectType.TypeName(), TypeName), false
										}
										name := arguments[1]
										if !name.IsTypeById(StringId) {
											return p.NewInvalidTypeError(context, name.TypeName(), StringName), false
										}
										self.String = fmt.Sprintf(" Value with name %s not Found inside object of type %s", name.String, objectType.Name)
										return p.GetNone(), true
									},
								),
							)
						},
					)
					return nil
				},
			),
		),
	)

	p.BuiltInContext.PeekSymbolTable().Set(InvalidNumberOfArgumentsError,
		p.NewType(p.BuiltInContext, true, InvalidNumberOfArgumentsError, p.BuiltInContext.PeekSymbolTable(), []*Value{exception},
			NewBuiltInConstructor(
				func(context *Context, object *Value) *Value {
					object.SetOnDemandSymbol(Initialize,
						func() *Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 2,
									func(self *Value, arguments ...*Value) (*Value, bool) {
										received := arguments[0]
										if !received.IsTypeById(IntegerId) {
											return p.NewInvalidTypeError(context, received.TypeName(), IntegerName), false
										}
										expecting := arguments[1]
										if !expecting.IsTypeById(IntegerId) {
											return p.NewInvalidTypeError(context, expecting.TypeName(), IntegerName), false
										}
										self.String = fmt.Sprintf("Received %d but expecting %d expecting", received.Integer, expecting.Integer)
										return p.GetNone(), true
									},
								),
							)
						},
					)
					return nil
				},
			),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set(GoRuntimeError,
		p.NewType(p.BuiltInContext, true, GoRuntimeError, p.BuiltInContext.PeekSymbolTable(), []*Value{exception},
			NewBuiltInConstructor(
				func(context *Context, object *Value) *Value {
					object.SetOnDemandSymbol(Initialize,
						func() *Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 1,
									func(self *Value, arguments ...*Value) (*Value, bool) {
										runtimeError := arguments[0]
										if !runtimeError.IsTypeById(StringId) {
											return p.NewInvalidTypeError(context, runtimeError.TypeName(), StringName), false
										}
										self.String = runtimeError.String
										return p.GetNone(), true
									},
								),
							)
						},
					)
					return nil
				},
			),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set(UnhashableTypeError,
		p.NewType(p.BuiltInContext, true, UnhashableTypeError, p.BuiltInContext.PeekSymbolTable(), []*Value{exception},
			NewBuiltInConstructor(
				func(context *Context, object *Value) *Value {
					object.SetOnDemandSymbol(Initialize,
						func() *Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 1,
									func(self *Value, arguments ...*Value) (*Value, bool) {
										objectType := arguments[0]
										if !objectType.IsTypeById(TypeId) {
											return p.NewInvalidTypeError(context, objectType.TypeName(), TypeName), false
										}
										self.String = fmt.Sprintf(" Value of type: %s is unhasable", objectType.Name)
										return p.GetNone(), true
									},
								),
							)
						},
					)
					return nil
				},
			),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set(IndexOutOfRangeError,
		p.NewType(p.BuiltInContext, true, IndexOutOfRangeError, p.BuiltInContext.PeekSymbolTable(), []*Value{exception},
			NewBuiltInConstructor(
				func(context *Context, object *Value) *Value {
					object.SetOnDemandSymbol(Initialize,
						func() *Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 2,
									func(self *Value, arguments ...*Value) (*Value, bool) {
										length := arguments[0]
										if !length.IsTypeById(IntegerId) {
											return p.NewInvalidTypeError(context, length.TypeName(), IntegerName), false
										}
										index := arguments[1]
										if !index.IsTypeById(IntegerId) {
											return p.NewInvalidTypeError(context, index.TypeName(), IntegerName), false
										}
										self.String = fmt.Sprintf("Index: %d, out of range (Length=%d)", index.Integer, length.Integer)
										return p.GetNone(), true
									},
								),
							)
						},
					)
					return nil
				},
			),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set(KeyNotFoundError,
		p.NewType(p.BuiltInContext, true, KeyNotFoundError, p.BuiltInContext.PeekSymbolTable(), []*Value{exception},
			NewBuiltInConstructor(
				func(context *Context, object *Value) *Value {
					object.SetOnDemandSymbol(Initialize,
						func() *Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 1,
									func(self *Value, arguments ...*Value) (*Value, bool) {
										key := arguments[0]
										keyToString, getError := key.Get(p, context, ToString)
										if getError != nil {
											return getError, false
										}
										keyString, success := p.CallFunction(context, keyToString)
										if !success {
											return keyString, false
										}
										if !keyString.IsTypeById(StringId) {
											return p.NewInvalidTypeError(context, keyString.TypeName(), StringName), false
										}
										self.String = fmt.Sprintf("Key %s not found", keyString.String)
										return p.GetNone(), true
									},
								),
							)
						},
					)
					return nil
				},
			),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set(IntegerParsingError,
		p.NewType(p.BuiltInContext, true, IntegerParsingError, p.BuiltInContext.PeekSymbolTable(), []*Value{exception},
			NewBuiltInConstructor(
				func(context *Context, object *Value) *Value {
					object.SetOnDemandSymbol(Initialize,
						func() *Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 0,
									func(self *Value, arguments ...*Value) (*Value, bool) {
										self.String = "Integer parsing error"
										return p.GetNone(), true
									},
								),
							)
						},
					)
					return nil
				},
			),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set(FloatParsingError,
		p.NewType(p.BuiltInContext, true, FloatParsingError, p.BuiltInContext.PeekSymbolTable(), []*Value{exception},
			NewBuiltInConstructor(
				func(context *Context, object *Value) *Value {
					object.SetOnDemandSymbol(Initialize,
						func() *Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 0,
									func(self *Value, arguments ...*Value) (*Value, bool) {
										self.String = "Float parsing error"
										return p.GetNone(), true
									},
								),
							)
						},
					)
					return nil
				},
			),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set(BuiltInSymbolProtectionError,
		p.NewType(p.BuiltInContext, true, BuiltInSymbolProtectionError, p.BuiltInContext.PeekSymbolTable(), []*Value{exception},
			NewBuiltInConstructor(
				func(context *Context, object *Value) *Value {
					object.SetOnDemandSymbol(Initialize,
						func() *Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 1,
									func(self *Value, arguments ...*Value) (*Value, bool) {
										symbolName := arguments[0]
										if !symbolName.IsTypeById(StringId) {
											return p.NewInvalidTypeError(context, symbolName.TypeName(), StringName), false
										}
										self.String = fmt.Sprintf("cannot assign/delete built-in symbol %s", symbolName.String)
										return p.GetNone(), true
									},
								),
							)
						},
					)
					return nil
				},
			),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set(ObjectNotCallableError,
		p.NewType(p.BuiltInContext, true, ObjectNotCallableError, p.BuiltInContext.PeekSymbolTable(), []*Value{exception},
			NewBuiltInConstructor(
				func(context *Context, object *Value) *Value {
					object.SetOnDemandSymbol(Initialize,
						func() *Value {
							return p.NewFunction(context, true, object.SymbolTable(),
								NewBuiltInClassFunction(object, 1,
									func(self *Value, arguments ...*Value) (*Value, bool) {
										receivedType := arguments[0]
										if !receivedType.IsTypeById(TypeId) {
											return p.NewInvalidTypeError(context, receivedType.TypeName(), TypeName), false
										}
										self.String = fmt.Sprintf(" Value of type %s object is not callable", receivedType.Name)
										return p.GetNone(), true
									},
								),
							)
						},
					)
					return nil
				},
			),
		),
	)
	//// Default Types
	p.BuiltInContext.PeekSymbolTable().Set(CallableName,
		p.NewType(p.BuiltInContext, true, CallableName, p.BuiltInContext.PeekSymbolTable(), nil,
			NewBuiltInConstructor(p.CallableInitialize(false)),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set(NoneName,
		p.NewType(p.BuiltInContext, true, NoneName, p.BuiltInContext.PeekSymbolTable(), nil,
			NewBuiltInConstructor(p.NoneInitialize(false)),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set(ModuleName,
		p.NewType(p.BuiltInContext, true, ModuleName, p.BuiltInContext.PeekSymbolTable(), nil,
			NewBuiltInConstructor(p.ObjectInitialize(false)),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set(BoolName,
		p.NewType(p.BuiltInContext, true, BoolName, p.BuiltInContext.PeekSymbolTable(), nil,
			NewBuiltInConstructor(p.BoolInitialize(false)),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set(IteratorName,
		p.NewType(p.BuiltInContext, true, IteratorName, p.BuiltInContext.PeekSymbolTable(), nil,
			NewBuiltInConstructor(p.IteratorInitialize(false)),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set(FloatName,
		p.NewType(p.BuiltInContext, true, FloatName, p.BuiltInContext.PeekSymbolTable(), nil,
			NewBuiltInConstructor(p.FloatInitialize(false)),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set(ValueName,
		p.NewType(p.BuiltInContext, true, ValueName,
			nil, nil,
			NewBuiltInConstructor(p.ObjectInitialize(false)),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set(FunctionName,
		p.NewType(p.BuiltInContext, true, FunctionName, p.BuiltInContext.PeekSymbolTable(), nil,
			NewBuiltInConstructor(p.ObjectInitialize(false)),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set(IntegerName,
		p.NewType(p.BuiltInContext, true, IntegerName, p.BuiltInContext.PeekSymbolTable(), nil,
			NewBuiltInConstructor(p.IntegerInitialize(false)),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set(StringName,
		p.NewType(p.BuiltInContext, true, StringName, p.BuiltInContext.PeekSymbolTable(), nil,
			NewBuiltInConstructor(p.StringInitialize(false)),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set(BytesName,
		p.NewType(p.BuiltInContext, true, BytesName, p.BuiltInContext.PeekSymbolTable(), nil,
			NewBuiltInConstructor(p.BytesInitialize(false)),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set(TupleName,
		p.NewType(p.BuiltInContext, true, TupleName, p.BuiltInContext.PeekSymbolTable(), nil,
			NewBuiltInConstructor(p.TupleInitialize(false)),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set(ArrayName,
		p.NewType(p.BuiltInContext, true, ArrayName, p.BuiltInContext.PeekSymbolTable(), nil,
			NewBuiltInConstructor(p.ArrayInitialize(false)),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set(HashName,
		p.NewType(p.BuiltInContext, true, HashName, p.BuiltInContext.PeekSymbolTable(), nil,
			NewBuiltInConstructor(p.HashTableInitialize(false)),
		),
	)
	// Names
	p.BuiltInContext.PeekSymbolTable().Set(TrueName, p.NewBool(p.BuiltInContext, true, true))
	p.BuiltInContext.PeekSymbolTable().Set(FalseName, p.NewBool(p.BuiltInContext, true, false))
	p.BuiltInContext.PeekSymbolTable().Set(None, p.NewNone(p.BuiltInContext, true, p.BuiltInContext.PeekSymbolTable()))
	// Functions
	p.BuiltInContext.PeekSymbolTable().Set("expand",
		p.NewFunction(p.BuiltInContext, true, p.BuiltInContext.PeekSymbolTable(),
			NewBuiltInFunction(2,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					receiver := arguments[0]
					for symbol, object := range arguments[1].SymbolTable().Symbols {
						receiver.Set(p, p.BuiltInContext, symbol, object)
					}
					return p.GetNone(), true
				},
			),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set("dir",
		p.NewFunction(p.BuiltInContext, true, p.BuiltInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					object := arguments[0]
					var symbols []*Value
					for symbol := range object.Dir() {
						symbols = append(symbols, p.NewString(p.BuiltInContext, false, symbol))
					}
					return p.NewTuple(p.BuiltInContext, false, symbols), true
				},
			),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set("set",
		p.NewFunction(p.BuiltInContext, true, p.BuiltInContext.PeekSymbolTable(),
			NewBuiltInFunction(3,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					source := arguments[0]
					symbol := arguments[1]
					value := arguments[2]
					if !symbol.IsTypeById(StringId) {
						return p.NewInvalidTypeError(p.BuiltInContext, symbol.TypeName(), StringName), false
					}
					if source.IsBuiltIn {
						return p.NewBuiltInSymbolProtectionError(p.BuiltInContext, symbol.String), false
					}
					source.Set(p, p.BuiltInContext, symbol.String, value)
					return p.GetNone(), true
				},
			),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set("get_from",
		p.NewFunction(p.BuiltInContext, true, p.BuiltInContext.PeekSymbolTable(),
			NewBuiltInFunction(2,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					source := arguments[0]
					symbol := arguments[1]
					if !symbol.IsTypeById(StringId) {
						return p.NewInvalidTypeError(p.BuiltInContext, symbol.TypeName(), StringName), false
					}
					value, getError := source.Get(p, p.BuiltInContext, symbol.String)
					if getError != nil {
						return getError, false
					}
					return value, true
				},
			),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set("delete_from",
		p.NewFunction(p.BuiltInContext, true, p.BuiltInContext.PeekSymbolTable(),
			NewBuiltInFunction(2,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					source := arguments[0]
					symbol := arguments[1]
					if !symbol.IsTypeById(StringId) {
						return p.NewInvalidTypeError(p.BuiltInContext, symbol.TypeName(), StringName), false
					}
					if source.IsBuiltIn {
						return p.NewBuiltInSymbolProtectionError(p.BuiltInContext, symbol.String), false
					}
					_, getError := source.SymbolTable().GetSelf(symbol.String)
					if getError != nil {
						return p.NewObjectWithNameNotFoundError(p.BuiltInContext, source.GetClass(p), symbol.String), false
					}
					delete(source.SymbolTable().Symbols, symbol.String)
					return p.GetNone(), true
				},
			),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set("input",
		p.NewFunction(p.BuiltInContext, true, p.BuiltInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					message := arguments[0]
					var messageString *Value
					if !message.IsTypeById(StringId) {
						messageToString, getError := message.Get(p, p.BuiltInContext, ToString)
						if getError != nil {
							return getError, false
						}
						toStringResult, success := p.CallFunction(p.BuiltInContext, messageToString)
						if !success {
							return toStringResult, false
						}
						if !toStringResult.IsTypeById(StringId) {
							return p.NewInvalidTypeError(p.BuiltInContext, toStringResult.TypeName(), StringName), false
						}
						messageString = toStringResult
					} else {
						messageString = message
					}
					_, writingError := p.StdOut().Write([]byte(messageString.String))
					if writingError != nil {
						return p.NewGoRuntimeError(p.BuiltInContext, writingError), false
					}
					if p.StdInScanner().Scan() {
						return p.NewString(p.BuiltInContext, false, p.StdInScanner().Text()), true
					}
					return p.NewGoRuntimeError(p.BuiltInContext, p.StdInScanner().Err()), false
				},
			),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set("print",
		p.NewFunction(p.BuiltInContext, true, p.BuiltInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					value := arguments[0]
					toString, getError := value.Get(p, p.BuiltInContext, ToString)
					if getError != nil {
						return getError, false
					}
					stringValue, success := p.CallFunction(p.BuiltInContext, toString)
					if !success {
						return stringValue, false
					}
					_, writeError := fmt.Fprintf(p.StdOut(), "%s", stringValue.String)
					if writeError != nil {
						return p.NewGoRuntimeError(p.BuiltInContext, writeError), false
					}
					return p.GetNone(), true
				},
			),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set("println",
		p.NewFunction(p.BuiltInContext, true, p.BuiltInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					value := arguments[0]
					toString, getError := value.Get(p, p.BuiltInContext, ToString)
					if getError != nil {
						return getError, false
					}
					stringValue, success := p.CallFunction(p.BuiltInContext, toString)
					if !success {
						return stringValue, false
					}
					_, writeError := fmt.Fprintf(p.StdOut(), "%s\n", stringValue.String)
					if writeError != nil {
						return p.NewGoRuntimeError(p.BuiltInContext, writeError), false
					}
					return p.GetNone(), true
				},
			),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set("id",
		p.NewFunction(p.BuiltInContext, true, p.BuiltInContext.PeekSymbolTable(),
			NewBuiltInFunction(1,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					object := arguments[0]
					return p.NewInteger(p.BuiltInContext, false, object.Id()), true
				},
			),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set("range",
		p.NewFunction(p.BuiltInContext, true, p.BuiltInContext.PeekSymbolTable(),
			NewBuiltInFunction(3,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					start := arguments[0]
					if !start.IsTypeById(IntegerId) {
						return p.NewInvalidTypeError(p.BuiltInContext, start.TypeName(), IntegerName), false
					}

					end := arguments[1]
					if !end.IsTypeById(IntegerId) {
						return p.NewInvalidTypeError(p.BuiltInContext, end.TypeName(), IntegerName), false
					}

					step := arguments[2]
					if !step.IsTypeById(IntegerId) {
						return p.NewInvalidTypeError(p.BuiltInContext, step.TypeName(), IntegerName), false
					}

					rangeInformation := struct {
						current int64
						end     int64
						step    int64
					}{
						current: start.Integer,
						end:     end.Integer,
						step:    step.Integer,
					}

					// This should return a iterator
					rangeIterator := p.NewIterator(p.BuiltInContext, false)

					rangeIterator.Set(p, p.BuiltInContext, HasNext,
						p.NewFunction(p.BuiltInContext, true, rangeIterator.SymbolTable(),
							NewBuiltInClassFunction(rangeIterator, 0,
								func(self *Value, _ ...*Value) (*Value, bool) {
									if rangeInformation.current < rangeInformation.end {
										return p.GetTrue(), true
									}
									return p.GetFalse(), true
								},
							),
						),
					)
					rangeIterator.Set(p, p.BuiltInContext, Next,
						p.NewFunction(p.BuiltInContext, true, rangeIterator.SymbolTable(),
							NewBuiltInClassFunction(rangeIterator, 0,
								func(self *Value, _ ...*Value) (*Value, bool) {
									number := rangeInformation.current
									rangeInformation.current += rangeInformation.step
									return p.NewInteger(p.BuiltInContext, false, number), true
								},
							),
						),
					)
					return rangeIterator, true
				},
			),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set("super",
		p.NewFunction(p.BuiltInContext, true, p.BuiltInContext.PeekSymbolTable(),
			NewBuiltInFunction(2,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					return p.NewProxyObject(p.BuiltInContext, arguments[0], arguments[1])
				},
			),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set("map",
		p.NewFunction(p.BuiltInContext, true, p.BuiltInContext.PeekSymbolTable(),
			NewBuiltInFunction(2,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					asIterator, interpretationSuccess := p.InterpretAsIterator(p.BuiltInContext, arguments[1])
					if !interpretationSuccess {
						return asIterator, false
					}
					hasNext, hasNextGetError := asIterator.Get(p, p.BuiltInContext, HasNext)
					if hasNextGetError != nil {
						return hasNextGetError, false
					}
					next, nextGetError := asIterator.Get(p, p.BuiltInContext, Next)
					if nextGetError != nil {
						return nextGetError, false
					}
					result := p.NewIterator(p.BuiltInContext, false)

					// Set Next
					result.Set(p, p.BuiltInContext, Next,
						p.NewFunction(p.BuiltInContext, false, result.symbols,
							NewBuiltInClassFunction(result, 0,
								func(value *Value, value2 ...*Value) (*Value, bool) {
									nextValue, success := p.CallFunction(p.BuiltInContext, next)
									if !success {
										return nextValue, false
									}
									return p.CallFunction(p.BuiltInContext, arguments[0], nextValue)
								},
							),
						),
					)

					// Set HasNext
					result.Set(p, p.BuiltInContext, HasNext, hasNext)

					return result, true
				},
			),
		),
	)
	p.BuiltInContext.PeekSymbolTable().Set("filter",
		p.NewFunction(p.BuiltInContext, true, p.BuiltInContext.PeekSymbolTable(),
			NewBuiltInFunction(2,
				func(_ *Value, arguments ...*Value) (*Value, bool) {
					asIterator, interpretationSuccess := p.InterpretAsIterator(p.BuiltInContext, arguments[1])
					if !interpretationSuccess {
						return asIterator, false
					}
					hasNext, hasNextGetError := asIterator.Get(p, p.BuiltInContext, HasNext)
					if hasNextGetError != nil {
						return hasNextGetError, false
					}
					next, nextGetError := asIterator.Get(p, p.BuiltInContext, Next)
					if nextGetError != nil {
						return nextGetError, false
					}
					result := p.NewIterator(p.BuiltInContext, false)

					// Set Next
					result.Set(p, p.BuiltInContext, Next,
						p.NewFunction(p.BuiltInContext, false, result.symbols,
							NewBuiltInClassFunction(result, 0,
								func(value *Value, value2 ...*Value) (*Value, bool) {
									var (
										nextValue    *Value
										filterResult *Value
									)
									for {
										doesHasNext, success := p.CallFunction(p.BuiltInContext, hasNext)
										if !success {
											return doesHasNext, false
										}
										nextValue, success = p.CallFunction(p.BuiltInContext, next)
										if !success {
											return nextValue, false
										}
										filterResult, success = p.CallFunction(p.BuiltInContext, arguments[0], nextValue)
										if !success {
											return filterResult, false
										}
										asBool, interpretationError := p.QuickGetBool(p.BuiltInContext, filterResult)
										if interpretationError != nil {
											return interpretationError, false
										}
										if asBool {
											return nextValue, true
										}
									}
								},
							),
						),
					)

					// Set HasNext
					result.Set(p, p.BuiltInContext, HasNext, hasNext)

					return result, true
				},
			),
		),
	)
}
