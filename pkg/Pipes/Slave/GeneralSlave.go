package Slave

import (
	"bufio"
	"crypto/tls"
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Pipes"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"github.com/shoriwe/FullProxy/pkg/Templates"
	"github.com/shoriwe/FullProxy/pkg/Templates/Types"
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

func (general *General) SetInboundFilter(_ Types.InboundFilter) error {
	return errors.New("This kind of PIPE doesn't support InboundFilters")
}

func (general *General) SetOutboundFilter(_ Types.OutboundFilter) error {
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
		_ = general.MasterConnection.SetReadDeadline(time.Now().Add(20 * time.Second))
		NumberOfReceivedBytes, buffer, connectionError := Sockets.Receive(general.MasterConnectionReader, 1024)
		if connectionError != nil {
			if parsedConnectionError, ok := connectionError.(net.Error); !(ok && parsedConnectionError.Timeout()) {
				finalError = connectionError
				break
			}
		}
		if NumberOfReceivedBytes != 1 {
			continue
		}
		if buffer[0] != Pipes.NewConnection[0] {
			continue
		}
		clientConnection, connectionError := Sockets.TLSConnect(&general.MasterHost, &general.MasterPort, general.TLSConfiguration)
		if connectionError != nil {
			finalError = connectionError
			break
		}
		Templates.LogData(general.LoggingMethod, "Client connection received from: ", clientConnection.RemoteAddr().String())
		clientConnectionReader, clientConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(clientConnection)
		go general.ProxyProtocol.Handle(clientConnection, clientConnectionReader, clientConnectionWriter)
	}
	return finalError
}
