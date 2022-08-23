package lexer

import (
	"errors"
	"github.com/shoriwe/gplasma/pkg/reader"
	"regexp"
)

var (
	identifierCheck = regexp.MustCompile("(?m)^[a-zA-Z_]+[a-zA-Z0-9_]*$")
	junkKindCheck   = regexp.MustCompile("(?m)^\\00+$")
)

type Lexer struct {
	currentToken *Token
	lastToken    *Token
	reader       reader.Reader
	complete     bool
}

var (
	CRLFInvalid          = errors.New("invalid CRLF")
	LineEscapeIncomplete = errors.New("incomplete line escape")
)

func (lexer *Lexer) HasNext() bool {
	return !lexer.complete
}

func (lexer *Lexer) next() (*Token, error) {
	lexer.currentToken = &Token{
		Contents:    nil,
		DirectValue: InvalidDirectValue,
		Kind:        EOF,
		Line:        lexer.reader.Line(),
		Index:       lexer.reader.Index(),
	}
	if !lexer.reader.HasNext() {
		lexer.complete = true
		return lexer.currentToken, nil
	}
	var tokenizingError error
	char := lexer.reader.Char()
	lexer.currentToken.append(char)
	lexer.reader.Next()
	switch char {
	case '\r':
		if lexer.reader.Char() != NewLineChar {
			return nil, CRLFInvalid
		}
		lexer.reader.Next()
		lexer.currentToken.append(char, '\n')
		lexer.currentToken.Kind = Separator
		lexer.currentToken.DirectValue = NewLine
	case NewLineChar:
		lexer.currentToken.Kind = Separator
		lexer.currentToken.DirectValue = NewLine
	case SemiColonChar:

		lexer.currentToken.Kind = Separator
		lexer.currentToken.DirectValue = SemiColon
	case ColonChar:

		lexer.currentToken.DirectValue = Colon
		lexer.currentToken.Kind = Punctuation
	case CommaChar:

		lexer.currentToken.DirectValue = Comma
		lexer.currentToken.Kind = Punctuation
	case OpenParenthesesChar:

		lexer.currentToken.DirectValue = OpenParentheses
		lexer.currentToken.Kind = Punctuation
	case CloseParenthesesChar:

		lexer.currentToken.DirectValue = CloseParentheses
		lexer.currentToken.Kind = Punctuation
	case OpenSquareBracketChar:

		lexer.currentToken.DirectValue = OpenSquareBracket
		lexer.currentToken.Kind = Punctuation
	case CloseSquareBracketChar:

		lexer.currentToken.DirectValue = CloseSquareBracket
		lexer.currentToken.Kind = Punctuation
	case OpenBraceChar:

		lexer.currentToken.DirectValue = OpenBrace
		lexer.currentToken.Kind = Punctuation
	case CloseBraceChar:

		lexer.currentToken.DirectValue = CloseBrace
		lexer.currentToken.Kind = Punctuation
	case DollarSignChar:

		lexer.currentToken.DirectValue = DollarSign
		lexer.currentToken.Kind = Punctuation
	case DotChar:

		lexer.currentToken.DirectValue = Dot
		lexer.currentToken.Kind = Punctuation
	case WhiteSpaceChar:

		lexer.currentToken.DirectValue = Blank
		lexer.currentToken.Kind = Whitespace
	case TabChar:

		lexer.currentToken.DirectValue = Blank
		lexer.currentToken.Kind = Whitespace
	case CommentChar:

		lexer.tokenizeComment()
	case '\'', '"', '`':

		tokenizingError = lexer.tokenizeStringLikeExpressions(char)
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':

		tokenizingError = lexer.tokenizeNumeric()
	case StarChar:

		lexer.tokenizeRepeatableOperator(Star, Operator, PowerOf, Operator, StarAssign, Assignment, PowerOfAssign, Assignment)
	case DivChar:

		lexer.tokenizeRepeatableOperator(Div, Operator, FloorDiv, Operator, DivAssign, Assignment, FloorDivAssign, Assignment)
	case LessThanChar:

		lexer.tokenizeRepeatableOperator(LessThan, Comparator, BitwiseLeft, Operator, LessOrEqualThan, Comparator, BitwiseLeftAssign, Assignment)
	case GreatThanChar:

		lexer.tokenizeRepeatableOperator(GreaterThan, Comparator, BitwiseRight, Operator, GreaterOrEqualThan, Comparator, BitwiseRightAssign, Assignment)
	case AddChar:

		lexer.tokenizeSingleOperator(Add, Operator, AddAssign, Assignment)
	case SubChar:

		lexer.tokenizeSingleOperator(Sub, Operator, SubAssign, Assignment)
	case ModulusChar:

		lexer.tokenizeSingleOperator(Modulus, Operator, ModulusAssign, Assignment)
	case BitwiseXorChar:

		lexer.tokenizeSingleOperator(BitwiseXor, Operator, BitwiseXorAssign, Assignment)
	case BitWiseAndChar:

		lexer.tokenizeSingleOperator(BitwiseAnd, Operator, BitwiseAndAssign, Assignment)
	case BitwiseOrChar:
		lexer.tokenizeSingleOperator(BitwiseOr, Operator, BitwiseOrAssign, Assignment)
	case SignNotChar:
		lexer.tokenizeSingleOperator(SignNot, Operator, NotEqual, Comparator)
	case NegateBitsChar:
		lexer.currentToken.Kind = Operator
		lexer.currentToken.DirectValue = NegateBits
	case EqualsChar:
		lexer.tokenizeSingleOperator(Assign, Assignment, Equals, Comparator)
	case BackSlashChar:

		if !lexer.reader.HasNext() {
			return nil, LineEscapeIncomplete
		}
		nextChar := lexer.reader.Char()
		if nextChar != '\n' {
			return nil, LineEscapeIncomplete
		}
		lexer.currentToken.append('\n')
		lexer.reader.Next()

	default:
		if char != 'b' || !lexer.reader.HasNext() {
			lexer.tokenizeWord()
			break
		}
		nextChar := lexer.reader.Char()
		if nextChar != '\'' && nextChar != '"' {
			lexer.tokenizeWord()
			break
		}
		lexer.reader.Next()
		lexer.currentToken.append(nextChar)
		tokenizingError = lexer.tokenizeStringLikeExpressions(nextChar)
		if lexer.currentToken.DirectValue != InvalidDirectValue {
			lexer.currentToken.DirectValue = ByteString
		}
		break
	}
	return lexer.currentToken, tokenizingError
}

/*
	This function will yield just the necessary token, this means not repeated separators
*/
func (lexer *Lexer) Next() (*Token, error) {
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
		reader:    codeReader,
		complete:  false,
	}
}
