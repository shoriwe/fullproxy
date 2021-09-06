package Slave

import (
	"bufio"
	"crypto/tls"
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Tools/Types"
	"net"
	"time"
)

type General struct {
	MasterConnection       net.Conn
	MasterConnectionReader *bufio.Reader
	MasterConnectionWriter *bufio.Writer
	MasterHost             string
	MasterPort             string
	TLSConfiguration       *tls.Config
	ProxyProtocol          Types.ProxyProtocol
	LoggingMethod          Types.LoggingMethod
}

func (general *General) SetInboundFilter(_ Types.IOFilter) error {
	return errors.New("This kind of PIPE doesn't support InboundFilters")
}

func (general *General) SetOutboundFilter(_ Types.IOFilter) error {
	return errors.New("This kind of PIPE doesn't support OutboundFilters")
}

func (general *General) SetTries(tries int) error {
	return general.ProxyProtocol.SetTries(tries)
}

func (general *General) SetTimeout(timeout time.Duration) error {
	return general.ProxyProtocol.SetTimeout(timeout)
}

func (general *General) SetLoggingMethod(loggingMethod Types.LoggingMethod) error {
	general.LoggingMethod = loggingMethod
	return nil
}

func (general *General) Serve() error {
	var finalError error
	for {

	}
	return finalError
}
