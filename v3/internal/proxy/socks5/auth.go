package socks5

import (
	"github.com/shoriwe/FullProxy/v3/internal/global"
	"net"
)

func (socks5 *Socks5) UsernamePasswordAuthentication(clientConnection net.Conn, context *Context) (bool, error) {
	_, connectionError := clientConnection.Read(context.Chunk[:])
	if connectionError != nil {
		return false, connectionError
	}
	if context.Chunk[0] != BasicNegotiation {
		return false, nil
	}
	userLength := context.Chunk[1]
	username := context.Chunk[2 : 2+userLength]
	passwordLength := context.Chunk[2+userLength]
	password := context.Chunk[2+userLength+1 : 2+userLength+1+passwordLength]
	loginSuccess, loginError := socks5.AuthenticationMethod(username, password)
	if loginError != nil || !loginSuccess {
		_, connectionError = clientConnection.Write([]byte{BasicNegotiation, FailedAuthentication})
		if connectionError != nil {
			return false, connectionError
		}
		return false, loginError
	}
	_, connectionError = clientConnection.Write([]byte{BasicNegotiation, SucceedAuthentication})
	if connectionError != nil {
		return false, connectionError
	}
	global.LogData(socks5.LoggingMethod, "Logging succeeded for: "+clientConnection.RemoteAddr().String())
	return true, nil
}

func (socks5 *Socks5) AuthenticateClient(clientConnection net.Conn, context *Context) (bool, error) {
	_, connectionError := clientConnection.Read(context.Chunk[:])
	if connectionError != nil {
		return false, connectionError
	}
	if context.Chunk[0] != SocksV5 {
		return false, nil
	}

	clientSupportedMethods := context.Chunk[2 : 2+context.Chunk[1]]

	if socks5.AuthenticationMethod == nil {
		for _, supportedMethod := range clientSupportedMethods {
			if supportedMethod == NoAuthRequired {
				_, connectionError = clientConnection.Write([]byte{SocksV5, NoAuthRequired})
				return true, connectionError
			}
		}
	} else {
		for _, supportedMethod := range clientSupportedMethods {
			if supportedMethod == UsernamePassword {
				_, connectionError = clientConnection.Write([]byte{SocksV5, UsernamePassword})
				if connectionError != nil {
					return false, connectionError
				}
				return socks5.UsernamePasswordAuthentication(clientConnection, context)
			}
		}
	}
	global.LogData(socks5.LoggingMethod, "Client doesn't support authentication methods: ", clientConnection.RemoteAddr().String())
	_, connectionError = clientConnection.Write([]byte{SocksV5, NoAcceptableMethods})
	return false, connectionError
}
