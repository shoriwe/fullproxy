package test

import (
	"bytes"
	"github.com/shoriwe/FullProxy/pkg/Pipes"
	"github.com/shoriwe/FullProxy/pkg/Proxies/SOCKS5"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
)

const (
	username     = "sulcud"
	password     = "password"
	testUrl      = "http://127.0.0.1:8080/big.txt"
	networkType  = "tcp"
	proxyAddress = "127.0.0.1:9050"
)

var (
	httpListener net.Listener
	bindPipe     *Pipes.Bind
)

func TestInitHTTPServer(t *testing.T) {
	var listenError error
	httpListener, listenError = net.Listen(networkType, "127.0.0.1:8080")
	if listenError != nil {
		t.Fatal(listenError)
	}
	http.HandleFunc("/big.txt",
		func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write(bytes.Repeat([]byte{'A'}, 0xFFFFFFF))
		},
	)

	go http.Serve(httpListener, nil)
}

func TestNoAuthInitialization(t *testing.T) {
	var pipeInitializationError error
	bindPipe, pipeInitializationError = Pipes.NewBindPipe(
		networkType,
		proxyAddress,
		SOCKS5.NewSocks5(nil, t.Log, nil),
		t.Log,
		nil,
	)
	if pipeInitializationError != nil {
		t.Fatal(pipeInitializationError)
	}
	go bindPipe.Serve()
}

func TestNoAuthHTTPRequest(t *testing.T) {
	// create a socks5 dialer
	dialer, err := proxy.SOCKS5(networkType, proxyAddress, nil, proxy.Direct)
	if err != nil {
		t.Fatal("can't connect to the proxy:", err)
	}
	// setup a http client
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	// set our socks5 as the dialer
	httpTransport.Dial = dialer.Dial
	// create a request
	req, err := http.NewRequest("GET", testUrl, nil)
	if err != nil {
		t.Fatal("can't create request:", err)
	}
	// use the http client to fetch the page
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatal("can't GET page:", err)
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("error reading body:", err)
	}
	// t.Log(string(b))
}

func TestCloseNoAuthPipe(t *testing.T) {
	closingError := bindPipe.Server.Close()
	if closingError != nil {
		t.Fatal(closingError)
	}
}

func TestUserPasswordAuthSocks5Init(t *testing.T) {
	var pipeInitializationError error
	bindPipe, pipeInitializationError = Pipes.NewBindPipe(
		networkType,
		proxyAddress,
		SOCKS5.NewSocks5(
			func(username []byte, password []byte) (bool, error) {
				if bytes.Equal(username, []byte("sulcud")) &&
					bytes.Equal(password, []byte("password")) {
					return true, nil
				}
				return false, nil
			},
			t.Log,
			nil,
		),
		t.Log,
		nil,
	)
	if pipeInitializationError != nil {
		t.Fatal(pipeInitializationError)
	}
	go bindPipe.Serve()
}

func TestUsernamePasswordHTTPRequest(t *testing.T) {
	// create a socks5 dialer
	dialer, err := proxy.SOCKS5(networkType, proxyAddress, &proxy.Auth{
		User:     username,
		Password: password,
	}, proxy.Direct)
	if err != nil {
		t.Fatal("can't connect to the proxy:", err)
	}
	// setup a http client
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	// set our socks5 as the dialer
	httpTransport.Dial = dialer.Dial
	// create a request
	req, err := http.NewRequest("GET", testUrl, nil)
	if err != nil {
		t.Fatal("can't create request:", err)
	}
	// use the http client to fetch the page
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatal("can't GET page:", err)
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("error reading body:", err)
	}
	// t.Log(string(b))
}

func TestUsernamePasswordWithNoAuthHTTPRequest(t *testing.T) {
	// create a socks5 dialer
	dialer, err := proxy.SOCKS5(networkType, proxyAddress, nil, proxy.Direct)
	if err != nil {
		t.Fatal("can't connect to the proxy:", err)
	}
	// setup a http client
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	// set our socks5 as the dialer
	httpTransport.Dial = dialer.Dial
	// create a request
	req, err := http.NewRequest("GET", testUrl, nil)
	if err != nil {
		t.Fatal("can't create request:", err)
	}
	// use the http client to fetch the page
	_, err = httpClient.Do(req)
	if err != nil {
		return
	}
	t.Fatal("Authentication bypassed")
}

func TestCloseUserPasswordAuthPipe(t *testing.T) {
	closingError := bindPipe.Server.Close()
	if closingError != nil {
		t.Fatal(closingError)
	}
}

func TestFinishHTTPServer(t *testing.T) {
	_ = httpListener.Close()
}
