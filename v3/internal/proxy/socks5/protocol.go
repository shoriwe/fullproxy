package socks5

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/shoriwe/FullProxy/v3/internal/global"
	"net"
)

const (
	DefaultContextChunkSize = 0xFFFF
)

type Socks5 struct {
	AuthenticationMethod global.AuthenticationMethod
	LoggingMethod        global.LoggingMethod
	OutboundFilter       global.IOFilter
	Dial                 global.DialFunc
	UDPRelay             *net.UDPConn
	relaySessions        map[string]*net.UDPConn
}

type Context struct {
	Chunk          [DefaultContextChunkSize]byte
	BNDAddressType int
	BNDHost        string
	BNDAddress     string
	BNDPort        int
}

func (c *Context) ParseAddress() error {
	var (
		rawHost, rawPort []byte
	)
	switch c.Chunk[3] {
	case IPv4:
		rawHost = c.Chunk[4 : 4+4]
		rawPort = c.Chunk[4+4 : 4+4+2]
		c.BNDHost = fmt.Sprintf("%d.%d.%d.%d", rawHost[0], rawHost[1], rawHost[2], rawHost[3])
	case DomainName:
		rawHost = c.Chunk[5 : 5+c.Chunk[4]]
		rawPort = c.Chunk[5+c.Chunk[4] : 5+c.Chunk[4]+2]
		c.BNDHost = string(rawHost)
	case IPv6:
		rawHost = c.Chunk[4 : 4+16]
		rawPort = c.Chunk[4+16 : 4+16+2]
		c.BNDHost = fmt.Sprintf("[%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x]",
			rawHost[0], rawHost[1], rawHost[2], rawHost[3],
			rawHost[4], rawHost[5], rawHost[6], rawHost[7],
			rawHost[8], rawHost[9], rawHost[10], rawHost[11],
			rawHost[12], rawHost[13], rawHost[14], rawHost[15],
		)
	default:
		// TODO: Return an error related to unknown address type
		panic("Implement me")
	}

	c.BNDPort = int(binary.BigEndian.Uint16(rawPort))
	c.BNDAddress = fmt.Sprintf("%s:%d", c.BNDHost, c.BNDPort)
	return nil
}

func (socks5 *Socks5) SetLoggingMethod(loggingMethod global.LoggingMethod) error {
	socks5.LoggingMethod = loggingMethod
	return nil
}

func (socks5 *Socks5) SetAuthenticationMethod(authenticationMethod global.AuthenticationMethod) error {
	socks5.AuthenticationMethod = authenticationMethod
	return nil
}

func (socks5 *Socks5) SetOutboundFilter(filter global.IOFilter) error {
	socks5.OutboundFilter = filter
	return nil
}

func (socks5 *Socks5) SetDial(dialFunc global.DialFunc) {
	socks5.Dial = dialFunc
}

func (socks5 *Socks5) Handle(clientConnection net.Conn) error {
	var context Context
	defer clientConnection.Close()
	authenticationSuccessful, connectionError := socks5.AuthenticateClient(clientConnection, &context)
	if connectionError != nil {
		return connectionError
	}
	if !authenticationSuccessful {
		errorMessage := "Authentication Failed with: " + clientConnection.RemoteAddr().String()
		_ = clientConnection.Close()
		// Templates.LogData(socks5.LoggingMethod, errorMessage)
		return errors.New(errorMessage)
	}
	_, connectionError = clientConnection.Read(context.Chunk[:])
	if connectionError != nil {
		return connectionError
	}
	version := context.Chunk[0]
	if version != SocksV5 {
		return SocksVersionNotSupported
	}

	// Cleanup the address
	addressParseError := context.ParseAddress()
	if addressParseError != nil {
		// TODO: Do something with the parsing error
		panic("Implement me")
	}
	switch context.Chunk[1] {
	case Connect:
		return socks5.Connect(clientConnection, &context)
	case UDPAssociate:
		// TODO: Return method not supported
		panic("Implement me")
	case Bind:
		return socks5.Bind(clientConnection, &context)
	default:
		// TODO: Do not panic, sending an non supported command could be used as DDoS
		panic("Implement me")
	}
}

func NewSocks5(
	authenticationMethod global.AuthenticationMethod,
	loggingMethod global.LoggingMethod,
	outboundFilter global.IOFilter,
	udpRelay *net.UDPConn,
) *Socks5 {
	return &Socks5{
		AuthenticationMethod: authenticationMethod,
		LoggingMethod:        loggingMethod,
		OutboundFilter:       outboundFilter,
		UDPRelay:             udpRelay,
	}
}
