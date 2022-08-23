package lexer

import "errors"

var (
	OctalInvalidSyntax = errors.New("invalid octal integer syntax")
)

func (lexer *Lexer) tokenizeOctal() error {
	if !lexer.reader.HasNext() {
		lexer.currentToken.Kind = Literal
		lexer.currentToken.DirectValue = InvalidDirectValue
		return OctalInvalidSyntax
	}
	nextDigit := lexer.reader.Char()
	if !('0' <= nextDigit && nextDigit <= '7') {
		lexer.currentToken.Kind = Literal
		lexer.currentToken.DirectValue = InvalidDirectValue
		return OctalInvalidSyntax
	}
	lexer.reader.Next()
	lexer.currentToken.append(nextDigit)
	for ; lexer.reader.HasNext(); lexer.reader.Next() {
		nextDigit = lexer.reader.Char()
		if !('0' <= nextDigit && nextDigit <= '7') && nextDigit != '_' {
			break
		}
		lexer.currentToken.append(nextDigit)
	}
	lexer.currentToken.Kind = Literal
	lexer.currentToken.DirectValue = OctalInteger
	return nil
}
