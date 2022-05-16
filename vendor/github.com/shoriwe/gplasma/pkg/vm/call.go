package vm

func (p *Plasma) CallableInitialize(isBuiltIn bool) ConstructorCallBack {
	return func(context *Context, object *Value) *Value {
		object.SetOnDemandSymbol(Call,
			func() *Value {
				return p.NewFunction(context, isBuiltIn, object.SymbolTable(),
					NewBuiltInClassFunction(object, 0,
						func(_ *Value, _ ...*Value) (*Value, bool) {
							return p.NewNotImplementedCallableError(context, Call), false
						},
					),
				)
			},
		)
		return nil
	}
}

func (p *Plasma) CallFunction(context *Context, function *Value, arguments ...*Value) (*Value, bool) {
	var (
		callFunction *Value
		result       *Value
		callError    *Value
		success      bool
	)
	switch function.BuiltInTypeId {
	case FunctionId:
		callFunction = function
	case TypeId:
		result, success = p.ConstructObject(context, function, function.SymbolTable().Parent)
		if !success {
			return result, false
		}
		resultInitialize, getError := result.Get(p, context, Initialize)
		if getError != nil {
			return getError, false
		}
		callError, success = p.CallFunction(context, resultInitialize, arguments...)
		// Construct the object and initialize it
		if !success {
			return callError, false
		}
		return result, true
	default:
		call, getError := function.Get(p, context, Call)
		if getError != nil {
			return getError, false
		}
		if !call.IsTypeById(FunctionId) {
			return p.NewInvalidTypeError(context, function.TypeName(), CallableName), false
		}
		callFunction = call
	}
	if callFunction.Callable.NumberOfArguments() != len(arguments) {
		//  Return Here a error related to number of arguments
		return p.NewInvalidNumberOfArgumentsError(context, len(arguments), callFunction.Callable.NumberOfArguments()), false
	}
	symbols := NewSymbolTable(function.SymbolTable().Parent)
	self, callback, code := callFunction.Callable.Call()
	if self != nil {
		symbols.Set(Self, self)
	} else {
		symbols.Set(Self, function)
	}
	context.PushSymbolTable(symbols)
	if callback != nil {
		result, success = callback(self, arguments...)
	} else if code != nil {
		// Load the arguments
		for i := len(arguments) - 1; i > -1; i-- {
			context.PushObject(arguments[i])
		}
		result, success = p.Execute(context, NewBytecodeFromArray(code))
	} else {
		panic("callback and code are nil")
	}
	context.PopSymbolTable()
	if !success {
		return result, false
	}
	return result, true
}

type FunctionCallback func(*Value, ...*Value) (*Value, bool)

type Callable interface {
	NumberOfArguments() int
	Call() (*Value, FunctionCallback, []*Code) // self should return directly the object or the code of the function
}

type PlasmaFunction struct {
	numberOfArguments int
	Code              []*Code
}

func (p *PlasmaFunction) NumberOfArguments() int {
	return p.numberOfArguments
}

func (p *PlasmaFunction) Call() (*Value, FunctionCallback, []*Code) {
	return nil, nil, p.Code
}

func NewPlasmaFunction(numberOfArguments int, code []*Code) *PlasmaFunction {
	return &PlasmaFunction{
		numberOfArguments: numberOfArguments,
		Code:              code,
	}
}

type PlasmaClassFunction struct {
	numberOfArguments int
	Code              []*Code
	Self              *Value
}

func (p *PlasmaClassFunction) NumberOfArguments() int {
	return p.numberOfArguments
}

func (p *PlasmaClassFunction) Call() (*Value, FunctionCallback, []*Code) {
	return p.Self, nil, p.Code
}

func NewPlasmaClassFunction(self *Value, numberOfArguments int, code []*Code) *PlasmaClassFunction {
	return &PlasmaClassFunction{
		numberOfArguments: numberOfArguments,
		Code:              code,
		Self:              self,
	}
}

type BuiltInFunction struct {
	numberOfArguments int
	callback          FunctionCallback
}

func (g *BuiltInFunction) NumberOfArguments() int {
	return g.numberOfArguments
}

func (g *BuiltInFunction) Call() (*Value, FunctionCallback, []*Code) {
	return nil, g.callback, nil
}

func NewBuiltInFunction(numberOfArguments int, callback FunctionCallback) *BuiltInFunction {
	return &BuiltInFunction{
		numberOfArguments: numberOfArguments,
		callback:          callback,
	}
}

type BuiltInClassFunction struct {
	numberOfArguments int
	callback          FunctionCallback
	Self              *Value
}

func (g *BuiltInClassFunction) NumberOfArguments() int {
	return g.numberOfArguments
}

func (g *BuiltInClassFunction) Call() (*Value, FunctionCallback, []*Code) {
	return g.Self, g.callback, nil
}

func NewBuiltInClassFunction(self *Value, numberOfArguments int, callback FunctionCallback) *BuiltInClassFunction {
	return &BuiltInClassFunction{
		numberOfArguments: numberOfArguments,
		callback:          callback,
		Self:              self,
	}
}
