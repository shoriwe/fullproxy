package ProxiesSetup

import (
	"github.com/shoriwe/FullProxy/internal/IOTools"
	"github.com/shoriwe/FullProxy/internal/PipesSetup"
	"github.com/shoriwe/FullProxy/internal/Templates"
	"github.com/shoriwe/FullProxy/pkg/Proxies/HTTP"
	"gopkg.in/elazarl/goproxy.v1"
	"log"
)

func SetupHTTP(
	host *string, port *string,
	slave *bool, tls *bool, username []byte,
	password []byte, inboundLists [2]string,
	outboundLists [2]string) {
	if *tls {
		log.Fatal("TLS is not implemented yet")
	}
	proxy := new(HTTP.HTTP)
	proxyController := goproxy.NewProxyHttpServer()
	proxy.ProxyController = proxyController
	_ = proxy.SetLoggingMethod(log.Print)
	if len(username) > 0 && len(password) > 0 {
		_ = proxy.SetAuthenticationMethod(Templates.BasicAuthentication(username, password))
	}
	filter, loadingError := IOTools.LoadList(outboundLists[0], outboundLists[1])
	if loadingError == nil {
		_ = proxy.SetOutboundFilter(filter)
	}
	if *slave {
		PipesSetup.GeneralSlave(host, port, proxy)
	} else {
		filter, loadingError := IOTools.LoadList(inboundLists[0], inboundLists[1])
		if loadingError == nil {
			_ = proxy.SetInboundFilter(filter)
		}
		PipesSetup.Bind(host, port, proxy)
	}
}
