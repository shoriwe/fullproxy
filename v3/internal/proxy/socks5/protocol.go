package socks5

import (
	"errors"
	"github.com/shoriwe/FullProxy/v3/internal/global"
	"net"
)

type Socks5 struct {
	AuthenticationMethod global.AuthenticationMethod
	LoggingMethod        global.LoggingMethod
	OutboundFilter       global.IOFilter
	Dial                 global.DialFunc
	UDPRelay             *net.UDPConn
	relaySessions        map[string]*net.UDPConn
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
	sessionChunk := make([]byte, 0xFFFF)
	defer clientConnection.Close()
	authenticationSuccessful, connectionError := socks5.AuthenticateClient(sessionChunk, clientConnection)
	if connectionError != nil {
		return connectionError
	}
	if !authenticationSuccessful {
		errorMessage := "Authentication Failed with: " + clientConnection.RemoteAddr().String()
		_ = clientConnection.Close()
		// Templates.LogData(socks5.LoggingMethod, errorMessage)
		return errors.New(errorMessage)
	}
	_, connectionError = clientConnection.Read(sessionChunk)
	if connectionError != nil {
		return connectionError
	}
	version := sessionChunk[0]
	if version != SocksV5 {
		return SocksVersionNotSupported
	}
	var (
		targetHostLength, hostStartIndex int
	)
	switch sessionChunk[3] {
	case IPv4:
		targetHostLength = 4
		hostStartIndex = 4
	case DomainName:
		targetHostLength = int(sessionChunk[4])
		hostStartIndex = 5
	case IPv6:
		targetHostLength = 16
		hostStartIndex = 4
	default:
		return protocolError
	}
	rawTargetHost := sessionChunk[hostStartIndex : hostStartIndex+targetHostLength]

	rawTargetPort := sessionChunk[hostStartIndex+targetHostLength : hostStartIndex+targetHostLength+2]

	// Cleanup the address
	port, host, hostPort := clean(sessionChunk[3], rawTargetHost, rawTargetPort)
	switch sessionChunk[1] {
	case Connect:
		return socks5.Connect(clientConnection, port, host, hostPort)
	case UDPAssociate:
		return socks5.UDPAssociate(sessionChunk, clientConnection, port, host, hostPort)
	case Bind:
		return socks5.Bind(sessionChunk, clientConnection, port, host, hostPort)
	default:
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
