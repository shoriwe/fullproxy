package Types

import "net"

type AuthenticationMethod func(username []byte, password []byte) bool

type LoggingMethod func(args ...interface{})

type InboundFilter func(address net.Addr) bool

type OutboundFilter func(address net.Addr) bool
