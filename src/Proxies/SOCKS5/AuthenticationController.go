package SOCKS5

import (
	"bufio"
	"github.com/shoriwe/FullProxy/src/Sockets"
	"log"
)

func GetClientAuthenticationImplementedMethods(
	clientAddress string,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer,
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
		_, connectionError := Sockets.Send(clientConnectionWriter, &UsernamePasswordSupported)
		if connectionError == nil {
			if success, authenticationProtocol := HandleUsernamePasswordAuthentication(clientConnectionReader, username, passwordHash); success && authenticationProtocol == UsernamePassword {
				_, connectionError = Sockets.Send(clientConnectionWriter, &UsernamePasswordSucceededResponse)
				authResult := connectionError == nil
				if connectionError == nil {
					log.Print("Login failed with invalid credentials from: ", clientAddress)
				}
				log.Print("Login succeeded from: ", clientAddress)
				return authResult
			}
			_, _ = Sockets.Send(clientConnectionWriter, &AuthenticationFailed)
		}
	case NoAuthRequired:
		_, connectionError := Sockets.Send(clientConnectionWriter, &NoAuthRequiredSupported)
		return connectionError == nil
	default:
		_, _ = Sockets.Send(clientConnectionWriter, &NoSupportedMethods)
	}
	return false
}
