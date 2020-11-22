package Slave

import (
	"crypto/tls"
	"github.com/shoriwe/FullProxy/pkg/ConnectionHandlers"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"log"
	"net"
	"time"
)

func ConnectToMasterServer(masterHost *string, masterPort *string) (net.Conn, *tls.Config) {
	tlsConfiguration, configurationCreationError := Sockets.CreateSlaveTLSConfiguration()
	if configurationCreationError == nil {
		log.Printf("Trying to connecto to %s:%s", *masterHost, *masterPort)
		masterConnection, connectionError := Sockets.TLSConnect(masterHost, masterPort, tlsConfiguration)
		log.Printf("Successfully connected to %s:%s", *masterHost, *masterPort)
		if connectionError == nil {
			return masterConnection, tlsConfiguration
		}
		log.Fatal("Could not connect to the master")
	}
	log.Fatal("Could not create the TLS certificate")
	return nil, nil
}

func GeneralSlave(masterHost *string, masterPort *string, protocol ConnectionHandlers.ProxyProtocol) {
	masterConnection, tlsConfiguration := ConnectToMasterServer(masterHost, masterPort)
	masterConnectionReader, _ := Sockets.CreateSocketConnectionReaderWriter(masterConnection)
	for {
		timeoutSetError := masterConnection.SetReadDeadline(time.Now().Add(20 * time.Second))
		if timeoutSetError != nil {
			log.Fatal(timeoutSetError)
		}
		NumberOfReceivedBytes, buffer, connectionError := Sockets.Receive(masterConnectionReader, 1024)
		if connectionError == nil {
			if NumberOfReceivedBytes == 1 {
				if buffer[0] == ConnectionHandlers.NewConnection[0] {
					clientConnection, connectionError := Sockets.TLSConnect(masterHost, masterPort, tlsConfiguration)

					if connectionError == nil {
						clientConnectionReader, clientConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(clientConnection)
						if clientConnectionReader != nil && clientConnectionWriter != nil {
							go protocol.Handle(clientConnection, clientConnectionReader, clientConnectionWriter)
						} else {
							_ = clientConnection.Close()
						}
					}
				}
			} else {
				continue
			}
		} else if parsedConnectionError, ok := connectionError.(net.Error); !(ok && parsedConnectionError.Timeout()) {
			log.Print(connectionError)
			break
		}
	}
}

func RemotePortForwardSlave(
	masterHost *string, masterPort *string,
	localHost *string, localPort *string,
	protocol ConnectionHandlers.ProxyProtocol) {
	localServer := Sockets.Bind(localHost, localPort)
	masterConnection, tlsConfiguration := ConnectToMasterServer(masterHost, masterPort)
	masterConnectionReader, masterConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(masterConnection)
	for {
		clientConnection, connectionError := localServer.Accept()
		if connectionError == nil {
			log.Print("Connection from: ", clientConnection.RemoteAddr().String())
			_, connectionError := Sockets.Send(masterConnectionWriter, &ConnectionHandlers.NewConnection)
			if connectionError == nil {
				_ = masterConnection.SetReadDeadline(time.Now().Add(3 * time.Second))
				numberOfBytesReceived, response, connectionError := Sockets.Receive(masterConnectionReader, 1)
				if connectionError == nil {
					if numberOfBytesReceived == 1 {
						switch response[0] {
						case ConnectionHandlers.NewConnection[0]:
							targetConnection, connectionError := Sockets.TLSConnect(masterHost, masterPort, tlsConfiguration)
							if connectionError == nil {
								go ConnectionHandlers.StartGeneralProxying(clientConnection, targetConnection)
							} else {
								_ = clientConnection.Close()
								log.Fatal("Connectivity issues with master server")
							}
						case ConnectionHandlers.FailToConnectToTarget[0]:
							_ = clientConnection.Close()
							log.Print("Something goes wrong when master connected to target")
							break
						case ConnectionHandlers.UnknownOperation[0]:
							_ = clientConnection.Close()
							log.Print("The master did not understood the message")
							break
						}
					}
				} else if parsedConnectionError, ok := connectionError.(net.Error); !(ok && parsedConnectionError.Timeout()) {
					log.Print(connectionError)
					break
				}
			}
		}
	}
}
