package pf_to_socks5

import (
	"github.com/shoriwe/FullProxy/v3/internal/global"
	"github.com/shoriwe/FullProxy/v3/internal/pipes"
	"golang.org/x/net/proxy"
	"net"
)

type customDialer struct {
	dialFunc global.DialFunc
}

func (c *customDialer) Dial(networkType string, address string) (net.Conn, error) {
	return c.dialFunc(networkType, address)
}

type ForwardToSocks5 struct {
	TargetAddress string
	socks5Dialer  *customDialer
	Socks5Dialer  proxy.Dialer
	LoggingMethod global.LoggingMethod
}

func (forwardToSocks5 *ForwardToSocks5) SetOutboundFilter(_ global.IOFilter) error {
	return nil
}

func (forwardToSocks5 *ForwardToSocks5) SetDial(dialFunc global.DialFunc) {
	forwardToSocks5.socks5Dialer.dialFunc = dialFunc
}

func NewForwardToSocks5(networkType, proxyAddress, username, password, targetAddress string, loggingMethod global.LoggingMethod) (*ForwardToSocks5, error) {
	fDialer := &customDialer{dialFunc: net.Dial}
	dialer, initializationError := proxy.SOCKS5(networkType, proxyAddress, &proxy.Auth{
		User:     username,
		Password: password,
	}, fDialer)
	if initializationError != nil {
		return nil, initializationError
	}
	return &ForwardToSocks5{
		TargetAddress: targetAddress,
		socks5Dialer:  fDialer,
		Socks5Dialer:  dialer,
		LoggingMethod: loggingMethod,
	}, nil
}

func (forwardToSocks5 *ForwardToSocks5) SetLoggingMethod(loggingMethod global.LoggingMethod) error {
	forwardToSocks5.LoggingMethod = loggingMethod
	return nil
}
func (forwardToSocks5 *ForwardToSocks5) SetAuthenticationMethod(_ global.AuthenticationMethod) error {
	return nil
}

func (forwardToSocks5 *ForwardToSocks5) Handle(clientConnection net.Conn) error {
	targetConnection, connectionError := forwardToSocks5.Socks5Dialer.Dial("tcp", forwardToSocks5.TargetAddress)
	if connectionError != nil {
		global.LogData(forwardToSocks5.LoggingMethod, connectionError)
		_ = clientConnection.Close()
		return connectionError
	}
	return pipes.ForwardTraffic(clientConnection, targetConnection)
}
