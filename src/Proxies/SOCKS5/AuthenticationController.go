package SOCKS5

import (
	"github.com/shoriwe/FullProxy/src/ConnectionStructures"
	"github.com/shoriwe/FullProxy/src/Sockets"
)

func GetClientAuthenticationImplementedMethods(clientConnectionReader ConnectionStructures.SocketReader,
	clientConnectionWriter ConnectionStructures.SocketWriter,
	username *[]byte,
	passwordHash *[]byte) bool {

	wantedMethod := UsernamePassword
	if *passwordHash == nil {
		wantedMethod = NoAuthRequired
	}
	var FoundMethod = InvalidMethod
	numberOfReceivedBytes, clientImplementedMethods, _ := Sockets.Receive(clientConnectionReader, 1024)
	if clientImplementedMethods == nil {
		_, _ = Sockets.Send(clientConnectionWriter, &NoSupportedMethods)
		return false
	} else if numberOfReceivedBytes >= 3 {
		if clientImplementedMethods[0] == Version && int(clientImplementedMethods[1]) == numberOfReceivedBytes-2 {
			for index := 2; index < numberOfReceivedBytes; index++ {
				if clientImplementedMethods[index] == wantedMethod {
					FoundMethod = wantedMethod
					break
				}
			}
		}
	}

	switch FoundMethod {
	case UsernamePassword:
		// Say to the client that we want to use the password protocol
		_, connectionError := Sockets.Send(clientConnectionWriter, &UsernamePasswordSupported)
		if connectionError == nil {
			if success, authenticationProtocol := HandleUsernamePasswordAuthentication(clientConnectionReader, username, passwordHash); success && authenticationProtocol == UsernamePassword {
				_, connectionError = Sockets.Send(clientConnectionWriter, &UsernamePasswordSucceededResponse)
				return connectionError == nil
			}
			_, _ = Sockets.Send(clientConnectionWriter, &AuthenticationFailed)
		}
	case NoAuthRequired:
		// Say to the client that he doesn't need to authenticate with us
		_, connectionError := Sockets.Send(clientConnectionWriter, &NoAuthRequiredSupported)
		return connectionError == nil
	default:
		_, _ = Sockets.Send(clientConnectionWriter, &InvalidMethodResponse)
	}
	return false
}
