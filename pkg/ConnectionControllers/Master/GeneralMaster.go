package Master

import (
	"bufio"
	"crypto/tls"
	"github.com/shoriwe/FullProxy/pkg/ConnectionControllers"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"net"
	"strings"
)

type General struct {
	MasterConnection       net.Conn
	MasterConnectionReader *bufio.Reader
	MasterConnectionWriter *bufio.Writer
	MasterHost             string
	TLSConfiguration       *tls.Config
	Server                 net.Listener
	LoggingMethod ConnectionControllers.LoggingMethod
}

func (general *General) SetLoggingMethod(loggingMethod ConnectionControllers.LoggingMethod) error {
	general.LoggingMethod = loggingMethod
	return nil
}

func (general *General) Serve() error {
	var finalError error
	for {
		clientConnection, connectionError := general.Server.Accept()
		if connectionError != nil {
			finalError = connectionError
			break
		}
		ConnectionControllers.LogData(general.LoggingMethod, "Client connection received from: ", clientConnection.RemoteAddr().String())
		_, connectionError = Sockets.Send(general.MasterConnectionWriter, &ConnectionControllers.NewConnection)
		if connectionError != nil {
			finalError = connectionError
			break
		}

		targetConnection, connectionError := general.Server.Accept()
		if connectionError != nil {
			ConnectionControllers.LogData(general.LoggingMethod, connectionError)
			continue
		}
		// Verify that the new connection is also from the slave
		if strings.Split(targetConnection.RemoteAddr().String(), ":")[0] == general.MasterHost {
			ConnectionControllers.LogData(general.LoggingMethod, "Target connection received from: ", targetConnection.RemoteAddr().String())
			targetConnection = Sockets.UpgradeServerToTLS(targetConnection, general.TLSConfiguration)
			go ConnectionControllers.StartGeneralProxying(clientConnection, targetConnection)
		} else {
			ConnectionControllers.LogData(general.LoggingMethod, "(Ignoring) Connection received from a non slave client: ", targetConnection.RemoteAddr().String())
		}
	}
	ConnectionControllers.LogData(general.LoggingMethod, finalError)
	return finalError
}
