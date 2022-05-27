package test

import (
	"bytes"
	"github.com/gorilla/websocket"
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers/reverse"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestRawReverseProxy(t *testing.T) {
	defer http.DefaultClient.CloseIdleConnections()
	httpServer := StartIPv4HTTPServer(t)
	defer httpServer.Close()
	reverseProxy := NewBindPipe(reverse.NewRaw([]*reverse.Host{
		{
			Network: "tcp",
			Address: "127.0.0.1:8080",
		},
	}), nil)
	defer reverseProxy.Close()

	result := GetRequestRaw("http://127.0.0.1:9050/big")
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
				Network: "tcp",
				Address: "127.0.0.1:8080",
			},
			{
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
				URI:           "/",
				Hosts: []*reverse.Host{
					{
						Scheme:    "http",
						URI:       "/",
						Network:   "tcp",
						Address:   "127.0.0.1:8080",
						TLSConfig: nil,
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
				URI:            "/",
				Hosts: []*reverse.Host{
					{
						Scheme:    "http",
						URI:       "/",
						Network:   "tcp",
						Address:   "127.0.0.1:8080",
						TLSConfig: nil,
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
					URI:            "/",
					CurrentHost:    0,
					Hosts: []*reverse.Host{
						{
							Scheme:    "http",
							URI:       "/",
							Network:   "tcp",
							Address:   "127.0.0.1:8080",
							TLSConfig: nil,
						},
						{
							Scheme:    "http",
							URI:       "/",
							Network:   "tcp",
							Address:   "127.0.0.1:8081",
							TLSConfig: nil,
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
				URI:            "/app",
				Hosts: []*reverse.Host{
					{
						Scheme:    "http",
						URI:       "/",
						Network:   "tcp",
						Address:   "127.0.0.1:8080",
						TLSConfig: nil,
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
				URI: "/app",
				Hosts: []*reverse.Host{
					{
						Scheme:    "http",
						URI:       "/",
						Network:   "tcp",
						Address:   "127.0.0.1:8080",
						TLSConfig: nil,
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

func TestHTTPReverseWebSocket(t *testing.T) {
	defer http.DefaultClient.CloseIdleConnections()
	httpServer := StartIPv4HTTPServer(t)
	defer httpServer.Close()
	h := reverse.NewHTTP(
		map[string]*reverse.Target{
			"127.0.0.1:9050": {
				RequestHeader:  http.Header{},
				ResponseHeader: http.Header{},
				URI:            "/",
				Hosts: []*reverse.Host{
					{
						WebsocketReadBufferSize:  1024,
						WebsocketWriteBufferSize: 1024,
						Scheme:                   "http",
						URI:                      "/",
						Network:                  "tcp",
						Address:                  "127.0.0.1:8080",
						TLSConfig:                nil,
					},
				},
			},
		},
	)
	h.SetSniffers(os.Stdout, os.Stdout)
	reverseProxy := NewBindHandler(h, nil)
	defer reverseProxy.Close()

	connection, _, dialError := websocket.DefaultDialer.Dial("ws://127.0.0.1:9050/ws", nil)

	if dialError != nil {
		t.Fatal(dialError)
	}
	defer connection.Close()
	msg := wsMessage{
		Succeed: true,
		MSG:     ClientMessage,
	}
	writeError := connection.WriteJSON(msg)
	if writeError != nil {
		t.Fatal(writeError)
	}
	readError := connection.ReadJSON(&msg)
	if readError != nil {
		t.Fatal(readError)
	}
	if msg.MSG != ServerMessage {
		t.Fatal(msg.MSG)
	}
}
