package SOCKS5

import (
	"encoding/binary"
	"fmt"
)

func clean(hostType uint8, rawHost []byte, rawPort []byte) string {
	var host string
	switch hostType {
	case IPv4:
		host = fmt.Sprintf("%d.%d.%d.%d", rawHost[0], rawHost[1], rawHost[2], rawHost[3])
	case DomainName:
		host = string(rawHost)
	case IPv6:
		host = fmt.Sprintf("[%d:%d:%d:%d:%d:%d:%d:%d:%d:%d:%d:%d:%d:%d:%d:%d]",
			rawHost[0], rawHost[1], rawHost[2], rawHost[3],
			rawHost[4], rawHost[5], rawHost[6], rawHost[7],
			rawHost[8], rawHost[9], rawHost[10], rawHost[11],
			rawHost[12], rawHost[13], rawHost[14], rawHost[15],
		)
	}
	port := binary.BigEndian.Uint16(rawPort)
	fmt.Println(fmt.Sprintf("%s:%d", host, port))
	return fmt.Sprintf("%s:%d", host, port)
}
