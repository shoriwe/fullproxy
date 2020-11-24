package Pipes

import (
	"github.com/shoriwe/FullProxy/pkg/Proxies/RawProxy"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"github.com/shoriwe/FullProxy/pkg/Templates/Types"
	"log"
	"net"
	"time"
)

var (
	NewConnection = []byte{1}
	// Shutdown              = []byte{2}
	FailToConnectToTarget = []byte{3}
	UnknownOperation      = []byte{4}
)

func Serve(pipe Types.Pipe) {
	log.Fatal(pipe.Serve())
}

func StartGeneralProxying(clientConnection net.Conn, targetConnection net.Conn, tries int, timeout time.Duration) {
	clientConnectionReader, clientConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(clientConnection)
	targetConnectionReader, targetConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(targetConnection)
	if targetConnectionReader != nil && targetConnectionWriter != nil {
		rawProxy := RawProxy.RawProxy{
			TargetConnection:       targetConnection,
			TargetConnectionReader: targetConnectionReader,
			TargetConnectionWriter: targetConnectionWriter,
		}
		_ = rawProxy.SetTries(tries)
		_ = rawProxy.SetTimeout(timeout)
		rawProxy.Handle(clientConnection, clientConnectionReader, clientConnectionWriter)
	} else {
		_ = clientConnection.Close()
		_ = targetConnection.Close()
	}
}
