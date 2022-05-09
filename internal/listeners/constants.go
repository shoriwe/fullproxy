package listeners

import (
	"errors"
)

var (
	SlaveConnectionRequestError = errors.New("new connection request error")
)

const (
	NewConnectionSucceeded byte = iota
	NewConnectionFailed
	RequestNewMasterConnectionCommand
	DialCommand
	UnknownCommand
)
