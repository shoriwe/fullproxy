package test

import (
	"bytes"
	socks52 "github.com/shoriwe/fullproxy/v3/internal/proxy/clients/socks5"
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers/socks5"
	"net"
	"strings"
	"testing"
	"time"
)

//// Test No auth

func TestSocks5NoAuthHTTPRequest(t *testing.T) {
	h := StartIPv4HTTPServer(t)
	defer h.Close()
	p := NewBindPipe(socks5.NewSocks5(nil), nil)
	defer p.Close()
	if GetRequestSocks5(testUrl, "", "") != Success {
		t.Fatal(testUrl)
	}
}

//// IPv6

func TestSocks5IPv6HTTPRequest(t *testing.T) {
	h := StartIPv6HTTPServer(t)
	defer h.Close()
	p := NewBindPipe(socks5.NewSocks5(nil), nil)
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
	p := NewBindPipe(socks5.NewSocks5(basicAuthFunc), nil)
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
	p := NewBindPipe(socks5.NewSocks5(basicAuthFunc), &BasicInbound{})
	defer p.Close()
	if GetRequestSocks5(testUrl, "sulcud", "password") != FailedRequest {
		t.Fatal("Bypassed inbound")
	}
}

//// Test outbound rules

func TestSocks5OutboundHTTPRequest(t *testing.T) {
	h := StartIPv4HTTPServer(t)
	defer h.Close()
	p := NewBindPipe(socks5.NewSocks5(basicAuthFunc), nil)
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
		socks5.NewSocks5(nil), nil)
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
		socks5.NewSocks5(nil), nil)
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
		socks5.NewSocks5(basicAuthFunc), nil)
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
		socks5.NewSocks5(basicAuthFunc), &BasicInbound{})
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
		socks5.NewSocks5(basicAuthFunc), &BasicOutbound{})
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
	Socks6BindSucceed(
		proxyAddress,
		nil, nil,
		"", "",
		t,
	)
}

func TestSocks5BasicAuthBind(t *testing.T) {
	Socks6BindSucceed(
		proxyAddress,
		basicAuthFunc,
		nil,
		"sulcud", "password",
		t,
	)
}

func TestSocks5BindListenFilter(t *testing.T) {
	time.Sleep(10 * time.Second)
	proxyServer := NewBindPipe(
		socks5.NewSocks5(basicAuthFunc),
		&BasicListen{},
	)
	defer proxyServer.Close()
	socksClient := socks52.Socks5{
		Network:  "tcp",
		Address:  proxyAddress,
		Username: "sulcud",
		Password: "password",
	}
	_, listenError := socksClient.Listen("tcp", SampleAddress.String())
	if listenError == nil {
		t.Fatal("the listen should fail")
	}
}

func TestSocks5BindAcceptFilter(t *testing.T) {
	time.Sleep(10 * time.Second)
	proxyServer := NewBindPipe(
		socks5.NewSocks5(basicAuthFunc),
		&BasicAccept{},
	)
	defer proxyServer.Close()
	socksClient := socks52.Socks5{
		Network:  "tcp",
		Address:  proxyAddress,
		Username: "sulcud",
		Password: "password",
	}
	var (
		listener                     net.Listener
		client, server               net.Conn
		connectionError, listenError error
	)
	go func() {
		listener, listenError = socksClient.Listen("tcp", SampleAddress.String())
		if listenError != nil {
			t.Fatal(listenError)
		}
		client, connectionError = listener.Accept()
		if connectionError != nil {
			t.Fatal(connectionError)
		}
		_, writeError := client.Write(SampleMessage)
		if writeError != nil {
			if !strings.Contains(writeError.Error(), "closed") {
				t.Fatal(writeError)
			}
		}
	}()
	time.Sleep(2 * time.Second)
	defer func(listener net.Listener) {
		if listener != nil {
			_ = listener.Close()
		}
	}(listener)
	defer func(client net.Conn) {
		if client != nil {
			_ = client.Close()
		}
	}(client)
	server, connectionError = net.DialTCP("tcp", SampleAddress, listener.Addr().(*net.TCPAddr))
	if connectionError != nil {
		t.Fatal(connectionError)
	}
	defer server.Close()
	response := make([]byte, len(SampleMessage))
	_, readError := server.Read(response)
	if readError == nil {
		t.Fatal("connection should close, expected EOF or connection closed")
	}
	if bytes.Equal(response, SampleMessage) {
		t.Fatal("empty array should not equal to targeted one")
	}
}
