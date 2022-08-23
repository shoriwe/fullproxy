package lexer

import "errors"

var (
	StringInvalidEscape = errors.New("invalid string escape sequence")
	StringNeverClosed   = errors.New("string never closed")
)

func (lexer *Lexer) tokenizeStringLikeExpressions(stringOpener rune) error {
	var target DirectValue
	switch stringOpener {
	case '\'':
		target = SingleQuoteString
	case '"':
		target = DoubleQuoteString
	case '`':
		target = CommandOutput
	}
	var directValue = InvalidDirectValue
	escaped := false
	finish := false
	for ; lexer.reader.HasNext() && !finish; lexer.reader.Next() {
		char := lexer.reader.Char()
		if escaped {
			switch char {
			case '\\', '\'', '"', '`', 'a', 'b', 'e', 'f', 'n', 'r', 't', '?', 'u', 'x':
				escaped = false
			default:
				return StringInvalidEscape
			}
		} else {
			switch char {
			case stringOpener:
				directValue = target
				finish = true
			case '\\':
				escaped = true
			}
		}
		lexer.currentToken.append(char)
	}
	if directValue != target {
		return StringNeverClosed
	}
	lexer.currentToken.Kind = Literal
	lexer.currentToken.DirectValue = directValue
	return nil
}
