package proxies

type Proxy interface {
	Serve() error
}
