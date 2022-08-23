package transformations_1

import (
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
)

func (transform *transformPass) Integer(integer *ast2.Integer) *ast3.Integer {
	return &ast3.Integer{
		Value: integer.Value,
	}
}

func (transform *transformPass) Float(float *ast2.Float) *ast3.Float {
	return &ast3.Float{
		Value: float.Value,
	}
}

func (transform *transformPass) String(s *ast2.String) *ast3.String {
	return &ast3.String{
		Contents: s.Contents,
	}
}

func (transform *transformPass) Bytes(bytes *ast2.Bytes) *ast3.Bytes {
	return &ast3.Bytes{
		Contents: bytes.Contents,
	}
}

func (transform *transformPass) True(t *ast2.True) *ast3.True {
	return &ast3.True{}
}

func (transform *transformPass) False(t *ast2.False) *ast3.False {
	return &ast3.False{}
}

func (transform *transformPass) None(t *ast2.None) *ast3.None {
	return &ast3.None{}
}
