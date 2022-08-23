package assembler

import (
	"github.com/shoriwe/gplasma/pkg/ast3"
	"github.com/shoriwe/gplasma/pkg/bytecode/opcodes"
)

func (a *assembler) Return(ret *ast3.Return) []byte {
	result := a.Expression(ret.Result)
	result = append(result, opcodes.Push, opcodes.Return)
	return result
}

func (a *assembler) Yield(yield *ast3.Yield) []byte {
	result := a.Expression(yield.Result)
	result = append(result, opcodes.Push, opcodes.Return)
	return result
}
