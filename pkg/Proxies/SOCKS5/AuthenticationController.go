package SOCKS5

import (
	"bufio"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
)

func (socks5 *Socks5) GetClientAuthenticationImplementedMethods(clientConnectionReader *bufio.Reader, clientConnectionWriter *bufio.Writer) bool {

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
			if success, authenticationProtocol := socks5.HandleUsernamePasswordAuthentication(clientConnectionReader); success && authenticationProtocol == UsernamePassword {
				_, connectionError = Sockets.Send(clientConnectionWriter, &UsernamePasswordSucceededResponse)
				authResult := connectionError == nil
				/*
					if connectionError == nil {
						log.Print("Login failed with invalid credentials from: ", clientHost)
					}
					log.Print("Login succeeded from: ", clientHost)
				*/
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
