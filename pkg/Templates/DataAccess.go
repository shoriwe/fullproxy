package Templates

import (
	"github.com/shoriwe/FullProxy/pkg/Templates/Types"
	"net"
	"time"
)

func LogData(loggingMethod Types.LoggingMethod, arguments ...interface{}) {
	if loggingMethod != nil {
		loggingMethod(arguments...)
	}
}

func FilterInbound(filter Types.IOFilter, address net.Addr) bool {
	if filter != nil {
		return filter(address)
	}
	return true
}

func FilterOutbound(filter Types.IOFilter, address net.Addr) bool {
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
