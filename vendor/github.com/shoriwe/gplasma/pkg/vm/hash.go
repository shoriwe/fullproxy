package vm

import magic_functions "github.com/shoriwe/gplasma/pkg/common/magic-functions"

func (plasma *Plasma) hashClass() *Value {
	class := plasma.NewValue(plasma.rootSymbols, BuiltInClassId, plasma.class)
	class.SetAny(Callback(func(argument ...*Value) (*Value, error) {
		return plasma.NewHash(argument[0].GetHash()), nil
	}))
	return class
}

/*
NewHash magic function:
In                  __in__
Length              __len__
Bool                __bool__
String              __string__
Bytes               __bytes__
Get                 __get__
Set                 __set__
Del				 	__del__
Copy                __copy__
*/
func (plasma *Plasma) NewHash(hash *Hash) *Value {
	result := plasma.NewValue(plasma.rootSymbols, HashId, plasma.hash)
	result.SetAny(hash)
	result.Set(magic_functions.In, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			in, inError := result.GetHash().In(argument[0])
			return plasma.NewBool(in), inError
		},
	))
	result.Set(magic_functions.Length, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewInt(result.GetHash().Size()), nil
		},
	))
	result.Set(magic_functions.Bool, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewBool(result.Bool()), nil
		},
	))
	result.Set(magic_functions.String, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewString([]byte(result.String())), nil
		},
	))
	result.Set(magic_functions.Bytes, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewBytes([]byte(result.String())), nil
		},
	))
	result.Set(magic_functions.Get, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return result.GetHash().Get(argument[0])
		},
	))
	result.Set(magic_functions.Set, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.none, result.GetHash().Set(argument[0], argument[1])
		},
	))
	result.Set(magic_functions.Del, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.none, result.GetHash().Del(argument[0])
		},
	))
	result.Set(magic_functions.Copy, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewHash(result.GetHash().Copy()), nil
		},
	))
	return result
}
