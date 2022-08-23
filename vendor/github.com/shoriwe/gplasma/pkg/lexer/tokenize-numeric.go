package lexer

func (lexer *Lexer) tokenizeNumeric() error {
	if !lexer.reader.HasNext() {
		lexer.currentToken.Kind = Literal
		lexer.currentToken.DirectValue = Integer
		return nil
	}
	nextChar := lexer.reader.Char()
	lexer.reader.Next()
	if lexer.currentToken.Contents[0] == '0' {
		switch nextChar {
		case 'x', 'X': // Hexadecimal
			lexer.currentToken.append(nextChar)
			return lexer.tokenizeHexadecimal()
		case 'b', 'B': // Binary
			lexer.currentToken.append(nextChar)
			return lexer.tokenizeBinary()
		case 'o', 'O': // Octal
			lexer.currentToken.append(nextChar)
			return lexer.tokenizeOctal()
		}
	}
	switch nextChar {
	case 'e', 'E': // Scientific float
		lexer.currentToken.append(nextChar)
		return lexer.tokenizeScientificFloat()
	case '.': // Maybe a float
		lexer.currentToken.append(nextChar)
		return lexer.tokenizeFloat() // Integer, Float Or Scientific Float
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '_':
		lexer.currentToken.append(nextChar)
		return lexer.tokenizeInteger() // Integer, Float or Scientific Float
	}
	lexer.reader.Redo()
	lexer.currentToken.Kind = Literal
	lexer.currentToken.DirectValue = Integer
	return nil
}
