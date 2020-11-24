package ProxiesSetup

import (
	"github.com/shoriwe/FullProxy/internal/IOTools"
	"github.com/shoriwe/FullProxy/internal/PipesSetup"
	"github.com/shoriwe/FullProxy/internal/Templates"
	"github.com/shoriwe/FullProxy/pkg/Proxies/SOCKS5"
	"log"
	"time"
)

func SetupSocks5(
	host *string, port *string,
	slave *bool, username []byte,
	password []byte, tries *int,
	timeout *time.Duration, inboundLists [2]string,
	outboundLists [2]string) {
	proxy := new(SOCKS5.Socks5)
	if len(username) > 0 && len(password) > 0 {
		proxy.WantedAuthMethod = SOCKS5.UsernamePassword
		_ = proxy.SetAuthenticationMethod(Templates.BasicAuthentication(username, password))
	} else {
		proxy.WantedAuthMethod = SOCKS5.NoAuthRequired
		_ = proxy.SetAuthenticationMethod(Templates.NoAuthentication)
	}
	_ = proxy.SetTries(*tries)
	_ = proxy.SetTimeout(*timeout)
	_ = proxy.SetLoggingMethod(log.Print)
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
