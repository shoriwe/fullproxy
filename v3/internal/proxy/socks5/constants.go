package socks5

import "errors"

const (
	BasicNegotiation      byte = 0x01
	SocksV5               byte = 0x05
	NoAuthRequired        byte = 0x00
	UsernamePassword      byte = 0x02
	Connect               byte = 0x01
	Bind                  byte = 0x02
	UDPAssociate          byte = 0x03
	IPv4                  byte = 0x01
	DomainName            byte = 0x03
	IPv6                  byte = 0x04
	NoAcceptableMethods   byte = 0xFF
	FailedAuthentication  byte = 0xFF
	SucceedAuthentication byte = 0x00
	ConnectionSucceed     byte = 0x00
)

var (
	protocolError            = errors.New("protocol error?!")
	SocksVersionNotSupported = errors.New("client requesting a non supported socks version")
)
