package circuit

import (
	"net"
	"net/http"
	"testing"

	"github.com/shoriwe/fullproxy/v4/proxies"
	"github.com/shoriwe/fullproxy/v4/sshd"
	"github.com/shoriwe/fullproxy/v4/utils/network"
	"github.com/stretchr/testify/assert"
)

func TestCircuit_Dial(t *testing.T) {
	// Prepare Socks5
	listener := network.ListenAny()
	defer listener.Close()
	socks := proxies.Socks5{
		Listener: listener,
		Dial:     net.Dial,
	}
	defer socks.Close()
	go socks.Serve()
	// Prepare chain
	chain := []Knot{
		&Socks5{
			Network: listener.Addr().Network(),
			Address: listener.Addr().String(),
		},
		&SSH{
			Network: "tcp",
			Address: "localhost:2222",
			Config:  *sshd.DefaultClientConfig(),
		},
	}
	// Run Tests
	t.Run("Basic", func(tt *testing.T) {
		circuit := &Circuit{
			Chain: chain,
		}
		expect := newExpect(tt, "http://echo:80", circuit.Dial)
		expect.GET("/").Expect().Status(http.StatusOK).Body().Contains("ECHO")
	})
	t.Run("Invalid Knot", func(tt *testing.T) {
		circuit := &Circuit{
			Chain: []Knot{&SSH{
				Network: "tcp",
				Address: "localhost:2222",
			}},
		}
		_, err := circuit.Dial("tcp", "echo:80")
		assert.NotNil(tt, err)
	})
}
