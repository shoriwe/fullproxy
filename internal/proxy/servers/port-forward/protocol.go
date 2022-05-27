package port_forward

import (
	"crypto/tls"
	"github.com/shoriwe/fullproxy/v3/internal/common"
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers"
	"io"
	"net"
)

type Forward struct {
	TargetNetwork                    string
	TargetAddress                    string
	DialFunc                         servers.DialFunc
	ListenAddress                    *net.TCPAddr
	IncomingSniffer, OutgoingSniffer io.Writer
	DialTLSConfig                    *tls.Config
}

func (f *Forward) SetSniffers(incoming, outgoing io.Writer) {
	f.IncomingSniffer = incoming
	f.OutgoingSniffer = outgoing
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
	if f.DialTLSConfig != nil {
		targetConnection = tls.Client(targetConnection, f.DialTLSConfig)
	}
	return common.ForwardTraffic(
		clientConnection,
		&common.Sniffer{
			WriteSniffer: f.OutgoingSniffer,
			ReadSniffer:  f.IncomingSniffer,
			Connection:   targetConnection,
		},
	)
}

func NewForward(targetNetwork, targetAddress string, dialTLSConfig *tls.Config) servers.Protocol {
	return &Forward{
		TargetNetwork: targetNetwork,
		TargetAddress: targetAddress,
		DialTLSConfig: dialTLSConfig,
	}
}
