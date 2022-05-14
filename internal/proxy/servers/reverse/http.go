package reverse

import (
	"context"
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers"
	"io"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type (
	Target struct {
		RequestHeader  http.Header
		ResponseHeader http.Header
		Path           string
		CurrentTarget  int
		Targets        []*Host
	}
	HTTP struct {
		Targets map[string]*Target
		Dial    servers.DialFunc
	}
)

func (target *Target) nextTarget() *Host {
	if target.CurrentTarget >= len(target.Targets) {
		target.CurrentTarget = 0
	}
	index := target.CurrentTarget
	target.CurrentTarget++
	return target.Targets[index]
}

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
	H.Dial = dialFunc
}

func createRequest(received *http.Request, reference *Target) (*http.Request, *Host, error) {
	host := reference.nextTarget()
	u, parseError := url.Parse(host.Url)
	if parseError != nil {
		return nil, nil, parseError
	}
	u.Path = path.Join(u.Path, strings.Replace(received.RequestURI, reference.Path, "/", 1))
	request, newRequestError := http.NewRequest(received.Method, u.String(), received.Body)
	if newRequestError != nil {
		return nil, nil, newRequestError
	}
	request.Header = reference.RequestHeader.Clone()
	for key, values := range received.Header {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}
	return request, host, nil
}

func (H *HTTP) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	target, found := H.Targets[request.Host]
	if found {
		if strings.Index(request.RequestURI, target.Path) == 0 {
			targetRequest, host, requestCreationError := createRequest(request, target)
			if requestCreationError != nil {
				return
			}
			client := &http.Client{
				Transport: &http.Transport{
					DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
						return H.Dial(host.Network, host.Address)
					},
				},
			}
			response, requestError := client.Do(targetRequest)
			if requestError != nil {
				return
			}
			for key, values := range response.Header {
				for _, value := range values {
					writer.Header().Add(key, value)
				}
			}
			for key, values := range target.ResponseHeader {
				for _, value := range values {
					writer.Header().Add(key, value)
				}
			}
			writer.WriteHeader(response.StatusCode)
			_, copyError := io.Copy(writer, response.Body)
			if copyError != nil {
				return
			}
			return
		}
	}
	writer.WriteHeader(http.StatusNotFound)
}

func NewHTTP(targets map[string]*Target) servers.HTTPHandler {
	return &HTTP{
		Targets: targets,
	}
}
