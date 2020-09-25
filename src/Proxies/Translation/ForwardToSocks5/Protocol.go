package ForwardToSocks5

import (
	"fmt"
	"github.com/shoriwe/FullProxy/src/BindServer"
	"github.com/shoriwe/FullProxy/src/ConnectionStructures"
	"github.com/shoriwe/FullProxy/src/Proxies/Basic"
	"golang.org/x/net/proxy"
	"log"
	"net"
)

func CreateTranslationSession(conn net.Conn, connReader  ConnectionStructures.SocketReader, connWriter ConnectionStructures.SocketWriter, args...interface{}) {
	targetConnection, connectionError := args[0].(proxy.Dialer).Dial("tcp", args[1].(string))
	if connectionError == nil {
		targetConnectionReader, targetConnectionWriter := ConnectionStructures.CreateSocketConnectionReaderWriter(targetConnection)
		Basic.Proxy(conn, targetConnection, connReader, connWriter, targetConnectionReader, targetConnectionWriter)
	} else {
		log.Print(connectionError)
		_ = conn.Close()
	}

}

func StartForwardToSocks5Translation(bindAddress *string, bindPort *string, socks5Address *string, socks5Port, username *string, password *string, targetAddress *string, targetPort *string) {
	if len(*targetAddress) > 0 && len(*targetPort) > 0 {
		var connectionDialer proxy.Dialer
		var connectionError error
		if len(*username) > 0 && len(*password) > 0 {
			auth := new(proxy.Auth)
			auth.User = *username
			auth.Password = *password
			connectionDialer, connectionError = proxy.SOCKS5("tcp", *socks5Address+":"+*socks5Port, auth, proxy.Direct)
		} else {
			connectionDialer, connectionError = proxy.SOCKS5("tcp", *socks5Address+":"+*socks5Port, nil, proxy.Direct)
		}

		if connectionError == nil {
			log.Print("Starting translation (Forward --> SOCKS5)")
			log.Printf("Targeting: %s:%s", *targetAddress, *targetPort)
			log.Printf("With the SOCKS5 tunnel: %s:%s", *socks5Address, *socks5Port)
			BindServer.Bind(bindAddress, bindPort, CreateTranslationSession, connectionDialer, *targetAddress+":"+*targetPort)
		} else {
			log.Print(connectionError)
		}
	} else {
		fmt.Println("You must specify a target and port address")
	}
}