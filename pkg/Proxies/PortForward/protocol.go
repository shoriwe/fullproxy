package PortForward

import (
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Pipes"
	"github.com/shoriwe/FullProxy/pkg/Tools/Types"
	"net"
)

type Forward struct {
	NetworkType   string
	TargetAddress string
	LoggingMethod Types.LoggingMethod
	DialFunc      Types.DialFunc
}

func NewForward(networkType string, targetAddress string, loggingMethod Types.LoggingMethod) *Forward {
	return &Forward{NetworkType: networkType, TargetAddress: targetAddress, LoggingMethod: loggingMethod}
}

func (localForward *Forward) SetAuthenticationMethod(_ Types.AuthenticationMethod) error {
	return errors.New("This kind of proxy doesn't support authentication methods")
}

func (localForward *Forward) SetLoggingMethod(loggingMethod Types.LoggingMethod) error {
	localForward.LoggingMethod = loggingMethod
	return nil
}

func (localForward *Forward) SetOutboundFilter(_ Types.IOFilter) error {
	return errors.New("This kind of proxy doesn't support OutboundFilters")
}

func (localForward *Forward) SetDial(dialFunc Types.DialFunc) {
	localForward.DialFunc = dialFunc
}

func (localForward *Forward) Handle(clientConnection net.Conn) error {
	defer clientConnection.Close()
	targetConnection, connectionError := localForward.DialFunc(localForward.NetworkType, localForward.TargetAddress)
	if connectionError != nil {
		return connectionError
	}
	return Pipes.ForwardTraffic(clientConnection, targetConnection)
}
