package lexer

func (lexer *Lexer) tokenizeSingleOperator(
	single DirectValue, singleKind Kind,
	assign DirectValue, assignKind Kind) {
	lexer.currentToken.Kind = singleKind
	lexer.currentToken.DirectValue = single
	if !lexer.reader.HasNext() {
		return
	}
	nextChar := lexer.reader.Char()
	if nextChar != '=' {
		return
	}
	lexer.currentToken.Kind = assignKind
	lexer.currentToken.DirectValue = assign
	lexer.currentToken.append(nextChar)
	lexer.reader.Next()
}
