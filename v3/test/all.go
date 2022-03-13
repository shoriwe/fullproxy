package test

import (
	"bytes"
	"fmt"
	"github.com/shoriwe/FullProxy/v3/internal/global"
	"github.com/shoriwe/FullProxy/v3/internal/pipes"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
)

const (
	testUrl      = "http://127.0.0.1:8080/big.txt"
	networkType  = "tcp"
	httpAddress  = "127.0.0.1:8080"
	proxyAddress = "127.0.0.1:9050"
	c2Address    = "127.0.0.1:9051"
	Success      = iota
	FailedProxySetup
	FailedRequest
)

func basicAuthFunc(username []byte, password []byte) (bool, error) {
	if bytes.Equal(username, []byte("sulcud")) &&
		bytes.Equal(password, []byte("password")) {
		return true, nil
	}
	return false, nil
}

func basicOutboundRule(host string) bool {
	if host == "google.com" {
		return false
	}
	return true
}

func basicInboundRule(host string) bool {
	if host == "127.0.0.1" {
		return false
	}
	return true
}

func NewBindPipe(protocol global.Protocol, inboundFilter global.IOFilter) net.Listener {
	bindPipe := pipes.NewBindPipe(
		networkType, proxyAddress,
		protocol,
		nil,
		inboundFilter,
	)
	go bindPipe.Serve()
	time.Sleep(2 * time.Second)
	return bindPipe.Server
}

func NewMasterSlave(inboundFilter global.IOFilter, protocol global.Protocol) (net.Listener, net.Listener) {
	cert, signError := pipes.SelfSignCertificate()
	if signError != nil {
		panic(signError)
	}
	masterPipe := pipes.NewMaster(
		networkType,
		c2Address,
		proxyAddress,
		nil,
		inboundFilter,
		protocol,
		cert,
	)
	go masterPipe.Serve()
	time.Sleep(1 * time.Second)
	slavePipe := pipes.NewSlave(
		networkType,
		c2Address,
		nil,
		true,
	)
	go slavePipe.Serve()
	time.Sleep(1 * time.Second)
	return masterPipe.ProxyListener, masterPipe.C2Listener
}

func StartHTTPServer(t *testing.T) net.Listener {
	httpListener, listenError := net.Listen(networkType, httpAddress)
	if listenError != nil {
		t.Fatal(listenError)
	}
	server := http.NewServeMux()
	server.HandleFunc("/big.txt",
		func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write(bytes.Repeat([]byte{'A'}, 0xFFFF))
		},
	)
	go http.Serve(httpListener, server)
	time.Sleep(1 * time.Second)
	return httpListener
}

func GetRequestRaw(url string) uint {
	_, e := http.Get(url)
	if e != nil {
		return FailedRequest
	}
	return Success
}

func GetRequestSocks5(url, username, password string) uint8 {
	dialer, err := proxy.SOCKS5(networkType, proxyAddress,
		&proxy.Auth{
			User:     username,
			Password: password,
		}, proxy.Direct)
	if err != nil {
		return FailedProxySetup
	}
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	httpTransport.Dial = dialer.Dial
	var req *http.Request
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return FailedRequest
	}
	_, err = httpClient.Do(req)
	if err != nil {
		return FailedRequest
	}
	return Success
}

func GetRequestHTTP(targetUrl, username, password string) uint8 {
	var (
		proxyUrl *url.URL
		err      error
	)
	if username != "" {
		proxyUrl, err = url.Parse(fmt.Sprintf("http://%s:%s@127.0.0.1:9050", username, password))
	} else {
		proxyUrl, err = url.Parse("http://127.0.0.1:9050")
	}
	httpClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	var req *http.Request
	req, err = http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		return FailedRequest
	}
	var response *http.Response
	response, err = httpClient.Do(req)
	if err != nil {
		return FailedRequest
	}
	if response.StatusCode == http.StatusForbidden {
		return FailedRequest
	}
	return Success
}
