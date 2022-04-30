package pf_to_socks5

import (
	"github.com/shoriwe/fullproxy/v3/internal/pipes"
	proxy2 "github.com/shoriwe/fullproxy/v3/internal/proxy"
	"golang.org/x/net/proxy"
	"net"
)

type customDialer struct {
	dialFunc proxy2.DialFunc
}

func (c *customDialer) Dial(networkType string, address string) (net.Conn, error) {
	return c.dialFunc(networkType, address)
}

type ForwardToSocks5 struct {
	TargetAddress string
	socks5Dialer  *customDialer
	Socks5Dialer  proxy.Dialer
	LoggingMethod proxy2.LoggingMethod
}

func (forwardToSocks5 *ForwardToSocks5) SetOutboundFilter(_ proxy2.IOFilter) error {
	return nil
}

func (forwardToSocks5 *ForwardToSocks5) SetDial(dialFunc proxy2.DialFunc) {
	forwardToSocks5.socks5Dialer.dialFunc = dialFunc
}

func NewForwardToSocks5(networkType, proxyAddress, username, password, targetAddress string, loggingMethod proxy2.LoggingMethod) (*ForwardToSocks5, error) {
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

func (forwardToSocks5 *ForwardToSocks5) SetLoggingMethod(loggingMethod proxy2.LoggingMethod) error {
	forwardToSocks5.LoggingMethod = loggingMethod
	return nil
}
func (forwardToSocks5 *ForwardToSocks5) SetAuthenticationMethod(_ proxy2.AuthenticationMethod) error {
	return nil
}

func (forwardToSocks5 *ForwardToSocks5) Handle(clientConnection net.Conn) error {
	defer clientConnection.Close()
	targetConnection, connectionError := forwardToSocks5.Socks5Dialer.Dial("tcp", forwardToSocks5.TargetAddress)
	if connectionError != nil {
		return connectionError
	}
	defer targetConnection.Close()
	return pipes.ForwardTraffic(clientConnection, targetConnection)
}
