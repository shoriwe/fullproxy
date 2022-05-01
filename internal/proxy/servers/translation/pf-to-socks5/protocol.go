package pf_to_socks5

import (
	"github.com/shoriwe/fullproxy/v3/internal/common"
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers"
	"golang.org/x/net/proxy"
	"net"
)

type customDialer struct {
	dialFunc servers.DialFunc
}

func (c *customDialer) Dial(networkType string, address string) (net.Conn, error) {
	return c.dialFunc(networkType, address)
}

type ForwardToSocks5 struct {
	TargetAddress string
	socks5Dialer  *customDialer
	Socks5Dialer  proxy.Dialer
}

func (forwardToSocks5 *ForwardToSocks5) SetDial(dialFunc servers.DialFunc) {
	forwardToSocks5.socks5Dialer.dialFunc = dialFunc
}

func NewForwardToSocks5(networkType, proxyAddress, username, password, targetAddress string) (*ForwardToSocks5, error) {
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
	}, nil
}

func (forwardToSocks5 *ForwardToSocks5) SetAuthenticationMethod(_ servers.AuthenticationMethod) error {
	return nil
}

func (forwardToSocks5 *ForwardToSocks5) Handle(clientConnection net.Conn) error {
	defer clientConnection.Close()
	targetConnection, connectionError := forwardToSocks5.Socks5Dialer.Dial("tcp", forwardToSocks5.TargetAddress)
	if connectionError != nil {
		return connectionError
	}
	defer targetConnection.Close()
	return common.ForwardTraffic(clientConnection, targetConnection)
}
