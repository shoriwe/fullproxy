package ast

import (
	"github.com/shoriwe/gplasma/pkg/compiler/lexer"
	"github.com/shoriwe/gplasma/pkg/errors"
	"github.com/shoriwe/gplasma/pkg/tools"
	"github.com/shoriwe/gplasma/pkg/vm"
	"strconv"
	"strings"
)

type IExpression interface {
	Node
	E()
}

type ArrayExpression struct {
	IExpression
	Values []IExpression
}

func (arrayExpression *ArrayExpression) Compile() ([]*vm.Code, *errors.Error) {
	valuesLength := len(arrayExpression.Values)
	var result []*vm.Code
	for i := valuesLength - 1; i > -1; i-- {
		childExpression, valueCompilationError := arrayExpression.Values[i].CompilePush(true)
		if valueCompilationError != nil {
			return nil, valueCompilationError
		}
		result = append(result, childExpression...)
	}
	return append(result, vm.NewCode(vm.NewArrayOP, errors.UnknownLine, len(arrayExpression.Values))), nil
}

func (arrayExpression *ArrayExpression) CompilePush(push bool) ([]*vm.Code, *errors.Error) {
	result, compilationError := arrayExpression.Compile()
	if compilationError != nil {
		return nil, compilationError
	}
	if push {
		result = append(result, vm.NewCode(vm.PushOP, errors.UnknownLine, nil))
	}
	return result, nil
}

type TupleExpression struct {
	IExpression
	Values []IExpression
}

func (tupleExpression *TupleExpression) Compile() ([]*vm.Code, *errors.Error) {
	valuesLength := len(tupleExpression.Values)
	var result []*vm.Code
	for i := valuesLength - 1; i > -1; i-- {
		childExpression, valueCompilationError := tupleExpression.Values[i].CompilePush(true)
		if valueCompilationError != nil {
			return nil, valueCompilationError
		}
		result = append(result, childExpression...)
	}
	return append(result, vm.NewCode(vm.NewTupleOP, errors.UnknownLine, len(tupleExpression.Values))), nil
}

func (tupleExpression *TupleExpression) CompilePush(push bool) ([]*vm.Code, *errors.Error) {
	result, compilationError := tupleExpression.Compile()
	if compilationError != nil {
		return nil, compilationError
	}
	if push {
		result = append(result, vm.NewCode(vm.PushOP, errors.UnknownLine, nil))
	}
	return result, nil
}

type KeyValue struct {
	Key   IExpression
	Value IExpression
}

type HashExpression struct {
	IExpression
	Values []*KeyValue
}

func (hashExpression *HashExpression) CompilePush(push bool) ([]*vm.Code, *errors.Error) {
	result, compilationError := hashExpression.Compile()
	if compilationError != nil {
		return nil, compilationError
	}
	if push {
		result = append(result, vm.NewCode(vm.PushOP, errors.UnknownLine, nil))
	}
	return result, nil
}

func (hashExpression *HashExpression) Compile() ([]*vm.Code, *errors.Error) {
	valuesLength := len(hashExpression.Values)
	var result []*vm.Code
	for i := valuesLength - 1; i > -1; i-- {
		key, valueCompilationError := hashExpression.Values[i].Value.CompilePush(true)
		if valueCompilationError != nil {
			return nil, valueCompilationError
		}
		result = append(result, key...)
		value, keyCompilationError := hashExpression.Values[i].Key.CompilePush(true)
		if keyCompilationError != nil {
			return nil, keyCompilationError
		}
		result = append(result, value...)
	}
	return append(result, vm.NewCode(vm.NewHashOP, errors.UnknownLine, len(hashExpression.Values))), nil
}

type Identifier struct {
	IExpression
	Token *lexer.Token
}

func (identifier *Identifier) CompilePush(push bool) ([]*vm.Code, *errors.Error) {
	result, compilationError := identifier.Compile()
	if compilationError != nil {
		return nil, compilationError
	}
	if push {
		result = append(result, vm.NewCode(vm.PushOP, errors.UnknownLine, nil))
	}
	return result, nil
}

func (identifier *Identifier) Compile() ([]*vm.Code, *errors.Error) {
	return []*vm.Code{vm.NewCode(vm.GetIdentifierOP, identifier.Token.Line, identifier.Token.String)}, nil
}

type BasicLiteralExpression struct {
	IExpression
	Token       *lexer.Token
	Kind        uint8
	DirectValue uint8
}

func (basicLiteral *BasicLiteralExpression) CompilePush(push bool) ([]*vm.Code, *errors.Error) {
	result, compilationError := basicLiteral.Compile()
	if compilationError != nil {
		return nil, compilationError
	}
	if push {
		result = append(result, vm.NewCode(vm.PushOP, errors.UnknownLine, nil))
	}
	return result, nil
}

func (basicLiteral *BasicLiteralExpression) Compile() ([]*vm.Code, *errors.Error) {
	switch basicLiteral.DirectValue {
	case lexer.Integer:
		numberString := basicLiteral.Token.String
		numberString = strings.ReplaceAll(numberString, "_", "")
		number, success := strconv.ParseInt(numberString, 10, 64)
		if success != nil {
			return nil, errors.New(basicLiteral.Token.Line, "Error parsing Integer", errors.GoRuntimeError)
		}
		return []*vm.Code{vm.NewCode(vm.NewIntegerOP, basicLiteral.Token.Line, number)}, nil
	case lexer.HexadecimalInteger:
		numberString := basicLiteral.Token.String
		numberString = strings.ReplaceAll(numberString, "_", "")
		numberString = numberString[2:]
		number, parsingError := strconv.ParseInt(numberString, 16, 64)
		if parsingError != nil {
			return nil, errors.New(basicLiteral.Token.Line, "Error parsing Hexadecimal Integer", errors.GoRuntimeError)
		}
		return []*vm.Code{vm.NewCode(vm.NewIntegerOP, basicLiteral.Token.Line, number)}, nil
	case lexer.OctalInteger:
		numberString := basicLiteral.Token.String
		numberString = strings.ReplaceAll(numberString, "_", "")
		numberString = numberString[2:]
		number, parsingError := strconv.ParseInt(numberString, 8, 64)
		if parsingError != nil {
			return nil, errors.New(basicLiteral.Token.Line, "Error parsing Octal Integer", errors.GoRuntimeError)
		}
		return []*vm.Code{vm.NewCode(vm.NewIntegerOP, basicLiteral.Token.Line, number)}, nil
	case lexer.BinaryInteger:
		numberString := basicLiteral.Token.String
		numberString = strings.ReplaceAll(numberString, "_", "")
		numberString = numberString[2:]
		number, parsingError := strconv.ParseInt(numberString, 2, 64)
		if parsingError != nil {
			return nil, errors.New(basicLiteral.Token.Line, "Error parsing Binary Integer", errors.GoRuntimeError)
		}
		return []*vm.Code{vm.NewCode(vm.NewIntegerOP, basicLiteral.Token.Line, number)}, nil
	case lexer.Float, lexer.ScientificFloat:
		numberString := basicLiteral.Token.String
		numberString = strings.ReplaceAll(numberString, "_", "")
		number, parsingError := strconv.ParseFloat(numberString, 64)
		if parsingError != nil {
			return nil, errors.New(basicLiteral.Token.Line, parsingError.Error(), errors.GoRuntimeError)
		}
		return []*vm.Code{vm.NewCode(vm.NewFloatOP, basicLiteral.Token.Line, number)}, nil
	case lexer.SingleQuoteString, lexer.DoubleQuoteString:
		return []*vm.Code{vm.NewCode(
			vm.NewStringOP, basicLiteral.Token.Line,

			string(
				tools.ReplaceEscaped(
					[]rune(basicLiteral.Token.String[1:len(basicLiteral.Token.String)-1])),
			),
		),
		}, nil
	case lexer.ByteString:
		return []*vm.Code{vm.NewCode(vm.NewBytesOP, basicLiteral.Token.Line,
			[]byte(
				string(
					tools.ReplaceEscaped(
						[]rune(basicLiteral.Token.String[2:len(basicLiteral.Token.String)-1]),
					),
				),
			),
		),
		}, nil
	case lexer.True:
		return []*vm.Code{vm.NewCode(vm.GetTrueOP, basicLiteral.Token.Line, nil)}, nil
	case lexer.False:
		return []*vm.Code{vm.NewCode(vm.GetFalseOP, basicLiteral.Token.Line, nil)}, nil
	case lexer.None:
		return []*vm.Code{vm.NewCode(vm.GetNoneOP, basicLiteral.Token.Line, nil)}, nil
	}
	panic(errors.NewUnknownVMOperationError(basicLiteral.Token.DirectValue))
}

type BinaryExpression struct {
	IExpression
	LeftHandSide  IExpression
	Operator      *lexer.Token
	RightHandSide IExpression
}

func (binaryExpression *BinaryExpression) CompilePush(push bool) ([]*vm.Code, *errors.Error) {
	result, compilationError := binaryExpression.Compile()
	if compilationError != nil {
		return nil, compilationError
	}
	if push {
		result = append(result, vm.NewCode(vm.PushOP, errors.UnknownLine, nil))
	}
	return result, nil
}

func (binaryExpression *BinaryExpression) Compile() ([]*vm.Code, *errors.Error) {
	var result []*vm.Code
	// CompilePush first right hand side
	right, rightHandSideCompileError := binaryExpression.RightHandSide.CompilePush(true)
	if rightHandSideCompileError != nil {
		return nil, rightHandSideCompileError
	}
	result = append(result, right...)
	// Then left hand side
	left, leftHandSideCompileError := binaryExpression.LeftHandSide.CompilePush(true)
	if leftHandSideCompileError != nil {
		return nil, leftHandSideCompileError
	}
	result = append(result, left...)
	var operation uint8
	// Finally decide the instruction to use
	switch binaryExpression.Operator.DirectValue {
	case lexer.Add:
		operation = vm.AddOP
	case lexer.Sub:
		operation = vm.SubOP
	case lexer.Star:
		operation = vm.MulOP
	case lexer.Div:
		operation = vm.DivOP
	case lexer.FloorDiv:
		operation = vm.FloorDivOP
	case lexer.Modulus:
		operation = vm.ModOP
	case lexer.PowerOf:
		operation = vm.PowOP
	case lexer.BitwiseXor:
		operation = vm.BitXorOP
	case lexer.BitWiseAnd:
		operation = vm.BitAndOP
	case lexer.BitwiseOr:
		operation = vm.BitOrOP
	case lexer.BitwiseLeft:
		operation = vm.BitLeftOP
	case lexer.BitwiseRight:
		operation = vm.BitRightOP
	case lexer.And:
		operation = vm.AndOP
	case lexer.Or:
		operation = vm.OrOP
	case lexer.Xor:
		operation = vm.XorOP
	case lexer.Equals:
		operation = vm.EqualsOP
	case lexer.NotEqual:
		operation = vm.NotEqualsOP
	case lexer.GreaterThan:
		operation = vm.GreaterThanOP
	case lexer.LessThan:
		operation = vm.LessThanOP
	case lexer.GreaterOrEqualThan:
		operation = vm.GreaterThanOrEqualOP
	case lexer.LessOrEqualThan:
		operation = vm.LessThanOrEqualOP
	case lexer.In:
		operation = vm.ContainsOP
	default:
		panic(errors.NewUnknownVMOperationError(binaryExpression.Operator.DirectValue))
	}
	result = append(result, vm.NewCode(vm.BinaryOP, binaryExpression.Operator.Line, operation))
	return result, nil
}

type UnaryExpression struct {
	IExpression
	Operator *lexer.Token
	X        IExpression
}

func (unaryExpression *UnaryExpression) CompilePush(push bool) ([]*vm.Code, *errors.Error) {
	result, compilationError := unaryExpression.Compile()
	if compilationError != nil {
		return nil, compilationError
	}
	if push {
		result = append(result, vm.NewCode(vm.PushOP, errors.UnknownLine, nil))
	}
	return result, nil
}

func (unaryExpression *UnaryExpression) Compile() ([]*vm.Code, *errors.Error) {
	result, expressionCompileError := unaryExpression.X.CompilePush(true)
	if expressionCompileError != nil {
		return nil, expressionCompileError
	}
	switch unaryExpression.Operator.DirectValue {
	case lexer.NegateBits:
		result = append(result, vm.NewCode(vm.UnaryOP, unaryExpression.Operator.Line, vm.NegateBitsOP))
	case lexer.Not, lexer.SignNot:
		result = append(result, vm.NewCode(vm.UnaryOP, unaryExpression.Operator.Line, vm.BoolNegateOP))
	case lexer.Sub:
		result = append(result, vm.NewCode(vm.UnaryOP, unaryExpression.Operator.Line, vm.NegativeOP))
	}
	return result, nil
}

type ParenthesesExpression struct {
	IExpression
	X IExpression
}

func (parenthesesExpression *ParenthesesExpression) CompilePush(push bool) ([]*vm.Code, *errors.Error) {
	result, compilationError := parenthesesExpression.Compile()
	if compilationError != nil {
		return nil, compilationError
	}
	if push {
		result = append(result, vm.NewCode(vm.PushOP, errors.UnknownLine, nil))
	}
	return result, nil
}

func (parenthesesExpression *ParenthesesExpression) Compile() ([]*vm.Code, *errors.Error) {
	return parenthesesExpression.X.CompilePush(false)
}

type LambdaExpression struct {
	IExpression
	Arguments []*Identifier
	Code      IExpression
}

func (lambdaExpression *LambdaExpression) CompilePush(push bool) ([]*vm.Code, *errors.Error) {
	result, compilationError := lambdaExpression.Compile()
	if compilationError != nil {
		return nil, compilationError
	}
	if push {
		result = append(result, vm.NewCode(vm.PushOP, errors.UnknownLine, nil))
	}
	return result, nil
}

func (lambdaExpression *LambdaExpression) Compile() ([]*vm.Code, *errors.Error) {
	functionCode, lambdaCodeCompilationError := lambdaExpression.Code.CompilePush(true)
	if lambdaCodeCompilationError != nil {
		return nil, lambdaCodeCompilationError
	}
	var arguments []string
	for _, argument := range lambdaExpression.Arguments {
		arguments = append(arguments, argument.Token.String)
	}
	functionBody := []*vm.Code{vm.NewCode(vm.LoadFunctionArgumentsOP, errors.UnknownLine, arguments)}
	functionBody = append(functionBody, functionCode...)
	functionBody = append(functionBody, vm.NewCode(vm.ReturnOP, errors.UnknownLine, 1))
	var result []*vm.Code
	result = append(result,
		vm.NewCode(
			vm.NewFunctionOP,
			errors.UnknownLine,
			vm.FunctionInformation{
				Name:              "",
				Body:              functionBody,
				NumberOfArguments: len(arguments),
			},
		),
	)
	return result, nil
}

type GeneratorExpression struct {
	IExpression
	Operation IExpression
	Receivers []*Identifier
	Source    IExpression
}

func (generatorExpression *GeneratorExpression) CompilePush(push bool) ([]*vm.Code, *errors.Error) {
	result, compilationError := generatorExpression.Compile()
	if compilationError != nil {
		return nil, compilationError
	}
	if push {
		result = append(result, vm.NewCode(vm.PushOP, errors.UnknownLine, nil))
	}
	return result, nil
}

func (generatorExpression *GeneratorExpression) Compile() ([]*vm.Code, *errors.Error) {
	source, sourceCompilationError := generatorExpression.Source.CompilePush(true)
	if sourceCompilationError != nil {
		return nil, sourceCompilationError
	}
	operationAsLambda := &LambdaExpression{
		Arguments: generatorExpression.Receivers,
		Code:      generatorExpression.Operation,
	}
	operation, operationCompilationError := operationAsLambda.CompilePush(true)
	if operationCompilationError != nil {
		return nil, operationCompilationError
	}

	var result []*vm.Code
	result = append(result, source...)
	result = append(result, operation...)
	result = append(result,
		vm.NewCode(vm.NewGeneratorOP, errors.UnknownLine, len(generatorExpression.Receivers)),
	)
	return result, nil
}

type SelectorExpression struct {
	IExpression
	X          IExpression
	Identifier *Identifier
}

func (selectorExpression *SelectorExpression) CompilePush(push bool) ([]*vm.Code, *errors.Error) {
	result, compilationError := selectorExpression.Compile()
	if compilationError != nil {
		return nil, compilationError
	}
	if push {
		result = append(result, vm.NewCode(vm.PushOP, errors.UnknownLine, nil))
	}
	return result, nil
}

func (selectorExpression *SelectorExpression) Compile() ([]*vm.Code, *errors.Error) {
	source, sourceCompilationError := selectorExpression.X.CompilePush(true)
	if sourceCompilationError != nil {
		return nil, sourceCompilationError
	}
	return append(source, vm.NewCode(vm.SelectNameFromObjectOP, selectorExpression.Identifier.Token.Line, selectorExpression.Identifier.Token.String)), nil
}

type MethodInvocationExpression struct {
	IExpression
	Function  IExpression
	Arguments []IExpression
}

func (methodInvocationExpression *MethodInvocationExpression) CompilePush(push bool) ([]*vm.Code, *errors.Error) {
	result, compilationError := methodInvocationExpression.Compile()
	if compilationError != nil {
		return nil, compilationError
	}
	if push {
		result = append(result, vm.NewCode(vm.PushOP, errors.UnknownLine, nil))
	}
	return result, nil
}

func (methodInvocationExpression *MethodInvocationExpression) Compile() ([]*vm.Code, *errors.Error) {
	numberOfArguments := len(methodInvocationExpression.Arguments)
	var result []*vm.Code
	for i := numberOfArguments - 1; i > -1; i-- {
		argument, argumentCompilationError := methodInvocationExpression.Arguments[i].CompilePush(true)
		if argumentCompilationError != nil {
			return nil, argumentCompilationError
		}
		result = append(result, argument...)
	}
	function, functionCompilationError := methodInvocationExpression.Function.CompilePush(true)
	if functionCompilationError != nil {
		return nil, functionCompilationError
	}
	result = append(result, function...)
	return append(result, vm.NewCode(vm.MethodInvocationOP, errors.UnknownLine, len(methodInvocationExpression.Arguments))), nil
}

type IndexExpression struct {
	IExpression
	Source IExpression
	Index  IExpression
}

func (indexExpression *IndexExpression) CompilePush(push bool) ([]*vm.Code, *errors.Error) {
	result, compilationError := indexExpression.Compile()
	if compilationError != nil {
		return nil, compilationError
	}
	if push {
		result = append(result, vm.NewCode(vm.PushOP, errors.UnknownLine, nil))
	}
	return result, nil
}

func (indexExpression *IndexExpression) Compile() ([]*vm.Code, *errors.Error) {
	source, sourceCompilationError := indexExpression.Source.CompilePush(true)
	if sourceCompilationError != nil {
		return nil, sourceCompilationError
	}
	index, indexCompilationError := indexExpression.Index.CompilePush(true)
	if indexCompilationError != nil {
		return nil, indexCompilationError
	}
	result := source
	result = append(result, index...)
	return append(result, vm.NewCode(vm.IndexOP, errors.UnknownLine, nil)), nil
}

type IfOneLinerExpression struct {
	IExpression
	Result     IExpression
	Condition  IExpression
	ElseResult IExpression
}

func (ifOneLinerExpression *IfOneLinerExpression) CompilePush(push bool) ([]*vm.Code, *errors.Error) {
	result, compilationError := ifOneLinerExpression.Compile()
	if compilationError != nil {
		return nil, compilationError
	}
	if push {
		result = append(result, vm.NewCode(vm.PushOP, errors.UnknownLine, nil))
	}
	return result, nil
}

func (ifOneLinerExpression *IfOneLinerExpression) Compile() ([]*vm.Code, *errors.Error) {
	condition, conditionCompilationError := ifOneLinerExpression.Condition.CompilePush(true)
	if conditionCompilationError != nil {
		return nil, conditionCompilationError
	}
	ifReturnResult := &ReturnStatement{
		Results: []IExpression{ifOneLinerExpression.Result},
	}
	ifResult, ifResultCompilationError := ifReturnResult.Compile()
	if ifResultCompilationError != nil {
		return nil, ifResultCompilationError
	}
	var elseResult []*vm.Code
	if ifOneLinerExpression.ElseResult != nil {
		elseReturnResult := &ReturnStatement{
			Results: []IExpression{ifOneLinerExpression.ElseResult},
		}
		var elseResultCompilationError *errors.Error
		elseResult, elseResultCompilationError = elseReturnResult.Compile()
		if elseResultCompilationError != nil {
			return nil, elseResultCompilationError
		}
	} else {
		elseResult = []*vm.Code{
			vm.NewCode(vm.GetNoneOP, errors.UnknownLine, nil),
			vm.NewCode(vm.PushOP, errors.UnknownLine, nil),
			vm.NewCode(vm.ReturnOP, errors.UnknownLine, 1),
		}
	}
	var result []*vm.Code
	result = append(result, condition...)
	result = append(result,
		vm.NewCode(vm.IfOneLinerOP,
			errors.UnknownLine,
			vm.ConditionInformation{
				Body:     ifResult,
				ElseBody: elseResult,
			},
		),
	)
	return result, nil
}

type UnlessOneLinerExpression struct {
	IExpression
	Result     IExpression
	Condition  IExpression
	ElseResult IExpression
}

func (unlessOneLinerExpression *UnlessOneLinerExpression) CompilePush(push bool) ([]*vm.Code, *errors.Error) {
	result, compilationError := unlessOneLinerExpression.Compile()
	if compilationError != nil {
		return nil, compilationError
	}
	if push {
		result = append(result, vm.NewCode(vm.PushOP, errors.UnknownLine, nil))
	}
	return result, nil
}

func (unlessOneLinerExpression *UnlessOneLinerExpression) Compile() ([]*vm.Code, *errors.Error) {
	condition, conditionCompilationError := unlessOneLinerExpression.Condition.CompilePush(true)
	if conditionCompilationError != nil {
		return nil, conditionCompilationError
	}
	ifReturnResult := &ReturnStatement{
		Results: []IExpression{unlessOneLinerExpression.Result},
	}
	ifResult, ifResultCompilationError := ifReturnResult.Compile()
	if ifResultCompilationError != nil {
		return nil, ifResultCompilationError
	}
	var elseResult []*vm.Code
	if unlessOneLinerExpression.ElseResult != nil {
		elseReturnResult := &ReturnStatement{
			Results: []IExpression{unlessOneLinerExpression.ElseResult},
		}
		var elseResultCompilationError *errors.Error
		elseResult, elseResultCompilationError = elseReturnResult.Compile()
		if elseResultCompilationError != nil {
			return nil, elseResultCompilationError
		}
	} else {
		elseResult = []*vm.Code{
			vm.NewCode(vm.GetNoneOP, errors.UnknownLine, nil),
			vm.NewCode(vm.PushOP, errors.UnknownLine, nil),
			vm.NewCode(vm.ReturnOP, errors.UnknownLine, 1),
		}
	}
	var result []*vm.Code
	result = append(result, condition...)
	result = append(result,
		vm.NewCode(vm.UnlessOneLinerOP,
			errors.UnknownLine,
			vm.ConditionInformation{
				Body:     ifResult,
				ElseBody: elseResult,
			},
		),
	)
	return result, nil
}
