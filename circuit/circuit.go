package circuit

import (
	"net"

	"github.com/shoriwe/fullproxy/v3/utils/network"
)

/*
Circuit are reduced to forwarding a end of chain service over multiple proxies

For example:

## The user knows these proxy servers

- Server 1: SOCKS-5
- Server 2: Forward
- Server 3: HTTP
- Server 4: Forward
- Server 5: SSH
- Server 6: SOCKS-5

- Service

## The user then design a circuit to connect to a service that is only accesible by "Server 6"

Client -> Server 1 -> Server 2 -> Server 3 -> Server 4 -> Server 5 -> Server 6 -> Service

# This way the user end Pivoting over the entire proxy chain to finally arrive it's requested service

## Server accesibility

# To allow Randomness and unpredictibility in the proxy chain outgoing connections, it is also possible to set Server accesibility

# Using the previus chain, but adding accesible annotation

- Server 1: SOCKS-5
-- accessible:
--- Server 3
--- Server 5
--- Server 6
- Server 2: Forward
- Server 3: HTTP
-- accessible:
--- Server 1
--- Server 5
--- Server 6
- Server 4: Forward
- Server 5: SSH
-- accessible:
--- Server 1
--- Server 3
--- Server 6
- Server 6: SOCKS-5
-- accessible:
--- Server 1
--- Server 5
- Service
-- accessible:
--- Server 1
--- Server 3
--- Server 5
--- Server 6

# Now the chain will be reconstructed every accepted connection, ensuring randomness in the knots

# First connection
Client -> Server 1 -> Server 2 -> Server 3 -> Server 4 -> Server 5 -> Server 6 -> Service

# Second connection

Client -> Server 3 -> Server 1 -> Server 4 -> Server 6 -> Server 2 -> Server 5 -> Service
*/
type Circuit struct {
	Chain []Knot
}

func (c *Circuit) Dial(n, addr string) (net.Conn, error) {
	var (
		result = &Conn{
			CloseFunctions: make([]network.CloseFunc, 0, len(c.Chain)),
		}
		closeFunc network.CloseFunc
		dial      network.DialFunc = net.Dial
		err       error
	)
	defer network.CloseOnError(&err, result)
	for _, knot := range c.Chain {
		closeFunc, dial, err = knot.Next(dial)
		if err != nil {
			return nil, err
		}
		result.CloseFunctions = append(result.CloseFunctions, closeFunc)
	}
	result.Conn, err = dial(n, addr)
	return result, err
}
