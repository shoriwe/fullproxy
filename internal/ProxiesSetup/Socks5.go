package ProxiesSetup

import (
	"github.com/shoriwe/FullProxy/pkg/ConnectionHandlers"
	"github.com/shoriwe/FullProxy/pkg/ConnectionHandlers/Slave"
	"github.com/shoriwe/FullProxy/pkg/Proxies/SOCKS5"
)

func SetupSocks5(host *string, port *string, slave *bool, username *[]byte, password *[]byte) {
	proxy := new(SOCKS5.Socks5)
	if len(*username) > 0 && len(*password) > 0 {
		proxy.WantedAuthMethod = SOCKS5.BasicNegotiation
		proxy.SetAuthenticationMethod(BasicAuthentication(username, password))
	} else {
		proxy.WantedAuthMethod = SOCKS5.NoAuthRequired
		proxy.SetAuthenticationMethod(NoAuthentication)
	}
	if *slave {
		Slave.GeneralSlave(host, port, proxy)
	} else {
		ConnectionHandlers.Bind(host, port, proxy)
	}
}
