package ProxiesSetup

import (
	"github.com/google/shlex"
	"github.com/shoriwe/FullProxy/internal/Authentication"
	"github.com/shoriwe/FullProxy/internal/IOTools"
	"github.com/shoriwe/FullProxy/internal/PipesSetup"
	"github.com/shoriwe/FullProxy/pkg/Proxies/SOCKS5"
	"log"
	"time"
)

func SetupSocks5(
	host *string, port *string,
	slave *bool, username []byte,
	password []byte, tries *int,
	timeout *time.Duration, inboundLists [2]string,
	outboundLists [2]string, commandAuth *string,
	databaseAuth *string) {
	proxy := new(SOCKS5.Socks5)
	if len(username) > 0 && len(password) > 0 {
		proxy.WantedAuthMethod = SOCKS5.UsernamePassword
		_ = proxy.SetAuthenticationMethod(Authentication.SingleUser(username, password))
	} else if len(*commandAuth) > 0 {
		splitProcess, splitError := shlex.Split(*commandAuth)
		if splitError != nil {
			log.Fatal(splitError)
		}
		proxy.WantedAuthMethod = SOCKS5.UsernamePassword
		_ = proxy.SetAuthenticationMethod(Authentication.CommandAuth(splitProcess[0], splitProcess[1:]))
	} else if len(*databaseAuth) > 0 {
		proxy.WantedAuthMethod = SOCKS5.UsernamePassword
		_ = proxy.SetAuthenticationMethod(Authentication.SQLite3Authentication(*databaseAuth))
	} else {
		proxy.WantedAuthMethod = SOCKS5.NoAuthRequired
		_ = proxy.SetAuthenticationMethod(Authentication.NoAuthentication)
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
