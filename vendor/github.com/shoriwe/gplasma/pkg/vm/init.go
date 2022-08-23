package vm

import (
	"bufio"
	magic_functions "github.com/shoriwe/gplasma/pkg/common/magic-functions"
	special_symbols "github.com/shoriwe/gplasma/pkg/common/special-symbols"
)

func (plasma *Plasma) init() {
	// On Demand values
	plasma.onDemand = map[string]func(*Value) *Value{
		magic_functions.Equal: func(self *Value) *Value {
			return plasma.NewBuiltInFunction(
				self.vtable,
				func(argument ...*Value) (*Value, error) {
					return plasma.NewBool(self == argument[0]), nil
				},
			)
		},
		magic_functions.NotEqual: func(self *Value) *Value {
			return plasma.NewBuiltInFunction(
				self.vtable,
				func(argument ...*Value) (*Value, error) {
					return plasma.NewBool(self != argument[0]), nil
				},
			)
		},
		magic_functions.And: func(self *Value) *Value {
			return plasma.NewBuiltInFunction(self.vtable,
				func(argument ...*Value) (*Value, error) {
					if self.Bool() && argument[0].Bool() {
						return plasma.true, nil
					}
					return plasma.false, nil
				})
		},
		magic_functions.Or: func(self *Value) *Value {
			return plasma.NewBuiltInFunction(self.vtable,
				func(argument ...*Value) (*Value, error) {
					if self.Bool() || argument[0].Bool() {
						return plasma.true, nil
					}
					return plasma.false, nil
				})
		},
		magic_functions.Xor: func(self *Value) *Value {
			return plasma.NewBuiltInFunction(self.vtable,
				func(argument ...*Value) (*Value, error) {
					if self.Bool() != argument[0].Bool() {
						return plasma.true, nil
					}
					return plasma.false, nil
				})
		},
		magic_functions.Is: func(self *Value) *Value {
			return plasma.NewBuiltInFunction(self.vtable,
				func(argument ...*Value) (*Value, error) {
					class := argument[0]
					switch class.TypeId() {
					case BuiltInClassId, ClassId:
						return plasma.NewBool(self.GetClass() == class), nil
					}
					return plasma.false, nil
				})
		},
		magic_functions.Implements: func(self *Value) *Value {
			return plasma.NewBuiltInFunction(self.vtable,
				func(argument ...*Value) (*Value, error) {
					class := argument[0]
					switch class.TypeId() {
					case BuiltInClassId, ClassId:
						return plasma.NewBool(self.GetClass().Implements(class)), nil
					}
					return plasma.false, nil
				})
		},
		magic_functions.Bool: func(self *Value) *Value {
			return plasma.NewBuiltInFunction(self.vtable,
				func(argument ...*Value) (*Value, error) {
					return plasma.NewBool(self.Bool()), nil
				})
		},
		magic_functions.Class: func(self *Value) *Value {
			return plasma.NewBuiltInFunction(self.vtable,
				func(argument ...*Value) (*Value, error) {
					return self.GetClass(), nil
				})
		},
		magic_functions.SubClasses: func(self *Value) *Value {
			return plasma.NewBuiltInFunction(self.vtable,
				func(argument ...*Value) (*Value, error) {
					return plasma.NewTuple(self.GetClass().GetClassInfo().Bases), nil
				})
		},
		magic_functions.Iter: func(self *Value) *Value {
			return plasma.NewBuiltInFunction(
				self.vtable,
				func(argument ...*Value) (*Value, error) {
					return self, nil
				},
			)
		},
	}
	// Init classes
	plasma.metaClass()
	plasma.value = plasma.valueClass()
	plasma.function = plasma.functionClass()
	plasma.string = plasma.stringClass()
	plasma.bytes = plasma.bytesClass()
	plasma.bool = plasma.boolClass()
	plasma.noneType = plasma.noneClass()
	plasma.int = plasma.integerClass()
	plasma.float = plasma.floatClass()
	plasma.array = plasma.arrayClass()
	plasma.tuple = plasma.tupleClass()
	plasma.hash = plasma.hashClass()
	// Init values
	plasma.true = plasma.NewBool(true)
	plasma.false = plasma.NewBool(false)
	plasma.none = plasma.NewNone()
	// Init symbols
	// -- Classes
	plasma.rootSymbols.Set(special_symbols.Value, plasma.value)
	plasma.rootSymbols.Set(special_symbols.String, plasma.string)
	plasma.rootSymbols.Set(special_symbols.Bytes, plasma.bytes)
	plasma.rootSymbols.Set(special_symbols.Bool, plasma.bool)
	plasma.rootSymbols.Set(special_symbols.None, plasma.noneType)
	plasma.rootSymbols.Set(special_symbols.Int, plasma.int)
	plasma.rootSymbols.Set(special_symbols.Float, plasma.float)
	plasma.rootSymbols.Set(special_symbols.Array, plasma.array)
	plasma.rootSymbols.Set(special_symbols.Tuple, plasma.tuple)
	plasma.rootSymbols.Set(special_symbols.Hash, plasma.hash)
	plasma.rootSymbols.Set(special_symbols.Function, plasma.function)
	plasma.rootSymbols.Set(special_symbols.Class, plasma.class)
	/*
		- input
		- print
		- println
		- range
	*/
	plasma.rootSymbols.Set(special_symbols.Input, plasma.NewBuiltInFunction(plasma.rootSymbols,
		func(argument ...*Value) (*Value, error) {
			_, writeError := plasma.Stdout.Write([]byte(argument[0].String()))
			if writeError != nil {
				panic(writeError)
			}
			scanner := bufio.NewScanner(plasma.Stdin)
			if scanner.Scan() {
				return plasma.NewString(scanner.Bytes()), nil
			}
			return plasma.none, nil
		},
	))
	plasma.rootSymbols.Set(special_symbols.Print, plasma.NewBuiltInFunction(plasma.rootSymbols,
		func(argument ...*Value) (*Value, error) {
			for index, arg := range argument {
				if index != 0 {
					_, writeError := plasma.Stdout.Write([]byte(" "))
					if writeError != nil {
						panic(writeError)
					}
				}
				_, writeError := plasma.Stdout.Write([]byte(arg.String()))
				if writeError != nil {
					panic(writeError)
				}

			}
			return plasma.none, nil
		},
	))
	plasma.rootSymbols.Set(special_symbols.Println, plasma.NewBuiltInFunction(plasma.rootSymbols,
		func(argument ...*Value) (*Value, error) {
			for index, arg := range argument {
				if index != 0 {
					_, writeError := plasma.Stdout.Write([]byte(" "))
					if writeError != nil {
						panic(writeError)
					}
				}
				_, writeError := plasma.Stdout.Write([]byte(arg.String()))
				if writeError != nil {
					panic(writeError)
				}
			}
			_, writeError := plasma.Stdout.Write([]byte("\n"))
			if writeError != nil {
				panic(writeError)
			}
			return plasma.none, nil
		},
	))
	plasma.rootSymbols.Set(special_symbols.Range, plasma.NewBuiltInFunction(plasma.rootSymbols,
		func(argument ...*Value) (*Value, error) {
			var (
				start              = argument[0]
				end                = argument[1]
				intStep      int64 = 1
				floatStep          = 1.0
				useFloatStep       = start.TypeId() == FloatId || end.TypeId() == FloatId
			)
			if len(argument) == 3 {
				step := argument[2]
				intStep = step.Int()
				floatStep = step.Float()
				if !useFloatStep {
					useFloatStep = step.TypeId() == FloatId
				}
			}
			iter := plasma.NewValue(plasma.rootSymbols, ValueId, plasma.value)
			if useFloatStep {
				iter.SetAny(start.Float())
				iter.Set(magic_functions.HasNext, plasma.NewBuiltInFunction(
					iter.vtable,
					func(_ ...*Value) (*Value, error) {
						return plasma.NewBool(iter.GetFloat64() < end.Float()), nil
					},
				))
				iter.Set(magic_functions.Next, plasma.NewBuiltInFunction(
					iter.vtable,
					func(_ ...*Value) (*Value, error) {
						current := iter.GetFloat64()
						// fmt.Println(current)
						iter.SetAny(current + floatStep)
						return plasma.NewFloat(current), nil
					},
				))
			} else {
				iter.SetAny(start.Int())
				iter.Set(magic_functions.HasNext, plasma.NewBuiltInFunction(
					iter.vtable,
					func(_ ...*Value) (*Value, error) {
						return plasma.NewBool(iter.GetInt64() < end.Int()), nil
					},
				))
				iter.Set(magic_functions.Next, plasma.NewBuiltInFunction(
					iter.vtable,
					func(_ ...*Value) (*Value, error) {
						current := iter.GetInt64()
						iter.SetAny(current + intStep)
						return plasma.NewInt(current), nil
					},
				))
			}
			return iter, nil
		},
	))
}
