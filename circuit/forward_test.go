package circuit

import (
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForward_Next(t *testing.T) {
	t.Run("Basic", func(tt *testing.T) {
		s := &Forward{
			Network: "tcp",
			Address: "localhost:8000",
		}
		closeFunc, dial, err := s.Next(net.Dial)
		assert.Nil(tt, err)
		defer closeFunc()
		expect := newExpect(tt, "http://localhost:8000", dial)
		expect.GET("/").Expect().Status(http.StatusOK).Body().Contains("ECHO")
	})
}

func Test_newForward(t *testing.T) {
	assert.NotNil(t, newForward())
}
