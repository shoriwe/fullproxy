package reverse

import (
	"context"
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers"
	"io"
	"net"
	"net/http"
)

type Pool struct {
	currentTarget int
	Targets       []string
	Client        *http.Client
}

func (pool *Pool) nextTarget() int {
	if pool.currentTarget >= len(pool.Targets) {
		pool.currentTarget = 0
	}
	result := pool.currentTarget
	pool.currentTarget++
	return result
}

func (pool *Pool) SetDial(dialFunc servers.DialFunc) {
	pool.Client = &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network string, _ string) (net.Conn, error) {
				return dialFunc(network, pool.Targets[pool.nextTarget()])
			},
		},
	}
}

func (pool *Pool) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	response, requestError := pool.Client.Do(request)
	if requestError != nil {
		// TODO: Do something with the error
		panic("DO SOMETHING?")
	}
	for key, values := range response.Header {
		// TODO: Remove headers
		writer.Header().Set(key, values[len(values)-1])
	}
	// TODO: Inject headers
	_, copyError := io.Copy(writer, response.Body)
	if copyError != nil {
		// TODO: Do something with the error
		panic("DO SOMETHING")
	}
	writer.WriteHeader(response.StatusCode)
}

func NewPool(targets []string) *Pool {
	return &Pool{
		currentTarget: 0,
		Targets:       targets,
	}
}
