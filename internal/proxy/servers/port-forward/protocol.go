package port_forward

import (
	"github.com/shoriwe/fullproxy/v3/internal/common"
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers"
	"net"
)

type Forward struct {
	TargetNetwork string
	TargetAddress string
	DialFunc      servers.DialFunc
	ListenAddress *net.TCPAddr
}

func (f *Forward) SetListen(_ servers.ListenFunc) {
}

func (f *Forward) SetListenAddress(address net.Addr) {
	f.ListenAddress = address.(*net.TCPAddr)
}

func (f *Forward) SetAuthenticationMethod(_ servers.AuthenticationMethod) {
}

func (f *Forward) SetDial(dialFunc servers.DialFunc) {
	f.DialFunc = dialFunc
}

func (f *Forward) Handle(clientConnection net.Conn) error {
	defer clientConnection.Close()
	targetConnection, connectionError := f.DialFunc(f.TargetNetwork, f.TargetAddress)
	if connectionError != nil {
		return connectionError
	}
	return common.ForwardTraffic(clientConnection, targetConnection)
}

func NewForward(targetNetwork, targetAddress string) servers.Protocol {
	return &Forward{
		TargetNetwork: targetNetwork,
		TargetAddress: targetAddress,
	}
}
