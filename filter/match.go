package filter

import (
	"net"
	"regexp"
	"strconv"
)

type Match struct {
	Host      *regexp.Regexp
	Port      int
	PortRange [2]int
}

func (m *Match) Match(addr net.Addr) bool {
	host, portS, _ := net.SplitHostPort(addr.String())
	port, _ := strconv.Atoi(portS)
	result := (m.Host == nil || (m.Host != nil && m.Host.MatchString(host))) &&
		(m.Port <= -1 || (m.Port > -1 && (m.Port == port))) &&
		(m.PortRange[0] <= -1 || (m.PortRange[0] > -1 && (m.PortRange[0] <= port))) &&
		(m.PortRange[1] <= -1 || (m.PortRange[1] > -1 && (port <= m.PortRange[1])))
	return result
}
