package assembler

import (
	"github.com/shoriwe/gplasma/pkg/ast3"
	"github.com/shoriwe/gplasma/pkg/bytecode/opcodes"
	"github.com/shoriwe/gplasma/pkg/common"
)

func (a *assembler) Defer(defer_ *ast3.Defer) []byte {
	expression := a.Expression(defer_.X)
	result := []byte{opcodes.Defer}
	result = append(result, common.IntToBytes(len(expression))...)
	result = append(result, expression...)
	return result
}
