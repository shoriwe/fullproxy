package socks5

import (
	"encoding/binary"
	"fmt"
)

func clean(hostType uint8, rawHost []byte, rawPort []byte) (port int, host, hostPort string) {
	switch hostType {
	case IPv4:
		host = fmt.Sprintf("%d.%d.%d.%d", rawHost[0], rawHost[1], rawHost[2], rawHost[3])
	case DomainName:
		host = string(rawHost)
	case IPv6:
		host = fmt.Sprintf("[%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x:%02x%02x]",
			rawHost[0], rawHost[1], rawHost[2], rawHost[3],
			rawHost[4], rawHost[5], rawHost[6], rawHost[7],
			rawHost[8], rawHost[9], rawHost[10], rawHost[11],
			rawHost[12], rawHost[13], rawHost[14], rawHost[15],
		)
	}
	port = int(binary.BigEndian.Uint16(rawPort))
	return port, host, fmt.Sprintf("%s:%d", host, port)
}
