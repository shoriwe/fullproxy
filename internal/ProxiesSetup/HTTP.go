package ProxiesSetup

import (
	"github.com/shoriwe/FullProxy/pkg/ConnectionHandlers"
	"github.com/shoriwe/FullProxy/pkg/ConnectionHandlers/Slave"
	"github.com/shoriwe/FullProxy/pkg/Proxies/HTTP"
	"gopkg.in/elazarl/goproxy.v1"
)

func SetupHTTP(host *string, port *string, slave *bool, tls *bool, username *[]byte, password *[]byte) {
	proxy := new(HTTP.HTTP)
	if len(*username) > 0 && len(*password) > 0 {
		proxy.SetAuthenticationMethod(BasicAuthentication(username, password))
	}
	proxyController := new(goproxy.ProxyHttpServer)
	proxy.ProxyController = proxyController
	if *slave {
		ConnectionHandlers.Bind(host, port, proxy)
	} else {
		Slave.GeneralSlave(host, port, proxy)
	}
}
