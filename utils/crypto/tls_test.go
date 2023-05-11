package crypto

import (
	"crypto/tls"
	"testing"

	"github.com/shoriwe/fullproxy/v3/utils/network"
	"github.com/stretchr/testify/assert"
)

const testMessage = "TESTING"

func TestDefaultTLSConfig(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		signal := make(chan struct{}, 1)
		l := network.ListenAny()
		defer l.Close()
		l = tls.NewListener(l, DefaultTLSConfig())
		go func() {
			conn := network.Dial(l.Addr().String())
			conn = tls.Client(conn, DefaultTLSConfig())
			_, wErr := conn.Write([]byte(testMessage))
			assert.Nil(tt, wErr)
			<-signal
		}()
		conn, aErr := l.Accept()
		assert.Nil(tt, aErr)
		buffer := make([]byte, len(testMessage))
		_, rErr := conn.Read(buffer)
		assert.Nil(tt, rErr)
		signal <- struct{}{}
	})
}
