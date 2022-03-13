package global

import "net"

type AuthenticationMethod func(username []byte, password []byte) (bool, error)

type LoggingMethod func(args ...interface{})

type IOFilter func(host string) bool

type DialFunc func(network, address string) (net.Conn, error)
