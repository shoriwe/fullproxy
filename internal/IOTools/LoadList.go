package IOTools

import (
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Templates/Types"
	"net"
)

func ParseAddresses(fileLines []string) []net.IP {
	var addresses []net.IP
	for _, wAddress := range fileLines {
		addresses = append(addresses, net.ParseIP(wAddress))
	}
	return addresses
}

func CreateWhitelistFilter(addresses []net.IP) Types.IOFilter {
	return func(address net.IP) bool {
		for _, wAddress := range addresses {
			if wAddress.Equal(address) {
				return true
			}
		}
		return false
	}
}

func CreateBlacklistFilter(addresses []net.IP) Types.IOFilter {
	return func(address net.IP) bool {
		for _, wAddress := range addresses {
			if wAddress.Equal(address) {
				return false
			}
		}
		return true
	}
}

func LoadList(l1 string, l2 string) (Types.IOFilter, error) {
	if len(l1) > 0 {
		if len(l2) > 0 {
			return nil, errors.New("Both parameters are set")
		}
		lines, readingError := ReadLines(l1)
		if readingError != nil {
			return nil, readingError
		}
		return CreateWhitelistFilter(ParseAddresses(lines)), nil
	}
	if len(l2) > 0 {
		lines, readingError := ReadLines(l2)
		if readingError != nil {
			return nil, readingError
		}
		return CreateBlacklistFilter(ParseAddresses(lines)), nil
	}
	return nil, errors.New("No files provided")
}
