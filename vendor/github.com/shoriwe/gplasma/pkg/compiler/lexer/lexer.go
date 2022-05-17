package lexer

import (
	"github.com/shoriwe/gplasma/pkg/errors"
	"github.com/shoriwe/gplasma/pkg/reader"
	"regexp"
)

var (
	identifierCheck = regexp.MustCompile("(?m)^[a-zA-Z_]+[a-zA-Z0-9_]*$")
	junkKindCheck   = regexp.MustCompile("(?m)^\\00+$")
)

type Lexer struct {
	lastToken *Token
	line      int
	reader    reader.Reader
	complete  bool
}

func (lexer *Lexer) HasNext() bool {
	return !lexer.complete
}

func (lexer *Lexer) tokenizeStringLikeExpressions(stringOpener rune) ([]rune, uint8, uint8, *errors.Error) {
	content := []rune{stringOpener}
	var target uint8
	switch stringOpener {
	case '\'':
		target = SingleQuoteString
	case '"':
		target = DoubleQuoteString
	case '`':
		target = CommandOutput
	}
	var directValue = Unknown
	var tokenizingError *errors.Error
	escaped := false
	finish := false
	for ; lexer.reader.HasNext() && !finish; lexer.reader.Next() {
		char := lexer.reader.Char()
		if escaped {
			switch char {
			case '\\', '\'', '"', '`', 'a', 'b', 'e', 'f', 'n', 'r', 't', '?', 'u', 'x':
				escaped = false
			default:
				tokenizingError = errors.New(lexer.line, "invalid escape sequence", errors.LexingError)
				finish = true
			}
		} else {
			switch char {
			case '\n':
				lexer.line++
			case stringOpener:
				directValue = target
				finish = true
			case '\\':
				escaped = true
			}
		}
		content = append(content, char)
	}
	if directValue != target {
		tokenizingError = errors.New(lexer.line, "string never closed", errors.LexingError)
	}
	return content, Literal, directValue, tokenizingError
}

func (lexer *Lexer) tokenizeHexadecimal(letterX rune) ([]rune, uint8, uint8, *errors.Error) {
	result := []rune{'0', letterX}
	if !lexer.reader.HasNext() {
		return result, Literal, Unknown, errors.NewUnknownTokenKindError(lexer.line)
	}
	nextDigit := lexer.reader.Char()
	if !(('0' <= nextDigit && nextDigit <= '9') ||
		('a' <= nextDigit && nextDigit <= 'f') ||
		('A' <= nextDigit && nextDigit <= 'F')) {
		return result, Literal, Unknown, errors.NewUnknownTokenKindError(lexer.line)
	}
	lexer.reader.Next()
	result = append(result, nextDigit)
	for ; lexer.reader.HasNext(); lexer.reader.Next() {
		nextDigit = lexer.reader.Char()
		if !(('0' <= nextDigit && nextDigit <= '9') || ('a' <= nextDigit && nextDigit <= 'f') || ('A' <= nextDigit && nextDigit <= 'F')) && nextDigit != '_' {
			return result, Literal, HexadecimalInteger, nil
		}
		result = append(result, nextDigit)
	}
	return result, Literal, HexadecimalInteger, nil
}

func (lexer *Lexer) tokenizeBinary(letterB rune) ([]rune, uint8, uint8, *errors.Error) {
	result := []rune{'0', letterB}
	if !lexer.reader.HasNext() {
		return result, Literal, Unknown, errors.NewUnknownTokenKindError(lexer.line)
	}
	nextDigit := lexer.reader.Char()
	if !(nextDigit == '0' || nextDigit == '1') {
		return result, Literal, Unknown, errors.NewUnknownTokenKindError(lexer.line)
	}
	lexer.reader.Next()
	result = append(result, nextDigit)
	for ; lexer.reader.HasNext(); lexer.reader.Next() {
		nextDigit = lexer.reader.Char()
		if !(nextDigit == '0' || nextDigit == '1') && nextDigit != '_' {
			return result, Literal, BinaryInteger, nil
		}
		result = append(result, nextDigit)
	}
	return result, Literal, BinaryInteger, nil
}

func (lexer *Lexer) tokenizeOctal(letterO rune) ([]rune, uint8, uint8, *errors.Error) {
	result := []rune{'0', letterO}
	if !lexer.reader.HasNext() {
		return result, Literal, Unknown, errors.NewUnknownTokenKindError(lexer.line)
	}
	nextDigit := lexer.reader.Char()
	if !('0' <= nextDigit && nextDigit <= '7') {
		return result, Literal, Unknown, errors.NewUnknownTokenKindError(lexer.line)
	}
	lexer.reader.Next()
	result = append(result, nextDigit)
	for ; lexer.reader.HasNext(); lexer.reader.Next() {
		nextDigit = lexer.reader.Char()
		if !('0' <= nextDigit && nextDigit <= '7') && nextDigit != '_' {
			return result, Literal, OctalInteger, nil
		}
		result = append(result, nextDigit)
	}
	return result, Literal, OctalInteger, nil
}

func (lexer *Lexer) tokenizeScientificFloat(base []rune) ([]rune, uint8, uint8, *errors.Error) {
	result := base
	if !lexer.reader.HasNext() {
		return result, Literal, Unknown, errors.NewUnknownTokenKindError(lexer.line)
	}
	direction := lexer.reader.Char()
	if (direction != '-') && (direction != '+') {
		return result, Literal, Unknown, errors.NewUnknownTokenKindError(lexer.line)
	}
	lexer.reader.Next()
	// Ensure next is a number
	if !lexer.reader.HasNext() {
		return result, Literal, Unknown, errors.NewUnknownTokenKindError(lexer.line)
	}
	result = append(result, direction)
	nextDigit := lexer.reader.Char()
	if !('0' <= nextDigit && nextDigit <= '9') {
		return result, Literal, Unknown, errors.NewUnknownTokenKindError(lexer.line)
	}
	result = append(result, nextDigit)
	lexer.reader.Next()
	for ; lexer.reader.HasNext(); lexer.reader.Next() {
		nextDigit = lexer.reader.Char()
		if !('0' <= nextDigit && nextDigit <= '9') && nextDigit != '_' {
			return result, Literal, ScientificFloat, nil
		}
		result = append(result, nextDigit)
	}
	return result, Literal, ScientificFloat, nil
}

func (lexer *Lexer) tokenizeFloat(base []rune) ([]rune, uint8, uint8, *errors.Error) {
	if !lexer.reader.HasNext() {
		lexer.reader.Redo()
		return base[:len(base)-1], Literal, Integer, nil
	}
	nextDigit := lexer.reader.Char()
	if !('0' <= nextDigit && nextDigit <= '9') {
		lexer.reader.Redo()
		return base[:len(base)-1], Literal, Integer, nil
	}
	lexer.reader.Next()
	result := append(base, nextDigit)
	for ; lexer.reader.HasNext(); lexer.reader.Next() {
		nextDigit = lexer.reader.Char()
		if ('0' <= nextDigit && nextDigit <= '9') || nextDigit == '_' {
			result = append(result, nextDigit)
		} else if (nextDigit == 'e') || (nextDigit == 'E') {
			lexer.reader.Next()
			return lexer.tokenizeScientificFloat(append(result, nextDigit))
		} else {
			break
		}
	}
	return result, Literal, Float, nil
}

func (lexer *Lexer) tokenizeInteger(base []rune) ([]rune, uint8, uint8, *errors.Error) {
	if !lexer.reader.HasNext() {
		return base, Literal, Integer, nil
	}
	nextDigit := lexer.reader.Char()
	if nextDigit == '.' {
		lexer.reader.Next()
		return lexer.tokenizeFloat(append(base, nextDigit))
	} else if nextDigit == 'e' || nextDigit == 'E' {
		lexer.reader.Next()
		return lexer.tokenizeScientificFloat(append(base, nextDigit))
	} else if !('0' <= nextDigit && nextDigit <= '9') {
		return base, Literal, Integer, nil
	}
	lexer.reader.Next()
	result := append(base, nextDigit)
	for ; lexer.reader.HasNext(); lexer.reader.Next() {
		nextDigit = lexer.reader.Char()
		if nextDigit == 'e' || nextDigit == 'E' {
			return lexer.tokenizeScientificFloat(append(result, nextDigit))
		} else if nextDigit == '.' {
			lexer.reader.Next()
			return lexer.tokenizeFloat(append(result, nextDigit))
		} else if ('0' <= nextDigit && nextDigit <= '9') || nextDigit == '_' {
			result = append(result, nextDigit)
		} else {
			break
		}
	}
	return result, Literal, Integer, nil
}

func (lexer *Lexer) tokenizeNumeric(firstDigit rune) ([]rune, uint8, uint8, *errors.Error) {
	if !lexer.reader.HasNext() {
		return []rune{firstDigit}, Literal, Integer, nil
	}
	nextChar := lexer.reader.Char()
	lexer.reader.Next()
	if firstDigit == '0' {
		switch nextChar {
		case 'x', 'X': // Hexadecimal
			return lexer.tokenizeHexadecimal(nextChar)
		case 'b', 'B': // Binary
			return lexer.tokenizeBinary(nextChar)
		case 'o', 'O': // Octal
			return lexer.tokenizeOctal(nextChar)
		case 'e', 'E': // Scientific float
			return lexer.tokenizeScientificFloat([]rune{firstDigit, nextChar})
		case '.': // Maybe a float
			return lexer.tokenizeFloat([]rune{firstDigit, nextChar}) // Integer, Float Or Scientific Float
		default:
			if ('0' <= nextChar && nextChar <= '9') || nextChar == '_' {
				return lexer.tokenizeInteger([]rune{firstDigit, nextChar}) // Integer, Float or Scientific Float
			}
		}
	} else {
		switch nextChar {
		case 'e', 'E': // Scientific float
			return lexer.tokenizeScientificFloat([]rune{firstDigit, nextChar})
		case '.': // Maybe a float
			return lexer.tokenizeFloat([]rune{firstDigit, nextChar}) // Integer, Float Or Scientific Float
		default:
			if ('0' <= nextChar && nextChar <= '9') || nextChar == '_' {
				return lexer.tokenizeInteger([]rune{firstDigit, nextChar}) // Integer, Float or Scientific Float
			}
		}
	}
	lexer.reader.Redo()
	return []rune{firstDigit}, Literal, Integer, nil
}

func guessKind(buffer []rune) (uint8, uint8) {

	switch string(buffer) {
	case PassString:
		return Keyboard, Pass
	case EndString:
		return Keyboard, End
	case IfString:
		return Keyboard, If
	case UnlessString:
		return Keyboard, Unless
	case ElseString:
		return Keyboard, Else
	case ElifString:
		return Keyboard, Elif
	case WhileString:
		return Keyboard, While
	case DoString:
		return Keyboard, Do
	case ForString:
		return Keyboard, For
	case UntilString:
		return Keyboard, Until
	case SwitchString:
		return Keyboard, Switch
	case CaseString:
		return Keyboard, Case
	case DefaultString:
		return Keyboard, Default
	case YieldString:
		return Keyboard, Yield
	case ReturnString:
		return Keyboard, Return
	case ContinueString:
		return Keyboard, Continue
	case BreakString:
		return Keyboard, Break
	case RedoString:
		return Keyboard, Redo
	case ModuleString:
		return Keyboard, Module
	case DefString:
		return Keyboard, Def
	case LambdaString:
		return Keyboard, Lambda
	case InterfaceString:
		return Keyboard, Interface
	case ClassString:
		return Keyboard, Class
	case TryString:
		return Keyboard, Try
	case ExceptString:
		return Keyboard, Except
	case FinallyString:
		return Keyboard, Finally
	case AndString:
		return Comparator, And
	case OrString:
		return Comparator, Or
	case XorString:
		return Comparator, Xor
	case InString:
		return Comparator, In
	case AsString:
		return Keyboard, As
	case RaiseString:
		return Keyboard, Raise
	case BEGINString:
		return Keyboard, BEGIN
	case ENDString:
		return Keyboard, END
	case NotString: // Unary operator
		return Operator, Not
	case TrueString:
		return Boolean, True
	case FalseString:
		return Boolean, False
	case NoneString:
		return NoneType, None
	case ContextString:
		return Keyboard, Context
	default:
		if identifierCheck.MatchString(string(buffer)) {
			return IdentifierKind, Unknown
		} else if junkKindCheck.MatchString(string(buffer)) {
			return JunkKind, Unknown
		}
	}
	return Unknown, Unknown
}

func (lexer *Lexer) tokenizeChars(startingChar rune) ([]rune, uint8, uint8, *errors.Error) {
	content := []rune{startingChar}
	for ; lexer.reader.HasNext(); lexer.reader.Next() {
		char := lexer.reader.Char()
		if ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || ('0' <= char && char <= '9') || (char == '_') {
			content = append(content, char)
		} else {
			break
		}
	}
	kind, directValue := guessKind(content)
	return content, kind, directValue, nil
}

func (lexer *Lexer) tokenizeComment() ([]rune, uint8, *errors.Error) {
	var content []rune
	for ; lexer.reader.HasNext(); lexer.reader.Next() {
		char := lexer.reader.Char()
		if char == '\n' {
			break
		}
		content = append(content, char)
	}
	return content, Comment, nil
}

func (lexer *Lexer) tokenizeRepeatableOperator(char rune, singleDirectValue uint8, singleKind uint8, doubleDirectValue uint8, doubleKind uint8, assignSingleDirectValue uint8, assignDoubleDirectValue uint8, assignSingleKind uint8, assignDoubleKind uint8) ([]rune, uint8, uint8) {
	content := []rune{char}
	kind := singleKind
	directValue := singleDirectValue
	if lexer.reader.HasNext() {
		nextChar := lexer.reader.Char()
		if nextChar == char {
			content = append(content, nextChar)
			lexer.reader.Next()
			kind = doubleKind
			directValue = doubleDirectValue
			if lexer.reader.HasNext() {
				nextNextChar := lexer.reader.Char()
				if nextNextChar == '=' {
					content = append(content, nextNextChar)
					kind = assignDoubleKind
					lexer.reader.Next()
					directValue = assignDoubleDirectValue
				}
			}
		} else if nextChar == '=' {
			kind = assignSingleKind
			content = append(content, nextChar)
			lexer.reader.Next()
			directValue = assignSingleDirectValue
		}
	}
	return content, kind, directValue
}

func (lexer *Lexer) tokenizeNotRepeatableOperator(char rune, single uint8, singleKind uint8, assign uint8, assignKind uint8) ([]rune, uint8, uint8) {
	content := []rune{char}
	kind := singleKind
	directValue := single
	if lexer.reader.HasNext() {
		nextChar := lexer.reader.Char()
		if nextChar == '=' {
			kind = assignKind
			directValue = assign
			content = append(content, nextChar)
			lexer.reader.Next()
		}
	}
	return content, kind, directValue
}
func (lexer *Lexer) next() (*Token, *errors.Error) {
	if !lexer.reader.HasNext() {
		lexer.complete = true
		return &Token{
			String:      "EOF",
			DirectValue: EOF,
			Kind:        EOF,
			Line:        lexer.line,
			Index:       lexer.reader.Index(),
		}, nil
	}
	var tokenizingError *errors.Error
	var kind uint8
	var content []rune
	directValue := Unknown
	line := lexer.line
	index := lexer.reader.Index()
	char := lexer.reader.Char()
	lexer.reader.Next()
	switch char {
	case '\r':
		if lexer.reader.Char() == NewLineChar {
			lexer.reader.Next()
			lexer.line++
			content = []rune{char}
			kind = Separator
			directValue = NewLine
		}
	case NewLineChar:
		lexer.line++
		content = []rune{char}
		kind = Separator
		directValue = NewLine
	case SemiColonChar:
		content = []rune{char}
		kind = Separator
		directValue = SemiColon
	case ColonChar:
		content = []rune{char}
		directValue = Colon
		kind = Punctuation
	case CommaChar:
		content = []rune{char}
		directValue = Comma
		kind = Punctuation
	case OpenParenthesesChar:
		content = []rune{char}
		directValue = OpenParentheses
		kind = Punctuation
	case CloseParenthesesChar:
		content = []rune{char}
		directValue = CloseParentheses
		kind = Punctuation
	case OpenSquareBracketChar:
		content = []rune{char}
		directValue = OpenSquareBracket
		kind = Punctuation
	case CloseSquareBracketChar:
		content = []rune{char}
		directValue = CloseSquareBracket
		kind = Punctuation
	case OpenBraceChar:
		content = []rune{char}
		directValue = OpenBrace
		kind = Punctuation
	case CloseBraceChar:
		content = []rune{char}
		directValue = CloseBrace
		kind = Punctuation
	case DollarSignChar:
		content = []rune{char}
		directValue = DollarSign
		kind = Punctuation
	case DotChar:
		content = []rune{char}
		directValue = Dot
		kind = Punctuation
	case WhiteSpaceChar:
		directValue = Whitespace
		kind = Whitespace
		content = []rune{char}
	case TabChar:
		directValue = Tab
		kind = Whitespace
		content = []rune{char}
	case CommentChar:
		content, kind, tokenizingError = lexer.tokenizeComment()
		content = append([]rune{'#'}, content...)
	case '\'', '"': // String1
		content, kind, directValue, tokenizingError = lexer.tokenizeStringLikeExpressions(char)
	case '`':
		content, kind, directValue, tokenizingError = lexer.tokenizeStringLikeExpressions(char)
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
		content, kind, directValue, tokenizingError = lexer.tokenizeNumeric(char)
	case StarChar:
		content, kind, directValue = lexer.tokenizeRepeatableOperator(char, Star, Operator, PowerOf, Operator, StarAssign, PowerOfAssign, Assignment, Assignment)
	case DivChar:
		content, kind, directValue = lexer.tokenizeRepeatableOperator(char, Div, Operator, FloorDiv, Operator, DivAssign, FloorDivAssign, Assignment, Assignment)
	case LessThanChar:
		content, kind, directValue = lexer.tokenizeRepeatableOperator(char, LessThan, Comparator, BitwiseLeft, Operator, LessOrEqualThan, BitwiseLeftAssign, Comparator, Assignment)
	case GreatThanChar:
		content, kind, directValue = lexer.tokenizeRepeatableOperator(char, GreaterThan, Comparator, BitwiseRight, Operator, GreaterOrEqualThan, BitwiseRightAssign, Comparator, Assignment)
	case AddChar:
		content, kind, directValue = lexer.tokenizeNotRepeatableOperator(char, Add, Operator, AddAssign, Assignment)
	case SubChar:
		content, kind, directValue = lexer.tokenizeNotRepeatableOperator(char, Sub, Operator, SubAssign, Assignment)
	case ModulusChar:
		content, kind, directValue = lexer.tokenizeNotRepeatableOperator(char, Modulus, Operator, ModulusAssign, Assignment)
	case BitwiseXorChar:
		content, kind, directValue = lexer.tokenizeNotRepeatableOperator(char, BitwiseXor, Operator, BitwiseXorAssign, Assignment)
	case BitWiseAndChar:
		content, kind, directValue = lexer.tokenizeNotRepeatableOperator(char, BitWiseAnd, Operator, BitWiseAndAssign, Assignment)
	case BitwiseOrChar:
		content, kind, directValue = lexer.tokenizeNotRepeatableOperator(char, BitwiseOr, Operator, BitwiseOrAssign, Assignment)
	case SignNotChar:
		content, kind, directValue = lexer.tokenizeNotRepeatableOperator(char, SignNot, Operator, NotEqual, Comparator)
	case NegateBitsChar:
		content, kind, directValue = lexer.tokenizeNotRepeatableOperator(char, NegateBits, Operator, NegateBitsAssign, Assignment)
	case EqualsChar:
		content, kind, directValue = lexer.tokenizeNotRepeatableOperator(char, Assign, Assignment, Equals, Comparator)
	case BackSlashChar:
		content = []rune{char}
		if lexer.reader.HasNext() {
			nextChar := lexer.reader.Char()
			if nextChar != '\n' {
				return nil, errors.New(lexer.line, "line escape not followed by a new line", errors.LexingError)
			}
			content = append(content, '\n')
			lexer.reader.Next()
		}
		kind = PendingEscape
	default:
		if char == 'b' {
			if lexer.reader.HasNext() {
				nextChar := lexer.reader.Char()
				if nextChar == '\'' || nextChar == '"' {
					var byteStringPart []rune
					lexer.reader.Next()
					byteStringPart, kind, directValue, tokenizingError = lexer.tokenizeStringLikeExpressions(nextChar)
					content = append([]rune{char}, byteStringPart...)
					if directValue != Unknown {
						directValue = ByteString
					}
					break
				}
			}
		}
		content, kind, directValue, tokenizingError = lexer.tokenizeChars(char)
	}
	return &Token{
		DirectValue: directValue,
		String:      string(content),
		Kind:        kind,
		Line:        line,
		Index:       index,
	}, tokenizingError
}

/*
	This function will yield just the necessary token, this means not repeated separators
*/
func (lexer *Lexer) Next() (*Token, *errors.Error) {
	token, tokenizingError := lexer.next()
	if tokenizingError != nil {
		return nil, tokenizingError
	}
	if token.Kind == JunkKind {
		return lexer.Next()
	}
	if token.Kind == Comment {
		return lexer.Next()
	}
	if token.Kind == Whitespace {
		return lexer.Next()
	}
	if token.Kind == Separator {
		if lexer.lastToken == nil {
			return lexer.Next()
		}
		switch lexer.lastToken.Kind {
		case Separator:
			return lexer.Next()
		case Operator, Comparator:
			return lexer.Next()
		default:
			break
		}
		switch lexer.lastToken.DirectValue {
		case Comma, OpenParentheses, OpenSquareBracket, OpenBrace:
			return lexer.Next()
		default:
			break
		}
	}
	lexer.lastToken = token
	return token, nil
}

func NewLexer(codeReader reader.Reader) *Lexer {
	return &Lexer{
		lastToken: nil,
		line:      1,
		reader:    codeReader,
		complete:  false,
	}
}
