package assembler

import (
	"github.com/shoriwe/gplasma/pkg/ast3"
	"github.com/shoriwe/gplasma/pkg/bytecode/opcodes"
	"github.com/shoriwe/gplasma/pkg/common"
)

func (a *assembler) Function(function *ast3.Function) []byte {
	arguments := make([]byte, 0, len(function.Arguments))
	for _, argument := range function.Arguments {
		arguments = append(arguments, common.IntToBytes(len(argument.Symbol))...)
		arguments = append(arguments, []byte(argument.Symbol)...)
	}
	body := make([]byte, 0, len(function.Body))
	for _, node := range function.Body {
		body = append(body, a.assemble(node)...)
	}
	var result []byte
	result = append(result, opcodes.NewFunction)
	result = append(result, common.IntToBytes(len(function.Arguments))...)
	result = append(result, arguments...)
	result = append(result, common.IntToBytes(len(body))...)
	result = append(result, body...)
	return result
}
