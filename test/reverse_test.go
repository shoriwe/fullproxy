package test

import (
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers/reverse"
	"testing"
)

func TestRawReverseProxy(t *testing.T) {
	httpServer := StartIPv4HTTPServer(t)
	defer httpServer.Close()
	reverseProxy := NewBindPipe(reverse.NewRaw([]string{"127.0.0.1:8080"}), nil)
	defer reverseProxy.Close()

	result := GetRequestRaw("http://127.0.0.1:9050")
	if result != Success {
		t.Fatal("Failed request")
	}
}

func TestMultipleRequestRawReverseProxy(t *testing.T) {
	httpServer := StartIPv4HTTPServer(t)
	defer httpServer.Close()
	reverseProxy := NewBindPipe(reverse.NewRaw([]string{"127.0.0.1:8080"}), nil)
	defer reverseProxy.Close()

	for i := 0; i < 100; i++ {
		if GetRequestRaw("http://127.0.0.1:9050/big") != Success {
			t.Fatal("Failed request")
		}
	}
}

func TestMultipleRequestPoolRawReverseProxy(t *testing.T) {
	server1, server2 := NewHTTPServers(t)
	defer server1.Close()
	defer server2.Close()

	reverseProxy := NewBindPipe(
		reverse.NewRaw([]string{"127.0.0.1:8080", "127.0.0.1:8081"}),
		nil,
	)
	defer reverseProxy.Close()

	if GetRequestToPool() != Success {
		t.Fatal("Failed")
	}
}
