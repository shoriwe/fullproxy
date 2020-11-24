package Types

import "net"

type AuthenticationMethod func(username []byte, password []byte) bool

type LoggingMethod func(args ...interface{})

type IOFilter func(address net.IP) bool
