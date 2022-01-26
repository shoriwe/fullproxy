package SOCKS5

import (
	"github.com/shoriwe/FullProxy/v2/pkg/Tools"
	"net"
)

func (socks5 *Socks5) UsernamePasswordAuthentication(clientConnection net.Conn) (bool, error) {
	_, connectionError := clientConnection.Write(UsernamePasswordSupported)
	if connectionError != nil {
		return false, connectionError
	}

	var numberOfBytesReceived int
	negotiationVersion := make([]byte, 1)
	numberOfBytesReceived, connectionError = clientConnection.Read(negotiationVersion)
	if connectionError != nil {
		return false, connectionError
	} else if numberOfBytesReceived != 1 {
		return false, nil
	}
	switch negotiationVersion[0] {
	case BasicNegotiation:
		break
	default:
		return false, nil
	}
	userLength := make([]byte, 1)
	numberOfBytesReceived, connectionError = clientConnection.Read(userLength)
	if connectionError != nil {
		return false, connectionError
	} else if numberOfBytesReceived != 1 {
		return false, nil
	}
	username := make([]byte, userLength[0])
	numberOfBytesReceived, connectionError = clientConnection.Read(username)
	if connectionError != nil {
		return false, connectionError
	} else if numberOfBytesReceived != int(userLength[0]) {
		return false, nil
	}
	passwordLength := make([]byte, 1)
	numberOfBytesReceived, connectionError = clientConnection.Read(passwordLength)
	if connectionError != nil {
		return false, connectionError
	} else if numberOfBytesReceived != 1 {
		return false, nil
	}
	password := make([]byte, passwordLength[0])
	numberOfBytesReceived, connectionError = clientConnection.Read(password)
	if connectionError != nil {
		return false, connectionError
	} else if numberOfBytesReceived != int(passwordLength[0]) {
		return false, nil
	}
	loginSuccess, loginError := socks5.AuthenticationMethod(username, password)
	if loginError != nil || !loginSuccess {
		_, connectionError = clientConnection.Write(AuthenticationFailed)
		if connectionError != nil {
			return false, connectionError
		}
		return false, loginError
	}
	_, connectionError = clientConnection.Write(AuthenticationSucceded)
	if connectionError != nil {
		return false, connectionError
	}
	Tools.LogData(socks5.LoggingMethod, "Logging succeeded for: "+clientConnection.RemoteAddr().String())
	return true, nil
}

func (socks5 *Socks5) AuthenticateClient(clientConnection net.Conn) (bool, error) {

	version := make([]byte, 1)
	bytesReceived, connectionError := clientConnection.Read(version)
	if connectionError != nil {
		return false, connectionError
	} else if bytesReceived != 1 {
		return false, nil
	}

	switch version[0] {
	case SocksV5, SocksV4:
		break
	default:
		return false, nil
	}

	numberOfMethods := make([]byte, 1)
	bytesReceived, connectionError = clientConnection.Read(numberOfMethods)
	if connectionError != nil {
		return false, connectionError
	} else if bytesReceived != 1 {
		return false, nil
	}
	clientSupportedMethods := make([]byte, numberOfMethods[0])
	bytesReceived, connectionError = clientConnection.Read(clientSupportedMethods)
	if connectionError != nil {
		return false, connectionError
	} else if bytesReceived != int(numberOfMethods[0]) {
		return false, nil
	}

	if socks5.AuthenticationMethod == nil {
		for _, supportedMethod := range clientSupportedMethods {
			if supportedMethod == NoAuthRequired {
				_, connectionError = clientConnection.Write(NoAuthRequiredSupported)
				return true, connectionError
			}
		}
	} else {
		for _, supportedMethod := range clientSupportedMethods {
			if supportedMethod == UsernamePassword {
				return socks5.UsernamePasswordAuthentication(clientConnection)
			}
		}
	}
	Tools.LogData(socks5.LoggingMethod, "Client doesn't support authentication methods: ", clientConnection.RemoteAddr().String())
	_, connectionError = clientConnection.Write(NoSupportedMethods)
	return false, connectionError
}
