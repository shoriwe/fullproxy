package assembler

import (
	"github.com/shoriwe/gplasma/pkg/ast3"
	magic_functions "github.com/shoriwe/gplasma/pkg/common/magic-functions"
)

func (a *assembler) Index(index *ast3.Index) []byte {
	return a.assemble(&ast3.Call{
		Expression: nil,
		Function: &ast3.Selector{
			Assignable: nil,
			X:          index.Source,
			Identifier: &ast3.Identifier{
				Assignable: nil,
				Symbol:     magic_functions.Get,
			},
		},
		Arguments: []ast3.Expression{index.Index},
	})
}
