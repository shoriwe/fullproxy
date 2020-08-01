package Sockets

import "strconv"


func GetIPv6(rawIPv6Address []byte) string{
	var ipv6 string = ""
	for index := 0; index < 16; index++{
		ipv6 += strconv.FormatInt(int64(rawIPv6Address[index]), 16)
		if (index != 15) && (index != 0) && (index % 2 != 0){
			ipv6 += ":"
		}
	}
	return "[" + ipv6 + "]"
}
