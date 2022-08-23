package transformations_1

import (
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
)

func (transform *transformPass) Continue(c *ast2.Continue) []ast3.Node {
	return []ast3.Node{&ast3.ContinueJump{}}
}

func (transform *transformPass) Break(b *ast2.Break) []ast3.Node {
	return []ast3.Node{&ast3.BreakJump{}}
}

func (transform *transformPass) Pass(p *ast2.Pass) []ast3.Node {
	return nil
}
