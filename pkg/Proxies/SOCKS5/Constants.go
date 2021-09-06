package SOCKS5

import "errors"

const (
	BasicNegotiation  byte = 1
	SocksV4           byte = 4
	SocksV5           byte = 5
	NoAuthRequired    byte = 0
	InvalidMethod     byte = 1
	UsernamePassword  byte = 2
	Connect           byte = 1
	Bind              byte = 2
	UDPAssociate      byte = 3
	IPv4              byte = 1
	DomainName        byte = 3
	IPv6              byte = 4
	ConnectionRefused byte = 5
	Succeeded         byte = 0
)

var protocolError = errors.New("Protocol error?!")

var (
	// SOCKS requests connection responses
	UsernamePasswordSupported = []byte{SocksV5, UsernamePassword}
	NoAuthRequiredSupported   = []byte{SocksV5, NoAuthRequired}
	NoSupportedMethods        = []byte{SocksV5, 255}
	AuthenticationSucceded    = []byte{BasicNegotiation, Succeeded}
	AuthenticationFailed      = []byte{BasicNegotiation, 1}
)
