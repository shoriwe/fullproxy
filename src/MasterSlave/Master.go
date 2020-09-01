package MasterSlave

import (
	"github.com/shoriwe/FullProxy/src/ConnectionStructures"
	"github.com/shoriwe/FullProxy/src/Proxies/Basic"
	"github.com/shoriwe/FullProxy/src/Sockets"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

type MasterFunction func(server net.Listener, masterConnection net.Conn, args interface{})

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
	clientConnectionReader, clientConnectionWriter := ConnectionStructures.CreateReaderWriter(clientConnection)
	targetConnectionReader, targetConnectionWriter := ConnectionStructures.CreateTunnelReaderWriter(targetConnection)
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

func masterHandler(address *string, port *string, masterFunction MasterFunction, args ...interface{}) {
	server := Sockets.Bind(address, port)
	log.Print("Waiting for proxy server connections...")
	masterConnection := receiveMasterConnectionFromSlave(server)
	setupControlCSignal(server, masterConnection)
	masterFunction(server, masterConnection, args)
}

func basicMasterServer(server net.Listener, masterConnection net.Conn, _ interface{}) {
	_, masterConnectionWriter := ConnectionStructures.CreateSocketConnectionReaderWriter(masterConnection)
	for {
		clientConnection, _ := server.Accept()
		_, connectionError := Sockets.Send(masterConnectionWriter, &NewConnection)
		if connectionError != nil {
			log.Print(connectionError)
			break
		}
		targetConnection, _ := server.Accept()
		go startGeneralProxying(clientConnection, targetConnection)
	}
	_ = masterConnection.Close()
	_ = server.Close()
}

func portForwardMasterServe(server net.Listener, masterConnection net.Conn, args interface{}) {
	remoteAddress := (args.([]interface{}))[0].(*string)
	remotePort := (args.([]interface{}))[1].(*string)
	masterConnectionReader, masterConnectionWriter := ConnectionStructures.CreateSocketConnectionReaderWriter(masterConnection)
	for {
		_ = masterConnection.SetReadDeadline(time.Now().Add(3 * time.Second))
		numberOfBytesReceived, _, connectionError := Sockets.Receive(masterConnectionReader, 1)
		if connectionError == nil {
			if numberOfBytesReceived == 1 {
				clientConnection := Sockets.Connect(remoteAddress, remotePort)
				if clientConnection != nil {
					_, _ = Sockets.Send(masterConnectionWriter, &NewConnection)
					targetConnection, connectionError := server.Accept()
					if connectionError == nil {
						go startGeneralProxying(clientConnection, targetConnection)
					} else {
						log.Print(connectionError)
					}
				} else {
					log.Print("Could not connect to target server")
					_, _ = Sockets.Send(masterConnectionWriter, &FailToConnectToTarget)
				}
			} else {
				log.Print("Error when interacting with slave")
				_, _ = Sockets.Send(masterConnectionWriter, &UnknownOperation)
			}
		} else if parsedConnectionError, ok := connectionError.(net.Error); !(ok && parsedConnectionError.Timeout()) {
			break
		}
	}
	_ = masterConnection.Close()
	_ = server.Close()
}

func Master(address *string, port *string, remoteAddress *string, remotePort *string) {
	if len(*remoteAddress) > 0 && len(*remotePort) > 0 {
		masterHandler(address, port, portForwardMasterServe, remoteAddress, remotePort)
	} else {
		masterHandler(address, port, basicMasterServer)
	}
}
