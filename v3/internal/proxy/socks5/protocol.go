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
	switch sessionChunk[1] {
	case Connect:
		return socks5.Connect(sessionChunk, clientConnection)
	case Bind:
		return socks5.Bind(sessionChunk, clientConnection)
	case UDPAssociate:
		return socks5.UDPAssociate(sessionChunk, clientConnection)
	default:
		return protocolError
	}
}

func NewSocks5(authenticationMethod global.AuthenticationMethod, loggingMethod global.LoggingMethod, outboundFilter global.IOFilter) *Socks5 {
	return &Socks5{
		AuthenticationMethod: authenticationMethod,
		LoggingMethod:        loggingMethod,
		OutboundFilter:       outboundFilter,
	}
}
