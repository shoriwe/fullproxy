package Proxies

import (
	"FullProxy/BindServer"
	"FullProxy/CryptoTools"
	"FullProxy/Interface"
	"FullProxy/Sockets"
	"bufio"
	"log"
	"math/big"
	"net"
	"time"
)

var (
	BasicNegotiation		byte = 1
	Version                 byte = 5
	NoAuthRequired          byte = 0
	InvalidMethod			byte = 1
	UsernamePassword		byte = 2
	// SOCKS requests connection response
	UsernamePasswordSupported	 = []byte{Version, UsernamePassword}
	NoAuthRequiredSupported      = []byte{Version, NoAuthRequired}
	NoSupportedMethods           = []byte{5, 255}
	// ConnectionNotAllowedByRuleset byte = 2
	// SOCKS commands
	Connect           byte = 1
	Bind              byte = 2
	UDPAssociate      byte = 3
	// SOCKS valid address types
	IPv4              byte = 1
	DomainName        byte = 3
	IPv6              byte = 4
	// SOCKS5 responses
	ConnectionRefused byte = 5
	Succeeded         byte = 0
)


func HandleUsernamePasswordAuthentication(clientConnectionReader *bufio.Reader,
										  username string,
										  passwordHash []byte) (bool, byte){
	numberOfReceivedBytes, credentials, connectionError := Sockets.Receive(clientConnectionReader, 1024)
	if connectionError != nil{
		return false, 0
	}
	if numberOfReceivedBytes < 4{
		return false, 0
	}

	if credentials[0] != BasicNegotiation{
		return false, 0
	}
	receivedUsernameLength := int(credentials[1])
	if receivedUsernameLength + 3  >= numberOfReceivedBytes{
		return false, 0
	}
	receivedUsername := credentials[2:2+receivedUsernameLength]
	if string(receivedUsername) == username{
		rawReceivedUsernamePassword := credentials[2+receivedUsernameLength+1:numberOfReceivedBytes]
		if string(CryptoTools.GetPasswordHashPasswordByteArray(username, rawReceivedUsernamePassword)) == string(passwordHash){
			return true, credentials[0]
		}
	}
	return false, 0
}


func GetClientImplementedMethods(clientConnectionReader *bufio.Reader,
	  							 clientConnectionWriter *bufio.Writer,
	  							 username string,
	  							 passwordHash []byte) bool{
	var wantedMethod = NoAuthRequired
	if passwordHash != nil{
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


func HandleConnectCommand(clientConnection net.Conn,
						  targetConnection net.Conn,
						  clientConnectionReader *bufio.Reader,
						  clientConnectionWriter *bufio.Writer,
						  targetConnectionReader *bufio.Reader,
						  targetConnectionWriter *bufio.Writer) {
	for {
		for state := 0; state < 2; state++ {
			switch state {
			case 1:
				_ = clientConnection.SetReadDeadline(time.Now().Add(100))
				NumberOfBytesReceived, buffer, ConnectionError := Sockets.Receive(clientConnectionReader, 20480)
				if ConnectionError != nil {
					if ConnectionError, ok := ConnectionError.(net.Error); !(ok && ConnectionError.Timeout()) {
						return
					}
				}
				_, ConnectionError = Sockets.Send(targetConnectionWriter, buffer[:NumberOfBytesReceived])
				if ConnectionError != nil {
					return
				}
			case 0:
				_ = targetConnection.SetReadDeadline(time.Now().Add(100))
				NumberOfBytesReceived, buffer, ConnectionError := Sockets.Receive(targetConnectionReader, 20480)
				if ConnectionError != nil {
					if ConnectionError, ok := ConnectionError.(net.Error); !(ok && ConnectionError.Timeout()) {
						return
					}
				}
				_, ConnectionError = Sockets.Send(clientConnectionWriter, buffer[:NumberOfBytesReceived])
				if ConnectionError != nil {
					return
				}
			}
		}
	}
}


func ReceiveTargetRequest(clientConnectionReader *bufio.Reader) (byte, byte, []byte, []byte) {
	numberOfBytesReceived, targetRequest, ConnectionError := Sockets.Receive(clientConnectionReader, 1024)
	if ConnectionError != nil {
		return 0, 0, nil, nil
	}
	if numberOfBytesReceived < 10 {
		return 0, 0, nil, nil
	}
	if targetRequest[0] == Version &&
		(targetRequest[1] == Connect || targetRequest[1] == Bind || targetRequest[1] == UDPAssociate) &&
		(targetRequest[3] == IPv4 || targetRequest[3] == IPv6 || targetRequest[3] == DomainName) {
		return targetRequest[1], targetRequest[3], targetRequest[4 : numberOfBytesReceived-2], targetRequest[numberOfBytesReceived-2 : numberOfBytesReceived]
	}
	return 0, 0, nil, nil
}


func CreateProxySession(clientConnection net.Conn, username string, passwordHash []byte) {

	var targetConnection net.Conn
	var targetAddress string
	clientConnectionReader := bufio.NewReader(clientConnection)
	clientConnectionWriter := bufio.NewWriter(clientConnection)

	// Receive connection
	clientHasCompatibleMethods := GetClientImplementedMethods(clientConnectionReader,
															  clientConnectionWriter,
															  username,
															  passwordHash)
	if !clientHasCompatibleMethods {
		_ = clientConnection.Close()
		return
	}
	// Receive and process connection request
	targetRequestedCommand, targetAddressType, rawTargetAddress, rawTargetPort := ReceiveTargetRequest(clientConnectionReader)
	if targetRequestedCommand == 0 && targetAddressType == 0 {
		_ = clientConnection.Close()
		return
	}
	switch targetAddressType {
	case IPv4:
		targetAddress = net.IPv4(rawTargetAddress[0], rawTargetAddress[1], rawTargetAddress[2], rawTargetAddress[3]).String()
	case IPv6:
		break
		// targetAddress = net.IP(rawTargetAddress...).String()
	case DomainName:
		targetAddress = string(rawTargetAddress[1:])
	}


	switch targetRequestedCommand {
	case Connect:
		targetConnection = Sockets.Connect(targetAddress, new(big.Int).SetBytes(rawTargetPort).String())
		if targetConnection == nil {
			_, _ = Sockets.Send(clientConnectionWriter, []byte{Version, ConnectionRefused, 0, targetAddressType, 0, 0})
			_ = clientConnection.Close()
			return
		}
		response := []byte{Version, Succeeded, 0, targetAddressType}
		response = append(response[:], rawTargetAddress[:]...)
		response = append(response[:], rawTargetPort[:]...)
		_, ConnectionError := Sockets.Send(clientConnectionWriter, response)
		if ConnectionError != nil {
			_ = clientConnection.Close()
			return
		}
		targetConnectionReader := bufio.NewReader(targetConnection)
		targetConnectionWriter := bufio.NewWriter(targetConnection)
		HandleConnectCommand(clientConnection, targetConnection, clientConnectionReader, clientConnectionWriter, targetConnectionReader, targetConnectionWriter)
	case Bind:
		break
	case UDPAssociate:
		break
	}
	if targetConnection != nil {
		_ = targetConnection.Close()
	}
	_ = clientConnection.Close()
}


func StartSocks5(ip string, port string, interfaceMode bool, username string, password string) {
	passwordHash := CryptoTools.GetPasswordHashPasswordString(username, password)
	switch interfaceMode {
	case true:
		log.Println("Starting SOCKS5 server in Interface Mode")
		Interface.Client(ip, port, username, passwordHash, CreateProxySession)
	case false:
		log.Println("Starting SOCKS5 server in Bind Mode")
		BindServer.BindServer(ip, port, username, passwordHash, CreateProxySession)
	}
}
