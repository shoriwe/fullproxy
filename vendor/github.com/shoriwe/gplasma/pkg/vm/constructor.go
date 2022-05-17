package vm

func (p *Plasma) constructSubClass(context *Context, subClass *Value, object *Value) *Value {
	for _, subSubClass := range subClass.subClasses {
		object.SymbolTable().Parent = subSubClass.symbols.Parent
		subSubClassConstructionError := p.constructSubClass(context, subSubClass, object)
		if subSubClassConstructionError != nil {
			return subSubClassConstructionError
		}
	}
	object.SymbolTable().Parent = subClass.symbols.Parent
	baseInitializationError := subClass.Constructor.Construct(context, p, object)
	if baseInitializationError != nil {
		return baseInitializationError
	}
	return nil
}

func (p *Plasma) ConstructObject(context *Context, type_ *Value, parent *SymbolTable) (*Value, bool) {
	object := p.NewValue(context, false, type_.Name, type_.subClasses, parent)
	for _, subclass := range object.subClasses {
		subClassConstructionError := p.constructSubClass(context, subclass, object)
		if subClassConstructionError != nil {
			return subClassConstructionError, false
		}
	}
	object.SymbolTable().Parent = parent
	object.class = type_
	baseInitializationError := type_.Constructor.Construct(context, p, object)
	if baseInitializationError != nil {
		return baseInitializationError, false
	}
	return object, true
}

type Constructor interface {
	Construct(*Context, *Plasma, *Value) *Value
}

type PlasmaConstructor struct {
	Constructor
	Code []*Code
}

func (c *PlasmaConstructor) Construct(context *Context, vm *Plasma, object *Value) *Value {
	context.PushSymbolTable(object.SymbolTable())
	context.PushObject(object)
	executionError, success := vm.Execute(context, NewBytecodeFromArray(c.Code))
	context.PopSymbolTable()
	context.PopObject()
	if !success {
		return executionError
	}
	return nil
}

func NewPlasmaConstructor(code []*Code) *PlasmaConstructor {
	return &PlasmaConstructor{
		Code: code,
	}
}

type ConstructorCallBack func(*Context, *Value) *Value

type BuiltInConstructor struct {
	Constructor
	callback ConstructorCallBack
}

func (c *BuiltInConstructor) Construct(context *Context, _ *Plasma, object *Value) *Value {
	return c.callback(context, object)
}

func NewBuiltInConstructor(callback ConstructorCallBack) *BuiltInConstructor {
	return &BuiltInConstructor{
		callback: callback,
	}
}
