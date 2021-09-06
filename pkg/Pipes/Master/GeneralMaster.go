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

	}
	Tools.LogData(general.LoggingMethod, finalError)
	return finalError
}
