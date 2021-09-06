package ProxiesSetup

import (
	"time"
)

func SetupForwardSocks5(
	bindHost string, bindPort string,
	socks5Host string, socks5Port string,
	username string, password string,
	targetHost string, targetPort string,
	tries *int, timeout *time.Duration,
	inboundLists [2]string) {

}
