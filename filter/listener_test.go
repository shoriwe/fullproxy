package filter

import (
	"regexp"
	"testing"

	"github.com/shoriwe/fullproxy/v4/utils/network"
	"github.com/stretchr/testify/assert"
)

func TestListener_Accept(t *testing.T) {
	t.Run("No Whitelist", func(tt *testing.T) {
		l := network.ListenAny()
		defer l.Close()
		filter := Listener{
			Listener:  l,
			Whitelist: []Match{},
		}
		go func() {
			conn := network.Dial(l.Addr().String())
			defer conn.Close()
		}()
		conn, err := filter.Accept()
		assert.Nil(tt, err)
		defer conn.Close()
	})
	t.Run("With Whitelist", func(tt *testing.T) {
		l := network.ListenAny()
		defer l.Close()
		filter := Listener{
			Listener: l,
			Whitelist: []Match{
				{
					Host:      regexp.MustCompile("127.0.0.1"),
					Port:      -1,
					PortRange: [2]int{-1, -1},
				},
			},
		}
		go func() {
			conn := network.Dial(l.Addr().String())
			defer conn.Close()
		}()
		conn, err := filter.Accept()
		assert.Nil(tt, err)
		defer conn.Close()
	})
	t.Run("With Whitelist Deny", func(tt *testing.T) {
		l := network.ListenAny()
		defer l.Close()
		filter := Listener{
			Listener: l,
			Whitelist: []Match{
				{
					Host:      regexp.MustCompile("9.9.9.9"),
					Port:      -1,
					PortRange: [2]int{-1, -1},
				},
			},
		}
		go func() {
			conn := network.Dial(l.Addr().String())
			defer conn.Close()
		}()
		_, err := filter.Accept()
		assert.NotNil(tt, err)
	})
	t.Run("Blacklist", func(tt *testing.T) {
		l := network.ListenAny()
		defer l.Close()
		filter := Listener{
			Listener: l,
			Blacklist: []Match{
				{
					Host:      regexp.MustCompile("127.0.0.1"),
					Port:      -1,
					PortRange: [2]int{-1, -1},
				},
			},
		}
		go func() {
			conn := network.Dial(l.Addr().String())
			defer conn.Close()
		}()
		_, err := filter.Accept()
		assert.NotNil(tt, err)
	})
}
