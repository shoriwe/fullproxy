package vm

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/bytecode/opcodes"
	"github.com/shoriwe/gplasma/pkg/common"
	magic_functions "github.com/shoriwe/gplasma/pkg/common/magic-functions"
	special_symbols "github.com/shoriwe/gplasma/pkg/common/special-symbols"
)

func (ctx *context) pushCode(bytecode []byte) {
	ctx.code.Push(
		&contextCode{
			bytecode: bytecode,
			rip:      0,
			onExit:   &common.ListStack[[]byte]{},
		},
	)
}
func (ctx *context) popCode() {
	// If there is defer code
	if ctxCode := ctx.code.Peek(); ctxCode.onExit.HasNext() {
		ctxCode.rip = int64(len(ctxCode.bytecode)) + 1
		if ctx.register != nil {
			ctx.stack.Push(ctx.register)
			ctx.pushCode([]byte{opcodes.Return})
			ctx.currentSymbols = NewSymbols(ctx.currentSymbols)
		}
		for ctxCode.onExit.HasNext() {
			ctx.pushCode(ctxCode.onExit.Pop())
			ctx.currentSymbols = NewSymbols(ctx.currentSymbols)
		}
		return
	}
	ctx.code.Pop()
	if ctx.currentSymbols.call != nil {
		ctx.currentSymbols = ctx.currentSymbols.call
	} else {
		ctx.currentSymbols = ctx.currentSymbols.Parent
	}
}

func (plasma *Plasma) prepareClassInitCode(classInfo *ClassInfo) {
	var result []byte
	for _, base := range classInfo.Bases {
		if base.TypeId() != ClassId {
			panic("no type received as base for class")
		}
		baseClassInfo := base.GetClassInfo()
		if !baseClassInfo.prepared {
			plasma.prepareClassInitCode(baseClassInfo)
		}
		result = append(result, baseClassInfo.Bytecode...)
	}
	result = append(result, classInfo.Bytecode...)
	classInfo.prepared = true
	classInfo.Bytecode = result
}

func (plasma *Plasma) do(ctx *context) {
	ctxCode := ctx.code.Peek()
	instruction := ctxCode.bytecode[ctxCode.rip]
	// fmt.Println(opcodes.OpCodes[instruction])
	// plasma.printStack(ctx)
	switch instruction {
	case opcodes.Push:
		ctxCode.rip++
		ctx.stack.Push(ctx.register)
	case opcodes.Pop:
		ctxCode.rip++
		ctx.register = ctx.stack.Pop()
	case opcodes.IdentifierAssign:
		ctxCode.rip++
		symbolLength := common.BytesToInt(ctxCode.bytecode[ctxCode.rip : ctxCode.rip+8])
		ctxCode.rip += 8
		symbol := string(ctxCode.bytecode[ctxCode.rip : ctxCode.rip+symbolLength])
		// fmt.Println(symbol)
		ctxCode.rip += symbolLength
		ctx.currentSymbols.Set(symbol, ctx.stack.Pop())
	case opcodes.SelectorAssign:
		ctxCode.rip++
		symbolLength := common.BytesToInt(ctxCode.bytecode[ctxCode.rip : ctxCode.rip+8])
		ctxCode.rip += 8
		symbol := string(ctxCode.bytecode[ctxCode.rip : ctxCode.rip+symbolLength])
		ctxCode.rip += symbolLength
		selector := ctx.stack.Pop()
		selector.Set(symbol, ctx.stack.Pop())
	case opcodes.Label:
		ctxCode.rip += 9 // OP + Label
	case opcodes.Jump:
		ctxCode.rip += common.BytesToInt(ctxCode.bytecode[1+ctxCode.rip : 9+ctxCode.rip])
	case opcodes.IfJump:
		if ctx.stack.Pop().Bool() {
			ctxCode.rip += common.BytesToInt(ctxCode.bytecode[1+ctxCode.rip : 9+ctxCode.rip])
		} else {
			ctxCode.rip += 9
		}
	case opcodes.Return:
		ctxCode.rip++
		ctx.register = ctx.stack.Pop()
		ctx.popCode()
	case opcodes.DeleteIdentifier:
		ctxCode.rip++
		symbolLength := common.BytesToInt(ctxCode.bytecode[ctxCode.rip : ctxCode.rip+8])
		ctxCode.rip += 8
		symbol := string(ctxCode.bytecode[ctxCode.rip : ctxCode.rip+symbolLength])
		ctxCode.rip += symbolLength
		delError := ctx.currentSymbols.Del(symbol)
		if delError != nil {
			panic(delError)
		}
	case opcodes.DeleteSelector:
		ctxCode.rip++
		symbolLength := common.BytesToInt(ctxCode.bytecode[ctxCode.rip : ctxCode.rip+8])
		ctxCode.rip += 8
		symbol := string(ctxCode.bytecode[ctxCode.rip : ctxCode.rip+symbolLength])
		ctxCode.rip += symbolLength
		selector := ctx.stack.Pop()
		delError := selector.Del(symbol)
		if delError != nil {
			panic(delError)
		}
	case opcodes.Defer:
		ctxCode.rip++
		exprLength := common.BytesToInt(ctxCode.bytecode[ctxCode.rip : ctxCode.rip+8])
		ctxCode.rip += 8
		onExitCode := ctxCode.bytecode[ctxCode.rip : ctxCode.rip+exprLength]
		ctxCode.rip += exprLength
		ctxCode.onExit.Push(onExitCode)
	case opcodes.NewFunction:
		ctxCode.rip++
		numberOfArgument := common.BytesToInt(ctxCode.bytecode[ctxCode.rip : ctxCode.rip+8])
		ctxCode.rip += 8
		arguments := make([]string, 0, numberOfArgument)
		for i := int64(0); i < numberOfArgument; i++ {
			symbolLength := common.BytesToInt(ctxCode.bytecode[ctxCode.rip : ctxCode.rip+8])
			ctxCode.rip += 8
			symbol := string(ctxCode.bytecode[ctxCode.rip : ctxCode.rip+symbolLength])
			ctxCode.rip += symbolLength
			arguments = append(arguments, symbol)
		}
		bytecodeLength := common.BytesToInt(ctxCode.bytecode[ctxCode.rip : ctxCode.rip+8])
		ctxCode.rip += 8
		bytecode := ctxCode.bytecode[ctxCode.rip : ctxCode.rip+bytecodeLength]
		ctxCode.rip += bytecodeLength
		funcInfo := FuncInfo{
			Arguments: arguments,
			Bytecode:  bytecode,
		}
		funcObject := plasma.NewValue(ctx.currentSymbols, FunctionId, plasma.function)
		funcObject.SetAny(funcInfo)
		ctx.register = funcObject
	case opcodes.NewClass:
		ctxCode.rip++
		numberOfBases := common.BytesToInt(ctxCode.bytecode[ctxCode.rip : ctxCode.rip+8])
		ctxCode.rip += 8
		bodyLength := common.BytesToInt(ctxCode.bytecode[ctxCode.rip : ctxCode.rip+8])
		ctxCode.rip += 8
		body := ctxCode.bytecode[ctxCode.rip : ctxCode.rip+bodyLength]
		ctxCode.rip += bodyLength
		// Get bases
		bases := make([]*Value, numberOfBases)
		for i := numberOfBases - 1; i >= 0; i-- {
			bases[i] = ctx.stack.Pop()
		}
		classInfo := &ClassInfo{
			Bases:    bases,
			Bytecode: body,
		}
		classObject := plasma.NewValue(ctx.currentSymbols, ClassId, plasma.class)
		classObject.SetAny(classInfo)
		ctx.register = classObject
	case opcodes.Call:
		ctxCode.rip++
		numberOfArguments := common.BytesToInt(ctxCode.bytecode[ctxCode.rip : ctxCode.rip+8])
		ctxCode.rip += 8
		function := ctx.stack.Pop()
		arguments := make([]*Value, numberOfArguments)
		for i := numberOfArguments - 1; i >= 0; i-- {
			arguments[i] = ctx.stack.Pop()
		}
		var callError error
		tries := 0
	doCall:
		if tries == MaxDoCallSearch {
			panic("infinite nested __call__")
		}
		switch function.TypeId() {
		case BuiltInFunctionId, BuiltInClassId:
			ctx.register, callError = function.Call(arguments...)
			if callError != nil {
				panic(callError)
			}
		case FunctionId:
			funcInfo := function.GetFuncInfo()
			// Push new symbol table based on the function
			newSymbols := NewSymbols(function.vtable)
			newSymbols.call = ctx.currentSymbols
			ctx.currentSymbols = newSymbols
			if int64(len(funcInfo.Arguments)) != numberOfArguments {
				panic("invalid number of argument for function call")
			}
			// Load arguments
			for index, argument := range funcInfo.Arguments {
				ctx.currentSymbols.Set(argument, arguments[index])
			}
			// Push code
			ctx.pushCode(funcInfo.Bytecode)
		case ClassId:
			classInfo := function.GetClassInfo()
			if !classInfo.prepared {
				plasma.prepareClassInitCode(classInfo)
			}
			// Instantiate object
			object := plasma.NewValue(function.vtable, ValueId, plasma.value)
			object.class = function
			object.Set(special_symbols.Self, object)
			// Push object
			ctx.stack.Push(object)
			for _, argument := range arguments {
				ctx.stack.Push(argument)
			}
			// Push class code
			classCode := make([]byte, 0, len(classInfo.Bytecode))
			classCode = append(classCode, classInfo.Bytecode...)
			// inject init code: object.__init__(arguments...)
			classCode = append(classCode, opcodes.Identifier)
			classCode = append(classCode, common.IntToBytes(len(magic_functions.Init))...)
			classCode = append(classCode, magic_functions.Init...)
			classCode = append(classCode, opcodes.Push)
			classCode = append(classCode, opcodes.Call)
			classCode = append(classCode, common.IntToBytes(numberOfArguments)...)
			// Inject pop object to register
			classCode = append(classCode, opcodes.Pop)
			// Load code
			ctx.pushCode(classCode)
			newSymbols := object.vtable
			newSymbols.call = ctx.currentSymbols
			ctx.currentSymbols = newSymbols
		default: // __call__
			call, getError := function.Get(magic_functions.Call)
			if getError != nil {
				panic(getError)
			}
			function = call
			tries++
			goto doCall
		}
	case opcodes.NewArray:
		ctxCode.rip++
		numberOfValues := common.BytesToInt(ctxCode.bytecode[ctxCode.rip : ctxCode.rip+8])
		ctxCode.rip += 8
		values := make([]*Value, numberOfValues)
		for i := numberOfValues - 1; i >= 0; i-- {
			values[i] = ctx.stack.Pop()
		}
		ctx.register = plasma.NewArray(values)
	case opcodes.NewTuple:
		ctxCode.rip++
		numberOfValues := common.BytesToInt(ctxCode.bytecode[ctxCode.rip : ctxCode.rip+8])
		ctxCode.rip += 8
		values := make([]*Value, numberOfValues)
		for i := numberOfValues - 1; i >= 0; i-- {
			values[i] = ctx.stack.Pop()
		}
		ctx.register = plasma.NewTuple(values)
	case opcodes.NewHash:
		ctxCode.rip++
		numberOfValues := common.BytesToInt(ctxCode.bytecode[ctxCode.rip : ctxCode.rip+8])
		ctxCode.rip += 8
		hash := plasma.NewInternalHash()
		for i := numberOfValues - 1; i >= 0; i-- {
			key := ctx.stack.Pop()
			value := ctx.stack.Pop()
			setError := hash.Set(key, value)
			if setError != nil {
				panic(setError)
			}
		}
		ctx.register = plasma.NewHash(hash)
	case opcodes.Identifier:
		ctxCode.rip++
		symbolLength := common.BytesToInt(ctxCode.bytecode[ctxCode.rip : ctxCode.rip+8])
		ctxCode.rip += 8
		symbol := string(ctxCode.bytecode[ctxCode.rip : ctxCode.rip+symbolLength])
		ctxCode.rip += symbolLength
		//fmt.Println(symbol)
		var getError error
		ctx.register, getError = ctx.currentSymbols.Get(symbol)
		if getError != nil {
			panic(getError)
		}
	case opcodes.Integer:
		ctxCode.rip++
		value := common.BytesToInt(ctxCode.bytecode[ctxCode.rip : ctxCode.rip+8])
		ctxCode.rip += 8
		ctx.register = plasma.NewInt(value)
	case opcodes.Float:
		ctxCode.rip++
		value := common.BytesToFloat(ctxCode.bytecode[ctxCode.rip : ctxCode.rip+8])
		ctxCode.rip += 8
		ctx.register = plasma.NewFloat(value)
	case opcodes.String:
		ctxCode.rip++
		stringLength := common.BytesToInt(ctxCode.bytecode[ctxCode.rip : ctxCode.rip+8])
		ctxCode.rip += 8
		contents := ctxCode.bytecode[ctxCode.rip : ctxCode.rip+stringLength]
		ctxCode.rip += stringLength
		ctx.register = plasma.NewString(contents)
	case opcodes.Bytes:
		ctxCode.rip++
		stringLength := common.BytesToInt(ctxCode.bytecode[ctxCode.rip : ctxCode.rip+8])
		ctxCode.rip += 8
		contents := ctxCode.bytecode[ctxCode.rip : ctxCode.rip+stringLength]
		ctxCode.rip += stringLength
		ctx.register = plasma.NewBytes(contents)
	case opcodes.True:
		ctxCode.rip++
		ctx.register = plasma.true
	case opcodes.False:
		ctxCode.rip++
		ctx.register = plasma.false
	case opcodes.None:
		ctxCode.rip++
		ctx.register = plasma.none
	case opcodes.Selector:
		ctxCode.rip++
		symbolLength := common.BytesToInt(ctxCode.bytecode[ctxCode.rip : ctxCode.rip+8])
		ctxCode.rip += 8
		symbol := string(ctxCode.bytecode[ctxCode.rip : ctxCode.rip+symbolLength])
		ctxCode.rip += symbolLength
		// fmt.Println(symbol)
		selector := ctx.stack.Pop()
		var getError error
		ctx.register, getError = selector.Get(symbol)
		if getError != nil {
			panic(getError)
		}
	case opcodes.Super:
		break // TODO: Implement me!
	default:
		panic(fmt.Sprintf("unknown opcode %d", instruction))
	}
}
