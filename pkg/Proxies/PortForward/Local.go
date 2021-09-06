package PortForward

import (
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Pipes"
	"github.com/shoriwe/FullProxy/pkg/Tools/Types"
	"net"
	"time"
)

type LocalForward struct {
	TargetAddress string
	NetworkType   string
	TargetHost    string
	TargetPort    string
	LoggingMethod Types.LoggingMethod
	Tries         int
	Timeout       time.Duration
	InboundFilter Types.IOFilter
}

func (localForward *LocalForward) SetAuthenticationMethod(_ Types.AuthenticationMethod) error {
	return errors.New("This kind of proxy doesn't support authentication methods")
}

func (localForward *LocalForward) SetLoggingMethod(loggingMethod Types.LoggingMethod) error {
	localForward.LoggingMethod = loggingMethod
	return nil
}

func (localForward *LocalForward) SetTries(tries int) error {
	localForward.Tries = tries
	return nil
}

func (localForward *LocalForward) SetTimeout(timeout time.Duration) error {
	localForward.Timeout = timeout
	return nil
}

func (localForward *LocalForward) SetInboundFilter(filter Types.IOFilter) error {
	localForward.InboundFilter = filter
	return nil
}

func (localForward *LocalForward) SetOutboundFilter(_ Types.IOFilter) error {
	return errors.New("This kind of proxy doesn't support OutboundFilters")
}

func (localForward *LocalForward) Handle(clientConnection net.Conn) error {
	defer clientConnection.Close()
	targetConnection, connectionError := net.Dial(localForward.NetworkType, localForward.TargetAddress)
	if connectionError != nil {
		return connectionError
	}
	return Pipes.ForwardTraffic(clientConnection, targetConnection)
}
