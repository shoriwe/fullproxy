package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simplify *simplifyPass) Return(ret *ast.ReturnStatement) *ast2.Return {
	switch len(ret.Results) {
	case 0:
		return &ast2.Return{
			Result: &ast2.None{},
		}
	case 1:
		return &ast2.Return{
			Result: simplify.Expression(ret.Results[0]),
		}
	}
	values := make([]ast2.Expression, 0, len(ret.Results))
	for _, result := range ret.Results {
		values = append(values, simplify.Expression(result))
	}
	return &ast2.Return{
		Result: &ast2.Tuple{
			Values: values,
		},
	}
}

func (simplify *simplifyPass) Yield(yield *ast.YieldStatement) *ast2.Yield {
	switch len(yield.Results) {
	case 0:
		return &ast2.Yield{
			Result: &ast2.None{},
		}
	case 1:
		return &ast2.Yield{
			Result: simplify.Expression(yield.Results[0]),
		}
	}
	values := make([]ast2.Expression, 0, len(yield.Results))
	for _, result := range yield.Results {
		values = append(values, simplify.Expression(result))
	}
	return &ast2.Yield{
		Result: &ast2.Tuple{
			Values: values,
		},
	}
}
