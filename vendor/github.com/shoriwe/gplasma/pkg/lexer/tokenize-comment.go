package lexer

func (lexer *Lexer) tokenizeComment() {
	lexer.currentToken.Kind = Comment
	lexer.currentToken.append('#')
	for ; lexer.reader.HasNext(); lexer.reader.Next() {
		char := lexer.reader.Char()
		if char == '\n' {
			break
		}
		lexer.currentToken.append(char)
	}
}
