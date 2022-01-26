package SOCKS5

import "errors"

const (
	BasicNegotiation byte = 0x01
	SocksV4          byte = 0x04
	SocksV5          byte = 0x05
	NoAuthRequired   byte = 0x00
	UsernamePassword byte = 0x02
	Connect          byte = 0x01
	Bind             byte = 0x02
	UDPAssociate     byte = 0x03
	IPv4             byte = 0x01
	DomainName       byte = 0x03
	IPv6             byte = 0x04
)

var protocolError = errors.New("Protocol error?!")

var (
	// SOCKS requests connection responses
	UsernamePasswordSupported = []byte{SocksV5, UsernamePassword}
	NoAuthRequiredSupported   = []byte{SocksV5, NoAuthRequired}
	NoSupportedMethods        = []byte{SocksV5, 0xFF}
	AuthenticationSucceded    = []byte{BasicNegotiation, 0x00}
	AuthenticationFailed      = []byte{SocksV5, 0xFF}
)
