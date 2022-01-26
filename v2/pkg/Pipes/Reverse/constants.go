package Reverse

const (
	NewConnectionSucceeded byte = iota
	NewConnectionFailed
	RequestNewMasterConnection
	Dial
	Bind
	UnknownCommand
)
