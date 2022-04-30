package socks5

import (
	"encoding/binary"
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers"
	"net"
	"strconv"
)

const (
	DefaultContextChunkSize = 0xFFFF
)

type Socks5 struct {
	AuthenticationMethod servers.AuthenticationMethod
	Dial                 servers.DialFunc
	Listen               servers.ListenFunc
	ListenAddress        *net.TCPAddr
}

func (socks5 *Socks5) SetListenAddress(address net.Addr) {
	socks5.ListenAddress = address.(*net.TCPAddr)
}

func (socks5 *Socks5) SetListen(listenFunc servers.ListenFunc) {
	socks5.Listen = listenFunc
}

type Context struct {
	ClientConnection net.Conn
	Chunk            [DefaultContextChunkSize]byte
	Version          byte
	Command          byte
	DST              string
	DSTAddressType   byte
	DSTAddress       string
	DSTRawAddress    []byte
	DSTPort          int
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
	switch c.Chunk[3] {
	case IPv4:
		c.DSTRawAddress = c.Chunk[4 : 4+4]
		c.DSTRawPort = c.Chunk[4+4 : 4+4+2]

		c.DSTAddress = net.IP(c.DSTRawAddress).To4().String()
	case DomainName:
		c.DSTRawAddress = c.Chunk[5 : 5+c.Chunk[4]]
		c.DSTRawPort = c.Chunk[5+c.Chunk[4] : 5+c.Chunk[4]+2]

		c.DSTAddress = string(c.Chunk[5 : 5+c.Chunk[4]])
	case IPv6:
		c.DSTRawAddress = c.Chunk[4 : 4+16]
		c.DSTRawPort = c.Chunk[4+16 : 4+16+2]

		c.DSTAddress = net.IP(c.DSTRawAddress).To16().String()
	default:
		return UnknownAddressType
	}

	c.DSTPort = int(binary.BigEndian.Uint16(c.DSTRawPort))
	c.DST = net.JoinHostPort(c.DSTAddress, strconv.Itoa(c.DSTPort))
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
	Version    byte
	StatusCode byte
	Address    net.IP
	Port       int
}

func (c CommandReply) Bytes() []byte {
	var (
		addressType byte
		address     []byte
	)
	if c.Address.To4() != nil {
		addressType = IPv4
		address = c.Address.To4()
	} else if c.Address.To16() != nil {
		addressType = IPv6
		address = c.Address.To16()
	} else {
		addressType = DomainName
		address = c.Address
	}
	result := make([]byte, 0, 7+len(address))
	result = append(result, c.Version, c.StatusCode, 0x00, addressType)
	if addressType == DomainName {
		result = append(result, byte(len(address)))
	}
	var portChunk [2]byte
	binary.BigEndian.PutUint16(portChunk[:], uint16(c.Port))
	result = append(result, address...)
	result = append(result, portChunk[0], portChunk[1])
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

func (socks5 *Socks5) SetAuthenticationMethod(authenticationMethod servers.AuthenticationMethod) error {
	socks5.AuthenticationMethod = authenticationMethod
	return nil
}

func (socks5 *Socks5) SetDial(dialFunc servers.DialFunc) {
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
			Version:    SocksV5,
			StatusCode: MethodNotSupportedCode,
			Address:    context.DSTRawAddress,
			Port:       context.DSTPort,
		})
		return MethodNotSupported
	}
}

func NewSocks5(
	authenticationMethod servers.AuthenticationMethod,
) servers.Protocol {
	return &Socks5{
		AuthenticationMethod: authenticationMethod,
	}
}
