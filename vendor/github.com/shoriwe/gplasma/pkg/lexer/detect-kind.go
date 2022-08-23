package lexer

func (lexer *Lexer) detectKindAndDirectValue() (Kind, DirectValue) {
	s := lexer.currentToken.String()
	switch s {
	case PassString:
		return Keyword, Pass
	case SuperString:
		return Keyword, Super
	case DeleteString:
		return Keyword, Delete
	case EndString:
		return Keyword, End
	case IfString:
		return Keyword, If
	case UnlessString:
		return Keyword, Unless
	case ElseString:
		return Keyword, Else
	case ElifString:
		return Keyword, Elif
	case WhileString:
		return Keyword, While
	case DoString:
		return Keyword, Do
	case ForString:
		return Keyword, For
	case UntilString:
		return Keyword, Until
	case SwitchString:
		return Keyword, Switch
	case CaseString:
		return Keyword, Case
	case DefaultString:
		return Keyword, Default
	case YieldString:
		return Keyword, Yield
	case ReturnString:
		return Keyword, Return
	case ContinueString:
		return Keyword, Continue
	case BreakString:
		return Keyword, Break
	case ModuleString:
		return Keyword, Module
	case DefString:
		return Keyword, Def
	case GeneratorString:
		return Keyword, Generator
	case LambdaString:
		return Keyword, Lambda
	case InterfaceString:
		return Keyword, Interface
	case ClassString:
		return Keyword, Class
	case AndString:
		return Comparator, And
	case OrString:
		return Comparator, Or
	case XorString:
		return Comparator, Xor
	case InString:
		return Comparator, In
	case IsString:
		return Comparator, Is
	case ImplementsString:
		return Comparator, Implements
	case BEGINString:
		return Keyword, BEGIN
	case ENDString:
		return Keyword, END
	case NotString: // Unary operator
		return Operator, Not
	case TrueString:
		return Boolean, True
	case FalseString:
		return Boolean, False
	case NoneString:
		return NoneType, None
	case DeferString:
		return Keyword, Defer
	default:
		if identifierCheck.MatchString(s) {
			return IdentifierKind, InvalidDirectValue
		} else if junkKindCheck.MatchString(s) {
			return JunkKind, InvalidDirectValue
		}
	}
	return Unknown, InvalidDirectValue
}
