package transformations_1

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
	"github.com/shoriwe/gplasma/pkg/common"
	"github.com/shoriwe/gplasma/pkg/common/magic-functions"
	"github.com/shoriwe/gplasma/pkg/common/special-symbols"
	"reflect"
)

func getRoot(expr ast3.Expression) ast3.Expression {
	current := expr
	for {
		switch c := current.(type) {
		case *ast3.Selector:
			current = c.X
		case *ast3.Index:
			current = c.Source
		default:
			return c
		}
	}
}

type generatorTransform struct {
	transform       *transformPass
	selfSymbols     map[string]struct{}
	hasNextVariable *ast3.Selector
	labelsOrder     []*ast3.Selector
	labels          map[*ast3.Selector]*ast3.Label
}

func newGeneratorTransform(transform *transformPass) *generatorTransform {
	return &generatorTransform{
		transform:   transform,
		selfSymbols: map[string]struct{}{},
		hasNextVariable: &ast3.Selector{
			X: &ast3.Identifier{
				Symbol: special_symbols.Self,
			},
			Identifier: &ast3.Identifier{
				Symbol: "____has_next",
			},
		},
		labelsOrder: nil,
		labels:      map[*ast3.Selector]*ast3.Label{},
	}
}

func (gt *generatorTransform) nextLabelIdentifier() (*ast3.Label, *ast3.Selector) {
	label := gt.transform.nextLabel()
	return label, &ast3.Selector{
		X: &ast3.Identifier{
			Symbol: special_symbols.Self,
		},
		Identifier: &ast3.Identifier{
			Symbol: fmt.Sprintf("_______label_%d", label.Code),
		},
	}
}

/*
	- Transform assignments to self.IDENTIFIER
	- Update Identifier access to self.IDENTIFIER
	- Add update jump condition variable before yield
	- Add label after yield
	- Update to has_next variable before return
*/
func (gt *generatorTransform) resolve(node ast3.Node, symbols map[string]struct{}) []ast3.Node {
	symbolsCopy := common.CopyMap(symbols)
	switch n := node.(type) {
	case *ast3.Assignment:
		if ident, ok := n.Left.(*ast3.Identifier); ok {
			_, found := gt.selfSymbols[ident.Symbol]
			if found {
				delete(symbolsCopy, ident.Symbol)
			}
		}
		return []ast3.Node{&ast3.Assignment{
			Left:  gt.resolve(n.Left, symbolsCopy)[0].(ast3.Assignable),
			Right: gt.resolve(n.Right, symbolsCopy)[0].(ast3.Expression),
		}}
	case *ast3.Label:
		return []ast3.Node{n}
	case *ast3.Jump:
		return []ast3.Node{n}
	case *ast3.ContinueJump:
		return []ast3.Node{n}
	case *ast3.BreakJump:
		return []ast3.Node{n}
	case *ast3.IfJump:
		return []ast3.Node{&ast3.IfJump{
			Condition: gt.resolve(n.Condition, symbolsCopy)[0].(ast3.Expression),
			Target:    n.Target,
		}}
	case *ast3.Return:
		return []ast3.Node{
			&ast3.Assignment{
				Left:  gt.hasNextVariable,
				Right: &ast3.False{},
			},
			&ast3.Return{
				Result: gt.resolve(n.Result, symbolsCopy)[0].(ast3.Expression),
			},
		}
	case *ast3.Yield:
		label, selector := gt.nextLabelIdentifier()
		gt.labels[selector] = label
		gt.labelsOrder = append(gt.labelsOrder, selector)
		return []ast3.Node{
			&ast3.Assignment{
				Left:  selector,
				Right: &ast3.True{},
			},
			&ast3.Yield{
				Result: gt.resolve(n.Result, symbolsCopy)[0].(ast3.Expression),
			},
			label,
		}
	case *ast3.Delete:
		return []ast3.Node{&ast3.Delete{
			Statement: nil,
			X:         gt.resolve(n.X, symbolsCopy)[0].(ast3.Assignable),
		}}
	case *ast3.Defer:
		return []ast3.Node{&ast3.Defer{
			Statement: nil,
			X:         gt.resolve(n.X, symbolsCopy)[0].(ast3.Expression),
		}}
	case *ast3.Function:
		for _, argument := range n.Arguments {
			if _, found := symbolsCopy[argument.Symbol]; found {
				delete(symbolsCopy, argument.Symbol)
			}
		}
		body := make([]ast3.Node, 0, len(n.Body))
		for _, child := range n.Body {
			body = append(body, gt.resolve(child, symbolsCopy)...)
		}
		return []ast3.Node{&ast3.Function{
			Arguments: n.Arguments,
			Body:      body,
		}}
	case *ast3.Class:
		bases := make([]ast3.Expression, 0, len(n.Bases))
		for _, base := range n.Bases {
			bases = append(bases, gt.resolve(base, symbolsCopy)[0].(ast3.Expression))
		}
		body := make([]ast3.Node, 0, len(n.Body))
		for _, child := range n.Body {
			body = append(body, gt.resolve(child, symbolsCopy)...)
		}
		return []ast3.Node{&ast3.Class{
			Bases: bases,
			Body:  body,
		}}
	case *ast3.Call:
		arguments := make([]ast3.Expression, 0, len(n.Arguments))
		for _, argument := range n.Arguments {
			arguments = append(arguments, gt.resolve(argument, symbolsCopy)[0].(ast3.Expression))
		}
		return []ast3.Node{&ast3.Call{
			Function:  gt.resolve(n.Function, symbolsCopy)[0].(ast3.Expression),
			Arguments: arguments,
		}}
	case *ast3.Array:
		values := make([]ast3.Expression, 0, len(n.Values))
		for _, value := range n.Values {
			values = append(values, gt.resolve(value, symbolsCopy)[0].(ast3.Expression))
		}
		return []ast3.Node{&ast3.Array{
			Values: values,
		}}
	case *ast3.Tuple:
		values := make([]ast3.Expression, 0, len(n.Values))
		for _, value := range n.Values {
			values = append(values, gt.resolve(value, symbolsCopy)[0].(ast3.Expression))
		}
		return []ast3.Node{&ast3.Tuple{
			Values: values,
		}}
	case *ast3.Hash:
		values := make([]*ast3.KeyValue, 0, len(n.Values))
		for _, keyValue := range n.Values {
			values = append(values, &ast3.KeyValue{
				Key:   gt.resolve(keyValue.Key, symbolsCopy)[0].(ast3.Expression),
				Value: gt.resolve(keyValue.Value, symbolsCopy)[0].(ast3.Expression),
			})
		}
		return []ast3.Node{&ast3.Hash{
			Values: values,
		}}
	case *ast3.Identifier:
		if _, found := symbolsCopy[n.Symbol]; found {
			return []ast3.Node{&ast3.Selector{
				X: &ast3.Identifier{
					Symbol: special_symbols.Self,
				},
				Identifier: n,
			}}
		}
		return []ast3.Node{n}
	case *ast3.Integer:
		return []ast3.Node{n}
	case *ast3.Float:
		return []ast3.Node{n}
	case *ast3.String:
		return []ast3.Node{n}
	case *ast3.Bytes:
		return []ast3.Node{n}
	case *ast3.True:
		return []ast3.Node{n}
	case *ast3.False:
		return []ast3.Node{n}
	case *ast3.None:
		return []ast3.Node{n}
	case *ast3.Selector:
		return []ast3.Node{&ast3.Selector{
			X:          gt.resolve(n.X, symbolsCopy)[0].(ast3.Expression),
			Identifier: n.Identifier,
		}}
	case *ast3.Index:
		return []ast3.Node{&ast3.Index{
			Source: gt.resolve(n.Source, symbolsCopy)[0].(ast3.Expression),
			Index:  gt.resolve(n.Index, symbolsCopy)[0].(ast3.Expression),
		}}
	case *ast3.Super:
		return []ast3.Node{&ast3.Super{
			X: gt.resolve(n.X, symbolsCopy)[0].(ast3.Expression),
		}}
	default:
		panic(fmt.Sprintf("unknown node type %s", reflect.TypeOf(node).String()))
	}
}

// Enumerate self symbols
func (gt *generatorTransform) enumerate(a ast3.Assignable) {
	switch left := a.(type) {
	case *ast3.Identifier:
		gt.selfSymbols[left.Symbol] = struct{}{}
	case *ast3.Selector:
		root := getRoot(left.X)
		ident, isIdent := root.(*ast3.Identifier)
		if !isIdent {
			return
		}
		if _, found := gt.selfSymbols[ident.Symbol]; !found {
			return
		}
		gt.selfSymbols[ident.Symbol] = struct{}{}
	case *ast3.Index:
		root := getRoot(left.Source)
		ident, isIdent := root.(*ast3.Identifier)
		if !isIdent {
			return
		}
		if _, found := gt.selfSymbols[ident.Symbol]; !found {
			return
		}
		gt.selfSymbols[ident.Symbol] = struct{}{}
	default:
		panic(fmt.Sprintf("unknown assignable type %s", reflect.TypeOf(a).String()))
	}
}

func (gt *generatorTransform) process(node ast3.Node) []ast3.Node {
	switch n := node.(type) {
	case *ast3.Assignment:
		gt.enumerate(n.Left)
		return []ast3.Node{&ast3.Assignment{
			Left:  gt.resolve(n.Left, gt.selfSymbols)[0].(ast3.Assignable),
			Right: gt.resolve(n.Right, gt.selfSymbols)[0].(ast3.Expression),
		}}
	default:
		return gt.resolve(n, gt.selfSymbols)
	}
}

/*
	- Prepend has_next jump
	- Prepend jump table
	- Append Body
	- Append On finish label
	- Append Return none
*/
func (gt *generatorTransform) setup(body []ast3.Node) []ast3.Node {
	result := make([]ast3.Node, 0, 3+len(gt.labelsOrder)+len(body))
	onFinishLabel := gt.transform.nextLabel()
	// Has next jump
	result = append(result, &ast3.IfJump{
		Condition: &ast3.Call{
			Function: &ast3.Selector{
				Assignable: nil,
				X:          gt.hasNextVariable,
				Identifier: &ast3.Identifier{
					Symbol: magic_functions.Not,
				},
			},
		},
		Target: onFinishLabel,
	})
	// Jump table
	for _, selector := range gt.labelsOrder {
		label := gt.labels[selector]
		result = append(result, &ast3.IfJump{
			Statement: nil,
			Condition: selector,
			Target:    label,
		})
	}
	// Body
	for _, node := range body {
		result = append(result, node)
	}
	// Finish label
	result = append(result, onFinishLabel)
	// Return none
	result = append(result, &ast3.Return{
		Result: &ast3.None{},
	})
	return result
}

func (gt *generatorTransform) next(rawBody []ast3.Node) *ast3.Assignment {
	body := make([]ast3.Node, 0, len(rawBody))
	for _, node := range rawBody {
		body = append(body, gt.process(node)...)
	}
	return &ast3.Assignment{
		Left: &ast3.Identifier{
			Symbol: magic_functions.Next,
		},
		Right: &ast3.Function{
			Body: gt.setup(body),
		},
	}
}

func (gt *generatorTransform) hasNext() *ast3.Assignment {
	body := []ast3.Node{
		&ast3.Return{
			Result: gt.hasNextVariable,
		},
	}
	return &ast3.Assignment{
		Statement: nil,
		Left: &ast3.Identifier{
			Symbol: magic_functions.HasNext,
		},
		Right: &ast3.Function{
			Body: body,
		},
	}
}

func (gt *generatorTransform) init(arguments []*ast3.Identifier) *ast3.Assignment {
	body := make([]ast3.Node, 0, len(arguments))
	for _, argument := range arguments {
		gt.selfSymbols[argument.Symbol] = struct{}{}
		body = append(body, &ast3.Assignment{
			Left: &ast3.Selector{
				Assignable: nil,
				X: &ast3.Identifier{
					Symbol: special_symbols.Self,
				},
				Identifier: argument,
			},
			Right: argument,
		})
	}
	return &ast3.Assignment{
		Left: &ast3.Identifier{
			Symbol: magic_functions.Init,
		},
		Right: &ast3.Function{
			Arguments: arguments,
			Body:      body,
		},
	}
}

func (gt *generatorTransform) class(rawFunctionBody []ast3.Node, arguments []*ast3.Identifier) *ast3.Class {
	initFunction := gt.init(arguments)
	nextFunction := gt.next(rawFunctionBody)
	hasNextFunction := gt.hasNext()
	body := make([]ast3.Node, 0, 3+len(gt.selfSymbols))
	body = append(body, &ast3.Assignment{
		Left:  gt.hasNextVariable,
		Right: &ast3.True{},
	})
	for selfSymbol := range gt.selfSymbols {
		body = append(body, &ast3.Assignment{
			Left: &ast3.Identifier{
				Symbol: selfSymbol,
			},
			Right: &ast3.None{},
		})
	}
	body = append(body, initFunction, hasNextFunction, nextFunction)
	for selector := range gt.labels {
		body = append(body, &ast3.Assignment{
			Statement: nil,
			Left:      selector.Identifier,
			Right:     &ast3.False{},
		})
	}
	return &ast3.Class{
		Body: body,
	}
}

func (transform *transformPass) GeneratorDef(generator *ast2.GeneratorDefinition) []ast3.Node {
	rawNextFunctionBody := make([]ast3.Node, 0, len(generator.Body))
	for _, node := range generator.Body {
		rawNextFunctionBody = append(rawNextFunctionBody, transform.Node(node)...)
	}
	arguments := make([]*ast3.Identifier, 0, len(generator.Arguments))
	for _, argument := range generator.Arguments {
		arguments = append(arguments, transform.Identifier(argument))
	}
	class := newGeneratorTransform(transform).class(rawNextFunctionBody, arguments)
	return []ast3.Node{&ast3.Assignment{
		Statement: nil,
		Left:      transform.Identifier(generator.Name),
		Right:     class,
	}}
}
