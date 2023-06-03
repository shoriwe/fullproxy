package compose

import (
	"crypto/tls"
	"fmt"
	"net"

	crypto2 "github.com/shoriwe/fullproxy/v3/utils/crypto"
)

const (
	CryptoTLS = "tls"
)

type Crypto struct {
	Mode               string  `yaml:"mode,omitempty" json:"mode,omitempty"`
	SelfSigned         bool    `yaml:"selfSigned" json:"selfSigned"`
	InsecureSkipVerify *bool   `yaml:"insecureSkipVerify,omitempty" json:"insecureSkipVerify,omitempty"`
	Cert               *string `yaml:"cert,omitempty" json:"cert,omitempty"`
	Key                *string `yaml:"key,omitempty" json:"key,omitempty"`
}

func (c *Crypto) tlsListener(l net.Listener) (net.Listener, error) {
	if c.SelfSigned {
		return tls.NewListener(l, crypto2.DefaultTLSConfig()), nil
	}
	if c.Cert == nil {
		return nil, fmt.Errorf("no cert provided")
	}
	if c.Key == nil {
		return nil, fmt.Errorf("no key provided")
	}
	cert, cErr := tls.LoadX509KeyPair(*c.Cert, *c.Key)
	if cErr != nil {
		return nil, cErr
	}
	config := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: c.InsecureSkipVerify != nil && *c.InsecureSkipVerify,
	}
	ll := tls.NewListener(l, config)
	return ll, nil
}

func (c *Crypto) WrapListener(l net.Listener) (net.Listener, error) {
	switch c.Mode {
	case CryptoTLS:
		return c.tlsListener(l)
	default:
		return nil, fmt.Errorf("unknown crypto channel %s", c.Mode)
	}
}

func (c *Crypto) tlsConn(conn net.Conn) (net.Conn, error) {
	if c.SelfSigned {
		return tls.Client(conn, crypto2.DefaultTLSConfig()), nil
	}
	config := &tls.Config{
		InsecureSkipVerify: c.InsecureSkipVerify != nil && *c.InsecureSkipVerify,
	}
	if c.Cert != nil && c.Key != nil {
		cert, cErr := tls.LoadX509KeyPair(*c.Cert, *c.Key)
		if cErr != nil {
			return nil, cErr
		}
		config.Certificates = append(config.Certificates, cert)
	}
	cc := tls.Client(conn, config)
	return cc, nil
}

func (c *Crypto) WrapConn(conn net.Conn) (net.Conn, error) {
	switch c.Mode {
	case CryptoTLS:
		return c.tlsConn(conn)
	default:
		return nil, fmt.Errorf("unknown crypto channel %s", c.Mode)
	}
}
