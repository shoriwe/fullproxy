package assembler

import (
	"github.com/shoriwe/gplasma/pkg/ast3"
	"github.com/shoriwe/gplasma/pkg/bytecode/opcodes"
	"github.com/shoriwe/gplasma/pkg/common"
)

func (a *assembler) Hash(hash *ast3.Hash) []byte {
	var result []byte
	for _, value := range hash.Values {
		result = append(result, a.Expression(value.Value)...)
		result = append(result, opcodes.Push)
		result = append(result, a.Expression(value.Key)...)
		result = append(result, opcodes.Push)
	}
	result = append(result, opcodes.NewHash)
	result = append(result, common.IntToBytes(len(hash.Values))...)
	return result
}
