package assembler

import (
	"github.com/shoriwe/gplasma/pkg/ast3"
	"github.com/shoriwe/gplasma/pkg/bytecode/opcodes"
	"github.com/shoriwe/gplasma/pkg/common"
)

func (a *assembler) Selector(selector *ast3.Selector) []byte {
	var result []byte
	result = append(result, a.Expression(selector.X)...)
	result = append(result, opcodes.Push)
	result = append(result, opcodes.Selector)
	result = append(result, common.IntToBytes(len(selector.Identifier.Symbol))...)
	result = append(result, []byte(selector.Identifier.Symbol)...)
	return result
}
