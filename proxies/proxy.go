package proxies

import "net"

type Proxy interface {
	Addr() net.Addr
	Close()
	Serve() error
}
