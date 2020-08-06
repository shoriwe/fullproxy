package MasterSlave

import (
	"github.com/shoriwe/FullProxy/src/ConnectionStructures"
	"github.com/shoriwe/FullProxy/src/Proxies/Basic"
	"github.com/shoriwe/FullProxy/src/Sockets"
	"log"
	"net"
	"os"
	"os/signal"
)


func setupControlCSignal(server net.Listener, masterConnection net.Conn){
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(serverConnection  net.Listener, clientConnection net.Conn){
		<-c
		_ = clientConnection.Close()
		_ = masterConnection.Close()
		os.Exit(0)
	}(server, masterConnection)
}


func startProxying(clientConnection net.Conn, targetConnection net.Conn){
	clientConnectionReader, clientConnectionWriter := ConnectionStructures.CreateReaderWriter(clientConnection)
	targetConnectionReader, targetConnectionWriter := ConnectionStructures.CreateTunnelReaderWriter(targetConnection)
	Basic.Proxy(
		clientConnection, targetConnection,
		clientConnectionReader, clientConnectionWriter,
		targetConnectionReader, targetConnectionWriter)
}


func Master(address *string, port *string){
	log.Print("Starting Master server")
	server, BindingError  := net.Listen("tcp", *address + ":" + *port)
	if BindingError != nil {
		log.Print("Something goes wrong: " + BindingError.Error())
		return
	}
	log.Printf("Bind successfully in %s:%s", *address, *port)
	log.Print("Waiting for proxy server connections...")
	masterConnection, connectionError := server.Accept()
	if connectionError != nil{
		_ = server.Close()
		return
	}
	log.Print("Reverse connection received from: ", masterConnection.RemoteAddr())
	setupControlCSignal(server, masterConnection)
	_, masterConnectionWriter := ConnectionStructures.CreateSocketConnectionReaderWriter(masterConnection)
	for {
		clientConnection, _ := server.Accept()
		_, connectionError := Sockets.Send(masterConnectionWriter, &NewConnection)
		if connectionError != nil{
			break
		}
		targetConnection, _ := server.Accept()
		go startProxying(clientConnection, targetConnection)
	}
	_ = masterConnection.Close()
	_ = server.Close()
}