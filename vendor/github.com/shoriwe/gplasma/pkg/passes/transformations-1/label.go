package transformations_1

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast3"
)

func (transform *transformPass) nextAnonIdentifier() *ast3.Identifier {
	identifier := transform.currentIdentifier
	transform.currentIdentifier++
	return &ast3.Identifier{
		Assignable: nil,
		Symbol:     fmt.Sprintf("_____transform_%d", identifier),
	}
}

func (transform *transformPass) nextLabel() *ast3.Label {
	label := transform.currentLabel
	transform.currentLabel++
	return &ast3.Label{
		Code: label,
	}
}
