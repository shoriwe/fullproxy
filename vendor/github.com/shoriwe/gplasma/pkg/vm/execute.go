package vm

func (p *Plasma) Execute(context *Context, bytecode *Bytecode) (*Value, bool) {
	if context == nil {
		context = p.NewContext()
	}
	var executionError *Value
	for bytecode.HasNext() {
		code := bytecode.Next()

		/*
			if code.Line != 0 {
				fmt.Println(color.GreenString(strconv.Itoa(code.Line)), instructionNames[code.Instruction.OpCode], code.Value)
			} else {
				fmt.Println(color.RedString("UL"), instructionNames[code.Instruction.OpCode], code.Value)
			}
			if context.ObjectStack.head != nil {
				current := context.ObjectStack.head
				for ; current != nil; current = current.next {
					fmt.Println(current.value.(*Value).GetClass(p).Name, current.value.(*Value).Integer)
				}
			}
		*/

		switch code.Instruction.OpCode {
		case GetFalseOP:
			context.LastObject = p.GetFalse()
		case GetTrueOP:
			context.LastObject = p.GetTrue()
		case GetNoneOP:
			context.LastObject = p.GetNone()
		case NewStringOP:
			executionError = p.newStringOP(context, code.Value.(string))
		case NewBytesOP:
			executionError = p.newBytesOP(context, code.Value.([]uint8))
		case NewIntegerOP:
			executionError = p.newIntegerOP(context, code.Value.(int64))
		case NewFloatOP:
			executionError = p.newFloatOP(context, code.Value.(float64))
		case NewArrayOP:
			executionError = p.newArrayOP(context, code.Value.(int))
		case NewTupleOP:
			executionError = p.newTupleOP(context, code.Value.(int))
		case NewHashOP:
			executionError = p.newHashTableOP(context, code.Value.(int))
		case UnaryOP:
			executionError = p.unaryOP(context, code.Value.(uint8))
		case BinaryOP:
			executionError = p.binaryOP(context, code.Value.(uint8))
		case MethodInvocationOP:
			executionError = p.methodInvocationOP(context, code.Value.(int))
		case GetIdentifierOP:
			executionError = p.getIdentifierOP(context, code.Value.(string))
		case SelectNameFromObjectOP:
			executionError = p.selectNameFromObjectOP(context, code.Value.(string))
		case IndexOP:
			executionError = p.indexOP(context)
		case PushOP:
			executionError = p.pushOP(context)
		case AssignIdentifierOP:
			executionError = p.assignIdentifierOP(context, code.Value.(string))
		case NewClassOP:
			executionError = p.newClassOP(context, code.Value.(ClassInformation))
		case NewClassFunctionOP:
			executionError = p.newClassFunctionOP(context, code.Value.(FunctionInformation))
		case NewFunctionOP:
			executionError = p.newFunctionOP(context, code.Value.(FunctionInformation))
		case LoadFunctionArgumentsOP:
			executionError = p.loadFunctionArgumentsOP(context, code.Value.([]string))
		case ReturnOP:
			returnResult := p.returnOP(context, code.Value.(int))
			context.ReturnState()
			return returnResult, true
		case IfOneLinerOP:
			executionError = p.ifOneLinerOP(context, code.Value.(ConditionInformation))
		case UnlessOneLinerOP:
			executionError = p.unlessOneLinerOP(context, code.Value.(ConditionInformation))
		case AssignSelectorOP:
			executionError = p.assignSelectorOP(context, code.Value.(string))
		case AssignIndexOP:
			executionError = p.assignIndexOP(context)
		case IfOP:
			executionError = p.ifOP(context, code.Value.(ConditionInformation))
			if executionError != nil {
				return executionError, false
			}
			switch context.LastState {
			case BreakState, RedoState, ContinueState:
				return p.GetNone(), true
			case ReturnState:
				result := context.LastObject
				context.LastObject = nil
				return result, true
			}
		case UnlessOP:
			executionError = p.unlessOP(context, code.Value.(ConditionInformation))
			if executionError != nil {
				return executionError, false
			}
			switch context.LastState {
			case BreakState, RedoState, ContinueState:
				return p.GetNone(), true
			case ReturnState:
				result := context.LastObject
				context.LastObject = nil
				return result, true
			}
		case ForLoopOP:
			executionError = p.forLoopOP(context, code.Value.(LoopInformation))
			if executionError != nil {
				return executionError, false
			} else if context.LastState == ReturnState {
				result := context.LastObject
				context.LastObject = nil
				return result, true
			}
		case NewGeneratorOP:
			executionError = p.newGeneratorOP(context, code.Value.(int))
		case WhileLoopOP:
			executionError = p.whileLoopOP(context, code.Value.(LoopInformation))
			if executionError != nil {
				return executionError, false
			} else if context.LastState == ReturnState {
				result := context.LastObject
				context.LastObject = nil
				return result, true
			}
		case DoWhileLoopOP:
			executionError = p.doWhileLoopOP(context, code.Value.(LoopInformation))
			if executionError != nil {
				return executionError, false
			} else if context.LastState == ReturnState {
				result := context.LastObject
				context.LastObject = nil
				return result, true
			}
		case UntilLoopOP:
			executionError = p.untilLoopOP(context, code.Value.(LoopInformation))
			if executionError != nil {
				return executionError, false
			} else if context.LastState == ReturnState {
				result := context.LastObject
				context.LastObject = nil
				return result, true
			}
		case BreakOP:
			context.BreakState()
			return p.GetNone(), true
		case ContinueOP:
			context.ContinueState()
			return p.GetNone(), true
		case RedoOP:
			context.RedoState()
			return p.GetNone(), true
		case NOP:
			break
		case TryOP:
			executionError = p.tryOP(context, code.Value.(TryInformation))
			if executionError != nil {
				return executionError, false
			}
			switch context.LastState {
			case BreakState, RedoState, ContinueState:
				return p.GetNone(), true
			case ReturnState:
				result := context.LastObject
				context.LastObject = nil
				return result, true
			}
		case RaiseOP:
			executionError = context.PopObject()
		case NewModuleOP:
			executionError = p.newModuleOP(context, code.Value.(ClassInformation))
		case SwitchOP:
			executionError = p.switchOP(context, code.Value.(SwitchInformation))
			if executionError != nil {
				return executionError, false
			}
			switch context.LastState {
			case BreakState, RedoState, ContinueState:
				return p.GetNone(), true
			case ReturnState:
				result := context.LastObject
				context.LastObject = nil
				return result, true
			}
		default:
			panic(instructionNames[code.Instruction.OpCode])
		}
		if executionError != nil {
			return executionError, false
		}
	}
	context.NoState()
	return p.GetNone(), true
}
