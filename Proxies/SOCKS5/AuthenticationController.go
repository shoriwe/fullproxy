package SOCKS5

import (
	"FullProxy/Sockets"
	"bufio"
	"bytes"
)


func GetClientAuthenticationImplementedMethods(clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer,
	username *[]byte,
	passwordHash *[]byte) bool{
	var wantedMethod = NoAuthRequired
	if !bytes.Equal(*passwordHash, []byte{}){
		wantedMethod = UsernamePassword
	}


	var FoundMethod = InvalidMethod
	numberOfReceivedBytes, clientImplementedMethods, _ := Sockets.Receive(clientConnectionReader, 1024)
	if clientImplementedMethods == nil {
		_, _ = Sockets.Send(clientConnectionWriter, NoSupportedMethods)
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


	var connectionError error
	success := false
	switch FoundMethod {
	case UsernamePassword:
		var negotiationVersion byte
		_, connectionError = Sockets.Send(clientConnectionWriter, UsernamePasswordSupported)
		if connectionError != nil{
			break
		}
		success, negotiationVersion = HandleUsernamePasswordAuthentication(clientConnectionReader, username, passwordHash)
		if success{
			_, connectionError = Sockets.Send(clientConnectionWriter, []byte{negotiationVersion, Succeeded})
			break
		}
		_, connectionError = Sockets.Send(clientConnectionWriter, []byte{Version, InvalidMethod})
	case NoAuthRequired:
		_, connectionError = Sockets.Send(clientConnectionWriter, NoAuthRequiredSupported)
		if connectionError != nil{
			break
		}
		success = true
	}
	if connectionError != nil {
		return false
	}
	return success
}
