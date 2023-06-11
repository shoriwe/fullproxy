package compose

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"testing"

	"github.com/gavv/httpexpect/v2"
	httputils "github.com/shoriwe/fullproxy/v3/utils/http"
	"github.com/shoriwe/fullproxy/v3/utils/network"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/proxy"
)

func TestProxy_getDialFunc(t *testing.T) {
	t.Run("no dialer", func(tt *testing.T) {
		p := Proxy{
			Listener: Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
		}
		*p.Listener.Network = "tcp"
		*p.Listener.Address = "localhost:0"
		dialFunc, err := p.getDialFunc()
		assert.Nil(tt, err)
		assert.NotNil(tt, dialFunc)
	})
	t.Run("with dialer", func(tt *testing.T) {
		p := Proxy{
			Listener: Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
			Dialer: &Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
		}
		*p.Listener.Network = "tcp"
		*p.Listener.Address = "localhost:0"
		*p.Dialer.Network = "tcp"
		*p.Dialer.Network = "localhost:0"
		dialFunc, err := p.getDialFunc()
		assert.Nil(tt, err)
		assert.NotNil(tt, dialFunc)
	})
}

func TestProxy_setupForward(t *testing.T) {
	t.Run("Succeed", func(tt *testing.T) {
		service := network.ListenAny()
		defer service.Close()
		httputils.NewMux(service)
		p := Proxy{
			Type:    ProxyForward,
			Network: new(string),
			Address: new(string),
			Listener: Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
		}
		*p.Network = service.Addr().Network()
		*p.Address = service.Addr().String()
		*p.Listener.Network = "tcp"
		*p.Listener.Address = "localhost:0"
		px, err := p.setupForward()
		assert.Nil(tt, err)
		defer px.Close()
		go px.Serve()
		expect := httpexpect.Default(tt, "http://"+px.Addr().String())
		expect.GET(httputils.EchoRoute).Expect().Status(http.StatusOK).Body().Contains(httputils.EchoMsg)
	})
}

func TestProxy_setupHTTP(t *testing.T) {
	t.Run("Succeed", func(tt *testing.T) {
		service := network.ListenAny()
		defer service.Close()
		httputils.NewMux(service)
		p := Proxy{
			Type:    ProxyHTTP,
			Network: new(string),
			Address: new(string),
			Listener: Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
		}
		*p.Network = service.Addr().Network()
		*p.Address = service.Addr().String()
		*p.Listener.Network = "tcp"
		*p.Listener.Address = "localhost:0"
		px, err := p.setupHTTP()
		assert.Nil(tt, err)
		defer px.Close()
		go px.Serve()
		proxyUrl, _ := url.Parse("http://" + px.Addr().String())
		expect := httpexpect.WithConfig(
			httpexpect.Config{
				BaseURL:  "http://" + service.Addr().String(),
				Reporter: httpexpect.NewAssertReporter(t),
				Client: &http.Client{
					Transport: &http.Transport{
						Proxy: http.ProxyURL(proxyUrl),
					},
				},
			},
		)
		expect.GET(httputils.EchoRoute).Expect().Status(http.StatusOK).Body().Contains(httputils.EchoMsg)
	})
}

func TestProxy_setupSocks5(t *testing.T) {
	t.Run("Succeed", func(tt *testing.T) {
		service := network.ListenAny()
		defer service.Close()
		httputils.NewMux(service)
		p := Proxy{
			Type:    ProxySocks5,
			Network: new(string),
			Address: new(string),
			Listener: Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
		}
		*p.Network = service.Addr().Network()
		*p.Address = service.Addr().String()
		*p.Listener.Network = "tcp"
		*p.Listener.Address = "localhost:0"
		px, err := p.setupSocks5()
		assert.Nil(tt, err)
		defer px.Close()
		go px.Serve()
		pxDialer, err := proxy.SOCKS5(px.Addr().Network(), px.Addr().String(), nil, nil)
		assert.Nil(tt, err)
		expect := httpexpect.WithConfig(
			httpexpect.Config{
				BaseURL:  "http://" + service.Addr().String(),
				Reporter: httpexpect.NewAssertReporter(t),
				Client: &http.Client{
					Transport: &http.Transport{
						DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
							return pxDialer.Dial(network, addr)
						},
					},
				},
			},
		)
		expect.GET(httputils.EchoRoute).Expect().Status(http.StatusOK).Body().Contains(httputils.EchoMsg)
	})
}

func TestProxy_setupProxy(t *testing.T) {
	t.Run(ProxyForward, func(tt *testing.T) {
		service := network.ListenAny()
		defer service.Close()
		httputils.NewMux(service)
		p := Proxy{
			Type:    ProxyForward,
			Network: new(string),
			Address: new(string),
			Listener: Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
		}
		*p.Network = service.Addr().Network()
		*p.Address = service.Addr().String()
		*p.Listener.Network = "tcp"
		*p.Listener.Address = "localhost:0"
		px, err := p.setupProxy()
		assert.Nil(tt, err)
		defer px.Close()
		go px.Serve()
		expect := httpexpect.Default(tt, "http://"+px.Addr().String())
		expect.GET(httputils.EchoRoute).Expect().Status(http.StatusOK).Body().Contains(httputils.EchoMsg)
	})
	t.Run(ProxyHTTP, func(tt *testing.T) {
		service := network.ListenAny()
		defer service.Close()
		httputils.NewMux(service)
		p := Proxy{
			Type:    ProxyHTTP,
			Network: new(string),
			Address: new(string),
			Listener: Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
		}
		*p.Network = service.Addr().Network()
		*p.Address = service.Addr().String()
		*p.Listener.Network = "tcp"
		*p.Listener.Address = "localhost:0"
		px, err := p.setupProxy()
		assert.Nil(tt, err)
		defer px.Close()
		go px.Serve()
		proxyUrl, _ := url.Parse("http://" + px.Addr().String())
		expect := httpexpect.WithConfig(
			httpexpect.Config{
				BaseURL:  "http://" + service.Addr().String(),
				Reporter: httpexpect.NewAssertReporter(t),
				Client: &http.Client{
					Transport: &http.Transport{
						Proxy: http.ProxyURL(proxyUrl),
					},
				},
			},
		)
		expect.GET(httputils.EchoRoute).Expect().Status(http.StatusOK).Body().Contains(httputils.EchoMsg)
	})
	t.Run(ProxySocks5, func(tt *testing.T) {
		service := network.ListenAny()
		defer service.Close()
		httputils.NewMux(service)
		p := Proxy{
			Type:    ProxySocks5,
			Network: new(string),
			Address: new(string),
			Listener: Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
		}
		*p.Network = service.Addr().Network()
		*p.Address = service.Addr().String()
		*p.Listener.Network = "tcp"
		*p.Listener.Address = "localhost:0"
		px, err := p.setupProxy()
		assert.Nil(tt, err)
		defer px.Close()
		go px.Serve()
		pxDialer, err := proxy.SOCKS5(px.Addr().Network(), px.Addr().String(), nil, nil)
		assert.Nil(tt, err)
		expect := httpexpect.WithConfig(
			httpexpect.Config{
				BaseURL:  "http://" + service.Addr().String(),
				Reporter: httpexpect.NewAssertReporter(t),
				Client: &http.Client{
					Transport: &http.Transport{
						DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
							return pxDialer.Dial(network, addr)
						},
					},
				},
			},
		)
		expect.GET(httputils.EchoRoute).Expect().Status(http.StatusOK).Body().Contains(httputils.EchoMsg)
	})
	t.Run("UNKNOWN", func(tt *testing.T) {
		service := network.ListenAny()
		defer service.Close()
		httputils.NewMux(service)
		p := Proxy{
			Type:    "UNKNOWN",
			Network: new(string),
			Address: new(string),
			Listener: Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
		}
		*p.Network = service.Addr().Network()
		*p.Address = service.Addr().String()
		*p.Listener.Network = "tcp"
		*p.Listener.Address = "localhost:0"
		_, err := p.setupProxy()
		assert.NotNil(tt, err)
	})
}

func TestProxy_Serve(t *testing.T) {
	t.Run("Succeed", func(tt *testing.T) {
		service := network.ListenAny()
		defer service.Close()
		httputils.NewMux(service)
		p := Proxy{
			Type:    ProxyForward,
			Network: new(string),
			Address: new(string),
			Listener: Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
		}
		*p.Network = service.Addr().Network()
		*p.Address = service.Addr().String()
		*p.Listener.Network = "tcp"
		*p.Listener.Address = "localhost:0"
		checkCh := make(chan struct{}, 1)
		go func() {
			go p.Serve()
			for p.proxy == nil {
			}
			checkCh <- struct{}{}
		}()
		<-checkCh
		expect := httpexpect.Default(tt, "http://"+p.proxy.Addr().String())
		expect.GET(httputils.EchoRoute).Expect().Status(http.StatusOK).Body().Contains(httputils.EchoMsg)
	})
}
