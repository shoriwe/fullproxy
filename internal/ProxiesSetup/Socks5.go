package ProxiesSetup

import (
	"github.com/shoriwe/FullProxy/internal/PipesSetup"
	"github.com/shoriwe/FullProxy/internal/Templates"
	"github.com/shoriwe/FullProxy/pkg/Proxies/SOCKS5"
	"log"
	"time"
)

func SetupSocks5(
	host *string, port *string,
	slave *bool, username []byte,
	password []byte, tries int,
	timeout time.Duration) {
	proxy := new(SOCKS5.Socks5)
	if len(username) > 0 && len(password) > 0 {
		proxy.WantedAuthMethod = SOCKS5.UsernamePassword
		proxy.SetAuthenticationMethod(Templates.BasicAuthentication(username, password))
	} else {
		proxy.WantedAuthMethod = SOCKS5.NoAuthRequired
		proxy.SetAuthenticationMethod(Templates.NoAuthentication)
	}
	proxy.SetTries(tries)
	proxy.SetTimeout(timeout)
	proxy.SetLoggingMethod(log.Print)
	if *slave {
		PipesSetup.GeneralSlave(host, port, proxy)
	} else {
		PipesSetup.Bind(host, port, proxy)
	}
}
