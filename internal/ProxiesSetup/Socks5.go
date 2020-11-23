package ProxiesSetup

import (
	"github.com/shoriwe/FullProxy/internal/SetupControllers"
	"github.com/shoriwe/FullProxy/pkg/Proxies/SOCKS5"
	"log"
)

func SetupSocks5(host *string, port *string, slave *bool, username []byte, password []byte) {
	proxy := new(SOCKS5.Socks5)
	if len(username) > 0 && len(password) > 0 {
		proxy.WantedAuthMethod = SOCKS5.UsernamePassword
		proxy.SetAuthenticationMethod(BasicAuthentication(username, password))
	} else {
		proxy.WantedAuthMethod = SOCKS5.NoAuthRequired
		proxy.SetAuthenticationMethod(NoAuthentication)
	}
	proxy.SetLoggingMethod(log.Print)
	if *slave {
		SetupControllers.GeneralSlave(host, port, proxy)
	} else {
		SetupControllers.Bind(host, port, proxy)
	}
}
