package errors

import (
	"errors"
	"fmt"
)

// Compiling Errors
const (
	LexingError  = "LexingError"
	ParsingError = "ParsingError"
	SyntaxError  = "SyntaxError"
)

type Error struct {
	type_   string
	message string
	line    int
}

func (plasmaError *Error) Type() string {
	return plasmaError.type_
}

func (plasmaError *Error) Message() string {
	return plasmaError.message
}

func (plasmaError *Error) String() string {
	return fmt.Sprintf("%s: %s at line %d", plasmaError.type_, plasmaError.message, plasmaError.line)
}

func (plasmaError *Error) Line() int {
	return plasmaError.line
}

func (plasmaError *Error) Error() error {
	return errors.New(plasmaError.String())
}

func New(line int, message string, type_ string) *Error {
	return &Error{
		type_:   type_,
		message: message,
		line:    line,
	}
}
