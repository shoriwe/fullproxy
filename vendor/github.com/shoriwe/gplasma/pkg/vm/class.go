package vm

import magic_functions "github.com/shoriwe/gplasma/pkg/common/magic-functions"

func (plasma *Plasma) metaClass() *Value {
	plasma.class = plasma.NewValue(plasma.rootSymbols, BuiltInClassId, plasma.class)
	plasma.class.class = plasma.class
	plasma.class.SetAny(Callback(func(argument ...*Value) (*Value, error) {
		return plasma.NewClass(), nil
	}))
	return plasma.class
}

/*
NewClass magic function:
Equal              __equal__
NotEqual            __not_equal__
*/
func (plasma *Plasma) NewClass() *Value {
	result := plasma.NewValue(plasma.rootSymbols, BuiltInClassId, plasma.class)
	result.Set(magic_functions.Equal, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewBool(result == argument[0]), nil
		},
	))
	result.Set(magic_functions.NotEqual, plasma.NewBuiltInFunction(
		result.vtable,
		func(argument ...*Value) (*Value, error) {
			return plasma.NewBool(result != argument[0]), nil
		},
	))
	return result
}
