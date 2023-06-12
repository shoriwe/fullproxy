package crypto

import (
	"crypto/tls"
	"os"
	"testing"

	"github.com/shoriwe/fullproxy/v4/utils/network"
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
			defer conn.Close()
			conn = tls.Client(conn, DefaultTLSConfig())
			_, wErr := conn.Write([]byte(testMessage))
			assert.Nil(tt, wErr)
			<-signal
		}()
		conn, aErr := l.Accept()
		assert.Nil(tt, aErr)
		defer conn.Close()
		buffer := make([]byte, len(testMessage))
		_, rErr := conn.Read(buffer)
		assert.Nil(tt, rErr)
		signal <- struct{}{}
	})
}

func TestTempCertKey(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		cert, key := TempCertKey()
		defer os.Remove(cert)
		defer os.Remove(key)
		_, err := tls.LoadX509KeyPair(cert, key)
		assert.Nil(tt, err)
	})
}
