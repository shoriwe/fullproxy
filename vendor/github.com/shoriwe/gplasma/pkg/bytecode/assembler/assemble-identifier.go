package assembler

import (
	"github.com/shoriwe/gplasma/pkg/ast3"
	"github.com/shoriwe/gplasma/pkg/bytecode/opcodes"
	"github.com/shoriwe/gplasma/pkg/common"
)

func (a *assembler) Identifier(ident *ast3.Identifier) []byte {
	var result []byte
	result = append(result, opcodes.Identifier)
	result = append(result, common.IntToBytes(len(ident.Symbol))...)
	result = append(result, []byte(ident.Symbol)...)
	return result
}
