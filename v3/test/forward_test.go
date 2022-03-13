package test

import (
	"github.com/shoriwe/FullProxy/v3/internal/proxy/port-forward"
	"testing"
)

// Test Bind

func TestPortForwardBindRequest(t *testing.T) {
	h := StartIPv4HTTPServer(t)
	defer h.Close()
	p := NewBindPipe(port_forward.NewForward(networkType, httpAddress, nil), nil)
	defer p.Close()
	result := GetRequestRaw("http://" + proxyAddress + "/big.txt")
	if result != Success {
		t.Fatal(proxyAddress, result)
	}
}

func TestPortForwardBindIPv6Request(t *testing.T) {
	h := StartIPv6HTTPServer(t)
	defer h.Close()
	p := NewBindPipe(port_forward.NewForward(networkType, httpIPv6Address, nil), nil)
	defer p.Close()
	result := GetRequestRaw("http://" + proxyAddress + "/big.txt")
	if result != Success {
		t.Fatal(proxyAddress, result)
	}
}

// Test master/slave

func TestPortForwardMasterSlaveRequest(t *testing.T) {
	h := StartIPv4HTTPServer(t)
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

func TestPortForwardMasterSlaveIPv6Request(t *testing.T) {
	h := StartIPv6HTTPServer(t)
	defer h.Close()
	a, b := NewMasterSlave(
		nil,
		port_forward.NewForward(networkType, httpIPv6Address, nil))
	defer func() {
		a.Close()
		b.Close()
	}()
	result := GetRequestRaw("http://" + proxyAddress + "/big.txt")
	if result != Success {
		t.Fatal(proxyAddress, result)
	}
}
