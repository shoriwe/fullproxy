package lexer

import "errors"

var (
	HexInvalidSyntax = errors.New("invalid hex integer syntax")
)

func (lexer *Lexer) tokenizeHexadecimal() error {
	if !lexer.reader.HasNext() {
		lexer.currentToken.Kind = Literal
		lexer.currentToken.DirectValue = InvalidDirectValue
		return HexInvalidSyntax
	}
	nextDigit := lexer.reader.Char()
	if !(('0' <= nextDigit && nextDigit <= '9') ||
		('a' <= nextDigit && nextDigit <= 'f') ||
		('A' <= nextDigit && nextDigit <= 'F')) {
		lexer.currentToken.Kind = Literal
		lexer.currentToken.DirectValue = InvalidDirectValue
		return HexInvalidSyntax
	}
	lexer.reader.Next()
	lexer.currentToken.append(nextDigit)
	for ; lexer.reader.HasNext(); lexer.reader.Next() {
		nextDigit = lexer.reader.Char()
		if !(('0' <= nextDigit && nextDigit <= '9') ||
			('a' <= nextDigit && nextDigit <= 'f') ||
			('A' <= nextDigit && nextDigit <= 'F')) &&
			nextDigit != '_' {
			break
		}
		lexer.currentToken.append(nextDigit)
	}
	lexer.currentToken.Kind = Literal
	lexer.currentToken.DirectValue = HexadecimalInteger
	return nil
}
