package port_forward

import (
	"errors"
	"github.com/shoriwe/fullproxy/v3/internal/pipes"
	"github.com/shoriwe/fullproxy/v3/internal/proxy"
	"net"
)

type Forward struct {
	TargetAddress string
	DialFunc      proxy.DialFunc
	ListenAddress *net.TCPAddr
}

func (localForward *Forward) SetListen(_ proxy.ListenFunc) {
}

func (localForward *Forward) SetListenAddress(address net.Addr) {
	localForward.ListenAddress = address.(*net.TCPAddr)
}

func NewForward(targetAddress string) proxy.Protocol {
	return &Forward{TargetAddress: targetAddress}
}

func (localForward *Forward) SetAuthenticationMethod(_ proxy.AuthenticationMethod) error {
	return errors.New("This kind of proxy doesn't support authentication methods")
}

func (localForward *Forward) SetDial(dialFunc proxy.DialFunc) {
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
