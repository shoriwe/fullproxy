package MasterSlave

import (
	"crypto/tls"
	"github.com/shoriwe/FullProxy/src/ConnectionStructures"
	"github.com/shoriwe/FullProxy/src/Proxies/Basic"
	"github.com/shoriwe/FullProxy/src/Sockets"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"
)

type MasterFunction func(server net.Listener, masterConnection net.Conn, tlsConfiguration *tls.Config, args interface{})

func setupControlCSignal(server net.Listener, masterConnection net.Conn) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(serverConnection net.Listener, clientConnection net.Conn) {
		<-c
		_ = clientConnection.Close()
		_ = masterConnection.Close()
		os.Exit(0)
	}(server, masterConnection)
}

func startGeneralProxying(clientConnection net.Conn, targetConnection net.Conn) {
	clientConnectionReader, clientConnectionWriter := ConnectionStructures.CreateSocketConnectionReaderWriter(clientConnection)
	targetConnectionReader, targetConnectionWriter := ConnectionStructures.CreateSocketConnectionReaderWriter(targetConnection)
	if targetConnectionReader != nil && targetConnectionWriter != nil {
		Basic.Proxy(
			clientConnection, targetConnection,
			clientConnectionReader, clientConnectionWriter,
			targetConnectionReader, targetConnectionWriter)
	} else {
		_ = clientConnection.Close()
		_ = targetConnection.Close()
	}
}

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

func masterHandler(address *string, port *string, masterFunction MasterFunction, args ...interface{}) {
	server := Sockets.Bind(address, port)
	log.Print("Waiting for proxy server connections...")
	masterConnection, tlsConfiguration := SetupMasterConnection(receiveMasterConnectionFromSlave(server))
	setupControlCSignal(server, masterConnection)
	masterFunction(server, masterConnection, tlsConfiguration, args)
}

func basicMasterServer(server net.Listener, masterConnection net.Conn, tlsConfiguration *tls.Config, _ interface{}) {
	masterAddress := strings.Split(masterConnection.RemoteAddr().String(), ":")[0]
	_, masterConnectionWriter := ConnectionStructures.CreateSocketConnectionReaderWriter(masterConnection)
	for {
		clientConnection, connectionError := server.Accept()
		if connectionError != nil {
			log.Fatal(connectionError)
		}
		_, connectionError = Sockets.Send(masterConnectionWriter, &NewConnection)
		if connectionError != nil {
			log.Print(connectionError)
			break
		}

		targetConnection, connectionError := server.Accept()
		if connectionError == nil {
			// Verify that the new connection is also from the slave
			if strings.Split(targetConnection.RemoteAddr().String(), ":")[0] == masterAddress {
				targetConnection = Sockets.UpgradeServerToTLS(targetConnection, tlsConfiguration)

				go startGeneralProxying(clientConnection, targetConnection)
			}
		}
	}
	_ = masterConnection.Close()
	_ = server.Close()
}

func RemotePortForwardMasterServer(server net.Listener, masterConnection net.Conn, tlsConfiguration *tls.Config, args interface{}) {
	masterAddress := strings.Split(masterConnection.RemoteAddr().String(), ":")[0]
	remoteAddress := (args.([]interface{}))[0].(*string)
	remotePort := (args.([]interface{}))[1].(*string)
	masterConnectionReader, masterConnectionWriter := ConnectionStructures.CreateSocketConnectionReaderWriter(masterConnection)
	for {
		_ = masterConnection.SetReadDeadline(time.Now().Add(20 * time.Second))
		numberOfBytesReceived, buffer, connectionError := Sockets.Receive(masterConnectionReader, 1)
		if connectionError == nil {
			if numberOfBytesReceived == 1 {
				if buffer[0] == NewConnection[0] {
					targetConnection := Sockets.Connect(remoteAddress, remotePort)
					if targetConnection != nil {
						_, _ = Sockets.Send(masterConnectionWriter, &NewConnection)

						clientConnection, connectionError := server.Accept()
						clientConnection = Sockets.UpgradeServerToTLS(clientConnection, tlsConfiguration)

						if connectionError == nil {
							if strings.Split(clientConnection.RemoteAddr().String(), ":")[0] == masterAddress {
								go startGeneralProxying(clientConnection, targetConnection)
							}
						} else {
							log.Print(connectionError)
						}
					} else {
						log.Print("Could not connect to target server")
						_, _ = Sockets.Send(masterConnectionWriter, &FailToConnectToTarget)
					}
				}
			} else {
				log.Print("Error when interacting with slave")
				_, _ = Sockets.Send(masterConnectionWriter, &UnknownOperation)
			}
		} else if parsedConnectionError, ok := connectionError.(net.Error); !(ok && parsedConnectionError.Timeout()) {
			log.Print(connectionError)
			break
		}
	}
	_ = masterConnection.Close()
	_ = server.Close()
}

func Master(address *string, port *string, remoteAddress *string, remotePort *string) {
	if len(*remoteAddress) > 0 && len(*remotePort) > 0 {
		masterHandler(address, port, RemotePortForwardMasterServer, remoteAddress, remotePort)
	} else {
		masterHandler(address, port, basicMasterServer)
	}
}
