package ProxiesSetup

import (
	"github.com/shoriwe/FullProxy/internal/ControllersSetup"
	"github.com/shoriwe/FullProxy/pkg/Proxies/HTTP"
	"gopkg.in/elazarl/goproxy.v1"
	"log"
)

func SetupHTTP(host *string, port *string, slave *bool, tls *bool, username []byte, password []byte) {
	proxy := new(HTTP.HTTP)
	proxyController := goproxy.NewProxyHttpServer()
	proxy.ProxyController = proxyController
	proxy.SetLoggingMethod(log.Print)
	if len(username) > 0 && len(password) > 0 {
		proxy.SetAuthenticationMethod(BasicAuthentication(username, password))
	}
	if *slave {
		ControllersSetup.GeneralSlave(host, port, proxy)
	} else {
		ControllersSetup.Bind(host, port, proxy)
	}
}
