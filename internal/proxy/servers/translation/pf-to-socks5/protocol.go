package pf_to_socks5

import (
	"github.com/shoriwe/fullproxy/v3/internal/common"
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"net/url"
)

type customDialer struct {
	dialFunc servers.DialFunc
}

func (c *customDialer) Dial(networkType string, address string) (net.Conn, error) {
	return c.dialFunc(networkType, address)
}

type ForwardToSocks5 struct {
	TargetNetwork                    string
	TargetAddress                    string
	socks5Dialer                     *customDialer
	Socks5Dialer                     proxy.Dialer
	IncomingSniffer, OutgoingSniffer io.Writer
}

func (forwardToSocks5 *ForwardToSocks5) SetSniffers(incoming, outgoing io.Writer) {
	forwardToSocks5.IncomingSniffer = incoming
	forwardToSocks5.OutgoingSniffer = outgoing
}

func (forwardToSocks5 *ForwardToSocks5) SetListen(_ servers.ListenFunc) {
}

func (forwardToSocks5 *ForwardToSocks5) SetListenAddress(_ net.Addr) {
}

func (forwardToSocks5 *ForwardToSocks5) SetDial(dialFunc servers.DialFunc) {
	forwardToSocks5.socks5Dialer.dialFunc = dialFunc
}

func NewForwardToSocks5(
	proxyNetwork, proxyAddress string,
	user *url.Userinfo,
	targetNetwork, targetAddress string,
) (servers.Protocol, error) {
	var (
		dialer              proxy.Dialer
		initializationError error
	)
	fDialer := &customDialer{dialFunc: net.Dial}
	if user != nil {
		password, _ := user.Password()
		dialer, initializationError = proxy.SOCKS5(proxyNetwork, proxyAddress, &proxy.Auth{
			User:     user.Username(),
			Password: password,
		}, fDialer)
	} else {
		dialer, initializationError = proxy.SOCKS5(proxyNetwork, proxyAddress, nil, fDialer)
	}
	if initializationError != nil {
		return nil, initializationError
	}
	return &ForwardToSocks5{
		TargetNetwork: targetNetwork,
		TargetAddress: targetAddress,
		socks5Dialer:  fDialer,
		Socks5Dialer:  dialer,
	}, nil
}

func (forwardToSocks5 *ForwardToSocks5) SetAuthenticationMethod(_ servers.AuthenticationMethod) {
}

func (forwardToSocks5 *ForwardToSocks5) Handle(clientConnection net.Conn) error {
	defer clientConnection.Close()
	targetConnection, connectionError := forwardToSocks5.Socks5Dialer.Dial(forwardToSocks5.TargetNetwork, forwardToSocks5.TargetAddress)
	if connectionError != nil {
		return connectionError
	}
	defer targetConnection.Close()
	return common.ForwardTraffic(
		clientConnection, targetConnection,
		forwardToSocks5.IncomingSniffer, forwardToSocks5.OutgoingSniffer)
}
