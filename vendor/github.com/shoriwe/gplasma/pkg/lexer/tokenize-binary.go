package lexer

import "errors"

var (
	BinaryInvalidSyntax = errors.New("invalid binary integer syntax")
)

func (lexer *Lexer) tokenizeBinary() error {
	if !lexer.reader.HasNext() {
		lexer.currentToken.Kind = Literal
		lexer.currentToken.DirectValue = InvalidDirectValue
		return BinaryInvalidSyntax
	}
	nextDigit := lexer.reader.Char()
	if !(nextDigit == '0' || nextDigit == '1') {
		lexer.currentToken.Kind = Literal
		lexer.currentToken.DirectValue = InvalidDirectValue
		return BinaryInvalidSyntax
	}
	lexer.reader.Next()
	lexer.currentToken.append(nextDigit)
	for ; lexer.reader.HasNext(); lexer.reader.Next() {
		nextDigit = lexer.reader.Char()
		if !(nextDigit == '0' || nextDigit == '1') && nextDigit != '_' {
			break
		}
		lexer.currentToken.append(nextDigit)
	}
	lexer.currentToken.Kind = Literal
	lexer.currentToken.DirectValue = BinaryInteger
	return nil
}
