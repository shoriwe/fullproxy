package config

import (
	"fmt"
	"net"

	"github.com/shoriwe/fullproxy/v3/proxies"
	"github.com/things-go/go-socks5"
)

const (
	Socks5 = "socks5"
)

type Proxy struct {
	Listener    *Listener         `yaml:"listener"`
	Reverse     *Reverse          `yaml:"reverse"`
	Credentials map[string]string `yaml:"credentials"`
	Protocol    string            `yaml:"protocol"`
}

func (p *Proxy) Create() (proxies.Proxy, error) {
	var (
		l    net.Listener
		lErr error
		dial func(string, string) (net.Conn, error)
	)
	switch {
	case p.Listener != nil:
		l, lErr = p.Listener.Listen()
		if lErr != nil {
			return nil, lErr
		}
		dial = net.Dial
	case p.Reverse != nil:
		m, lErr := p.Reverse.Master()
		if lErr != nil {
			return nil, lErr
		}
		l = m
		dial = m.Dial
	default:
		return nil, fmt.Errorf("no listener provided")
	}
	var proxy proxies.Proxy
	switch p.Protocol {
	case Socks5:
		var authMethods []socks5.Authenticator
		if p.Credentials != nil {
			authMethods = append(authMethods, socks5.UserPassAuthenticator{
				Credentials: socks5.StaticCredentials(p.Credentials),
			})
		}
		proxy = &proxies.Socks5{
			Listener:    l,
			Dial:        dial,
			AuthMethods: authMethods,
		}
	default:
		l.Close()
		return nil, fmt.Errorf("unknown protocol: %s", p.Protocol)
	}
	return proxy, nil
}
