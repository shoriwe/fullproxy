package SOCKS5

import (
	"bufio"
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"github.com/shoriwe/FullProxy/pkg/Templates"
	"github.com/shoriwe/FullProxy/pkg/Templates/Types"
	"math/big"
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

func ReceiveTargetRequest(clientConnectionReader *bufio.Reader) (byte, byte, []byte, []byte) {
	numberOfBytesReceived, targetRequest, ConnectionError := Sockets.Receive(clientConnectionReader, 1024)
	if ConnectionError != nil {
		return 0, 0, nil, nil
	}
	if targetRequest[0] != Version {
		return 0, 0, nil, nil
	}
	if !(targetRequest[1] == Connect || targetRequest[1] == Bind || targetRequest[1] == UDPAssociate) {
		return 0, 0, nil, nil
	}
	if !(targetRequest[3] == IPv4 || targetRequest[3] == IPv6 || targetRequest[3] == DomainName) {
		return 0, 0, nil, nil
	}
	return targetRequest[1], targetRequest[3], targetRequest[4 : numberOfBytesReceived-2], targetRequest[numberOfBytesReceived-2 : numberOfBytesReceived]
}

func GetTargetHostPort(targetRequestedCommand *byte, targetHostType *byte, rawTargetHost []byte, rawTargetPort []byte) (byte, string, string) {
	if *targetRequestedCommand != 0 && *targetHostType != 0 {
		switch *targetHostType {
		case IPv4:
			return *targetRequestedCommand, net.IPv4(rawTargetHost[0], rawTargetHost[1], rawTargetHost[2], rawTargetHost[3]).String(), new(big.Int).SetBytes(rawTargetPort).String()
		case IPv6:
			return *targetRequestedCommand, Sockets.GetIPv6(rawTargetHost), new(big.Int).SetBytes(rawTargetPort).String()
		case DomainName:
			return *targetRequestedCommand, string(rawTargetHost[1:]), new(big.Int).SetBytes(rawTargetPort).String()
		}
	}
	return ConnectionRefused, "", ""
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

func (socks5 *Socks5) Handle(
	clientConnection net.Conn,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer) error {
	if !Templates.FilterInbound(socks5.InboundFilter, Templates.ParseIP(clientConnection.RemoteAddr().String())) {
		errorMessage := "Unwanted connection received from " + clientConnection.RemoteAddr().String()
		_ = clientConnection.Close()
		Templates.LogData(socks5.LoggingMethod, errorMessage)
		return errors.New(errorMessage)
	}
	var targetRequestedCommand byte
	// Receive connection
	authenticationSuccessful := socks5.AuthenticateClient(clientConnection, clientConnectionReader, clientConnectionWriter)
	if !authenticationSuccessful {
		errorMessage := "Authentication Failed with: " + clientConnection.RemoteAddr().String()
		_ = clientConnection.Close()
		Templates.LogData(socks5.LoggingMethod, errorMessage)
		return errors.New(errorMessage)
	}
	Templates.LogData(socks5.LoggingMethod, "Login succeeded from: ", clientConnection.RemoteAddr().String())
	var targetHost string
	var targetPort string
	rawTargetRequestedCommand, targetHostType, rawTargetHost, rawTargetPort := ReceiveTargetRequest(
		clientConnectionReader)
	targetRequestedCommand, targetHost, targetPort = GetTargetHostPort(
		&rawTargetRequestedCommand, &targetHostType,
		rawTargetHost, rawTargetPort)
	if targetRequestedCommand == ConnectionRefused {
		errorMessage := "Target connection refused: " + targetHost + ":" + targetPort
		_ = clientConnection.Close()
		Templates.LogData(socks5.LoggingMethod, errorMessage)
		return errors.New(errorMessage)
	}
	if !Templates.FilterOutbound(socks5.OutboundFilter, net.ParseIP(targetHost)) {
		errorMessage := "Unwanted outbound connection requested by: " + clientConnection.RemoteAddr().String() + " To: " + targetHost
		_ = clientConnection.Close()
		Templates.LogData(socks5.LoggingMethod, errorMessage)
		return errors.New(errorMessage)
	}
	return socks5.ExecuteCommand(clientConnection, clientConnectionReader, clientConnectionWriter,
		&targetRequestedCommand, &targetHostType, &targetHost, &targetPort)
}
