package filter

import (
	"net"
)

type Filter interface {
	GetBlacklist() []Match
	GetWhitelist() []Match
}

func VerifyConn(filter Filter, remote net.Addr) bool {
	for _, black := range filter.GetBlacklist() {
		if black.Match(remote) {
			return false
		}
	}
	if len(filter.GetWhitelist()) == 0 {
		return true
	}
	for _, white := range filter.GetWhitelist() {
		if white.Match(remote) {
			return true
		}
	}
	return false
}
