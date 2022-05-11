package reverse

import (
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers"
	"net"
	"net/http"
	"regexp"
)

type (
	Target struct {
		Host regexp.Regexp
		Path regexp.Regexp
		Pool *Pool
	}
	HTTP struct {
		Cache   map[string]map[string]*Pool
		Targets []Target
	}
)

func (H *HTTP) SetListenAddress(_ net.Addr) {
}

func (H *HTTP) SetAuthenticationMethod(_ servers.AuthenticationMethod) {
}

func (H *HTTP) SetListen(_ servers.ListenFunc) {
}

func (H *HTTP) Handle(_ net.Conn) error {
	panic("this should not be used")
}

func (H *HTTP) SetDial(dialFunc servers.DialFunc) {
	for _, value := range H.Targets {
		value.Pool.SetDial(dialFunc)
	}
}

func (H *HTTP) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if cached, hostFound := H.Cache[request.Host]; hostFound {
		if pool, poolFound := cached[request.RequestURI]; poolFound {
			pool.ServeHTTP(writer, request)
			return
		}
	}
	for _, target := range H.Targets {
		if target.Host.MatchString(request.Host) && target.Path.MatchString(request.RequestURI) {
			target.Pool.ServeHTTP(writer, request)
			return
		}
	}
	writer.WriteHeader(http.StatusNotFound)
}

func NewHTTP(targets []Target) servers.HTTPHandler {
	return &HTTP{
		Cache:   map[string]map[string]*Pool{},
		Targets: targets,
	}
}
