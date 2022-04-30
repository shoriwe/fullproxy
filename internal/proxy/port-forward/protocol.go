package port_forward

import (
	"errors"
	"github.com/shoriwe/fullproxy/v3/internal/global"
	"github.com/shoriwe/fullproxy/v3/internal/pipes"
	"net"
)

type Forward struct {
	TargetAddress string
	LoggingMethod global.LoggingMethod
	DialFunc      global.DialFunc
	ListenAddress *net.TCPAddr
}

func (localForward *Forward) SetListen(_ global.ListenFunc) {
}

func (localForward *Forward) SetListenAddress(address net.Addr) {
	localForward.ListenAddress = address.(*net.TCPAddr)
}

func NewForward(targetAddress string, loggingMethod global.LoggingMethod) global.Protocol {
	return &Forward{TargetAddress: targetAddress, LoggingMethod: loggingMethod}
}

func (localForward *Forward) SetAuthenticationMethod(_ global.AuthenticationMethod) error {
	return errors.New("This kind of proxy doesn't support authentication methods")
}

func (localForward *Forward) SetLoggingMethod(loggingMethod global.LoggingMethod) error {
	localForward.LoggingMethod = loggingMethod
	return nil
}

func (localForward *Forward) SetOutboundFilter(_ global.IOFilter) error {
	return errors.New("This kind of proxy doesn't support OutboundFilters")
}

func (localForward *Forward) SetDial(dialFunc global.DialFunc) {
	localForward.DialFunc = dialFunc
}

func (localForward *Forward) Handle(clientConnection net.Conn) error {
	defer clientConnection.Close()
	targetConnection, connectionError := localForward.DialFunc("tcp", localForward.TargetAddress)
	if connectionError != nil {
		return connectionError
	}
	return pipes.ForwardTraffic(clientConnection, targetConnection)
}
