package ProxiesSetup

import (
	"github.com/shoriwe/FullProxy/internal/ControllersSetup"
	"github.com/shoriwe/FullProxy/pkg/Proxies/Translation/ForwardToSocks5"
	"golang.org/x/net/proxy"
	"log"
)

func SetupForwardSocks5(
	bindHost *string, bindPort *string,
	socks5Host *string, socks5Port *string,
	username *string, password *string,
	targetHost *string, targetPort *string) {
	proxyProtocol := new(ForwardToSocks5.ForwardToSocks5)
	proxyProtocol.TargetHost = *targetHost
	proxyProtocol.TargetPort = *targetPort
	proxyProtocol.SetLoggingMethod(log.Print)
	proxyAuth := new(proxy.Auth)
	if len(*username) > 0 && len(*password) > 0 {
		proxyAuth.User = *username
		proxyAuth.Password = *password
	} else {
		proxyAuth = nil
	}
	proxyDialer, dialerCreationError := proxy.SOCKS5("tcp", *socks5Host+":"+*socks5Port, proxyAuth, proxy.Direct)
	if dialerCreationError != nil {
		log.Fatal(dialerCreationError)
	}
	proxyProtocol.Socks5Dialer = proxyDialer
	ControllersSetup.Bind(bindHost, bindPort, proxyProtocol)
}
