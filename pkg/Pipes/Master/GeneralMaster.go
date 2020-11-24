package Master

import (
	"bufio"
	"crypto/tls"
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Pipes"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"github.com/shoriwe/FullProxy/pkg/Templates"
	"github.com/shoriwe/FullProxy/pkg/Templates/Types"
	"net"
	"strings"
	"time"
)

type General struct {
	MasterConnection       net.Conn
	MasterConnectionReader *bufio.Reader
	MasterConnectionWriter *bufio.Writer
	SlaveHost              string
	TLSConfiguration       *tls.Config
	Server                 net.Listener
	LoggingMethod          Types.LoggingMethod
	Tries                  int
	Timeout                time.Duration
	InboundFilter          Types.IOFilter
}

func (general *General) SetInboundFilter(filter Types.IOFilter) error {
	general.InboundFilter = filter
	return nil
}

func (general *General) SetOutboundFilter(_ Types.IOFilter) error {
	return errors.New("This kind of PIPE doesn't support OutboundFilters")
}

func (general *General) SetTries(tries int) error {
	general.Tries = tries
	return nil
}

func (general *General) SetTimeout(timeout time.Duration) error {
	general.Timeout = timeout
	return nil
}

func (general *General) SetLoggingMethod(loggingMethod Types.LoggingMethod) error {
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
		if !Templates.FilterInbound(general.InboundFilter, Templates.ParseIP(clientConnection.RemoteAddr().String())) {
			_ = clientConnection.Close()
			Templates.LogData(general.LoggingMethod, "Unwanted connection received from "+clientConnection.RemoteAddr().String())
			continue
		}
		Templates.LogData(general.LoggingMethod, "Client connection received from: ", clientConnection.RemoteAddr().String())
		_, connectionError = Sockets.Send(general.MasterConnectionWriter, &Pipes.NewConnection)
		if connectionError != nil {
			finalError = connectionError
			break
		}

		targetConnection, connectionError := general.Server.Accept()
		if connectionError != nil {
			Templates.LogData(general.LoggingMethod, connectionError)
			continue
		}
		// Verify that the new connection is also from the slave
		if strings.Split(targetConnection.RemoteAddr().String(), ":")[0] != general.SlaveHost {
			Templates.LogData(general.LoggingMethod, "(Ignoring) Connection received from a non slave client: ", targetConnection.RemoteAddr().String())
			_ = targetConnection.Close()
			_ = clientConnection.Close()
			continue
		}
		Templates.LogData(general.LoggingMethod, "Target connection received from: ", targetConnection.RemoteAddr().String())
		targetConnection = Sockets.UpgradeServerToTLS(targetConnection, general.TLSConfiguration)
		go Pipes.StartGeneralProxying(
			clientConnection, targetConnection,
			Templates.GetTries(general.Tries), Templates.GetTimeout(general.Timeout))
	}
	Templates.LogData(general.LoggingMethod, finalError)
	return finalError
}
