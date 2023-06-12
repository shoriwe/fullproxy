package filter

import (
	"fmt"
	"net"
	"strconv"

	"github.com/shoriwe/fullproxy/v4/utils/network"
)

type DialFunc struct {
	network.DialFunc
	Whitelist []Match
	Blacklist []Match
}

func (df *DialFunc) GetWhitelist() []Match { return df.Whitelist }
func (df *DialFunc) GetBlacklist() []Match { return df.Blacklist }

func (df *DialFunc) Dial(network, address string) (conn net.Conn, err error) {
	var host, portS string
	host, portS, err = net.SplitHostPort(address)
	ip := net.ParseIP(host)
	if err == nil {
		port, _ := strconv.Atoi(portS)
		remote := &net.TCPAddr{
			IP:   ip,
			Port: port,
		}
		if VerifyConn(df, remote) {
			conn, err = df.DialFunc(network, address)
		} else {
			err = fmt.Errorf("connection to %s denied by rule", remote)
		}
	}
	return conn, err
}
