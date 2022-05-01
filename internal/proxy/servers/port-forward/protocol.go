package port_forward

import (
	"errors"
	"github.com/shoriwe/fullproxy/v3/internal/common"
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers"
	"net"
)

type Forward struct {
	TargetAddress string
	DialFunc      servers.DialFunc
	ListenAddress *net.TCPAddr
}

func (localForward *Forward) SetListen(_ servers.ListenFunc) {
}

func (localForward *Forward) SetListenAddress(address net.Addr) {
	localForward.ListenAddress = address.(*net.TCPAddr)
}

func NewForward(targetAddress string) servers.Protocol {
	return &Forward{TargetAddress: targetAddress}
}

func (localForward *Forward) SetAuthenticationMethod(_ servers.AuthenticationMethod) error {
	return errors.New("this kind of proxy doesn't support authentication methods")
}

func (localForward *Forward) SetDial(dialFunc servers.DialFunc) {
	localForward.DialFunc = dialFunc
}

func (localForward *Forward) Handle(clientConnection net.Conn) error {
	defer clientConnection.Close()
	targetConnection, connectionError := localForward.DialFunc("tcp", localForward.TargetAddress)
	if connectionError != nil {
		return connectionError
	}
	return common.ForwardTraffic(clientConnection, targetConnection)
}
