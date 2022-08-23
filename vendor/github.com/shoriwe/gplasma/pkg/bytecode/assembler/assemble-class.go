package assembler

import (
	"github.com/shoriwe/gplasma/pkg/ast3"
	"github.com/shoriwe/gplasma/pkg/bytecode/opcodes"
	"github.com/shoriwe/gplasma/pkg/common"
)

func (a *assembler) Class(class *ast3.Class) []byte {
	var bases []byte
	for _, base := range class.Bases {
		bases = append(bases, a.Expression(base)...)
		bases = append(bases, opcodes.Push)
	}
	var body []byte
	for _, node := range class.Body {
		body = append(body, a.assemble(node)...)
	}
	var result []byte
	result = append(result, bases...)
	result = append(result, opcodes.NewClass)
	result = append(result, common.IntToBytes(len(class.Bases))...)
	result = append(result, common.IntToBytes(len(body))...)
	result = append(result, body...)
	return result
}
