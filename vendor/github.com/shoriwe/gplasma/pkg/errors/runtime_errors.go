package errors

import (
	"fmt"
)

// This should be changed
const (
	UnknownLine = 0
)

// Errors Types
const (
	UnknownVmOperationError = "UnknownVMOperationError"
	NameNotFoundError       = "NameNotFoundError"
	IndexError              = "IndexError"
	GoRuntimeError          = "GoRuntimeError"
)

// Errors Messages
const (
	UnknownTokenKindMessage = "Unknown token kind"
	NameNotFoundMessage     = "Name not found"
)

func NewUnknownTokenKindError(line int) *Error {
	return New(line, UnknownTokenKindMessage, LexingError)
}

func NewIndexOutOfRangeError(line int, length int, index int) *Error {
	return New(line, fmt.Sprintf("index %d out of bound for a %d container", index, length), IndexError)
}

func NewNameNotFoundError() *Error {
	return New(UnknownLine, NameNotFoundMessage, NameNotFoundError)
}

func NewUnknownVMOperationError(operation uint8) *Error {
	return New(UnknownLine, fmt.Sprintf("unknown operation with value %d", operation), UnknownVmOperationError)
}
