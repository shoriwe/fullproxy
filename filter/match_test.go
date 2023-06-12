package filter

import (
	"net"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch_Match(t *testing.T) {
	t.Run("127.0.0.1:443", func(tt *testing.T) {
		m := Match{
			Host:      regexp.MustCompile("127.0.0.1"),
			Port:      443,
			PortRange: [2]int{-1, -1},
		}
		addr1 := net.TCPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: 443,
		}
		addr2 := net.TCPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: 20,
		}
		assert.True(tt, m.Match(&addr1))
		assert.False(tt, m.Match(&addr2))
	})
	t.Run("127.0.0.1:80-443", func(tt *testing.T) {
		m := Match{
			Host:      regexp.MustCompile("127.0.0.1"),
			Port:      -1,
			PortRange: [2]int{80, 443},
		}
		addr1 := net.TCPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: 443,
		}
		addr2 := net.TCPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: 1000,
		}
		assert.True(tt, m.Match(&addr1))
		assert.False(tt, m.Match(&addr2))
	})
}
