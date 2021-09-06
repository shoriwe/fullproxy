package Types

import "net"

type AuthenticationMethod func(username []byte, password []byte) (bool, error)

type LoggingMethod func(args ...interface{})

type IOFilter func(address net.IP) bool
