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
}

func (general *General) Serve() error {
	var finalError error
	for {
		clientConnection, connectionError := general.Server.Accept()
		if connectionError != nil {
			finalError = connectionError
			break
		}
		_, connectionError = Sockets.Send(general.MasterConnectionWriter, &ConnectionControllers.NewConnection)
		if connectionError != nil {
			finalError = connectionError
			break
		}

		targetConnection, connectionError := general.Server.Accept()
		if connectionError != nil {
			continue
		}
		// Verify that the new connection is also from the slave
		if strings.Split(targetConnection.RemoteAddr().String(), ":")[0] == general.MasterHost {
			targetConnection = Sockets.UpgradeServerToTLS(targetConnection, general.TLSConfiguration)
			go ConnectionControllers.StartGeneralProxying(clientConnection, targetConnection)
		}
	}
	return finalError
}
