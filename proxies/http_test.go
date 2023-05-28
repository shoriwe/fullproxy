package proxies

import (
	"net"
	"net/http"
	"net/url"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/shoriwe/fullproxy/v3/reverse"
	"github.com/shoriwe/fullproxy/v3/utils/network"
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
	testMessage := "TEST"
	echoRoute := "/echo"
	setupHTTP := func(tt *testing.T, proxy, service net.Listener) *httpexpect.Expect {
		mux := http.NewServeMux()
		mux.HandleFunc(echoRoute, func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(testMessage))
		})
		go http.Serve(service, mux)
		proxyUrl, _ := url.Parse("http://" + proxy.Addr().String())
		return httpexpect.WithConfig(
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
	}
	t.Run("Basic", func(tt *testing.T) {
		service := network.ListenAny()
		defer service.Close()
		listener := network.ListenAny()
		defer listener.Close()
		expect := setupHTTP(tt, listener, service)
		h := HTTP{
			Listener: listener,
			Dial:     net.Dial,
		}
		defer h.Close()
		go h.Serve()
		expect.GET(echoRoute).Expect().Status(http.StatusOK).Body().Contains(testMessage)
	})
	t.Run("Reverse", func(tt *testing.T) {
		service := network.ListenAny()
		defer service.Close()
		data := network.ListenAny()
		defer data.Close()
		control := network.ListenAny()
		defer control.Close()
		master := network.Dial(control.Addr().String())
		defer master.Close()
		controlCh := make(chan struct{}, 1)
		defer close(controlCh)
		go func() {
			s := &reverse.Slave{
				Master: master,
			}
			defer s.Close()
			go s.Serve()
			<-controlCh
		}()
		m := &reverse.Master{
			Data:    data,
			Control: control,
		}
		expect := setupHTTP(tt, data, service)
		h := HTTP{
			Listener: m,
			Dial:     net.Dial,
		}
		defer h.Close()
		go h.Serve()
		expect.GET(echoRoute).Expect().Status(http.StatusOK).Body().Contains(testMessage)
		controlCh <- struct{}{}
	})

}
