package filter

import (
	"net"
	"regexp"
	"testing"

	"github.com/shoriwe/fullproxy/v4/utils/network"
	"github.com/stretchr/testify/assert"
)

func TestDialFunc_Dial(t *testing.T) {
	t.Run("No Whitelist", func(tt *testing.T) {
		l := network.ListenAny()
		defer l.Close()
		filter := DialFunc{
			DialFunc:  net.Dial,
			Whitelist: []Match{},
		}
		go l.Accept()
		conn, err := filter.Dial(l.Addr().Network(), l.Addr().String())
		assert.Nil(tt, err)
		defer conn.Close()
	})
	t.Run("With Whitelist", func(tt *testing.T) {
		l := network.ListenAny()
		defer l.Close()
		filter := DialFunc{
			DialFunc: net.Dial,
			Whitelist: []Match{
				{
					Host:      regexp.MustCompile("127.0.0.1"),
					Port:      -1,
					PortRange: [2]int{-1, -1},
				},
				{
					Host:      regexp.MustCompile("localhost"),
					Port:      -1,
					PortRange: [2]int{-1, -1},
				},
			},
		}
		go l.Accept()
		conn, err := filter.Dial(l.Addr().Network(), l.Addr().String())
		assert.Nil(tt, err)
		defer conn.Close()
	})
	t.Run("With Whitelist Deny", func(tt *testing.T) {
		l := network.ListenAny()
		defer l.Close()
		filter := DialFunc{
			DialFunc: net.Dial,
			Whitelist: []Match{
				{
					Host:      regexp.MustCompile("9.9.9.9"),
					Port:      -1,
					PortRange: [2]int{-1, -1},
				},
			},
		}
		_, err := filter.Dial(l.Addr().Network(), l.Addr().String())
		assert.NotNil(tt, err)
	})
	t.Run("Blacklist", func(tt *testing.T) {
		l := network.ListenAny()
		defer l.Close()
		filter := DialFunc{
			DialFunc: net.Dial,
			Blacklist: []Match{
				{
					Host:      regexp.MustCompile("127.0.0.1"),
					Port:      -1,
					PortRange: [2]int{-1, -1},
				},
			},
		}
		_, err := filter.Dial(l.Addr().Network(), l.Addr().String())
		assert.NotNil(tt, err)
	})
}
