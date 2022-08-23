package assembler

import (
	"github.com/shoriwe/gplasma/pkg/ast3"
	"github.com/shoriwe/gplasma/pkg/bytecode/opcodes"
	"github.com/shoriwe/gplasma/pkg/common"
)

func (a *assembler) Integer(integer *ast3.Integer) []byte {
	var result []byte
	result = append(result, opcodes.Integer)
	result = append(result, common.IntToBytes(integer.Value)...)
	return result
}

func (a *assembler) Float(float *ast3.Float) []byte {
	var result []byte
	result = append(result, opcodes.Float)
	result = append(result, common.FloatToBytes(float.Value)...)
	return result
}

func (a *assembler) String(s *ast3.String) []byte {
	var result []byte
	result = append(result, opcodes.String)
	result = append(result, common.IntToBytes(len(s.Contents))...)
	result = append(result, s.Contents...)
	return result
}

func (a *assembler) Bytes(bytes *ast3.Bytes) []byte {
	var result []byte
	result = append(result, opcodes.Bytes)
	result = append(result, common.IntToBytes(len(bytes.Contents))...)
	result = append(result, bytes.Contents...)
	return result
}

func (a *assembler) True(t *ast3.True) []byte {
	return []byte{opcodes.True}
}

func (a *assembler) False(f *ast3.False) []byte {
	return []byte{opcodes.False}
}

func (a *assembler) None(none *ast3.None) []byte {
	return []byte{opcodes.None}
}
