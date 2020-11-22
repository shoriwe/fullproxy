package Master

import (
	"crypto/tls"
	"github.com/shoriwe/FullProxy/pkg/ConnectionHandlers"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"log"
	"net"
	"strings"
	"time"
)

func receiveMasterConnectionFromSlave(server net.Listener) net.Conn {
	masterConnection, connectionError := server.Accept()
	if connectionError != nil {
		_ = server.Close()
		log.Fatal(connectionError)
	}
	log.Print("Reverse connection received from: ", masterConnection.RemoteAddr())
	return masterConnection
}

func SetupMasterConnection(masterConnection net.Conn) (net.Conn, *tls.Config) {
	tlsConfiguration, configurationError := Sockets.CreateMasterTLSConfiguration()
	if configurationError == nil {
		masterConnection = Sockets.UpgradeServerToTLS(masterConnection, tlsConfiguration)
		return masterConnection, tlsConfiguration
	}
	log.Fatal(configurationError)
	return nil, nil
}

func masterHandler(host *string, port *string, masterFunction ConnectionHandlers.MasterFunction, args ...interface{}) {
	server := Sockets.Bind(host, port)
	log.Print("Waiting for proxy server connections...")
	masterConnection, tlsConfiguration := SetupMasterConnection(receiveMasterConnectionFromSlave(server))
	masterFunction(server, masterConnection, tlsConfiguration, args)
}

func basicMasterServer(server net.Listener, masterConnection net.Conn, tlsConfiguration *tls.Config, _ interface{}) {
	masterHost := strings.Split(masterConnection.RemoteAddr().String(), ":")[0]
	_, masterConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(masterConnection)
	for {
		clientConnection, connectionError := server.Accept()
		if connectionError != nil {
			log.Fatal(connectionError)
		}
		_, connectionError = Sockets.Send(masterConnectionWriter, &ConnectionHandlers.NewConnection)
		if connectionError != nil {
			log.Print(connectionError)
			break
		}

		targetConnection, connectionError := server.Accept()
		if connectionError == nil {
			// Verify that the new connection is also from the slave
			if strings.Split(targetConnection.RemoteAddr().String(), ":")[0] == masterHost {
				targetConnection = Sockets.UpgradeServerToTLS(targetConnection, tlsConfiguration)

				go ConnectionHandlers.StartGeneralProxying(clientConnection, targetConnection)
			}
		}
	}
	_ = masterConnection.Close()
	_ = server.Close()
}

func RemotePortForwardMasterServer(server net.Listener, masterConnection net.Conn, tlsConfiguration *tls.Config, args interface{}) {
	masterHost := strings.Split(masterConnection.RemoteAddr().String(), ":")[0]
	remoteHost := (args.([]interface{}))[0].(*string)
	remotePort := (args.([]interface{}))[1].(*string)
	masterConnectionReader, masterConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(masterConnection)
	for {
		_ = masterConnection.SetReadDeadline(time.Now().Add(20 * time.Second))
		numberOfBytesReceived, buffer, connectionError := Sockets.Receive(masterConnectionReader, 1)
		if connectionError == nil {
			if numberOfBytesReceived == 1 {
				if buffer[0] == ConnectionHandlers.NewConnection[0] {
					targetConnection, connectionError := Sockets.Connect(remoteHost, remotePort)
					if connectionError == nil {
						_, _ = Sockets.Send(masterConnectionWriter, &ConnectionHandlers.NewConnection)

						clientConnection, connectionError := server.Accept()
						clientConnection = Sockets.UpgradeServerToTLS(clientConnection, tlsConfiguration)

						if connectionError == nil {
							if strings.Split(clientConnection.RemoteAddr().String(), ":")[0] == masterHost {
								go ConnectionHandlers.StartGeneralProxying(clientConnection, targetConnection)
							}
						} else {
							log.Print(connectionError)
						}
					} else {
						log.Print("Could not connect to target server")
						_, _ = Sockets.Send(masterConnectionWriter, &ConnectionHandlers.FailToConnectToTarget)
					}
				}
			} else {
				log.Print("Error when interacting with slave")
				_, _ = Sockets.Send(masterConnectionWriter, &ConnectionHandlers.UnknownOperation)
			}
		} else if parsedConnectionError, ok := connectionError.(net.Error); !(ok && parsedConnectionError.Timeout()) {
			log.Print(connectionError)
			break
		}
	}
	_ = masterConnection.Close()
	_ = server.Close()
}

func Master(host *string, port *string, remoteHost *string, remotePort *string) {
	if len(*remoteHost) > 0 && len(*remotePort) > 0 {
		masterHandler(host, port, RemotePortForwardMasterServer, remoteHost, remotePort)
	} else {
		masterHandler(host, port, basicMasterServer)
	}
}
