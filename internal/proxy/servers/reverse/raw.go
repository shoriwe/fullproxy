package reverse

import (
	"crypto/tls"
	"github.com/shoriwe/fullproxy/v3/internal/common"
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers"
	"io"
	"net"
)

type Raw struct {
	CurrentHost                      int
	Hosts                            []*Host
	Dial                             servers.DialFunc
	IncomingSniffer, OutgoingSniffer io.Writer
}

func (r *Raw) SetSniffers(incoming, outgoing io.Writer) {
	r.IncomingSniffer = incoming
	r.OutgoingSniffer = outgoing
}

func (r *Raw) nextHost() *Host {
	if r.CurrentHost >= len(r.Hosts) {
		r.CurrentHost = 0
	}
	index := r.CurrentHost
	r.CurrentHost++
	return r.Hosts[index]
}

func (r *Raw) SetAuthenticationMethod(_ servers.AuthenticationMethod) {
}

func (r *Raw) SetListen(_ servers.ListenFunc) {
}

func (r *Raw) SetListenAddress(_ net.Addr) {
}

func (r *Raw) SetDial(dialFunc servers.DialFunc) {
	r.Dial = dialFunc
}

func (r *Raw) Handle(conn net.Conn) error {
	host := r.nextHost()
	targetConnection, connectionError := r.Dial(host.Network, host.Address)
	if connectionError != nil {
		return connectionError
	}
	if host.TLSConfig != nil {
		// TODO: Do something to control tls config
		targetConnection = tls.Client(targetConnection, host.TLSConfig)
	}
	return common.ForwardTraffic(
		conn,
		&common.Sniffer{
			WriteSniffer: r.OutgoingSniffer,
			ReadSniffer:  r.IncomingSniffer,
			Connection:   targetConnection,
		},
	)
}

func NewRaw(targets []*Host) servers.Protocol {
	return &Raw{
		CurrentHost: 0,
		Hosts:       targets,
	}
}
