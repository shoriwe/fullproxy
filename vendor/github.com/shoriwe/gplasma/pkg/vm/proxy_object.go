package vm

func (p *Plasma) NewProxyObject(context *Context, self, class *Value) (*Value, bool) {
	symbolTable := NewSymbolTable(context.PeekSymbolTable())
	proxyObject, success := p.ConstructObject(context, class, symbolTable)
	if !success {
		return proxyObject, false
	}
	for _, value := range proxyObject.SymbolTable().Symbols {
		if value.IsTypeById(FunctionId) {
			switch value.Callable.(type) {
			case *BuiltInClassFunction:
				value.Callable.(*BuiltInClassFunction).Self = self
			case *PlasmaClassFunction:
				value.Callable.(*PlasmaClassFunction).Self = self
			}
		}
	}
	return proxyObject, true
}
