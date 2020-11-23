package SOCKS5

import (
	"bufio"
	"github.com/shoriwe/FullProxy/pkg/ConnectionControllers"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"net"
)

func (socks5 *Socks5) UsernamePasswordAuthentication(clientConnectionReader *bufio.Reader) (bool, byte) {
	numberOfReceivedBytes, credentials, connectionError := Sockets.Receive(clientConnectionReader, 1024)
	if connectionError != nil {
		return false, 0
	}
	if numberOfReceivedBytes < 4 {
		return false, 0
	}
	if credentials[0] != BasicNegotiation {
		return false, 0
	}
	receivedUsernameLength := int(credentials[1])
	if receivedUsernameLength+3 >= numberOfReceivedBytes {
		return false, 0
	}
	receivedUsername := credentials[2 : 2+receivedUsernameLength]
	rawReceivedUsernamePassword := credentials[2+receivedUsernameLength+1 : numberOfReceivedBytes]
	if socks5.AuthenticationMethod(receivedUsername, rawReceivedUsernamePassword) {
		return true, UsernamePassword
	}
	return false, 0
}

func (socks5 *Socks5) AuthenticateClient(clientConnection net.Conn, clientConnectionReader *bufio.Reader, clientConnectionWriter *bufio.Writer) bool {

	var foundMethod = InvalidMethod
	numberOfReceivedBytes, clientImplementedMethods, _ := Sockets.Receive(clientConnectionReader, 1024)
	if clientImplementedMethods == nil {
		_, _ = Sockets.Send(clientConnectionWriter, &NoSupportedMethods)
		return false
	} else if numberOfReceivedBytes >= 3 {
		if clientImplementedMethods[0] == Version && int(clientImplementedMethods[1]) == numberOfReceivedBytes-2 {
			for index := 2; index < numberOfReceivedBytes; index++ {
				if clientImplementedMethods[index] == socks5.WantedAuthMethod {
					foundMethod = socks5.WantedAuthMethod
					break
				}
			}
		}
	}

	switch foundMethod {
	case UsernamePassword:
		_, connectionError := Sockets.Send(clientConnectionWriter, &UsernamePasswordSupported)
		if connectionError == nil {
			if success, authenticationProtocol := socks5.UsernamePasswordAuthentication(clientConnectionReader); success && authenticationProtocol == UsernamePassword {
				_, connectionError = Sockets.Send(clientConnectionWriter, &UsernamePasswordSucceededResponse)
				return connectionError == nil
			}
			_, _ = Sockets.Send(clientConnectionWriter, &AuthenticationFailed)
		}
	case NoAuthRequired:
		_, connectionError := Sockets.Send(clientConnectionWriter, &NoAuthRequiredSupported)
		return connectionError == nil
	default:
		ConnectionControllers.LogData(socks5.LoggingMethod, "Client doesn't support authentication methods: ", clientConnection.RemoteAddr().String())
		_, _ = Sockets.Send(clientConnectionWriter, &NoSupportedMethods)
	}
	return false
}
