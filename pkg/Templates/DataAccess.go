package Templates

import (
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Templates/Types"
	"net"
	"time"
)
// Taken from https://play.golang.org/p/dAoV99_7iPY
func ParseIP(s string) (net.IP, error) {
	ip, _, err := net.SplitHostPort(s)
	if err == nil {
		return net.ParseIP(ip), nil
	}

	ip2 := net.ParseIP(s)
	if ip2 == nil {
		return nil, errors.New("invalid IP")
	}

	return ip2, nil
}

func LogData(loggingMethod Types.LoggingMethod, arguments ...interface{}) {
	if loggingMethod != nil {
		loggingMethod(arguments...)
	}
}

func FilterInbound(filter Types.IOFilter, address net.IP) bool {
	if filter != nil {
		return filter(address)
	}
	return true
}

func FilterOutbound(filter Types.IOFilter, address net.IP) bool {
	if filter != nil {
		return filter(address)
	}
	return true
}

func GetTries(tries int) int {
	if tries != 0 {
		return tries
	}
	return 5
}

func GetTimeout(timeout time.Duration) time.Duration {
	if timeout != 0 {
		return timeout
	}
	return 10 * time.Second
}
