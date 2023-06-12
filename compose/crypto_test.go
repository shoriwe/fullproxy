package compose

import (
	"os"
	"testing"

	"github.com/shoriwe/fullproxy/v4/utils/crypto"
	"github.com/shoriwe/fullproxy/v4/utils/network"
	"github.com/stretchr/testify/assert"
)

func TestCrypto_tlsListener(t *testing.T) {
	t.Run("SelfSigned", func(tt *testing.T) {
		c := Crypto{
			Mode:       CryptoTLS,
			SelfSigned: true,
		}
		l := network.ListenAny()
		defer l.Close()
		_, err := c.tlsListener(l)
		assert.Nil(tt, err)
	})
	t.Run("No Cert", func(tt *testing.T) {
		c := Crypto{
			Mode: CryptoTLS,
		}
		l := network.ListenAny()
		defer l.Close()
		_, err := c.tlsListener(l)
		assert.NotNil(tt, err)
	})
	t.Run("No Key", func(tt *testing.T) {
		c := Crypto{
			Mode: CryptoTLS,
			Cert: new(string),
		}
		l := network.ListenAny()
		defer l.Close()
		_, err := c.tlsListener(l)
		assert.NotNil(tt, err)
	})
	t.Run("Valid Certs", func(tt *testing.T) {
		c := Crypto{
			Mode: CryptoTLS,
			Cert: new(string),
			Key:  new(string),
		}
		cert, key := crypto.TempCertKey()
		defer os.Remove(cert)
		defer os.Remove(key)
		*c.Cert = cert
		*c.Key = key
		l := network.ListenAny()
		defer l.Close()
		_, err := c.tlsListener(l)
		assert.Nil(tt, err)
	})
	t.Run("Invalid Certs", func(tt *testing.T) {
		c := Crypto{
			Mode: CryptoTLS,
			Cert: new(string),
			Key:  new(string),
		}
		cert, key := crypto.TempCertKey()
		defer os.Remove(cert)
		defer os.Remove(key)
		*c.Cert = ""
		*c.Key = key
		l := network.ListenAny()
		defer l.Close()
		_, err := c.tlsListener(l)
		assert.NotNil(tt, err)
	})
}

func TestCrypto_tlsConn(t *testing.T) {
	t.Run("SelfSigned", func(tt *testing.T) {
		c := Crypto{
			Mode:       CryptoTLS,
			SelfSigned: true,
		}
		l := network.ListenAny()
		defer l.Close()
		go l.Accept()
		conn := network.Dial(l.Addr().String())
		defer conn.Close()
		_, err := c.tlsConn(conn)
		assert.Nil(tt, err)
	})
	t.Run("InsecureSkipVerify", func(tt *testing.T) {
		c := Crypto{
			Mode:               CryptoTLS,
			InsecureSkipVerify: new(bool),
		}
		*c.InsecureSkipVerify = true
		l := network.ListenAny()
		defer l.Close()
		go l.Accept()
		conn := network.Dial(l.Addr().String())
		defer conn.Close()
		_, err := c.tlsConn(conn)
		assert.Nil(tt, err)
	})
	t.Run("Cert", func(tt *testing.T) {
		// Listener
		l := network.ListenAny()
		defer l.Close()
		cc := Crypto{
			Mode:       CryptoTLS,
			SelfSigned: true,
		}
		_, err := cc.tlsListener(l)
		assert.Nil(tt, err)
		go l.Accept()
		// Client
		c := Crypto{
			Mode: CryptoTLS,
			Cert: new(string),
			Key:  new(string),
		}
		cert, key := crypto.TempCertKey()
		defer os.Remove(cert)
		defer os.Remove(key)
		*c.Cert = cert
		*c.Key = key
		conn := network.Dial(l.Addr().String())
		defer conn.Close()
		_, err = c.tlsConn(conn)
		assert.Nil(tt, err)
	})
	t.Run("Invalid Cert", func(tt *testing.T) {
		// Listener
		l := network.ListenAny()
		defer l.Close()
		cc := Crypto{
			Mode:       CryptoTLS,
			SelfSigned: true,
		}
		_, err := cc.tlsListener(l)
		assert.Nil(tt, err)
		go l.Accept()
		// Client
		c := Crypto{
			Mode: CryptoTLS,
			Cert: new(string),
			Key:  new(string),
		}
		cert, key := crypto.TempCertKey()
		defer os.Remove(cert)
		defer os.Remove(key)
		*c.Cert = ""
		*c.Key = key
		conn := network.Dial(l.Addr().String())
		defer conn.Close()
		_, err = c.tlsConn(conn)
		assert.NotNil(tt, err)
	})
}

func TestWrapListener(t *testing.T) {
	t.Run(CryptoTLS, func(tt *testing.T) {
		c := Crypto{
			Mode:       CryptoTLS,
			SelfSigned: true,
		}
		l := network.ListenAny()
		defer l.Close()
		_, err := c.WrapListener(l)
		assert.Nil(tt, err)
	})
	t.Run("Invalid", func(tt *testing.T) {
		c := Crypto{
			Mode:       "INVALID",
			SelfSigned: true,
		}
		l := network.ListenAny()
		defer l.Close()
		_, err := c.WrapListener(l)
		assert.NotNil(tt, err)
	})
}

func TestWrapConn(t *testing.T) {
	t.Run(CryptoTLS, func(tt *testing.T) {
		c := Crypto{
			Mode:       CryptoTLS,
			SelfSigned: true,
		}
		l := network.ListenAny()
		defer l.Close()
		go l.Accept()
		conn := network.Dial(l.Addr().String())
		defer conn.Close()
		_, err := c.WrapConn(conn)
		assert.Nil(tt, err)
	})
	t.Run("Invalid", func(tt *testing.T) {
		c := Crypto{
			Mode:       "INVALID",
			SelfSigned: true,
		}
		l := network.ListenAny()
		defer l.Close()
		go l.Accept()
		conn := network.Dial(l.Addr().String())
		defer conn.Close()
		_, err := c.WrapConn(conn)
		assert.NotNil(tt, err)
	})
}
