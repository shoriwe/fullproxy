package test

import (
	"bytes"
	"github.com/shoriwe/FullProxy/pkg/Pipes"
	"github.com/shoriwe/FullProxy/pkg/Proxies/SOCKS5"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"log"
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

func init() {
	var listenError error
	httpListener, listenError = net.Listen(networkType, "127.0.0.1:8080")
	if listenError != nil {
		log.Fatal(listenError)
	}
	http.HandleFunc("/big.txt",
		func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write(bytes.Repeat([]byte{'A'}, 0xFFFF))
		},
	)

	go http.Serve(httpListener, nil)
}

// Test Bind

func TestNoAuthInitialization(t *testing.T) {
	var pipeInitializationError error
	bindPipe, pipeInitializationError = Pipes.NewBindPipe(
		networkType,
		proxyAddress,
		SOCKS5.NewSocks5(nil, log.Print, nil),
		log.Print,
		nil,
	)
	if pipeInitializationError != nil {
		t.Fatal(pipeInitializationError)
	}
	go bindPipe.Serve()
}

//// Test No auth

func TestNoAuthHTTPRequest(t *testing.T) {
	dialer, err := proxy.SOCKS5(networkType, proxyAddress, nil, proxy.Direct)
	if err != nil {
		t.Fatal("can't connect to the proxy:", err)
	}
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	httpTransport.Dial = dialer.Dial
	req, err := http.NewRequest("GET", testUrl, nil)
	if err != nil {
		t.Fatal("can't create request:", err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatal("can't GET page:", err)
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("error reading body:", err)
	}
}

func TestCloseNoAuthPipe(t *testing.T) {
	closingError := bindPipe.Server.Close()
	if closingError != nil {
		t.Fatal(closingError)
	}
}

//// Test Auth

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
			log.Print,
			nil,
		),
		log.Print,
		nil,
	)
	if pipeInitializationError != nil {
		t.Fatal(pipeInitializationError)
	}
	go bindPipe.Serve()
}

func TestUsernamePasswordHTTPRequest(t *testing.T) {
	dialer, err := proxy.SOCKS5(networkType, proxyAddress, &proxy.Auth{
		User:     username,
		Password: password,
	}, proxy.Direct)
	if err != nil {
		t.Fatal("can't connect to the proxy:", err)
	}
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	httpTransport.Dial = dialer.Dial
	req, err := http.NewRequest("GET", testUrl, nil)
	if err != nil {
		t.Fatal("can't create request:", err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatal("can't GET page:", err)
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("error reading body:", err)
	}
}

func TestUsernamePasswordWithNoAuthHTTPRequest(t *testing.T) {
	dialer, err := proxy.SOCKS5(networkType, proxyAddress, nil, proxy.Direct)
	if err != nil {
		t.Fatal("can't connect to the proxy:", err)
	}
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	httpTransport.Dial = dialer.Dial
	req, err := http.NewRequest("GET", testUrl, nil)
	if err != nil {
		t.Fatal("can't create request:", err)
	}
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

//// Test inbound rules
func TestInboundRulesSocks5Init(t *testing.T) {
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
			log.Print,
			nil,
		),
		log.Print,
		func(host string) bool {
			if host == "127.0.0.1" {
				return false
			}
			return true
		},
	)
	if pipeInitializationError != nil {
		t.Fatal(pipeInitializationError)
	}
	go bindPipe.Serve()
}

func TestInvalidInboundHTTPRequest(t *testing.T) {
	dialer, err := proxy.SOCKS5(networkType, proxyAddress, &proxy.Auth{
		User:     username,
		Password: password,
	}, proxy.Direct)
	if err != nil {
		t.Fatal("can't connect to the proxy:", err)
	}
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	httpTransport.Dial = dialer.Dial
	req, err := http.NewRequest("GET", testUrl, nil)
	if err != nil {
		t.Fatal("can't create request:", err)
	}
	_, err = httpClient.Do(req)
	if err != nil {
		return
	}
	log.Fatal("Bypassed inbound rules")
}

func TestCloseInboundRulesPipe(t *testing.T) {
	closingError := bindPipe.Server.Close()
	if closingError != nil {
		t.Fatal(closingError)
	}
}

//// Test outbound rules

func TestOutboundRulesSocks5Init(t *testing.T) {
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
			log.Print,
			func(host string) bool {
				if host == "google.com" {
					return false
				}
				return true
			},
		),
		log.Print,
		nil,
	)
	if pipeInitializationError != nil {
		t.Fatal(pipeInitializationError)
	}
	go bindPipe.Serve()
}

func TestInvalidOutboundHTTPRequest(t *testing.T) {
	dialer, err := proxy.SOCKS5(networkType, proxyAddress, &proxy.Auth{
		User:     username,
		Password: password,
	}, proxy.Direct)
	if err != nil {
		t.Fatal("can't connect to the proxy:", err)
	}
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	httpTransport.Dial = dialer.Dial
	req, err := http.NewRequest("GET", "https://google.com", nil)
	if err != nil {
		t.Fatal("can't create request:", err)
	}
	_, err = httpClient.Do(req)
	if err != nil {
		return
	}
	log.Fatal("Bypassed outbound rules")
}

func TestOutboundSuccessHTTPRequest(t *testing.T) {
	dialer, err := proxy.SOCKS5(networkType, proxyAddress, &proxy.Auth{
		User:     username,
		Password: password,
	}, proxy.Direct)
	if err != nil {
		t.Fatal("can't connect to the proxy:", err)
	}
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	httpTransport.Dial = dialer.Dial
	req, err := http.NewRequest("GET", testUrl, nil)
	if err != nil {
		t.Fatal("can't create request:", err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatal("can't GET page:", err)
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("error reading body:", err)
	}
}

func TestCloseOutboundRulesPipe(t *testing.T) {
	closingError := bindPipe.Server.Close()
	if closingError != nil {
		t.Fatal(closingError)
	}
}

// Test Master Slave

//// Test No auth

//// Test Auth

//// Test inbound rules

//// Test outbound rules

// Finally, close the HTTP server

func TestClose(t *testing.T) {
	_ = httpListener.Close()
}
