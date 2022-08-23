package transformations_1

import (
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
	magic_functions "github.com/shoriwe/gplasma/pkg/common/magic-functions"
	special_symbols "github.com/shoriwe/gplasma/pkg/common/special-symbols"
)

func (transform *transformPass) GeneratorExpr(generator *ast2.Generator) *ast3.Call {
	selfSource := &ast3.Selector{
		X: &ast3.Identifier{
			Symbol: special_symbols.Self,
		},
		Identifier: &ast3.Identifier{
			Symbol: "_____source",
		},
	}
	initSourceArgument := &ast3.Identifier{
		Symbol: "__source__",
	}
	initFunction := &ast3.Assignment{
		Left: &ast3.Identifier{
			Symbol: magic_functions.Init,
		},
		Right: &ast3.Function{
			Arguments: []*ast3.Identifier{initSourceArgument},
			Body: []ast3.Node{
				&ast3.Assignment{
					Left: selfSource,
					Right: &ast3.Call{
						Function: &ast3.Selector{
							X: initSourceArgument,
							Identifier: &ast3.Identifier{
								Symbol: magic_functions.Iter,
							},
						},
					},
				},
			},
		},
	}
	hasNextFunction := &ast3.Assignment{
		Left: &ast3.Identifier{
			Symbol: magic_functions.HasNext,
		},
		Right: &ast3.Function{
			Body: []ast3.Node{
				&ast3.Return{
					Result: &ast3.Call{
						Function: &ast3.Selector{
							X: selfSource,
							Identifier: &ast3.Identifier{
								Symbol: magic_functions.HasNext,
							},
						},
					},
				},
			},
		},
	}
	anonReceiver := transform.nextAnonIdentifier()
	sourceNext := &ast3.Assignment{
		Left: anonReceiver,
		Right: &ast3.Call{
			Function: &ast3.Selector{
				X: selfSource,
				Identifier: &ast3.Identifier{
					Symbol: magic_functions.Next,
				},
			},
		},
	}
	expand := make([]ast3.Node, 0, len(generator.Receivers))
	if len(generator.Receivers) > 1 {
		for index, receiver := range generator.Receivers {
			expand = append(expand, &ast3.Assignment{
				Left: transform.Identifier(receiver),
				Right: &ast3.Index{
					Source: anonReceiver,
					Index: &ast3.Integer{
						Value: int64(index),
					},
				},
			})
		}
	} else {
		expand = append(expand, &ast3.Assignment{
			Left: &ast3.Identifier{
				Symbol: generator.Receivers[0].Symbol,
			},
			Right: anonReceiver,
		})
	}
	result := &ast3.Return{
		Statement: nil,
		Result:    transform.Expression(generator.Operation),
	}
	nextBody := make([]ast3.Node, 0, 2+len(expand))
	nextBody = append(nextBody, sourceNext)
	nextBody = append(nextBody, expand...)
	nextBody = append(nextBody, result)
	nextFunction := &ast3.Assignment{
		Left: &ast3.Identifier{
			Symbol: magic_functions.Next,
		},
		Right: &ast3.Function{
			Body: nextBody,
		},
	}
	class := &ast3.Class{
		Expression: nil,
		Bases:      nil,
		Body:       []ast3.Node{initFunction, hasNextFunction, nextFunction},
	}
	sourceExpr := transform.Expression(generator.Source)
	return &ast3.Call{
		Function:  class,
		Arguments: []ast3.Expression{sourceExpr},
	}
}
