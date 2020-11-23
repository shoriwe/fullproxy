package SOCKS5

import (
	"bufio"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
)

func (socks5 *Socks5) HandleUsernamePasswordAuthentication(clientConnectionReader *bufio.Reader) (bool, byte) {
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
