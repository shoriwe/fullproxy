package lexer

func (lexer *Lexer) tokenizeFloat() error {
	if !lexer.reader.HasNext() {
		lexer.reader.Redo()
		lexer.currentToken.Contents = lexer.currentToken.Contents[:len(lexer.currentToken.Contents)-1]
		lexer.currentToken.Kind = Literal
		lexer.currentToken.DirectValue = Integer
		return nil
	}
	nextDigit := lexer.reader.Char()
	if !('0' <= nextDigit && nextDigit <= '9') {
		lexer.reader.Redo()
		lexer.currentToken.Contents = lexer.currentToken.Contents[:len(lexer.currentToken.Contents)-1]
		lexer.currentToken.Kind = Literal
		lexer.currentToken.DirectValue = Integer
		return nil
	}
	lexer.reader.Next()
	lexer.currentToken.append(nextDigit)
	for ; lexer.reader.HasNext(); lexer.reader.Next() {
		nextDigit = lexer.reader.Char()
		if ('0' <= nextDigit && nextDigit <= '9') || nextDigit == '_' {
			lexer.currentToken.append(nextDigit)
		} else if (nextDigit == 'e') || (nextDigit == 'E') {
			lexer.reader.Next()
			lexer.currentToken.append(nextDigit)
			return lexer.tokenizeScientificFloat()
		} else {
			break
		}
	}
	lexer.currentToken.Kind = Literal
	lexer.currentToken.DirectValue = Float
	return nil
}
