package proxies

import (
	"net"
	"net/http"
	"net/url"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/shoriwe/fullproxy/v4/reverse"
	httputils "github.com/shoriwe/fullproxy/v4/utils/http"
	"github.com/shoriwe/fullproxy/v4/utils/network"
	"github.com/stretchr/testify/assert"
)

func TestHTTP_Addr(t *testing.T) {
	listener := network.ListenAny()
	defer listener.Close()
	h := HTTP{
		Listener: listener,
		Dial:     net.Dial,
	}
	defer h.Close()
	assert.NotNil(t, h.Addr())
}

func TestHTTP_Serve(t *testing.T) {
	doProxyTest := func(tt *testing.T, proxy, targetService net.Listener) {
		// Setup expect
		httputils.NewMux(targetService)
		proxyUrl, _ := url.Parse("http://" + proxy.Addr().String())
		expect := httpexpect.WithConfig(
			httpexpect.Config{
				BaseURL:  "http://" + targetService.Addr().String(),
				Reporter: httpexpect.NewAssertReporter(t),
				Client: &http.Client{
					Transport: &http.Transport{
						Proxy: http.ProxyURL(proxyUrl),
					},
				},
			},
		)
		// HTTP
		h := HTTP{
			Listener: proxy,
			Dial:     net.Dial,
		}
		defer h.Close()
		go h.Serve()
		expect.GET(httputils.EchoRoute).Expect().Status(http.StatusOK).Body().Contains(httputils.EchoMsg)
	}
	t.Run("Basic", func(tt *testing.T) {
		// - Proxy
		listener := network.ListenAny()
		defer listener.Close()
		// - Service
		service := network.ListenAny()
		defer service.Close()
		// Run test
		doProxyTest(tt, listener, service)
	})
	t.Run("Reverse", func(tt *testing.T) {
		// - Proxy
		data := network.ListenAny()
		defer data.Close()
		control := network.ListenAny()
		defer control.Close()
		master := network.Dial(control.Addr().String())
		defer master.Close()
		// - Service
		service := network.ListenAny()
		defer service.Close()
		// Slave
		s := &reverse.Slave{
			Master: master,
		}
		defer s.Close()
		go s.Serve()
		// Master
		m := &reverse.Master{
			Data:    data,
			Control: control,
		}
		defer m.Close()
		// Run test
		doProxyTest(tt, data, service)
	})

}
