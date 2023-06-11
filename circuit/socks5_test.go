package circuit

import (
	"net"
	"net/http"
	"testing"

	"github.com/shoriwe/fullproxy/v4/proxies"
	"github.com/shoriwe/fullproxy/v4/utils/network"
	"github.com/stretchr/testify/assert"
)

func TestSocks5_Next(t *testing.T) {
	listener := network.ListenAny()
	defer listener.Close()
	socks := proxies.Socks5{
		Listener: listener,
		Dial:     net.Dial,
	}
	defer socks.Close()
	go socks.Serve()
	t.Run("Basic", func(tt *testing.T) {
		s := &Socks5{
			Network: listener.Addr().Network(),
			Address: listener.Addr().String(),
		}
		closeFunc, dial, err := s.Next(net.Dial)
		assert.Nil(tt, err)
		defer closeFunc()
		expect := newExpect(tt, "http://localhost:8000", dial)
		expect.GET("/").Expect().Status(http.StatusOK).Body().Contains("ECHO")
	})
}

func Test_newSocks5(t *testing.T) {
	assert.NotNil(t, newSocks5())
}
