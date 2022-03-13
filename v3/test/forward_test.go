package test

import (
	"github.com/shoriwe/FullProxy/v3/internal/proxy/port-forward"
	"testing"
)

// Test Bind

func TestPortForwardBindRequest(t *testing.T) {
	h := StartHTTPServer(t)
	defer h.Close()
	p := NewBindPipe(port_forward.NewForward(networkType, httpAddress, nil), nil)
	defer p.Close()
	result := GetRequestRaw("http://" + proxyAddress + "/big.txt")
	if result != Success {
		t.Fatal(proxyAddress, result)
	}
}

// Test master/slave

func TestPortForwardMasterSlaveRequest(t *testing.T) {
	h := StartHTTPServer(t)
	defer h.Close()
	a, b := NewMasterSlave(
		nil,
		port_forward.NewForward(networkType, httpAddress, nil))
	defer func() {
		a.Close()
		b.Close()

	}()
	result := GetRequestRaw("http://" + proxyAddress + "/big.txt")
	if result != Success {
		t.Fatal(proxyAddress, result)
	}
}
