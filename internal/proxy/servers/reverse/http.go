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
		Header        http.Header
		Path          string
		CurrentTarget int
		Targets       []string
	}
	HTTP struct {
		Targets map[string]*Target
		Client  *http.Client
	}
)

func (target *Target) nextTarget() string {
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
	H.Client.Transport = &http.Transport{
		DialContext: func(_ context.Context, network, address string) (net.Conn, error) {
			return dialFunc(network, address)
		},
	}
}

func createRequest(received *http.Request, reference *Target) (*http.Request, error) {
	u, parseError := url.Parse(reference.nextTarget())
	if parseError != nil {
		return nil, parseError
	}
	u.Path = path.Join(u.Path, strings.Replace(received.RequestURI, reference.Path, "/", 1))
	request, newRequestError := http.NewRequest(received.Method, u.String(), received.Body)
	if newRequestError != nil {
		return nil, newRequestError
	}
	request.Header = reference.Header.Clone()
	for key, value := range received.Header {
		request.Header.Set(key, value[0])
	}
	return request, nil
}

func (H *HTTP) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	target, found := H.Targets[request.Host]
	if found {
		if strings.Index(request.RequestURI, target.Path) == 0 {
			targetRequest, requestCreationError := createRequest(request, target)
			if requestCreationError != nil {
				return
			}
			response, requestError := H.Client.Do(targetRequest)
			if requestError != nil {
				return
			}
			for key, value := range response.Header {
				writer.Header().Set(key, value[0])
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
		Client:  &http.Client{},
	}
}
