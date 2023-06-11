package compose

import (
	"fmt"
	"net"

	"github.com/shoriwe/fullproxy/v3/proxies"
	"github.com/shoriwe/fullproxy/v3/utils/network"
)

const (
	ProxyForward = "forward"
	ProxyHTTP    = "http"
	ProxySocks5  = "socks5"
)

type Proxy struct {
	Type     string   `yaml:"type" json:"type"`
	Listener Network  `yaml:"listener" json:"listener"`
	Dialer   *Network `yaml:"dialer,omitempty" json:"dialer,omitempty"`
	Network  *string  `yaml:"network,omitempty" json:"network,omitempty"`
	Address  *string  `yaml:"address,omitempty" json:"address,omitempty"`
	proxy    proxies.Proxy
}

func (p *Proxy) getDialFunc() (network.DialFunc, error) {
	switch p.Dialer {
	case nil:
		return p.Listener.DialFunc()
	default:
		return p.Dialer.DialFunc()
	}
}

func (p *Proxy) setupForward() (proxy proxies.Proxy, err error) {
	var (
		l        net.Listener
		dialFunc network.DialFunc
	)
	l, err = p.Listener.Listen()
	if err == nil {
		dialFunc, err = p.getDialFunc()
		if err == nil {
			network.CloseOnError(&err, l)
			proxy = &proxies.Forward{
				Network:  *p.Network,
				Address:  *p.Address,
				Listener: l,
				Dial:     dialFunc,
			}
		}
	}
	return proxy, err
}

func (p *Proxy) setupHTTP() (proxy proxies.Proxy, err error) {
	var (
		l        net.Listener
		dialFunc network.DialFunc
	)
	l, err = p.Listener.Listen()
	if err == nil {
		dialFunc, err = p.getDialFunc()
		if err == nil {
			network.CloseOnError(&err, l)
			proxy = &proxies.HTTP{
				Listener: l,
				Dial:     dialFunc,
			}
		}
	}
	return proxy, err
}

func (p *Proxy) setupSocks5() (proxy proxies.Proxy, err error) {
	var (
		l        net.Listener
		dialFunc network.DialFunc
	)
	l, err = p.Listener.Listen()
	if err == nil {
		dialFunc, err = p.getDialFunc()
		if err == nil {
			network.CloseOnError(&err, l)
			proxy = &proxies.Socks5{
				Listener: l,
				Dial:     dialFunc,
			}
		}
	}
	return proxy, err
}

func (p *Proxy) setupProxy() (proxies.Proxy, error) {
	switch p.Type {
	case ProxyForward:
		return p.setupForward()
	case ProxyHTTP:
		return p.setupHTTP()
	case ProxySocks5:
		return p.setupSocks5()
	default:
		return nil, fmt.Errorf("unknown proxy type: %s", p.Type)
	}
}

func (p *Proxy) Serve() (err error) {
	p.proxy, err = p.setupProxy()
	if err == nil {
		err = p.proxy.Serve()
	}
	return err
}
