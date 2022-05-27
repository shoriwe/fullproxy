package reverse

import (
	"crypto/tls"
)

type Host struct {
	WebsocketReadBufferSize  int
	WebsocketWriteBufferSize int
	Scheme                   string
	URI                      string
	Network                  string
	Address                  string
	TLSConfig                *tls.Config
}
