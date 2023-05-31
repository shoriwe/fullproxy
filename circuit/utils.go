package circuit

import (
	"context"
	"net"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/shoriwe/fullproxy/v3/utils/network"
)

func newExpect(t *testing.T, baseUrl string, dial network.DialFunc) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  baseUrl,
		Reporter: httpexpect.NewAssertReporter(t),
		Client: &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return dial(network, addr)
				},
			},
		},
	})
}
