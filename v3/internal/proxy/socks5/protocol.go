package socks5

import (
	"encoding/binary"
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
	ClientConnection net.Conn
	Chunk            [DefaultContextChunkSize]byte
	Version          byte
	Command          byte
	DSTAddressType   byte
	DSTHost          string
	DSTAddress       string
	DSTPort          int
	DSTRawAddress    []byte
	DSTRawPort       []byte
}

func (c *Context) ReadCommand() error {
	_, connectionError := c.ClientConnection.Read(c.Chunk[:])
	if connectionError != nil {
		return connectionError
	}
	c.Version = c.Chunk[0]
	c.Command = c.Chunk[1]
	if c.Version != SocksV5 {
		return SocksVersionNotSupported
	}
	// Cleanup the address
	var (
		rawHost, rawPort []byte
	)
	switch c.Chunk[3] {
	case IPv4:
		rawHost = c.Chunk[4 : 4+4]
		rawPort = c.Chunk[4+4 : 4+4+2]
		c.DSTHost = fmt.Sprintf("%d.%d.%d.%d", rawHost[0], rawHost[1], rawHost[2], rawHost[3])
	case DomainName:
		rawHost = c.Chunk[5 : 5+c.Chunk[4]]
		rawPort = c.Chunk[5+c.Chunk[4] : 5+c.Chunk[4]+2]
		c.DSTHost = string(rawHost)
	case IPv6:
		rawHost = c.Chunk[4 : 4+16]
		rawPort = c.Chunk[4+16 : 4+16+2]
		c.DSTHost = fmt.Sprintf("[%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x]",
			rawHost[0], rawHost[1], rawHost[2], rawHost[3],
			rawHost[4], rawHost[5], rawHost[6], rawHost[7],
			rawHost[8], rawHost[9], rawHost[10], rawHost[11],
			rawHost[12], rawHost[13], rawHost[14], rawHost[15],
		)
	default:
		return UnknownAddressType
	}

	c.DSTPort = int(binary.BigEndian.Uint16(rawPort))
	c.DSTAddress = fmt.Sprintf("%s:%d", c.DSTHost, c.DSTPort)
	return nil
}

type Reply interface {
	Bytes() []byte
}

type BasicReply struct {
	Version    byte
	StatusCode byte
}

func (b BasicReply) Bytes() []byte {
	return []byte{b.Version, b.StatusCode}
}

type CommandReply struct {
	Version     byte
	StatusCode  byte
	AddressType byte
	Address     []byte
	Port        []byte
}

func (c CommandReply) Bytes() []byte {
	result := make([]byte, 0, 7+len(c.Address))
	result = append(result, c.Version, c.StatusCode, 0x00, c.AddressType)
	if c.AddressType == DomainName {
		result = append(result, byte(len(c.Address)))
	}
	result = append(result, c.Address...)
	result = append(result, c.Port...)
	return result
}

func (c *Context) Reply(reply Reply) error {
	_, connectionError := c.ClientConnection.Write(reply.Bytes())
	return connectionError
}

func NewContext(conn net.Conn) *Context {
	return &Context{
		ClientConnection: conn,
	}
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
	context := NewContext(clientConnection)
	defer clientConnection.Close()
	authenticationError := socks5.AuthenticateClient(context)
	if authenticationError != nil {
		return authenticationError
	}
	readCommandError := context.ReadCommand()
	if readCommandError != nil {
		return readCommandError
	}

	switch context.Command {
	case Connect:
		return socks5.Connect(context)
	case Bind:
		return socks5.Bind(context)
	default:
		_ = context.Reply(CommandReply{
			Version:     SocksV5,
			StatusCode:  MethodNotSupportedCode,
			AddressType: context.DSTAddressType,
			Address:     context.DSTRawAddress,
			Port:        context.DSTRawPort,
		})
		return MethodNotSupported
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
