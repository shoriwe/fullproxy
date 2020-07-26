package SOCKS5

import (
	"FullProxy/BindServer"
	"FullProxy/CryptoTools"
	"FullProxy/Interface"
	"FullProxy/Sockets"
	"bufio"
	"log"
	"math/big"
	"net"
)


const (
	BasicNegotiation		byte = 1
	Version                 byte = 5
	NoAuthRequired          byte = 0
	InvalidMethod			byte = 1
	UsernamePassword		byte = 2
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

var (
	// SOCKS requests connection responses
	UsernamePasswordSupported	 = []byte{Version, UsernamePassword}
	NoAuthRequiredSupported      = []byte{Version, NoAuthRequired}
	NoSupportedMethods           = []byte{5, 255}
)


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


func GetTargetAddressPort(targetRequestedCommand *byte, targetAddressType *byte, rawTargetAddress []byte, rawTargetPort []byte) (byte, string, string){
	if !(*targetRequestedCommand == 0 && *targetAddressType == 0) {
		switch *targetAddressType {
		case IPv4:
			return *targetRequestedCommand, net.IPv4(rawTargetAddress[0], rawTargetAddress[1], rawTargetAddress[2], rawTargetAddress[3]).String(), new(big.Int).SetBytes(rawTargetPort).String()
		case IPv6:
			return *targetRequestedCommand, Sockets.GetIPv6(rawTargetAddress), new(big.Int).SetBytes(rawTargetPort).String()
		case DomainName:
			return *targetRequestedCommand, string(rawTargetAddress[1:]), new(big.Int).SetBytes(rawTargetPort).String()
		}
	}
	return ConnectionRefused, "", ""
}


func CreateProxySession(clientConnection net.Conn, username *[]byte, passwordHash *[]byte) {

	var targetConnection net.Conn
	var targetAddress string
	clientConnectionReader := bufio.NewReader(clientConnection)
	clientConnectionWriter := bufio.NewWriter(clientConnection)

	// Receive connection
	clientHasCompatibleMethods := GetClientAuthenticationImplementedMethods(clientConnectionReader,
															  clientConnectionWriter,
															  username,
															  passwordHash)
	if !clientHasCompatibleMethods {
		_ = clientConnection.Close()
		return
	}
	// Receive and process connection request
	rawTargetRequestedCommand, targetAddressType, rawTargetAddress, rawTargetPort := ReceiveTargetRequest(clientConnectionReader)
	targetRequestedCommand, targetAddress, targetPort := GetTargetAddressPort(
		&rawTargetRequestedCommand, &targetAddressType, rawTargetAddress, rawTargetPort)
	if targetRequestedCommand == ConnectionRefused {
		_ = clientConnection.Close()
		return
	}

	targetConnection = HandleCommandExecution(
		clientConnection, clientConnectionReader, clientConnectionWriter, &targetRequestedCommand,
		&targetAddressType, &targetAddress, &targetPort, rawTargetAddress, rawTargetPort)
	if targetConnection != nil {
		_ = targetConnection.Close()
	}
	_ = clientConnection.Close()
}


func StartSocks5(ip string, port string, interfaceMode bool, username []byte, password []byte) {
	var passwordHash []byte
	rawPasswordHash := CryptoTools.GetPasswordHashPasswordByteArray(&username, &password)
	if rawPasswordHash == nil{
		passwordHash = []byte{}
	} else {
		passwordHash = rawPasswordHash
	}
	switch interfaceMode {
	case true:
		log.Println("Starting SOCKS5 server in Interface Mode")
		Interface.Client(ip, port, &username, &passwordHash, CreateProxySession)
	case false:
		log.Println("Starting SOCKS5 server in Bind Mode")
		BindServer.BindServer(ip, port, &username, &passwordHash, CreateProxySession)
	}
}
