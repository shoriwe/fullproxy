package test

import (
	http2 "github.com/shoriwe/fullproxy/v3/internal/proxy/servers/http"
	"testing"
)

//// Test No auth

func TestHTTPNoAuthHTTPRequest(t *testing.T) {
	h := StartIPv4HTTPServer(t)
	defer h.Close()
	p := NewBindHandler(http2.NewHTTP(nil), nil)
	defer p.Close()
	if GetRequestHTTP(testUrl, "", "") != Success {
		t.Fatal(testUrl)
	}
}

//// Test Auth

func TestHTTPUsernamePasswordHTTPRequest(t *testing.T) {
	h := StartIPv4HTTPServer(t)
	defer h.Close()
	p := NewBindHandler(http2.NewHTTP(basicAuthFunc), nil)
	defer p.Close()
	if GetRequestHTTP(testUrl, "sulcud", "password") != Success {
		t.Fatal(testUrl)
	}
	if GetRequestHTTP(testUrl, "", "") != FailedRequest {
		t.Fatal("Auth bypassed")
	}
}

//// Test inbound rules

func TestHTTPInvalidInboundHTTPRequest(t *testing.T) {
	h := StartIPv4HTTPServer(t)
	defer h.Close()
	p := NewBindHandler(http2.NewHTTP(basicAuthFunc), &BasicInbound{})
	defer p.Close()
	if GetRequestHTTP(testUrl, "sulcud", "password") != FailedRequest {
		t.Fatal("Bypassed inbound")
	}
}

//// Test outbound rules

func TestHTTPOutboundHTTPRequest(t *testing.T) {
	h := StartIPv4HTTPServer(t)
	defer h.Close()
	p := NewBindHandler(http2.NewHTTP(basicAuthFunc), nil)
	defer p.Close()
	if GetRequestHTTP("google.com", "sulcud", "password") == Success {
		t.Fatal("Bypassed outbound")
	}
	if GetRequestHTTP(testUrl, "sulcud", "password") != Success {
		t.Fatal(testUrl)
	}
}

// Test master slave

//// Test No auth

func TestHTTPNoAuthMasterSlaveHTTPRequest(t *testing.T) {
	h := StartIPv4HTTPServer(t)
	defer h.Close()
	a, b := NewMasterSlaveHandler(
		http2.NewHTTP(nil), nil)
	defer func() {
		a.Close()
		b.Close()

	}()
	if GetRequestHTTP(testUrl, "", "") != Success {
		t.Fatal(testUrl)
	}
}

//// Test Auth

func TestHTTPUsernamePasswordMasterSlaveHTTPRequest(t *testing.T) {
	h := StartIPv4HTTPServer(t)
	defer h.Close()
	a, b := NewMasterSlaveHandler(
		http2.NewHTTP(basicAuthFunc), nil)
	defer func() {
		a.Close()
		b.Close()

	}()
	if GetRequestHTTP(testUrl, "sulcud", "password") != Success {
		t.Fatal(testUrl)
	}
}

//// Test inbound rules

func TestHTTPInboundMasterSlaveHTTPRequest(t *testing.T) {
	h := StartIPv4HTTPServer(t)
	defer h.Close()
	a, b := NewMasterSlaveHandler(
		http2.NewHTTP(basicAuthFunc), &BasicInbound{})
	defer func() {
		a.Close()
		b.Close()
	}()
	if GetRequestHTTP(testUrl, "sulcud", "password") != FailedRequest {
		t.Fatal("Bypassed inbound rule")
	}
}

//// Test outbound rules

func TestHTTPInvalidOutboundMasterSlaveHTTPRequest(t *testing.T) {
	h := StartIPv4HTTPServer(t)
	a, b := NewMasterSlaveHandler(
		http2.NewHTTP(basicAuthFunc), nil)
	defer func() {
		_ = h.Close()
		_ = a.Close()
		_ = b.Close()
	}()
	if GetRequestHTTP("https://google.com", "sulcud", "password") != FailedRequest {
		t.Fatal("Bypassed outbound rule")
	}
}

func TestHTTPOutboundMasterSlaveHTTPRequest(t *testing.T) {
	h := StartIPv4HTTPServer(t)
	defer h.Close()
	a, b := NewMasterSlaveHandler(
		http2.NewHTTP(basicAuthFunc), nil)
	defer func() {
		a.Close()
		b.Close()

	}()
	if GetRequestHTTP(testUrl, "sulcud", "password") != Success {
		t.Fatal(testUrl)
	}
}
