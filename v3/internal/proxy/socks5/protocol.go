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
	defer clientConnection.Close()
	// if !Templates.FilterInbound(socks5.InboundFilter, Templates.ParseIP(clientConnection.RemoteAddr().String())) {
	// 	errorMessage := "Connection denied to: " + clientConnection.RemoteAddr().String()
	// 	_ = clientConnection.Close()
	// 	Templates.LogData(socks5.LoggingMethod, errorMessage)
	// 	return errors.New(errorMessage)
	// }
	// Templates.LogData(socks5.LoggingMethod, "Connection Received from: ", clientConnection.RemoteAddr().String())
	// Receive connection
	authenticationSuccessful, connectionError := socks5.AuthenticateClient(clientConnection)
	if connectionError != nil {
		return connectionError
	}
	if !authenticationSuccessful {
		errorMessage := "Authentication Failed with: " + clientConnection.RemoteAddr().String()
		_ = clientConnection.Close()
		// Templates.LogData(socks5.LoggingMethod, errorMessage)
		return errors.New(errorMessage)
	}

	version := make([]byte, 1)
	var numberOfBytesReceived int
	numberOfBytesReceived, connectionError = clientConnection.Read(version)
	if connectionError != nil {
		return connectionError
	} else if numberOfBytesReceived != 1 {
		return protocolError
	}
	command := make([]byte, 1)
	numberOfBytesReceived, connectionError = clientConnection.Read(command)
	if connectionError != nil {
		return connectionError
	} else if numberOfBytesReceived != 1 {
		return protocolError
	}
	switch command[0] {
	case Connect:
		return socks5.Connect(clientConnection)
	case Bind:
		return socks5.Bind(clientConnection)
	case UDPAssociate:
		return socks5.UDPAssociate(clientConnection)
	default:
		return protocolError
	}
}

func NewSocks5(authenticationMethod global.AuthenticationMethod, loggingMethod global.LoggingMethod, outboundFilter global.IOFilter) *Socks5 {
	return &Socks5{AuthenticationMethod: authenticationMethod, LoggingMethod: loggingMethod, OutboundFilter: outboundFilter}
}
