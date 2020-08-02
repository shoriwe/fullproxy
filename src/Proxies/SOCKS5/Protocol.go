package SOCKS5

import (
	"bufio"
	"github.com/shoriwe/FullProxy/src/BindServer"
	"github.com/shoriwe/FullProxy/src/CryptoTools"
	"github.com/shoriwe/FullProxy/src/Interface"
	"github.com/shoriwe/FullProxy/src/Sockets"
	"log"
	"math/big"
	"net"
)


func ReceiveTargetRequest(clientConnectionReader *bufio.Reader) (byte, byte, []byte, []byte) {
	numberOfBytesReceived, targetRequest, ConnectionError := Sockets.Receive(clientConnectionReader, 1024)
	if ConnectionError == nil{
		if targetRequest[0] == Version {
			if targetRequest[1] == Connect || targetRequest[1] == Bind || targetRequest[1] == UDPAssociate {
				if targetRequest[3] == IPv4 || targetRequest[3] == IPv6 || targetRequest[3] == DomainName {
					return targetRequest[1], targetRequest[3], targetRequest[4 : numberOfBytesReceived-2], targetRequest[numberOfBytesReceived-2 : numberOfBytesReceived]
				}
			}
		}
	}
	return 0, 0, nil, nil
}


func GetTargetAddressPort(targetRequestedCommand *byte, targetAddressType *byte, rawTargetAddress []byte, rawTargetPort []byte) (byte, string, string){
	if *targetRequestedCommand != 0 && *targetAddressType != 0 {
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
	var targetRequestedCommand byte
	clientConnectionReader := bufio.NewReader(clientConnection)
	clientConnectionWriter := bufio.NewWriter(clientConnection)

	// Receive connection
	clientHasCompatibleMethods := GetClientAuthenticationImplementedMethods(
		clientConnectionReader,
		clientConnectionWriter,
		username,
		passwordHash)
	if clientHasCompatibleMethods{
		var targetAddress string
		var targetPort string
		rawTargetRequestedCommand, targetAddressType, rawTargetAddress, rawTargetPort := ReceiveTargetRequest(
			clientConnectionReader)
		targetRequestedCommand, targetAddress, targetPort = GetTargetAddressPort(
			&rawTargetRequestedCommand, &targetAddressType,
			rawTargetAddress, rawTargetPort)
		if targetRequestedCommand != ConnectionRefused{
			HandleCommandExecution(
				clientConnection, clientConnectionReader, clientConnectionWriter, &targetRequestedCommand,
				&targetAddressType, &targetAddress, &targetPort)
		}
	}
	if (!clientHasCompatibleMethods) || (targetRequestedCommand == ConnectionRefused){
		_ = clientConnection.Close()
	}
}


func StartSocks5(ip string, port string, interfaceMode bool, username []byte, password []byte) {
	passwordHash := CryptoTools.GetPasswordHashPasswordByteArray(&username, &password)
	switch interfaceMode {
	case true:
		log.Println("Starting SOCKS5 server in Interface Mode")
		Interface.Client(ip, port, &username, &passwordHash, CreateProxySession)
	case false:
		log.Println("Starting SOCKS5 server in Bind Mode")
		BindServer.BindServer(ip, port, &username, &passwordHash, CreateProxySession)
	}
}
