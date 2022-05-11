package reverse

import (
	"github.com/shoriwe/fullproxy/v3/internal/common"
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers"
	"net"
)

type Raw struct {
	currentTarget int
	Targets       []string
	Dial          servers.DialFunc
}

func (r *Raw) nextTarget() int {
	if r.currentTarget >= len(r.Targets) {
		r.currentTarget = 0
	}
	result := r.currentTarget
	r.currentTarget++
	return result
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
	targetConnection, connectionError := r.Dial("tcp", r.Targets[r.nextTarget()])
	if connectionError != nil {
		return connectionError
	}
	return common.ForwardTraffic(conn, targetConnection)
}

func NewRaw(targets []string) servers.Protocol {
	return &Raw{
		currentTarget: 0,
		Targets:       targets,
	}
}
