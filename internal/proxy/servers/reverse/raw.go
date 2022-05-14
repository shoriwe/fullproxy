package reverse

import (
	"github.com/shoriwe/fullproxy/v3/internal/common"
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers"
	"net"
)

type Raw struct {
	currentTarget int
	Targets       []*Host
	Dial          servers.DialFunc
}

func (r *Raw) nextTarget() *Host {
	if r.currentTarget >= len(r.Targets) {
		r.currentTarget = 0
	}
	index := r.currentTarget
	r.currentTarget++
	return r.Targets[index]
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
	host := r.nextTarget()
	targetConnection, connectionError := r.Dial(host.Network, host.Address)
	if connectionError != nil {
		return connectionError
	}
	return common.ForwardTraffic(conn, targetConnection)
}

func NewRaw(targets []*Host) servers.Protocol {
	return &Raw{
		currentTarget: 0,
		Targets:       targets,
	}
}
