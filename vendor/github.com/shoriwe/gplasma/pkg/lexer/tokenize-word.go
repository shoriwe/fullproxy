package lexer

func (lexer *Lexer) tokenizeWord() {
	for ; lexer.reader.HasNext(); lexer.reader.Next() {
		char := lexer.reader.Char()
		if ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || ('0' <= char && char <= '9') || (char == '_') {
			lexer.currentToken.append(char)
		} else {
			break
		}
	}
	lexer.currentToken.Kind, lexer.currentToken.DirectValue = lexer.detectKindAndDirectValue()
}
