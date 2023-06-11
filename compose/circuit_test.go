package compose

import (
	"testing"

	"net/http"

	utilshttp "github.com/shoriwe/fullproxy/v3/utils/http"

	"github.com/gavv/httpexpect/v2"
	"github.com/shoriwe/fullproxy/v3/utils/network"
	"github.com/stretchr/testify/assert"
)

func TestCircuit_setupCircuit(t *testing.T) {
	t.Run("Succeed", func(tt *testing.T) {
		service := network.ListenAny()
		defer service.Close()
		utilshttp.NewMux(service)
		c := Circuit{
			Network: service.Addr().Network(),
			Address: service.Addr().String(),
			Listener: Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
			Knots: []Knot{
				{
					Type:    KnotForward,
					Network: service.Addr().Network(),
					Address: service.Addr().String(),
				},
			},
		}
		*c.Listener.Network = "tcp"
		*c.Listener.Address = "localhost:0"
		assert.Nil(tt, c.setupCircuit())
	})
	t.Run("Invalid Knot", func(tt *testing.T) {
		service := network.ListenAny()
		defer service.Close()
		utilshttp.NewMux(service)
		c := Circuit{
			Network: service.Addr().Network(),
			Address: service.Addr().String(),
			Listener: Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
			Knots: []Knot{
				{
					Type:    KnotSSH,
					Network: service.Addr().Network(),
					Address: service.Addr().String(),
				},
			},
		}
		*c.Listener.Network = "tcp"
		*c.Listener.Address = "localhost:0"
		assert.NotNil(tt, c.setupCircuit())
	})
}

func TestCircuit_handle(t *testing.T) {
	t.Run("Succeed", func(tt *testing.T) {
		service := network.ListenAny()
		defer service.Close()
		utilshttp.NewMux(service)
		c := Circuit{
			Network: service.Addr().Network(),
			Address: service.Addr().String(),
			Listener: Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
			Knots: []Knot{
				{
					Type:    KnotForward,
					Network: service.Addr().Network(),
					Address: service.Addr().String(),
				},
			},
		}
		*c.Listener.Network = "tcp"
		*c.Listener.Address = "localhost:0"
		assert.Nil(tt, c.setupCircuit())
		ll, err := c.Listener.Listen()
		assert.Nil(tt, err)
		defer ll.Close()
		go func() {
			conn, err := ll.Accept()
			assert.Nil(tt, err)
			c.handle(conn)
		}()
		expect := httpexpect.Default(tt, "http://"+ll.Addr().String())
		expect.GET(utilshttp.EchoRoute).Expect().Status(http.StatusOK)
	})
}

func TestCircuit_serve(t *testing.T) {
	t.Run("Succeed", func(tt *testing.T) {
		service := network.ListenAny()
		defer service.Close()
		utilshttp.NewMux(service)
		c := Circuit{
			Network: service.Addr().Network(),
			Address: service.Addr().String(),
			Listener: Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
			Knots: []Knot{
				{
					Type:    KnotForward,
					Network: service.Addr().Network(),
					Address: service.Addr().String(),
				},
			},
		}
		*c.Listener.Network = "tcp"
		*c.Listener.Address = "localhost:0"
		assert.Nil(tt, c.setupCircuit())
		ll, err := c.Listener.Listen()
		assert.Nil(tt, err)
		defer ll.Close()
		go c.serve(ll)
		expect := httpexpect.Default(tt, "http://"+ll.Addr().String())
		expect.GET(utilshttp.EchoRoute).Expect().Status(http.StatusOK)
	})
}

func TestCircuit_Serve(t *testing.T) {
	t.Run("Succeed", func(tt *testing.T) {
		service := network.ListenAny()
		defer service.Close()
		utilshttp.NewMux(service)
		c := Circuit{
			Network: service.Addr().Network(),
			Address: service.Addr().String(),
			Listener: Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
			Knots: []Knot{
				{
					Type:    KnotForward,
					Network: service.Addr().Network(),
					Address: service.Addr().String(),
				},
			},
		}
		*c.Listener.Network = "tcp"
		*c.Listener.Address = "localhost:0"
		checkChn := make(chan struct{}, 1)
		defer close(checkChn)
		go func() {
			go c.Serve()
			for c.listener == nil {
			}
			checkChn <- struct{}{}
		}()
		<-checkChn
		defer c.listener.Close()
	})
	t.Run("Invalid Setup", func(tt *testing.T) {
		service := network.ListenAny()
		defer service.Close()
		utilshttp.NewMux(service)
		c := Circuit{
			Network: service.Addr().Network(),
			Address: service.Addr().String(),
			Listener: Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
			Knots: []Knot{
				{
					Type:    KnotSSH,
					Network: service.Addr().Network(),
					Address: service.Addr().String(),
				},
			},
		}
		*c.Listener.Network = "tcp"
		*c.Listener.Address = "localhost:0"
		assert.NotNil(tt, c.Serve())
	})
}
