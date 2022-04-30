package test

import (
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers/socks5"
	"testing"
)

//// Test No auth

func TestSocks5NoAuthHTTPRequest(t *testing.T) {
	h := StartIPv4HTTPServer(t)
	defer h.Close()
	p := NewBindPipe(socks5.NewSocks5(nil), nil, nil)
	defer p.Close()
	if GetRequestSocks5(testUrl, "", "") != Success {
		t.Fatal(testUrl)
	}
}

//// IPv6

func TestSocks5IPv6HTTPRequest(t *testing.T) {
	h := StartIPv6HTTPServer(t)
	defer h.Close()
	p := NewBindPipe(socks5.NewSocks5(nil), nil, nil)
	defer p.Close()
	result := GetRequestSocks5(testUrlIPv6, "", "")
	if result != Success {
		t.Fatal(testUrlIPv6, result)
	}
}

//// Test Auth

func TestSocks5UsernamePasswordHTTPRequest(t *testing.T) {
	h := StartIPv4HTTPServer(t)
	defer h.Close()
	p := NewBindPipe(socks5.NewSocks5(basicAuthFunc), nil, nil)
	defer p.Close()
	if GetRequestSocks5(testUrl, "sulcud", "password") != Success {
		t.Fatal(testUrl)
	}
	if GetRequestSocks5(testUrl, "", "") != FailedRequest {
		t.Fatal("Auth bypassed")
	}
}

//// Test inbound rules

func TestSocks5InvalidInboundHTTPRequest(t *testing.T) {
	h := StartIPv4HTTPServer(t)
	defer h.Close()
	p := NewBindPipe(socks5.NewSocks5(basicAuthFunc), basicInboundRule, nil)
	defer p.Close()
	if GetRequestSocks5(testUrl, "sulcud", "password") != FailedRequest {
		t.Fatal("Bypassed inbound")
	}
}

//// Test outbound rules

func TestSocks5OutboundHTTPRequest(t *testing.T) {
	h := StartIPv4HTTPServer(t)
	defer h.Close()
	p := NewBindPipe(socks5.NewSocks5(basicAuthFunc),
		nil, basicOutboundRule,
	)
	defer p.Close()
	if GetRequestSocks5("google.com", "sulcud", "password") == Success {
		t.Fatal("Bypassed outbound")
	}
	if GetRequestSocks5(testUrl, "sulcud", "password") != Success {
		t.Fatal(testUrl)
	}
}

// Test master slave

//// Test No auth

func TestSocks5NoAuthMasterSlaveHTTPRequest(t *testing.T) {
	h := StartIPv4HTTPServer(t)
	defer h.Close()
	a, b := NewMasterSlave(
		socks5.NewSocks5(nil), nil, nil)
	defer func() {
		a.Close()
		b.Close()

	}()
	if GetRequestSocks5(testUrl, "", "") != Success {
		t.Fatal(testUrl)
	}
}

func TestSocks5NoAuthIPv6MasterSlaveHTTPRequest(t *testing.T) {
	h := StartIPv6HTTPServer(t)
	defer h.Close()
	a, b := NewMasterSlave(
		socks5.NewSocks5(nil), nil, nil)
	defer func() {
		a.Close()
		b.Close()
	}()
	if GetRequestSocks5(testUrlIPv6, "", "") != Success {
		t.Fatal(testUrl)
	}
}

//// Test Auth

func TestSocks5UsernamePasswordMasterSlaveHTTPRequest(t *testing.T) {
	h := StartIPv4HTTPServer(t)
	defer h.Close()
	a, b := NewMasterSlave(
		socks5.NewSocks5(basicAuthFunc), nil, nil)
	defer func() {
		a.Close()
		b.Close()

	}()
	if GetRequestSocks5(testUrl, "sulcud", "password") != Success {
		t.Fatal(testUrl)
	}
}

//// Test inbound rules

func TestSocks5InboundMasterSlaveHTTPRequest(t *testing.T) {
	h := StartIPv4HTTPServer(t)
	defer h.Close()
	a, b := NewMasterSlave(
		socks5.NewSocks5(basicAuthFunc), basicInboundRule, nil)
	defer func() {
		a.Close()
		b.Close()

	}()
	if GetRequestSocks5(testUrl, "sulcud", "password") != FailedRequest {
		t.Fatal("Bypassed inbound rule")
	}
}

//// Test outbound rules

func TestSocks5OutboundMasterSlaveHTTPRequest(t *testing.T) {
	h := StartIPv4HTTPServer(t)
	defer h.Close()
	a, b := NewMasterSlave(
		socks5.NewSocks5(basicAuthFunc), nil, basicOutboundRule)
	defer func() {
		a.Close()
		b.Close()

	}()
	if GetRequestSocks5("https://google.com", "sulcud", "password") != FailedRequest {
		t.Fatal("Bypassed outbound rule")
	}
	if GetRequestSocks5(testUrl, "sulcud", "password") != Success {
		t.Fatal(testUrl)
	}
}

// BIND

func TestSocks5NoAuthBind(t *testing.T) {
	Socks5BindTest(
		proxyAddress,
		nil,
		"", "",
		t,
	)
}

func TestSocks5BasicAuthBind(t *testing.T) {
	Socks5BindTest(
		proxyAddress,
		basicAuthFunc,
		"sulcud", "password",
		t,
	)
}
