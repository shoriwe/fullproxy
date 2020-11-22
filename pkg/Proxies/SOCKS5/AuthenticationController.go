package SOCKS5

import (
	"github.com/shoriwe/FullProxy/pkg/Sockets"
)

func (socks5 *Socks5)GetClientAuthenticationImplementedMethods() bool {

	var foundMethod = InvalidMethod
	numberOfReceivedBytes, clientImplementedMethods, _ := Sockets.Receive(socks5.ClientConnectionReader, 1024)
	if clientImplementedMethods == nil {
		_, _ = Sockets.Send(socks5.ClientConnectionWriter, &NoSupportedMethods)
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
		_, connectionError := Sockets.Send(socks5.ClientConnectionWriter, &UsernamePasswordSupported)
		if connectionError == nil {
			if success, authenticationProtocol := socks5.HandleUsernamePasswordAuthentication(); success && authenticationProtocol == UsernamePassword {
				_, connectionError = Sockets.Send(socks5.ClientConnectionWriter, &UsernamePasswordSucceededResponse)
				authResult := connectionError == nil
				/*
				if connectionError == nil {
					log.Print("Login failed with invalid credentials from: ", clientAddress)
				}
				log.Print("Login succeeded from: ", clientAddress)
				 */
				return authResult
			}
			_, _ = Sockets.Send(socks5.ClientConnectionWriter, &AuthenticationFailed)
		}
	case NoAuthRequired:
		_, connectionError := Sockets.Send(socks5.ClientConnectionWriter, &NoAuthRequiredSupported)
		return connectionError == nil
	default:
		_, _ = Sockets.Send(socks5.ClientConnectionWriter, &NoSupportedMethods)
	}
	return false
}
