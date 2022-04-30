package test

import (
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers/port-forward"
	"testing"
)

// Test Bind

func TestPortForwardBindRequest(t *testing.T) {
	h := StartIPv4HTTPServer(t)
	defer h.Close()
	p := NewBindPipe(port_forward.NewForward(httpAddress), nil, nil)
	defer p.Close()
	result := GetRequestRaw("http://" + proxyAddress + "/big.txt")
	if result != Success {
		t.Fatal(proxyAddress, result)
	}
}

func TestPortForwardBindIPv6Request(t *testing.T) {
	h := StartIPv6HTTPServer(t)
	defer h.Close()
	p := NewBindPipe(port_forward.NewForward(httpIPv6Address), nil, nil)
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
		port_forward.NewForward(httpAddress), nil, nil)
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
		port_forward.NewForward(httpIPv6Address), nil, nil)
	defer func() {
		a.Close()
		b.Close()
	}()
	result := GetRequestRaw("http://" + proxyAddress + "/big.txt")
	if result != Success {
		t.Fatal(proxyAddress, result)
	}
}
