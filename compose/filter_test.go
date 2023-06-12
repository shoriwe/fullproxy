package compose

import (
	"net"
	"regexp"
	"testing"

	"github.com/shoriwe/fullproxy/v4/filter"
	"github.com/shoriwe/fullproxy/v4/utils/network"
	"github.com/stretchr/testify/assert"
)

func TestMatch_Compile(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		m := Match{
			Host: new(string),
			Port: new(int),
			Range: &PortRange{
				From: new(int),
				To:   new(int),
			},
		}
		*m.Host = "127.0.0.1"
		*m.Port = 80
		*m.Range.From = 0
		*m.Range.To = 8000
		match, err := m.Compile()
		assert.Nil(tt, err)
		expect := filter.Match{
			Host:      regexp.MustCompile("127.0.0.1"),
			Port:      80,
			PortRange: [2]int{0, 8000},
		}
		assert.Equal(tt, expect, match)
	})
}

func TestFilter_Listener(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		m := Match{
			Host: new(string),
			Port: new(int),
			Range: &PortRange{
				From: new(int),
				To:   new(int),
			},
		}
		*m.Host = "127.0.0.1"
		*m.Port = 80
		*m.Range.From = 0
		*m.Range.To = 8000
		f := Filter{
			Whitelist: []Match{m},
			Blacklist: []Match{m},
		}
		l := network.ListenAny()
		defer l.Close()
		ll, err := f.Listener(l)
		assert.Nil(tt, err)
		defer ll.Close()
	})
	t.Run("Invalid Whitelist", func(tt *testing.T) {
		m := Match{
			Host: new(string),
			Port: new(int),
			Range: &PortRange{
				From: new(int),
				To:   new(int),
			},
		}
		*m.Host = ")"
		*m.Port = 80
		*m.Range.From = 0
		*m.Range.To = 8000
		f := Filter{
			Whitelist: []Match{m},
			// Blacklist: []Match{m},
		}
		l := network.ListenAny()
		defer l.Close()
		_, err := f.Listener(l)
		assert.NotNil(tt, err)
	})
	t.Run("Invalid Blacklist", func(tt *testing.T) {
		m := Match{
			Host: new(string),
			Port: new(int),
			Range: &PortRange{
				From: new(int),
				To:   new(int),
			},
		}
		*m.Host = ")"
		*m.Port = 80
		*m.Range.From = 0
		*m.Range.To = 8000
		f := Filter{
			// Whitelist: []Match{m},
			Blacklist: []Match{m},
		}
		l := network.ListenAny()
		defer l.Close()
		_, err := f.Listener(l)
		assert.NotNil(tt, err)
	})
}

func TestFilter_DialFunc(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		m := Match{
			Host: new(string),
			Port: new(int),
			Range: &PortRange{
				From: new(int),
				To:   new(int),
			},
		}
		*m.Host = "127.0.0.1"
		*m.Port = 80
		*m.Range.From = 0
		*m.Range.To = 8000
		f := Filter{
			Whitelist: []Match{m},
			Blacklist: []Match{m},
		}
		df, err := f.DialFunc(net.Dial)
		assert.Nil(tt, err)
		assert.NotNil(tt, df)
	})
	t.Run("Invalid Whitelist", func(tt *testing.T) {
		m := Match{
			Host: new(string),
			Port: new(int),
			Range: &PortRange{
				From: new(int),
				To:   new(int),
			},
		}
		*m.Host = ")"
		*m.Port = 80
		*m.Range.From = 0
		*m.Range.To = 8000
		f := Filter{
			Whitelist: []Match{m},
			// Blacklist: []Match{m},
		}
		_, err := f.DialFunc(net.Dial)
		assert.NotNil(tt, err)
	})
	t.Run("Invalid Blacklist", func(tt *testing.T) {
		m := Match{
			Host: new(string),
			Port: new(int),
			Range: &PortRange{
				From: new(int),
				To:   new(int),
			},
		}
		*m.Host = ")"
		*m.Port = 80
		*m.Range.From = 0
		*m.Range.To = 8000
		f := Filter{
			// Whitelist: []Match{m},
			Blacklist: []Match{m},
		}
		_, err := f.DialFunc(net.Dial)
		assert.NotNil(tt, err)
	})
}
