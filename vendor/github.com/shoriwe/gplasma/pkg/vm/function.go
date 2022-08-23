package vm

func (plasma *Plasma) functionClass() *Value {
	class := plasma.NewValue(plasma.rootSymbols, BuiltInClassId, plasma.class)
	class.SetAny(Callback(func(argument ...*Value) (*Value, error) {
		return plasma.NewBuiltInFunction(
			plasma.rootSymbols,
			func(argument ...*Value) (*Value, error) {
				return plasma.none, nil
			}), nil
	}))
	return class
}

func (plasma *Plasma) NewBuiltInFunction(parent *Symbols, callback Callback) *Value {
	function := plasma.NewValue(parent, BuiltInFunctionId, plasma.function)
	function.SetAny(callback)
	return function
}
