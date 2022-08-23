package lexer

func (lexer *Lexer) tokenizeRepeatableOperator(
	singleDirectValue DirectValue, singleKind Kind,
	doubleDirectValue DirectValue, doubleKind Kind,
	assignSingleDirectValue DirectValue, assignSingleKind Kind,
	assignDoubleDirectValue DirectValue, assignDoubleKind Kind,
) {
	lexer.currentToken.Kind = singleKind
	lexer.currentToken.DirectValue = singleDirectValue
	if !lexer.reader.HasNext() {
		return
	}
	nextChar := lexer.reader.Char()
	if nextChar == '=' {
		lexer.currentToken.Kind = assignSingleKind
		lexer.currentToken.append(nextChar)
		lexer.reader.Next()
		lexer.currentToken.DirectValue = assignSingleDirectValue
		return
	}
	if nextChar != lexer.currentToken.Contents[0] {
		return
	}
	lexer.currentToken.append(nextChar)
	lexer.reader.Next()
	lexer.currentToken.Kind = doubleKind
	lexer.currentToken.DirectValue = doubleDirectValue
	if !lexer.reader.HasNext() {
		return
	}
	nextNextChar := lexer.reader.Char()
	if nextNextChar != '=' {
		return
	}
	lexer.currentToken.append(nextNextChar)
	lexer.currentToken.Kind = assignDoubleKind
	lexer.reader.Next()
	lexer.currentToken.DirectValue = assignDoubleDirectValue
}
