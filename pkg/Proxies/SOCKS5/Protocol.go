package SOCKS5

import (
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Tools/Types"
	"net"
	"time"
)

type Socks5 struct {
	AuthenticationMethod Types.AuthenticationMethod
	WantedAuthMethod     byte
	LoggingMethod        Types.LoggingMethod
	Tries                int
	Timeout              time.Duration
	InboundFilter        Types.IOFilter
	OutboundFilter       Types.IOFilter
}

func (socks5 *Socks5) SetLoggingMethod(loggingMethod Types.LoggingMethod) error {
	socks5.LoggingMethod = loggingMethod
	return nil
}

func (socks5 *Socks5) SetAuthenticationMethod(authenticationMethod Types.AuthenticationMethod) error {
	socks5.AuthenticationMethod = authenticationMethod
	return nil
}

func (socks5 *Socks5) SetTries(tries int) error {
	socks5.Tries = tries
	return nil
}

func (socks5 *Socks5) SetTimeout(timeout time.Duration) error {
	socks5.Timeout = timeout
	return nil
}

func (socks5 *Socks5) SetInboundFilter(filter Types.IOFilter) error {
	socks5.InboundFilter = filter
	return nil
}

func (socks5 *Socks5) SetOutboundFilter(filter Types.IOFilter) error {
	socks5.OutboundFilter = filter
	return nil
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
