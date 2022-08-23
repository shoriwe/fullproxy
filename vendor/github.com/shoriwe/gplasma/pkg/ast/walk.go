package ast

import (
	"fmt"
	"reflect"
)

type Visitor interface {
	Visit(node Node) Visitor
}

func walk(visitor Visitor, node Node) {
	// Visit the node
	if visitor = visitor.Visit(node); visitor == nil {
		return
	}
	// Visit its children
	switch n := node.(type) {
	case *Identifier, *BasicLiteralExpression:
		break
	case *Program:
		if n.Begin != nil {
			walk(visitor, n.Begin)
		}
		for _, child := range n.Body {
			walk(visitor, child)
		}
		if n.End != nil {
			walk(visitor, n.End)
		}
	case *ArrayExpression:
		for _, value := range n.Values {
			walk(visitor, value)
		}
	case *TupleExpression:
		for _, value := range n.Values {
			walk(visitor, value)
		}
	case *HashExpression:
		for _, keyValue := range n.Values {
			walk(visitor, keyValue.Key)
			walk(visitor, keyValue.Value)
		}
	case *BinaryExpression:
		walk(visitor, n.LeftHandSide)
		walk(visitor, n.RightHandSide)
	case *UnaryExpression:
		walk(visitor, n.X)
	case *ParenthesesExpression:
		walk(visitor, n.X)
	case *LambdaExpression:
		for _, argument := range n.Arguments {
			walk(visitor, argument)
		}
		walk(visitor, n.Code)
	case *GeneratorExpression:
		walk(visitor, n.Operation)
		for _, identifier := range n.Receivers {
			walk(visitor, identifier)
		}
		walk(visitor, n.Source)
	case *SelectorExpression:
		walk(visitor, n.X)
		walk(visitor, n.Identifier)
	case *MethodInvocationExpression:
		walk(visitor, n.Function)
		for _, argument := range n.Arguments {
			walk(visitor, argument)
		}
	case *IndexExpression:
		walk(visitor, n.Source)
		walk(visitor, n.Index)
	case *IfOneLinerExpression:
		walk(visitor, n.Result)
		walk(visitor, n.Condition)
		walk(visitor, n.ElseResult)
	case *UnlessOneLinerExpression:
		walk(visitor, n.Result)
		walk(visitor, n.Condition)
		walk(visitor, n.ElseResult)
	case *AssignStatement:
		walk(visitor, n.LeftHandSide)
		walk(visitor, n.RightHandSide)
	case *DoWhileStatement:
		walk(visitor, n.Condition)
		for _, bodyNode := range n.Body {
			walk(visitor, bodyNode)
		}
	case *WhileLoopStatement:
		walk(visitor, n.Condition)
		for _, bodyNode := range n.Body {
			walk(visitor, bodyNode)
		}
	case *UntilLoopStatement:
		walk(visitor, n.Condition)
		for _, bodyNode := range n.Body {
			walk(visitor, bodyNode)
		}
	case *ForLoopStatement:
		for _, receiver := range n.Receivers {
			walk(visitor, receiver)
		}
		walk(visitor, n.Source)
		for _, bodyNode := range n.Body {
			walk(visitor, bodyNode)
		}
	case *SwitchStatement:
		walk(visitor, n.Target)
		for _, caseBlock := range n.CaseBlocks {
			for _, case_ := range caseBlock.Cases {
				walk(visitor, case_)
			}
			for _, bodyNode := range caseBlock.Body {
				walk(visitor, bodyNode)
			}
		}
		for _, bodyNode := range n.Default {
			walk(visitor, bodyNode)
		}
	case *ModuleStatement:
		walk(visitor, n.Name)
		for _, bodyNode := range n.Body {
			walk(visitor, bodyNode)
		}
	case *FunctionDefinitionStatement:
		walk(visitor, n.Name)
		for _, argument := range n.Arguments {
			walk(visitor, argument)
		}
		for _, bodyNode := range n.Body {
			walk(visitor, bodyNode)
		}
	case *InterfaceStatement:
		walk(visitor, n.Name)
		for _, base := range n.Bases {
			walk(visitor, base)
		}
		for _, methodDefinition := range n.MethodDefinitions {
			walk(visitor, methodDefinition)
		}
	case *ClassStatement:
		walk(visitor, n.Name)
		for _, base := range n.Bases {
			walk(visitor, base)
		}
		for _, bodyNode := range n.Body {
			walk(visitor, bodyNode)
		}
	case *BeginStatement:
		for _, bodyNode := range n.Body {
			walk(visitor, bodyNode)
		}
	case *EndStatement:
		for _, bodyNode := range n.Body {
			walk(visitor, bodyNode)
		}
	case *ReturnStatement:
		for _, value := range n.Results {
			walk(visitor, value)
		}
	case *YieldStatement:
		for _, value := range n.Results {
			walk(visitor, value)
		}
	case *IfStatement:
		walk(visitor, n.Condition)
		for _, bodyChild := range n.Body {
			walk(visitor, bodyChild)
		}
		for _, elifBlock := range n.ElifBlocks {
			walk(visitor, elifBlock.Condition)
			for _, elifBlockChild := range elifBlock.Body {
				walk(visitor, elifBlockChild)
			}
		}
	case *UnlessStatement:
		walk(visitor, n.Condition)
		for _, bodyChild := range n.Body {
			walk(visitor, bodyChild)
		}
		for _, elifBlock := range n.ElifBlocks {
			walk(visitor, elifBlock.Condition)
			for _, elifBlockChild := range elifBlock.Body {
				walk(visitor, elifBlockChild)
			}
		}
	case *DeleteStatement:
		walk(visitor, n.X)
	case *PassStatement:
		return
	case nil:
		break // Ignore nil
	default:
		panic(fmt.Sprintf("unknown node type %s", reflect.TypeOf(node).String()))
	}
}

func Walk(visitor Visitor, node Node) {
	walk(visitor, node)
}
