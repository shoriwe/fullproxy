package ForwardToSocks5

import (
	"github.com/shoriwe/FullProxy/v2/pkg/Pipes"
	"github.com/shoriwe/FullProxy/v2/pkg/Tools"
	"github.com/shoriwe/FullProxy/v2/pkg/Tools/Types"
	"golang.org/x/net/proxy"
	"net"
)

type customDialer struct {
	dialFunc Types.DialFunc
}

func (c *customDialer) Dial(networkType string, address string) (net.Conn, error) {
	return c.dialFunc(networkType, address)
}

type ForwardToSocks5 struct {
	TargetAddress string
	socks5Dialer  *customDialer
	Socks5Dialer  proxy.Dialer
	LoggingMethod Types.LoggingMethod
}

func (forwardToSocks5 *ForwardToSocks5) SetOutboundFilter(_ Types.IOFilter) error {
	return nil
}

func (forwardToSocks5 *ForwardToSocks5) SetDial(dialFunc Types.DialFunc) {
	forwardToSocks5.socks5Dialer.dialFunc = dialFunc
}

func NewForwardToSocks5(networkType, proxyAddress, username, password, targetAddress string, loggingMethod Types.LoggingMethod) (*ForwardToSocks5, error) {
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

func (forwardToSocks5 *ForwardToSocks5) SetLoggingMethod(loggingMethod Types.LoggingMethod) error {
	forwardToSocks5.LoggingMethod = loggingMethod
	return nil
}
func (forwardToSocks5 *ForwardToSocks5) SetAuthenticationMethod(_ Types.AuthenticationMethod) error {
	return nil
}

func (forwardToSocks5 *ForwardToSocks5) Handle(clientConnection net.Conn) error {
	targetConnection, connectionError := forwardToSocks5.Socks5Dialer.Dial("tcp", forwardToSocks5.TargetAddress)
	if connectionError != nil {
		Tools.LogData(forwardToSocks5.LoggingMethod, connectionError)
		_ = clientConnection.Close()
		return connectionError
	}
	return Pipes.ForwardTraffic(clientConnection, targetConnection)
}
