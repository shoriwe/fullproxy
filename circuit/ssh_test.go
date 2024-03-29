package circuit

import (
	"net"
	"net/http"
	"testing"

	"github.com/shoriwe/fullproxy/v4/sshd"
	"github.com/stretchr/testify/assert"
)

func TestSSH_Next(t *testing.T) {
	t.Run("Basic", func(tt *testing.T) {
		s := &SSH{
			Network: "tcp",
			Address: sshd.DefaultAddr,
			Config:  *sshd.DefaultClientConfig(),
		}
		closeFunc, dial, err := s.Next(net.Dial)
		assert.Nil(tt, err)
		defer closeFunc()
		expect := newExpect(tt, "http://echo:80", dial)
		expect.GET("/").Expect().Status(http.StatusOK).Body().Contains("ECHO")
	})
}

func TestSSH_newSSH(t *testing.T) {
	assert.NotNil(t, newSSH())
}
