package global

import (
	"net"
)

// Taken from https://play.golang.org/p/dAoV99_7iPY
func ParseIP(s string) net.IP {
	ip, _, err := net.SplitHostPort(s)
	if err == nil {
		return net.ParseIP(ip)
	}
	ip2 := net.ParseIP(s)
	if ip2 == nil {
		return nil
	}
	return ip2
}

func LogData(loggingMethod LoggingMethod, arguments ...interface{}) {
	if loggingMethod != nil {
		loggingMethod(arguments...)
	}
}

func FilterInbound(filter IOFilter, host string) bool {
	if filter != nil {
		return filter(host)
	}
	return true
}

func FilterOutbound(filter IOFilter, host string) bool {
	if filter != nil {
		return filter(host)
	}
	return true
}
