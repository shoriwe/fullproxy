package ast

import (
	"github.com/shoriwe/gplasma/pkg/compiler/lexer"
	"github.com/shoriwe/gplasma/pkg/errors"
	"github.com/shoriwe/gplasma/pkg/vm"
	"reflect"
)

type Statement interface {
	S()
	Node
}

func compileClassBody(body []Node) ([]*vm.Code, *errors.Error) {
	foundInitialize := false
	var isInitialize bool
	var nodeCode []*vm.Code
	var compilationError *errors.Error
	var result []*vm.Code
	for _, node := range body {
		switch node.(type) {
		case IExpression:
			nodeCode, compilationError = node.(IExpression).CompilePush(true)
		case Statement:
			if _, ok := node.(*FunctionDefinitionStatement); ok {
				nodeCode, compilationError, isInitialize = node.(*FunctionDefinitionStatement).CompileAsClassFunction()
				if isInitialize && !foundInitialize {
					foundInitialize = true
				}
			} else {
				nodeCode, compilationError = node.(Statement).Compile()
			}
		}
		if compilationError != nil {
			return nil, compilationError
		}
		result = append(result, nodeCode...)
	}
	if !foundInitialize {
		initFunction := &FunctionDefinitionStatement{
			Name: &Identifier{
				Token: &lexer.Token{
					String: vm.Initialize,
				},
			},
			Arguments: nil,
			Body:      nil,
		}
		nodeCode, _, _ = initFunction.CompileAsClassFunction()
		result = append(result, nodeCode...)
	}
	return result, nil
}

type AssignStatement struct {
	Statement
	LeftHandSide   IExpression // Identifiers or Selectors
	AssignOperator *lexer.Token
	RightHandSide  IExpression
}

func compileAssignStatementMiddleBinaryExpression(leftHandSide IExpression, assignOperator *lexer.Token) ([]*vm.Code, *errors.Error) {
	result, leftHandSideCompilationError := leftHandSide.CompilePush(true)
	if leftHandSideCompilationError != nil {
		return nil, leftHandSideCompilationError
	}
	// Finally decide the instruction to use
	var operation uint8
	switch assignOperator.DirectValue {
	case lexer.AddAssign:
		operation = vm.AddOP
	case lexer.SubAssign:
		operation = vm.SubOP
	case lexer.StarAssign:
		operation = vm.MulOP
	case lexer.DivAssign:
		operation = vm.DivOP
	case lexer.ModulusAssign:
		operation = vm.ModOP
	case lexer.PowerOfAssign:
		operation = vm.PowOP
	case lexer.BitwiseXorAssign:
		operation = vm.BitXorOP
	case lexer.BitWiseAndAssign:
		operation = vm.BitAndOP
	case lexer.BitwiseOrAssign:
		operation = vm.BitOrOP
	case lexer.BitwiseLeftAssign:
		operation = vm.BitLeftOP
	case lexer.BitwiseRightAssign:
		operation = vm.BitRightOP
	default:
		panic(errors.NewUnknownVMOperationError(operation))
	}
	return append(result, vm.NewCode(vm.BinaryOP, assignOperator.Line, operation)), nil
}

func compileIdentifierAssign(identifier *Identifier) ([]*vm.Code, *errors.Error) {
	return []*vm.Code{vm.NewCode(vm.AssignIdentifierOP, identifier.Token.Line, identifier.Token.String)}, nil
}

func compileSelectorAssign(selectorExpression *SelectorExpression) ([]*vm.Code, *errors.Error) {
	result, sourceCompilationError := selectorExpression.X.CompilePush(true)
	if sourceCompilationError != nil {
		return nil, sourceCompilationError
	}
	return append(result, vm.NewCode(vm.AssignSelectorOP, selectorExpression.Identifier.Token.Line, selectorExpression.Identifier.Token.String)), nil
}

func compileIndexAssign(indexExpression *IndexExpression) ([]*vm.Code, *errors.Error) {
	result, sourceCompilationError := indexExpression.Source.CompilePush(true)
	if sourceCompilationError != nil {
		return nil, sourceCompilationError
	}
	index, indexCompilationError := indexExpression.Index.CompilePush(true)
	if indexCompilationError != nil {
		return nil, indexCompilationError
	}
	result = append(result, index...)
	return append(result, vm.NewCode(vm.AssignIndexOP, errors.UnknownLine, nil)), nil
}

func (assignStatement *AssignStatement) Compile() ([]*vm.Code, *errors.Error) {
	result, valueCompilationError := assignStatement.RightHandSide.CompilePush(true)
	if valueCompilationError != nil {
		return nil, valueCompilationError
	}
	if assignStatement.AssignOperator.DirectValue != lexer.Assign {
		// Do something here to evaluate the operation
		assignOperation, middleOperationCompilationError := compileAssignStatementMiddleBinaryExpression(assignStatement.LeftHandSide, assignStatement.AssignOperator)
		if middleOperationCompilationError != nil {
			return nil, middleOperationCompilationError
		}
		result = append(result, assignOperation...)
		result = append(result, vm.NewCode(vm.PushOP, errors.UnknownLine, nil))
	}
	var leftHandSide []*vm.Code
	var leftHandSideCompilationError *errors.Error
	switch assignStatement.LeftHandSide.(type) {
	case *Identifier:
		leftHandSide, leftHandSideCompilationError = compileIdentifierAssign(assignStatement.LeftHandSide.(*Identifier))
	case *SelectorExpression:
		leftHandSide, leftHandSideCompilationError = compileSelectorAssign(assignStatement.LeftHandSide.(*SelectorExpression))
	case *IndexExpression:
		leftHandSide, leftHandSideCompilationError = compileIndexAssign(assignStatement.LeftHandSide.(*IndexExpression))
	default:
		panic(reflect.TypeOf(assignStatement.LeftHandSide))
	}
	if leftHandSideCompilationError != nil {
		return nil, leftHandSideCompilationError
	}
	return append(result, leftHandSide...), nil
}

type DoWhileStatement struct {
	Statement
	Condition IExpression
	Body      []Node
}

func (doWhileStatement *DoWhileStatement) Compile() ([]*vm.Code, *errors.Error) {
	conditionReturn := &ReturnStatement{
		Results: []IExpression{doWhileStatement.Condition},
	}
	condition, conditionCompilationError := conditionReturn.Compile()
	if conditionCompilationError != nil {
		return nil, conditionCompilationError
	}
	body, bodyCompilationError := compileBody(doWhileStatement.Body)
	if bodyCompilationError != nil {
		return nil, bodyCompilationError
	}
	var result []*vm.Code
	result = append(result,
		vm.NewCode(vm.DoWhileLoopOP, errors.UnknownLine,
			vm.LoopInformation{
				Body:      body,
				Condition: condition,
				Receivers: nil,
			},
		),
	)

	return result, nil
}

type WhileLoopStatement struct {
	Statement
	Condition IExpression
	Body      []Node
}

func (whileStatement *WhileLoopStatement) Compile() ([]*vm.Code, *errors.Error) {
	conditionReturn := &ReturnStatement{
		Results: []IExpression{whileStatement.Condition},
	}
	condition, conditionCompilationError := conditionReturn.Compile()
	if conditionCompilationError != nil {
		return nil, conditionCompilationError
	}
	body, bodyCompilationError := compileBody(whileStatement.Body)
	if bodyCompilationError != nil {
		return nil, bodyCompilationError
	}
	var result []*vm.Code
	result = append(result,
		vm.NewCode(vm.WhileLoopOP, errors.UnknownLine,
			vm.LoopInformation{
				Body:      body,
				Condition: condition,
				Receivers: nil,
			},
		),
	)
	return result, nil
}

type UntilLoopStatement struct {
	Statement
	Condition IExpression
	Body      []Node
}

func (untilLoop *UntilLoopStatement) Compile() ([]*vm.Code, *errors.Error) {
	conditionReturn := &ReturnStatement{
		Results: []IExpression{untilLoop.Condition},
	}
	condition, conditionCompilationError := conditionReturn.Compile()
	if conditionCompilationError != nil {
		return nil, conditionCompilationError
	}
	body, bodyCompilationError := compileBody(untilLoop.Body)
	if bodyCompilationError != nil {
		return nil, bodyCompilationError
	}
	var result []*vm.Code
	result = append(result,
		vm.NewCode(vm.UntilLoopOP, errors.UnknownLine,
			vm.LoopInformation{
				Body:      body,
				Condition: condition,
				Receivers: nil,
			},
		),
	)
	return result, nil
}

type ForLoopStatement struct {
	Statement
	Receivers []*Identifier
	Source    IExpression
	Body      []Node
}

func (forStatement *ForLoopStatement) Compile() ([]*vm.Code, *errors.Error) {
	source, compilationError := forStatement.Source.CompilePush(true)
	if compilationError != nil {
		return nil, compilationError
	}
	body, bodyCompilationError := compileBody(forStatement.Body)
	if bodyCompilationError != nil {
		return nil, bodyCompilationError
	}
	var receivers []string
	for _, receiver := range forStatement.Receivers {
		receivers = append(
			receivers,
			receiver.Token.String,
		)
	}
	var result []*vm.Code
	result = append(result, source...)
	result = append(
		result,
		vm.NewCode(
			vm.ForLoopOP,
			errors.UnknownLine,
			vm.LoopInformation{
				Body:      body,
				Condition: nil,
				Receivers: receivers,
			},
		),
	)
	return result, nil
}

type IfStatement struct {
	Statement
	Condition IExpression
	Body      []Node
	Else      []Node
}

func (ifStatement *IfStatement) Compile() ([]*vm.Code, *errors.Error) {
	condition, conditionCompilationError := ifStatement.Condition.CompilePush(true)
	if conditionCompilationError != nil {
		return nil, conditionCompilationError
	}
	body, bodyCompilationError := compileBody(ifStatement.Body)
	if bodyCompilationError != nil {
		return nil, bodyCompilationError
	}
	var elseBody []*vm.Code
	if ifStatement.Else != nil {
		var elseBodyCompilationError *errors.Error
		elseBody, elseBodyCompilationError = compileBody(ifStatement.Else)
		if elseBodyCompilationError != nil {
			return nil, elseBodyCompilationError
		}
	}
	var result []*vm.Code
	result = append(result, condition...)
	result = append(result,
		vm.NewCode(vm.IfOP, errors.UnknownLine,
			vm.ConditionInformation{
				Body:     body,
				ElseBody: elseBody,
			},
		),
	)
	return result, nil
}

type UnlessStatement struct {
	Statement
	Condition IExpression
	Body      []Node
	Else      []Node
}

func (unlessStatement *UnlessStatement) Compile() ([]*vm.Code, *errors.Error) {
	condition, conditionCompilationError := unlessStatement.Condition.CompilePush(true)
	if conditionCompilationError != nil {
		return nil, conditionCompilationError
	}
	body, bodyCompilationError := compileBody(unlessStatement.Body)
	if bodyCompilationError != nil {
		return nil, bodyCompilationError
	}
	var elseBody []*vm.Code
	if unlessStatement.Body != nil {
		var elseBodyCompilationError *errors.Error
		elseBody, elseBodyCompilationError = compileBody(unlessStatement.Else)
		if elseBodyCompilationError != nil {
			return nil, elseBodyCompilationError
		}
	}
	var result []*vm.Code
	result = append(result, condition...)
	result = append(result,
		vm.NewCode(vm.UnlessOP, errors.UnknownLine,
			vm.ConditionInformation{
				Body:     body,
				ElseBody: elseBody,
			},
		),
	)
	return result, nil
}

type CaseBlock struct {
	Cases []IExpression
	Body  []Node
}

type SwitchStatement struct {
	Statement
	Target     IExpression
	CaseBlocks []*CaseBlock
	Default    []Node
}

func (switchStatement *SwitchStatement) Compile() ([]*vm.Code, *errors.Error) {
	reference, referenceCompilationError := switchStatement.Target.CompilePush(true)
	if referenceCompilationError != nil {
		return nil, referenceCompilationError
	}

	var caseBlocks []vm.CaseInformation
	for _, caseChunk := range switchStatement.CaseBlocks {
		targetsAsTuple := &ReturnStatement{
			Results: []IExpression{
				&TupleExpression{
					Values: caseChunk.Cases,
				},
			},
		}
		targets, targetCompilationError := targetsAsTuple.Compile()
		if targetCompilationError != nil {
			return nil, targetCompilationError
		}
		body, bodyCompilationError := compileBody(caseChunk.Body)
		if bodyCompilationError != nil {
			return nil, bodyCompilationError
		}
		caseBlocks = append(caseBlocks,
			vm.CaseInformation{
				Targets: targets,
				Body:    body,
			},
		)
	}

	var defaultBlock []*vm.Code
	if switchStatement.Default != nil {
		var defaultBlockCompilationError *errors.Error
		defaultBlock, defaultBlockCompilationError = compileBody(switchStatement.Default)
		if defaultBlockCompilationError != nil {
			return nil, defaultBlockCompilationError
		}
	}

	var result []*vm.Code
	result = append(result, reference...)
	result = append(result,
		vm.NewCode(vm.SwitchOP, errors.UnknownLine, vm.SwitchInformation{
			Cases:   caseBlocks,
			Default: defaultBlock,
		}))
	return result, nil
}

type ModuleStatement struct {
	Statement
	Name *Identifier
	Body []Node
}

func (moduleStatement *ModuleStatement) Compile() ([]*vm.Code, *errors.Error) {
	body, bodyCompilationError := compileBody(moduleStatement.Body)
	if bodyCompilationError != nil {
		return nil, bodyCompilationError
	}
	return []*vm.Code{
		vm.NewCode(
			vm.NewModuleOP,
			errors.UnknownLine,
			vm.ClassInformation{
				Name: moduleStatement.Name.Token.String,
				Body: body,
			},
		),
		vm.NewCode(vm.PushOP, errors.UnknownLine, nil),
		vm.NewCode(vm.AssignIdentifierOP, errors.UnknownLine, moduleStatement.Name.Token.String),
	}, nil
}

type FunctionDefinitionStatement struct {
	Statement
	Name      *Identifier
	Arguments []*Identifier
	Body      []Node
}

func (functionDefinition *FunctionDefinitionStatement) Compile() ([]*vm.Code, *errors.Error) {
	functionCode, functionDefinitionBodyCompilationError := compileBody(functionDefinition.Body)
	if functionDefinitionBodyCompilationError != nil {
		return nil, functionDefinitionBodyCompilationError
	}
	var arguments []string
	for _, argument := range functionDefinition.Arguments {
		arguments = append(arguments, argument.Token.String)
	}

	var body []*vm.Code
	body = append(body, vm.NewCode(vm.LoadFunctionArgumentsOP, errors.UnknownLine, arguments))
	body = append(body, functionCode...)
	body = append(body, vm.NewCode(vm.ReturnOP, errors.UnknownLine, 0))

	var result []*vm.Code
	result = append(result, vm.NewCode(vm.NewFunctionOP, errors.UnknownLine,
		vm.FunctionInformation{
			Name:              functionDefinition.Name.Token.String,
			Body:              body,
			NumberOfArguments: len(functionDefinition.Arguments),
		},
	))
	result = append(result,
		vm.NewCode(vm.PushOP, errors.UnknownLine, nil),
	)
	result = append(
		result,
		vm.NewCode(
			vm.AssignIdentifierOP,
			functionDefinition.Name.Token.Line,
			functionDefinition.Name.Token.String,
		),
	)
	return result, nil
}

func (functionDefinition *FunctionDefinitionStatement) CompileAsClassFunction() ([]*vm.Code, *errors.Error, bool) {
	functionCode, functionDefinitionBodyCompilationError := compileBody(functionDefinition.Body)
	if functionDefinitionBodyCompilationError != nil {
		return nil, functionDefinitionBodyCompilationError, false
	}

	var arguments []string
	for _, argument := range functionDefinition.Arguments {
		arguments = append(arguments, argument.Token.String)
	}
	var body []*vm.Code
	body = append(body, vm.NewCode(vm.LoadFunctionArgumentsOP, errors.UnknownLine, arguments))
	body = append(body, functionCode...)
	body = append(body, vm.NewCode(vm.ReturnOP, errors.UnknownLine, 0))
	var result []*vm.Code
	result = append(
		result,
		vm.NewCode(
			vm.NewClassFunctionOP,
			errors.UnknownLine,
			vm.FunctionInformation{
				Name:              functionDefinition.Name.Token.String,
				Body:              body,
				NumberOfArguments: len(arguments),
			},
		),
	)

	result = append(result,
		vm.NewCode(vm.PushOP, errors.UnknownLine, nil),
	)
	result = append(
		result,
		vm.NewCode(
			vm.AssignIdentifierOP,
			functionDefinition.Name.Token.Line,
			functionDefinition.Name.Token.String,
		),
	)
	return result, nil, functionDefinition.Name.Token.String == vm.Initialize
}

type InterfaceStatement struct {
	Statement
	Name              *Identifier
	Bases             []IExpression
	MethodDefinitions []*FunctionDefinitionStatement
}

func (interfaceStatement *InterfaceStatement) Compile() ([]*vm.Code, *errors.Error) {
	var body []*vm.Code
	for _, functionDefinition := range interfaceStatement.MethodDefinitions {
		function, functionCompilationError := functionDefinition.Compile()
		if functionCompilationError != nil {
			return nil, functionCompilationError
		}
		body = append(body, function...)
	}

	basesAsTuple := &TupleExpression{
		Values: interfaceStatement.Bases,
	}
	bases, basesCompilationError := basesAsTuple.CompilePush(true)
	if basesCompilationError != nil {
		return nil, basesCompilationError
	}
	var result []*vm.Code
	result = append(result, bases...)
	result = append(result, vm.NewCode(vm.NewClassOP, errors.UnknownLine,
		vm.ClassInformation{
			Name: interfaceStatement.Name.Token.String,
			Body: body,
		},
	))

	result = append(result, vm.NewCode(vm.PushOP, errors.UnknownLine, nil))
	result = append(result, vm.NewCode(vm.AssignIdentifierOP, errors.UnknownLine, interfaceStatement.Name.Token.String))
	return result, nil
}

type ClassStatement struct {
	Statement
	Name  *Identifier
	Bases []IExpression // Identifiers and selectors
	Body  []Node
}

func (classStatement *ClassStatement) Compile() ([]*vm.Code, *errors.Error) {
	basesAsTuple := &TupleExpression{
		Values: classStatement.Bases,
	}
	bases, basesCompilationError := basesAsTuple.CompilePush(true)
	if basesCompilationError != nil {
		return nil, basesCompilationError
	}
	body, bodyCompilationError := compileClassBody(classStatement.Body)
	if bodyCompilationError != nil {
		return nil, bodyCompilationError
	}
	var result []*vm.Code
	result = append(result, bases...)
	result = append(result, vm.NewCode(vm.NewClassOP, errors.UnknownLine,
		vm.ClassInformation{
			Name: classStatement.Name.Token.String,
			Body: body,
		},
	))

	result = append(result, vm.NewCode(vm.PushOP, errors.UnknownLine, nil))
	result = append(result, vm.NewCode(vm.AssignIdentifierOP, errors.UnknownLine, classStatement.Name.Token.String))
	return result, nil
}

type ExceptBlock struct {
	Targets     []IExpression
	CaptureName *Identifier
	Body        []Node
}

type RaiseStatement struct {
	Statement
	X IExpression
}

func (raise *RaiseStatement) Compile() ([]*vm.Code, *errors.Error) {
	result, expressionCompilationError := raise.X.CompilePush(true)
	if expressionCompilationError != nil {
		return nil, expressionCompilationError
	}
	result = append(result, vm.NewCode(vm.RaiseOP, errors.UnknownLine, nil))
	return result, nil
}

type TryStatement struct {
	Statement
	Body         []Node
	ExceptBlocks []*ExceptBlock
	Else         []Node
	Finally      []Node
}

func (tryStatement *TryStatement) Compile() ([]*vm.Code, *errors.Error) {
	body, bodyCompilationError := compileBody(tryStatement.Body)
	if bodyCompilationError != nil {
		return nil, bodyCompilationError
	}
	var exceptInformation []vm.ExceptInformation
	for _, except := range tryStatement.ExceptBlocks {
		captureName := ""
		if except.CaptureName != nil {
			captureName = except.CaptureName.Token.String
		}
		targetsReturn := &ReturnStatement{
			Results: []IExpression{
				&TupleExpression{
					Values: except.Targets,
				},
			},
		}
		targets, targetsCompilationError := targetsReturn.Compile()
		if targetsCompilationError != nil {
			return nil, targetsCompilationError
		}
		exceptBody, exceptCompilationError := compileBody(except.Body)
		if exceptCompilationError != nil {
			return nil, exceptCompilationError
		}

		exceptBlock := vm.ExceptInformation{
			CaptureName: captureName,
			Targets:     targets,
			Body:        exceptBody,
		}
		exceptInformation = append(exceptInformation, exceptBlock)
	}
	if tryStatement.Else != nil {
		elseBody, elseCompilationError := compileBody(tryStatement.Else)
		if elseCompilationError != nil {
			return nil, elseCompilationError
		}
		exceptInformation = append(exceptInformation,
			vm.ExceptInformation{
				CaptureName: "",
				Targets:     nil,
				Body:        elseBody,
			},
		)
	}
	var finally []*vm.Code
	if tryStatement.Finally != nil {
		var finallyCompilationError *errors.Error
		finally, finallyCompilationError = compileBody(tryStatement.Finally)
		if finallyCompilationError != nil {
			return nil, finallyCompilationError
		}
	}
	var result []*vm.Code
	result = append(result,
		vm.NewCode(vm.TryOP, errors.UnknownLine, vm.TryInformation{
			Body:    body,
			Excepts: exceptInformation,
			Finally: finally,
		}))
	return result, nil
}

type BeginStatement struct {
	Statement
	Body []Node
}

func (beginStatement *BeginStatement) Compile() ([]*vm.Code, *errors.Error) {
	return compileBody(beginStatement.Body)
}

type EndStatement struct {
	Statement
	Body []Node
}

func (endStatement *EndStatement) Compile() ([]*vm.Code, *errors.Error) {
	return compileBody(endStatement.Body)
}

type ReturnStatement struct {
	Statement
	Results []IExpression
}

func (returnStatement *ReturnStatement) Compile() ([]*vm.Code, *errors.Error) {
	numberOfResults := len(returnStatement.Results)
	var result []*vm.Code
	for i := numberOfResults - 1; i > -1; i-- {
		returnResult, resultCompilationError := returnStatement.Results[i].CompilePush(true)
		if resultCompilationError != nil {
			return nil, resultCompilationError
		}
		result = append(result, returnResult...)
	}
	return append(result, vm.NewCode(vm.ReturnOP, errors.UnknownLine, numberOfResults)), nil
}

type YieldStatement struct {
	Statement
	Results []IExpression
}

type ContinueStatement struct {
	Statement
}

func (_ *ContinueStatement) Compile() ([]*vm.Code, *errors.Error) {
	return []*vm.Code{vm.NewCode(vm.ContinueOP, errors.UnknownLine, nil)}, nil
}

type BreakStatement struct {
	Statement
}

func (_ *BreakStatement) Compile() ([]*vm.Code, *errors.Error) {
	return []*vm.Code{vm.NewCode(vm.BreakOP, errors.UnknownLine, nil)}, nil
}

type RedoStatement struct {
	Statement
}

func (_ *RedoStatement) Compile() ([]*vm.Code, *errors.Error) {
	return []*vm.Code{vm.NewCode(vm.RedoOP, errors.UnknownLine, nil)}, nil
}

type PassStatement struct {
	Statement
}

func (_ *PassStatement) Compile() ([]*vm.Code, *errors.Error) {
	return []*vm.Code{vm.NewCode(vm.NOP, errors.UnknownLine, nil)}, nil
}
