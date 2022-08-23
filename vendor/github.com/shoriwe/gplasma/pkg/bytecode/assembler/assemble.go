package assembler

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast3"
	"github.com/shoriwe/gplasma/pkg/bytecode/opcodes"
	"github.com/shoriwe/gplasma/pkg/common"
	"reflect"
)

type (
	assembler struct{}
)

func newAssembler() *assembler {
	return &assembler{}
}

func (a *assembler) assemble(node ast3.Node) []byte {
	if node == nil {
		return nil
	}
	switch n := node.(type) {
	case ast3.Statement:
		return a.Statement(n)
	case ast3.Expression:
		return a.Expression(n)
	default:
		panic(fmt.Sprintf("unknown node type %s", reflect.TypeOf(n).String()))
	}
}

func (a *assembler) enumLabels(bytecode []byte) map[int64]int64 {
	bytecodeLength := int64(len(bytecode))
	labels := map[int64]int64{}
	for index := int64(0); index < bytecodeLength; {
		op := bytecode[index]
		switch op {
		case opcodes.Push:
			index++
		case opcodes.Pop:
			index++
		case opcodes.IdentifierAssign:
			index++
			symbolLength := common.BytesToInt(bytecode[index : index+8])
			index += 8 + symbolLength
		case opcodes.SelectorAssign:
			index++
			symbolLength := common.BytesToInt(bytecode[index : index+8])
			index += 8 + symbolLength
		case opcodes.Label:
			index++
			labelCode := common.BytesToInt(bytecode[index : index+8])
			labels[labelCode] = index - 1
			index += 8
		case opcodes.Jump:
			index++
			index += 8
		case opcodes.IfJump:
			index++
			index += 8
		case opcodes.Return:
			index++
		case opcodes.DeleteIdentifier:
			index++
			symbolLength := common.BytesToInt(bytecode[index : index+8])
			index += 8 + symbolLength
		case opcodes.DeleteSelector:
			index++
			symbolLength := common.BytesToInt(bytecode[index : index+8])
			index += 8 + symbolLength
		case opcodes.Defer:
			index++
			index += 8
		case opcodes.NewFunction:
			index++
			argsNumber := common.BytesToInt(bytecode[index : index+8])
			index += 8
			for arg := int64(0); arg < argsNumber; arg++ {
				argSymbolLength := common.BytesToInt(bytecode[index : index+8])
				index += 8 + argSymbolLength
			}
			index += 8
		case opcodes.NewClass:
			index++
			index += 8
			index += 8
		case opcodes.Call:
			index++
			index += 8
		case opcodes.NewArray:
			index++
			index += 8
		case opcodes.NewTuple:
			index++
			index += 8
		case opcodes.NewHash:
			index++
			index += 8
		case opcodes.Identifier:
			index++
			symbolLength := common.BytesToInt(bytecode[index : index+8])
			index += 8 + symbolLength
		case opcodes.Integer:
			index++
			index += 8
		case opcodes.Float:
			index++
			index += 8
		case opcodes.String:
			index++
			stringLength := common.BytesToInt(bytecode[index : index+8])
			index += 8 + stringLength
		case opcodes.Bytes:
			index++
			bytesLength := common.BytesToInt(bytecode[index : index+8])
			index += 8 + bytesLength
		case opcodes.True:
			index++
		case opcodes.False:
			index++
		case opcodes.None:
			index++
		case opcodes.Selector:
			index++
			symbolLength := common.BytesToInt(bytecode[index : index+8])
			index += 8 + symbolLength
		case opcodes.Super:
			index++
		default:
			panic(fmt.Sprintf("unknown opcode %d in %v", op, bytecode[index-5:]))
		}
	}
	return labels
}

func (a *assembler) resolveLabels(bytecode []byte, labels map[int64]int64) []byte {
	bytecodeLength := int64(len(bytecode))
	for index := int64(0); index < bytecodeLength; {
		op := bytecode[index]
		switch op {
		case opcodes.Push:
			index++
		case opcodes.Pop:
			index++
		case opcodes.IdentifierAssign:
			index++
			symbolLength := common.BytesToInt(bytecode[index : index+8])
			index += 8 + symbolLength
		case opcodes.SelectorAssign:
			index++
			symbolLength := common.BytesToInt(bytecode[index : index+8])
			index += 8 + symbolLength
		case opcodes.Label:
			index++
			index += 8
		case opcodes.Jump:
			labelCode := common.BytesToInt(bytecode[index+1 : index+9])
			jump := labels[labelCode] - index
			index++
			copy(bytecode[index:index+8], common.IntToBytes(jump))
			index += 8
		case opcodes.IfJump:
			labelCode := common.BytesToInt(bytecode[index+1 : index+9])
			jump := labels[labelCode] - index
			index++
			copy(bytecode[index:index+8], common.IntToBytes(jump))
			index += 8
		case opcodes.Return:
			index++
		case opcodes.DeleteIdentifier:
			index++
			symbolLength := common.BytesToInt(bytecode[index : index+8])
			index += 8 + symbolLength
		case opcodes.DeleteSelector:
			index++
			symbolLength := common.BytesToInt(bytecode[index : index+8])
			index += 8 + symbolLength
		case opcodes.Defer:
			index++
			index += 8
		case opcodes.NewFunction:
			index++
			argsNumber := common.BytesToInt(bytecode[index : index+8])
			index += 8
			for arg := int64(0); arg < argsNumber; arg++ {
				argSymbolLength := common.BytesToInt(bytecode[index : index+8])
				index += 8 + argSymbolLength
			}
			index += 8
		case opcodes.NewClass:
			index++
			index += 8
			index += 8
		case opcodes.Call:
			index++
			index += 8
		case opcodes.NewArray:
			index++
			index += 8
		case opcodes.NewTuple:
			index++
			index += 8
		case opcodes.NewHash:
			index++
			index += 8
		case opcodes.Identifier:
			index++
			symbolLength := common.BytesToInt(bytecode[index : index+8])
			index += 8 + symbolLength
		case opcodes.Integer:
			index++
			index += 8
		case opcodes.Float:
			index++
			index += 8
		case opcodes.String:
			index++
			stringLength := common.BytesToInt(bytecode[index : index+8])
			index += 8 + stringLength
		case opcodes.Bytes:
			index++
			bytesLength := common.BytesToInt(bytecode[index : index+8])
			index += 8 + bytesLength
		case opcodes.True:
			index++
		case opcodes.False:
			index++
		case opcodes.None:
			index++
		case opcodes.Selector:
			index++
			symbolLength := common.BytesToInt(bytecode[index : index+8])
			index += 8 + symbolLength
		case opcodes.Super:
			index++
		default:
			panic(fmt.Sprintf("unknown opcode %d in %v", op, bytecode[index-5:]))
		}
	}
	return bytecode
}

func (a *assembler) Assemble(program ast3.Program) ([]byte, error) {
	resultChan := make(chan []byte, 1)
	errorChan := make(chan error, 1)
	go func(rChan chan []byte, eChan chan error) {
		defer func() {
			err := recover()
			if err != nil {
				rChan <- nil
				eChan <- err.(error)
			}
		}()
		bytecode := make([]byte, 0, len(program))
		for _, node := range program {
			chunk := a.assemble(node)
			bytecode = append(bytecode, chunk...)
		}
		labels := a.enumLabels(bytecode)
		rChan <- a.resolveLabels(bytecode, labels)
		eChan <- nil
	}(resultChan, errorChan)
	return <-resultChan, <-errorChan
}

func AssembleAny(node ast3.Node) ([]byte, error) {
	a := newAssembler()
	return a.Assemble(ast3.Program{node})
}

func Assemble(program ast3.Program) ([]byte, error) {
	a := newAssembler()
	return a.Assemble(program)
}
