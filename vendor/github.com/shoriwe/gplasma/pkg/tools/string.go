package tools

import (
	"strconv"
)

var directCharEscapeValue = map[rune][]rune{
	'a': {7},
	'b': {8},
	'e': {'\\', 'e'},
	'f': {12},
	'n': {10},
	'r': {13},
	't': {9},
	'?': {'\\', '?'},
}

func Repeat(s string, times int64) string {
	result := ""
	for i := int64(0); i < times; i++ {
		result += s
	}
	return result
}

func ReplaceEscaped(s []rune) []rune {
	sLength := len(s)
	escaped := false
	var result []rune
	for index := 0; index < sLength; index++ {
		char := s[index]
		if escaped {
			switch char {
			case 'a', 'b', 'e', 'f', 'n', 'r', 't', '?':
				// Replace char based
				result = append(result, directCharEscapeValue[char]...)
			case '\\', '\'', '"', '`':
				// Replace escaped literals
				result = append(result, char)
			case 'x':
				// Replace hex with numbers
				index++
				a := s[index]
				index++
				b := s[index]
				number, parsingError := strconv.ParseUint(string([]rune{a, b}), 16, 32)
				if parsingError != nil {
					panic(parsingError)
				}
				result = append(result, rune(number))
			case 'u':
				// Replace unicode with numbers
				index++
				a := s[index]
				index++
				b := s[index]
				index++
				c := s[index]
				index++
				d := s[index]
				number, parsingError := strconv.ParseUint(string([]rune{a, b, c, d}), 16, 32)
				if parsingError != nil {
					panic(parsingError)
				}
				result = append(result, rune(number))
			}
			escaped = false
		} else if char == '\\' {
			escaped = true
		} else {
			result = append(result, char)
		}
	}
	return result
}
