package assembler

import (
	"github.com/shoriwe/gplasma/pkg/ast3"
	"github.com/shoriwe/gplasma/pkg/bytecode/opcodes"
	"github.com/shoriwe/gplasma/pkg/common"
)

func (a *assembler) Label(label *ast3.Label) []byte {
	result := []byte{opcodes.Label}
	result = append(result, common.IntToBytes(label.Code)...)
	return result
}
