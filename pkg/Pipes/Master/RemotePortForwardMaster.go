package Master

import (
	"bufio"
	"crypto/tls"
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Tools"
	"github.com/shoriwe/FullProxy/pkg/Tools/Types"
	"net"
	"time"
)

type RemotePortForward struct {
	Server                 net.Listener
	MasterConnection       net.Conn
	MasterConnectionReader *bufio.Reader
	MasterConnectionWriter *bufio.Writer
	TLSConfiguration       *tls.Config
	SlaveHost              string
	SlavePort              string
	LoggingMethod          Types.LoggingMethod
	Tries                  int
	Timeout                time.Duration
}

func (remotePortForward *RemotePortForward) SetInboundFilter(_ Types.IOFilter) error {
	return errors.New("This kind of PIPE doesn't support InboundFilters")
}

func (remotePortForward *RemotePortForward) SetOutboundFilter(_ Types.IOFilter) error {
	return errors.New("This kind of PIPE doesn't support OutboundFilters")
}

func (remotePortForward *RemotePortForward) SetTries(tries int) error {
	remotePortForward.Tries = tries
	return nil
}

func (remotePortForward *RemotePortForward) SetTimeout(timeout time.Duration) error {
	remotePortForward.Timeout = timeout
	return nil
}

func (remotePortForward *RemotePortForward) SetLoggingMethod(loggingMethod Types.LoggingMethod) error {
	remotePortForward.LoggingMethod = loggingMethod
	return nil
}

func (remotePortForward *RemotePortForward) Serve() error {
	var finalError error
	for {

	}
	Tools.LogData(remotePortForward.LoggingMethod, finalError)
	return finalError
}
