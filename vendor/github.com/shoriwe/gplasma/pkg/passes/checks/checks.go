package checks

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/common"
)

/*
Check verifies:
- Yield statement is only on generator statements
- Break/Continue/Redo are only in loop statements
- Return statement is only on functions and generator statements
*/
type Check struct {
	InvalidFunctionNodesStack  common.ListStack[ast.Node]
	InvalidGeneratorNodesStack common.ListStack[ast.Node]
	InvalidLoopNodesStack      common.ListStack[ast.Node]
	insideFunction             bool
	insideGenerator            bool
	insideLoop                 bool
}

func (c *Check) Visit(node ast.Node) ast.Visitor {
	defer func(oldInsideFunction, oldInsideGenerator, oldInsideLoop bool) {
		c.insideFunction, c.insideGenerator, c.insideLoop = oldInsideFunction, oldInsideGenerator, oldInsideLoop
	}(c.insideFunction, c.insideGenerator, c.insideLoop)
	switch n := node.(type) {
	case *ast.FunctionDefinitionStatement:
		c.insideFunction = true
		c.insideGenerator = false
		c.insideLoop = false
		for _, child := range n.Body {
			c.Visit(child)
		}
		return nil
	case *ast.GeneratorDefinitionStatement:
		c.insideFunction = false
		c.insideGenerator = true
		c.insideLoop = false
		for _, child := range n.Body {
			c.Visit(child)
		}
		return nil
	case *ast.ForLoopStatement:
		c.insideLoop = true
		for _, child := range n.Body {
			c.Visit(child)
		}
		return nil
	case *ast.WhileLoopStatement:
		c.insideLoop = true
		for _, child := range n.Body {
			c.Visit(child)
		}
		return nil
	case *ast.DoWhileStatement:
		c.insideLoop = true
		for _, child := range n.Body {
			c.Visit(child)
		}
		return nil
	case *ast.UntilLoopStatement:
		c.insideLoop = true
		for _, child := range n.Body {
			c.Visit(child)
		}
		return nil
	case *ast.ReturnStatement:
		if !c.insideFunction && !c.insideGenerator {
			c.InvalidFunctionNodesStack.Push(n)
		}
		return nil
	case *ast.YieldStatement:
		if !c.insideGenerator {
			c.InvalidGeneratorNodesStack.Push(n)
		}
		return nil
	case *ast.DeferStatement:
		if !c.insideFunction && !c.insideGenerator {
			c.InvalidFunctionNodesStack.Push(n)
		}
		return nil
	case *ast.BreakStatement, *ast.ContinueStatement:
		if !c.insideLoop {
			c.InvalidLoopNodesStack.Push(n)
		}
		return nil
	}
	return c
}

func (c *Check) CountInvalidFunctionNodes() int {
	result := 0
	for current := c.InvalidFunctionNodesStack.Top; current != nil; current = current.Next {
		result++
	}
	return result
}

func (c *Check) CountInvalidGeneratorNodes() int {
	result := 0
	for current := c.InvalidGeneratorNodesStack.Top; current != nil; current = current.Next {
		result++
	}
	return result
}

func (c *Check) CountInvalidLoopNodes() int {
	result := 0
	for current := c.InvalidLoopNodesStack.Top; current != nil; current = current.Next {
		result++
	}
	return result
}

func NewCheckPass() *Check {
	return &Check{
		InvalidFunctionNodesStack:  common.ListStack[ast.Node]{},
		InvalidGeneratorNodesStack: common.ListStack[ast.Node]{},
		InvalidLoopNodesStack:      common.ListStack[ast.Node]{},
		insideFunction:             false,
		insideGenerator:            false,
		insideLoop:                 false,
	}
}
