package Slave

import (
	"bufio"
	"crypto/tls"
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Tools/Types"
	"net"
	"time"
)

type RemotePortForward struct {
	LocalServer            net.Listener
	MasterConnection       net.Conn
	MasterConnectionReader *bufio.Reader
	MasterConnectionWriter *bufio.Writer
	MasterHost             string
	MasterPort             string
	TLSConfiguration       *tls.Config
	LoggingMethod          Types.LoggingMethod
	Tries                  int
	Timeout                time.Duration
	InboundFilter          Types.IOFilter
}

func (remotePortForward *RemotePortForward) SetInboundFilter(filter Types.IOFilter) error {
	remotePortForward.InboundFilter = filter
	return nil
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
	return finalError
}
