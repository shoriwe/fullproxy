package reverse

import (
	"crypto/tls"
)

type Host struct {
	WebsocketReadBufferSize  int         `yaml:"websocket-read-buffer-size"`
	WebsocketWriteBufferSize int         `yaml:"websocket-write-buffer-size"`
	Scheme                   string      `yaml:"scheme"`
	URI                      string      `yaml:"uri"`
	Network                  string      `yaml:"network"`
	Address                  string      `yaml:"address"`
	TLSConfig                *tls.Config `yaml:"tls-config"`
}
