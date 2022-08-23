package assembler

import (
	"github.com/shoriwe/gplasma/pkg/ast3"
	"github.com/shoriwe/gplasma/pkg/bytecode/opcodes"
	"github.com/shoriwe/gplasma/pkg/common"
)

func (a *assembler) Jump(jump *ast3.Jump) []byte {
	result := []byte{opcodes.Jump}
	result = append(result, common.IntToBytes(jump.Target.Code)...)
	return result
}

func (a *assembler) ContinueJump(jump *ast3.ContinueJump) []byte {
	result := []byte{opcodes.Jump}
	result = append(result, common.IntToBytes(jump.Target.Code)...)
	return result
}

func (a *assembler) BreakJump(jump *ast3.BreakJump) []byte {
	result := []byte{opcodes.Jump}
	result = append(result, common.IntToBytes(jump.Target.Code)...)
	return result
}

func (a *assembler) IfJump(jump *ast3.IfJump) []byte {
	result := a.Expression(jump.Condition)
	result = append(result, opcodes.Push, opcodes.IfJump)
	result = append(result, common.IntToBytes(jump.Target.Code)...)
	return result
}
