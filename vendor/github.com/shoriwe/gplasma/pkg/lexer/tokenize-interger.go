package lexer

func (lexer *Lexer) tokenizeInteger() error {
	if !lexer.reader.HasNext() {
		lexer.currentToken.Kind = Literal
		lexer.currentToken.DirectValue = Integer
		return nil
	}
	nextDigit := lexer.reader.Char()
	switch nextDigit {
	case '.':
		lexer.reader.Next()
		lexer.currentToken.append(nextDigit)
		return lexer.tokenizeFloat()
	case 'e', 'E':
		lexer.reader.Next()
		lexer.currentToken.append(nextDigit)
		return lexer.tokenizeScientificFloat()
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		// If no of this match return
	default:
		lexer.currentToken.Kind = Literal
		lexer.currentToken.DirectValue = Integer
		return nil
	}
	lexer.reader.Next()
	lexer.currentToken.append(nextDigit)
loop:
	for ; lexer.reader.HasNext(); lexer.reader.Next() {
		nextDigit = lexer.reader.Char()
		switch nextDigit {
		case 'e', 'E':
			lexer.currentToken.append(nextDigit)
			return lexer.tokenizeScientificFloat()
		case '.':
			lexer.reader.Next()
			lexer.currentToken.append(nextDigit)
			return lexer.tokenizeFloat()
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '_':
			lexer.currentToken.append(nextDigit)
		default:
			break loop
		}
	}
	lexer.currentToken.Kind = Literal
	lexer.currentToken.DirectValue = Integer
	return nil
}
