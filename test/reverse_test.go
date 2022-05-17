package test

import (
	"bytes"
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers/reverse"
	"io"
	"net/http"
	"testing"
)

func TestRawReverseProxy(t *testing.T) {
	defer http.DefaultClient.CloseIdleConnections()
	httpServer := StartIPv4HTTPServer(t)
	defer httpServer.Close()
	reverseProxy := NewBindPipe(reverse.NewRaw([]*reverse.Host{
		{
			Url:     "http://127.0.0.1:8080",
			Network: "tcp",
			Address: "127.0.0.1:8080",
		},
	}), nil)
	defer reverseProxy.Close()

	result := GetRequestRaw("http://127.0.0.1:9050")
	if result != Success {
		t.Fatal("Failed request")
	}
}

func TestMultipleRequestRawReverseProxy(t *testing.T) {
	defer http.DefaultClient.CloseIdleConnections()
	httpServer := StartIPv4HTTPServer(t)
	defer httpServer.Close()
	reverseProxy := NewBindPipe(reverse.NewRaw([]*reverse.Host{
		{
			Url:     "http://127.0.0.1:8080",
			Network: "tcp",
			Address: "127.0.0.1:8080",
		},
	}), nil)
	defer reverseProxy.Close()

	for i := 0; i < 100; i++ {
		if GetRequestRaw("http://127.0.0.1:9050/big") != Success {
			t.Fatal("Failed request")
		}
	}
}

func TestMultipleRequestPoolRawReverseProxy(t *testing.T) {
	defer http.DefaultClient.CloseIdleConnections()
	server1, server2 := NewHTTPServers(t)
	defer server1.Close()
	defer server2.Close()

	reverseProxy := NewBindPipe(
		reverse.NewRaw([]*reverse.Host{
			{
				Url:     "http://127.0.0.1:8080",
				Network: "tcp",
				Address: "127.0.0.1:8080",
			},
			{
				Url:     "http://127.0.0.1:8081",
				Network: "tcp",
				Address: "127.0.0.1:8081",
			},
		}),
		nil,
	)
	defer reverseProxy.Close()

	if GetRequestToPool() != Success {
		t.Fatal("Failed")
	}
}

func TestHTTPReverseProxy(t *testing.T) {
	defer http.DefaultClient.CloseIdleConnections()
	httpServer := StartIPv4HTTPServer(t)
	defer httpServer.Close()
	reverseProxy := NewBindHandler(reverse.NewHTTP(
		map[string]*reverse.Target{
			"127.0.0.1:9050": {
				RequestHeader: http.Header{},
				Path:          "/",
				Targets: []*reverse.Host{
					{
						Url:     "http://127.0.0.1:8080",
						Network: "tcp",
						Address: "127.0.0.1:8080",
					},
				},
			},
		},
	), nil)
	defer reverseProxy.Close()

	response, requestError := http.Get("http://127.0.0.1:9050/big")
	if requestError != nil {
		t.Fatal(requestError)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatal("Failed")
	}
	contents, readError := io.ReadAll(response.Body)
	if readError != nil {
		t.Fatal(readError)
	}
	if !bytes.Equal(contents, bytes.Repeat([]byte{'A'}, DefaultChunkSize)) {
		t.Fatal(string(contents))
	}
}

func TestMultipleRequestHTTPReverseProxy(t *testing.T) {
	defer http.DefaultClient.CloseIdleConnections()
	httpServer := StartIPv4HTTPServer(t)
	defer httpServer.Close()
	reverseProxy := NewBindHandler(reverse.NewHTTP(
		map[string]*reverse.Target{
			"127.0.0.1:9050": {
				RequestHeader:  http.Header{},
				ResponseHeader: http.Header{},
				Path:           "/",
				Targets: []*reverse.Host{
					{
						Url:     "http://127.0.0.1:8080",
						Network: "tcp",
						Address: "127.0.0.1:8080",
					},
				},
			},
		},
	), nil)
	defer reverseProxy.Close()

	for i := 0; i < 100; i++ {
		response, requestError := http.Get("http://127.0.0.1:9050/big")
		if requestError != nil {
			t.Fatal(requestError)
		}
		if response.StatusCode != http.StatusOK {
			t.Fatal("Failed")
		}
		contents, readError := io.ReadAll(response.Body)
		if readError != nil {
			t.Fatal(readError)
		}
		if !bytes.Equal(contents, bytes.Repeat([]byte{'A'}, DefaultChunkSize)) {
			t.Fatal(string(contents))
		}
	}
}

func TestMultipleRequestPoolHTTPReverseProxy(t *testing.T) {
	defer http.DefaultClient.CloseIdleConnections()
	httpServer1, httpServer2 := NewHTTPServers(t)
	defer httpServer1.Close()
	defer httpServer2.Close()
	reverseProxy := NewBindHandler(
		reverse.NewHTTP(
			map[string]*reverse.Target{
				"127.0.0.1:9050": {
					RequestHeader:  http.Header{},
					ResponseHeader: http.Header{},
					Path:           "/",
					CurrentTarget:  0,
					Targets: []*reverse.Host{
						{
							Url:     "http://127.0.0.1:8080",
							Network: "tcp",
							Address: "127.0.0.1:8080",
						},
						{
							Url:     "http://127.0.0.1:8081",
							Network: "tcp",
							Address: "127.0.0.1:8081",
						},
					},
				},
			},
		),
		nil,
	)
	defer reverseProxy.Close()

	response, requestError := http.Get("http://127.0.0.1:9050/big")
	if requestError != nil {
		t.Fatal(requestError)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatal("Failed")
	}
	contents, readError := io.ReadAll(response.Body)
	if readError != nil {
		t.Fatal(readError)
	}
	if !bytes.Equal(contents, bytes.Repeat([]byte{'A'}, DefaultChunkSize)) {
		t.Fatal(string(contents))
	}

	response, requestError = http.Get("http://127.0.0.1:9050/big")
	if requestError != nil {
		t.Fatal(requestError)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatal("Failed")
	}
	contents, readError = io.ReadAll(response.Body)
	if readError != nil {
		t.Fatal(readError)
	}
	if !bytes.Equal(contents, bytes.Repeat([]byte{'B'}, DefaultChunkSize)) {
		t.Fatal(string(contents))
	}
}

func TestHTTPReversePathBasedProxy(t *testing.T) {
	defer http.DefaultClient.CloseIdleConnections()
	httpServer := StartIPv4HTTPServer(t)
	defer httpServer.Close()
	reverseProxy := NewBindHandler(reverse.NewHTTP(
		map[string]*reverse.Target{
			"127.0.0.1:9050": {
				RequestHeader:  http.Header{},
				ResponseHeader: http.Header{},
				Path:           "/app",
				Targets: []*reverse.Host{
					{
						Url:     "http://127.0.0.1:8080",
						Network: "tcp",
						Address: "127.0.0.1:8080",
					},
				},
			},
		},
	), nil)
	defer reverseProxy.Close()

	response, requestError := http.Get("http://127.0.0.1:9050/app/big")
	if requestError != nil {
		t.Fatal(requestError)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatal("Failed")
	}
	contents, readError := io.ReadAll(response.Body)
	if readError != nil {
		t.Fatal(readError)
	}
	if !bytes.Equal(contents, bytes.Repeat([]byte{'A'}, DefaultChunkSize)) {
		t.Fatal(string(contents))
	}
}

func TestHTTPReverseInjectHeadersProxy(t *testing.T) {
	defer http.DefaultClient.CloseIdleConnections()
	httpServer := StartIPv4HTTPServer(t)
	defer httpServer.Close()
	reverseProxy := NewBindHandler(reverse.NewHTTP(
		map[string]*reverse.Target{
			"127.0.0.1:9050": {
				RequestHeader: http.Header{},
				ResponseHeader: http.Header{
					"Name": []string{"sulcud"},
				},
				Path: "/app",
				Targets: []*reverse.Host{
					{
						Url:     "http://127.0.0.1:8080",
						Network: "tcp",
						Address: "127.0.0.1:8080",
					},
				},
			},
		},
	), nil)
	defer reverseProxy.Close()

	response, requestError := http.Get("http://127.0.0.1:9050/app/big")
	if requestError != nil {
		t.Fatal(requestError)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatal("Failed")
	}
	contents, readError := io.ReadAll(response.Body)
	if readError != nil {
		t.Fatal(readError)
	}
	if !bytes.Equal(contents, bytes.Repeat([]byte{'A'}, DefaultChunkSize)) {
		t.Fatal(string(contents))
	}

	if response.Header.Get("Name") != "sulcud" {
		t.Fatal("invalid Name header")
	}
}
