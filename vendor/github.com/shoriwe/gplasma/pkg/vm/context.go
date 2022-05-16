package vm

type Context struct {
	ObjectStack *ObjectStack
	SymbolStack *SymbolStack
	LastObject  *Value
	LastState   uint8
}

func (c *Context) PushObject(object *Value) {
	c.ObjectStack.Push(object)
}
func (c *Context) PeekObject() *Value {
	return c.ObjectStack.Peek()
}

func (c *Context) PopObject() *Value {
	return c.ObjectStack.Pop()
}

func (c *Context) PushSymbolTable(table *SymbolTable) {
	c.SymbolStack.Push(table)
}

func (c *Context) PopSymbolTable() *SymbolTable {
	return c.SymbolStack.Pop()
}

func (c *Context) PeekSymbolTable() *SymbolTable {
	return c.SymbolStack.Peek()
}

func (c *Context) ReturnState() {
	c.LastState = ReturnState
}

func (c *Context) BreakState() {
	c.LastState = BreakState
}

func (c *Context) RedoState() {
	c.LastState = RedoState
}

func (c *Context) ContinueState() {
	c.LastState = ContinueState
}

func (c *Context) NoState() {
	c.LastState = NoState
}
