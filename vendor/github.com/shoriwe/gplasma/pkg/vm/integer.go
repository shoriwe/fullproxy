package vm

import (
	"bytes"
	"encoding/binary"
	magic_functions "github.com/shoriwe/gplasma/pkg/common/magic-functions"
	"math"
)

func (plasma *Plasma) integerClass() *Value {
	class := plasma.NewValue(plasma.rootSymbols, BuiltInClassId, plasma.class)
	class.SetAny(Callback(func(argument ...*Value) (*Value, error) {
		return plasma.NewInt(argument[0].Int()), nil
	}))
	return class
}

/*
NewInt magic function:
Positive:           __positive__
Negative:           __negative__
NegateBits:         __negate_its__
Equal:             __equal__
NotEqual:           __not_equal__
GreaterThan:        __greater_than__
GreaterOrEqualThan: __greater_or_equal_than__
LessThan:           __less_than__
LessOrEqualThan:    __less_or_equal_than__
BitwiseOr:          __bitwise_or__
BitwiseXor:         __bitwise_xor__
BitwiseAnd:         __bitwise_and__
BitwiseLeft:        __bitwise_left__
BitwiseRight:       __bitwise_right__
Add:                __add__
Sub:                __sub__
Mul:                __mul__
Div:                __div__
FloorDiv:           __floor_div__
Modulus:            __mod__
PowerOf:            __pow__
Bool:               __bool__
String             __string__
Int                __int__
Float              __float__
Copy:               __copy__
BigEndian			big_endian
LittleEndian		little_endian
FromBig				from_big
FromLittle			from_little
*/
func (plasma *Plasma) NewInt(i int64) *Value {
	result := plasma.NewValue(plasma.rootSymbols, IntId, plasma.int)
	result.SetAny(i)
	result.Set(magic_functions.Positive,
		plasma.NewBuiltInFunction(
			result.vtable,
			func(argument ...*Value) (*Value, error) {
				return result, nil
			},
		))
	result.Set(magic_functions.Negative,
		plasma.NewBuiltInFunction(
			result.vtable,
			func(argument ...*Value) (*Value, error) {
				return plasma.NewInt(-result.GetInt64()), nil
			},
		))
	result.Set(magic_functions.NegateBits,
		plasma.NewBuiltInFunction(
			result.vtable,
			func(argument ...*Value) (*Value, error) {
				return plasma.NewInt(^result.GetInt64()), nil
			},
		))
	result.Set(magic_functions.Equal,
		plasma.NewBuiltInFunction(
			result.vtable,
			func(argument ...*Value) (*Value, error) {
				switch argument[0].TypeId() {
				case IntId, FloatId:
					return plasma.NewBool(result.Equal(argument[0])), nil
				}
				return plasma.false, nil
			},
		))
	result.Set(magic_functions.NotEqual,
		plasma.NewBuiltInFunction(
			result.vtable,
			func(argument ...*Value) (*Value, error) {
				switch argument[0].TypeId() {
				case IntId, FloatId:
					return plasma.NewBool(!result.Equal(argument[0])), nil
				}
				return plasma.true, nil
			},
		))
	result.Set(magic_functions.GreaterThan,
		plasma.NewBuiltInFunction(
			result.vtable,
			func(argument ...*Value) (*Value, error) {
				switch argument[0].TypeId() {
				case IntId:
					return plasma.NewBool(result.Int() > argument[0].Int()), nil
				case FloatId:
					return plasma.NewBool(result.Float() > argument[0].Float()), nil
				}
				return nil, NotComparable
			},
		))
	result.Set(magic_functions.GreaterOrEqualThan,
		plasma.NewBuiltInFunction(
			result.vtable,
			func(argument ...*Value) (*Value, error) {
				switch argument[0].TypeId() {
				case IntId:
					return plasma.NewBool(result.Int() >= argument[0].Int()), nil
				case FloatId:
					return plasma.NewBool(result.Float() >= argument[0].Float()), nil
				}
				return nil, NotComparable
			},
		))
	result.Set(magic_functions.LessThan,
		plasma.NewBuiltInFunction(
			result.vtable,
			func(argument ...*Value) (*Value, error) {
				switch argument[0].TypeId() {
				case IntId:
					return plasma.NewBool(result.Int() < argument[0].Int()), nil
				case FloatId:
					return plasma.NewBool(result.Float() < argument[0].Float()), nil
				}
				return nil, NotComparable
			},
		))
	result.Set(magic_functions.LessOrEqualThan,
		plasma.NewBuiltInFunction(
			result.vtable,
			func(argument ...*Value) (*Value, error) {
				switch argument[0].TypeId() {
				case IntId:
					return plasma.NewBool(result.Int() <= argument[0].Int()), nil
				case FloatId:
					return plasma.NewBool(result.Float() <= argument[0].Float()), nil
				}
				return nil, NotComparable
			},
		))
	result.Set(magic_functions.BitwiseOr,
		plasma.NewBuiltInFunction(
			result.vtable,
			func(argument ...*Value) (*Value, error) {
				switch argument[0].TypeId() {
				case IntId:
					return plasma.NewInt(result.Int() | argument[0].Int()), nil
				case FloatId:
					return plasma.NewInt(int64(uint64(result.Int()) | math.Float64bits(argument[0].Float()))), nil
				}
				return nil, NotOperable
			},
		))
	result.Set(magic_functions.BitwiseXor,
		plasma.NewBuiltInFunction(
			result.vtable,
			func(argument ...*Value) (*Value, error) {
				switch argument[0].TypeId() {
				case IntId:
					return plasma.NewInt(result.Int() ^ argument[0].Int()), nil
				case FloatId:
					return plasma.NewInt(int64(uint64(result.Int()) ^ math.Float64bits(argument[0].Float()))), nil
				}
				return nil, NotOperable
			},
		))
	result.Set(magic_functions.BitwiseAnd,
		plasma.NewBuiltInFunction(
			result.vtable,
			func(argument ...*Value) (*Value, error) {
				switch argument[0].TypeId() {
				case IntId:
					return plasma.NewInt(result.Int() & argument[0].Int()), nil
				case FloatId:
					return plasma.NewInt(int64(uint64(result.Int()) & math.Float64bits(argument[0].Float()))), nil
				}
				return nil, NotOperable
			},
		))
	result.Set(magic_functions.BitwiseLeft,
		plasma.NewBuiltInFunction(
			result.vtable,
			func(argument ...*Value) (*Value, error) {
				switch argument[0].TypeId() {
				case IntId:
					return plasma.NewInt(result.Int() << argument[0].Int()), nil
				case FloatId:
					return plasma.NewInt(int64(uint64(result.Int()) << math.Float64bits(argument[0].Float()))), nil
				}
				return nil, NotOperable
			},
		))
	result.Set(magic_functions.BitwiseRight,
		plasma.NewBuiltInFunction(
			result.vtable,
			func(argument ...*Value) (*Value, error) {
				switch argument[0].TypeId() {
				case IntId:
					return plasma.NewInt(result.Int() >> argument[0].Int()), nil
				case FloatId:
					return plasma.NewInt(int64(uint64(result.Int()) >> math.Float64bits(argument[0].Float()))), nil
				}
				return nil, NotOperable
			},
		))
	result.Set(magic_functions.Add,
		plasma.NewBuiltInFunction(
			result.vtable,
			func(argument ...*Value) (*Value, error) {
				switch argument[0].TypeId() {
				case IntId:
					return plasma.NewInt(result.Int() + argument[0].Int()), nil
				case FloatId:
					return plasma.NewFloat(result.Float() + argument[0].Float()), nil
				}
				return nil, NotOperable
			},
		))
	result.Set(magic_functions.Sub,
		plasma.NewBuiltInFunction(
			result.vtable,
			func(argument ...*Value) (*Value, error) {
				switch argument[0].TypeId() {
				case IntId:
					return plasma.NewInt(result.Int() - argument[0].Int()), nil
				case FloatId:
					return plasma.NewFloat(result.Float() - argument[0].Float()), nil
				}
				return nil, NotOperable
			},
		))
	result.Set(magic_functions.Mul,
		plasma.NewBuiltInFunction(
			result.vtable,
			func(argument ...*Value) (*Value, error) {
				switch argument[0].TypeId() {
				case IntId:
					return plasma.NewInt(result.Int() * argument[0].Int()), nil
				case FloatId:
					return plasma.NewFloat(result.Float() * argument[0].Float()), nil
				case StringId:
					s := argument[0].GetBytes()
					times := result.GetInt64()
					return plasma.NewString(bytes.Repeat(s, int(times))), nil
				case BytesId:
					s := argument[0].GetBytes()
					times := result.GetInt64()
					return plasma.NewBytes(bytes.Repeat(s, int(times))), nil
				case ArrayId:
					times := result.GetInt64()
					currentValues := argument[0].GetValues()
					newValues := make([]*Value, 0, int64(len(currentValues))*times)
					for t := int64(0); t < times; t++ {
						for _, value := range currentValues {
							newValues = append(newValues, value)
						}
					}
					return plasma.NewArray(newValues), nil
				}
				return nil, NotOperable
			},
		))
	result.Set(magic_functions.Div,
		plasma.NewBuiltInFunction(
			result.vtable,
			func(argument ...*Value) (*Value, error) {
				switch argument[0].TypeId() {
				case IntId, FloatId:
					return plasma.NewFloat(result.Float() / argument[0].Float()), nil
				}
				return nil, NotOperable
			},
		))
	result.Set(magic_functions.FloorDiv,
		plasma.NewBuiltInFunction(
			result.vtable,
			func(argument ...*Value) (*Value, error) {
				switch argument[0].TypeId() {
				case IntId, FloatId:
					return plasma.NewInt(result.Int() / argument[0].Int()), nil
				}
				return nil, NotOperable
			},
		))
	result.Set(magic_functions.Modulus,
		plasma.NewBuiltInFunction(
			result.vtable,
			func(argument ...*Value) (*Value, error) {
				switch argument[0].TypeId() {
				case IntId:
					return plasma.NewInt(result.Int() % argument[0].Int()), nil
				case FloatId:
					return plasma.NewFloat(math.Mod(result.Float(), argument[0].Float())), nil
				}
				return nil, NotOperable
			},
		))
	result.Set(magic_functions.PowerOf,
		plasma.NewBuiltInFunction(
			result.vtable,
			func(argument ...*Value) (*Value, error) {
				switch argument[0].TypeId() {
				case IntId:
					times := argument[0].Int()
					value := result.Int()
					if times >= 0 {
						v := int64(1)
						for t := int64(0); t < times; t++ {
							v *= value
						}
						return plasma.NewInt(v), nil
					}
					v := int64(1)
					for t := int64(0); times < t; t-- {
						v *= value
					}
					return plasma.NewFloat(1 / float64(v)), nil
				case FloatId:
					return plasma.NewFloat(math.Pow(result.Float(), argument[0].Float())), nil
				}
				return nil, NotOperable
			},
		))
	result.Set(magic_functions.Bool,
		plasma.NewBuiltInFunction(
			result.vtable,
			func(argument ...*Value) (*Value, error) {
				return plasma.NewBool(result.GetInt64() != 0), nil
			},
		))
	result.Set(magic_functions.String,
		plasma.NewBuiltInFunction(
			result.vtable,
			func(argument ...*Value) (*Value, error) {
				return plasma.NewString([]byte(result.String())), nil
			},
		))
	result.Set(magic_functions.Int,
		plasma.NewBuiltInFunction(
			result.vtable,
			func(argument ...*Value) (*Value, error) {
				return result, nil
			},
		))
	result.Set(magic_functions.Float,
		plasma.NewBuiltInFunction(
			result.vtable,
			func(argument ...*Value) (*Value, error) {
				return plasma.NewFloat(result.Float()), nil
			},
		))
	result.Set(magic_functions.Copy,
		plasma.NewBuiltInFunction(
			result.vtable,
			func(argument ...*Value) (*Value, error) {
				return plasma.NewInt(result.GetInt64()), nil
			},
		))
	result.Set(magic_functions.BigEndian, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			b := make([]byte, 8)
			binary.BigEndian.PutUint64(b, uint64(result.Int()))
			return plasma.NewBytes(b), nil
		},
	))
	result.Set(magic_functions.LittleEndian, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			b := make([]byte, 8)
			binary.LittleEndian.PutUint64(b, uint64(result.Int()))
			return plasma.NewBytes(b), nil
		},
	))
	result.Set(magic_functions.FromBig, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewInt(int64(binary.BigEndian.Uint64(argument[0].GetBytes()))), nil
		},
	))
	result.Set(magic_functions.FromLittle, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewInt(int64(binary.LittleEndian.Uint64(argument[0].GetBytes()))), nil
		},
	))
	return result
}
