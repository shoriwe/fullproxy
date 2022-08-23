package vm

import "fmt"

func (plasma *Plasma) printStack(ctx *context) {
	current := ctx.stack.Top
	for current != nil {
		fmt.Println("\t", current.Value.(*Value).TypeId(), current.Value)
		current = current.Next
	}
}
