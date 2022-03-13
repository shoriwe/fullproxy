package httpprotocol

import (
	"bytes"
	"fmt"
	"github.com/shoriwe/FullProxy/v3/internal/pipes"
	http2 "github.com/shoriwe/FullProxy/v3/internal/proxy/http"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"testing"
)

const (
	testUrl      = "http://127.0.0.1:8080/big.txt"
	networkType  = "tcp"
	proxyAddress = "127.0.0.1:9050"
)

var (
	httpListener net.Listener
	bindPipe     *pipes.Bind
	masterPipe   *pipes.Master
	slavePipe    *pipes.Slave
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

const (
	Success = iota
	FailedRequest
)

func getRequest(targetUrl string, username, password string) uint8 {
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
	content, _ := io.ReadAll(response.Body)
	fmt.Println(string(content))
	return Success
}

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
	bindPipe = pipes.NewBindPipe(
		networkType,
		proxyAddress,
		http2.NewHTTP(nil, nil, nil),
		nil,
		nil,
	)
	go bindPipe.Serve()
}

//// Test No auth

func TestNoAuthHTTPRequest(t *testing.T) {
	if getRequest(testUrl, "", "") != Success {
		t.Fatal(testUrl)
	}
}

func TestCloseNoAuthPipe(t *testing.T) {
	closingError := bindPipe.Server.Close()
	if closingError != nil {
		t.Fatal(closingError)
	}
}

//// Test Auth

func TestUserPasswordAuthHTTPInit(t *testing.T) {
	bindPipe = pipes.NewBindPipe(
		networkType,
		proxyAddress,
		http2.NewHTTP(
			basicAuthFunc,
			nil,
			nil,
		),
		nil,
		nil,
	)
	go bindPipe.Serve()
}

func TestUsernamePasswordHTTPRequest(t *testing.T) {
	if getRequest(testUrl, "sulcud", "password") != Success {
		t.Fatal(testUrl)
	}
}

func TestUsernamePasswordWithNoAuthHTTPRequest(t *testing.T) {
	if getRequest(testUrl, "", "") != FailedRequest {
		t.Fatal("Auth bypassed")
	}
}

func TestCloseUserPasswordAuthPipe(t *testing.T) {
	closingError := bindPipe.Server.Close()
	if closingError != nil {
		t.Fatal(closingError)
	}
}

//// Test inbound rules
func TestInboundRulesHTTPInit(t *testing.T) {
	bindPipe = pipes.NewBindPipe(
		networkType,
		proxyAddress,
		http2.NewHTTP(
			basicAuthFunc,
			nil,
			nil,
		),
		nil,
		basicInboundRule,
	)
	go bindPipe.Serve()
}

func TestInvalidInboundHTTPRequest(t *testing.T) {
	if getRequest(testUrl, "sulcud", "password") != FailedRequest {
		t.Fatal("Bypassed inbound")
	}
}

func TestCloseInboundRulesPipe(t *testing.T) {
	closingError := bindPipe.Server.Close()
	if closingError != nil {
		t.Fatal(closingError)
	}
}

//// Test outbound rules

func TestOutboundRulesHTTPInit(t *testing.T) {
	bindPipe = pipes.NewBindPipe(
		networkType,
		proxyAddress,
		http2.NewHTTP(
			basicAuthFunc,
			nil,
			basicOutboundRule,
		),
		nil,
		nil,
	)
	go bindPipe.Serve()
}

func TestInvalidOutboundHTTPRequest(t *testing.T) {
	if getRequest("google.com", "sulcud", "password") == Success {
		t.Fatal("Bypassed outbound")
	}
}

func TestOutboundSuccessHTTPRequest(t *testing.T) {
	if getRequest(testUrl, "sulcud", "password") != Success {
		t.Fatal(testUrl)
	}
}

func TestCloseOutboundRulesPipe(t *testing.T) {
	closingError := bindPipe.Server.Close()
	if closingError != nil {
		t.Fatal(closingError)
	}
}

// Test master slave

//// Test No auth

func TestNoAuthMasterSlaveInitialization(t *testing.T) {
	masterPipe = pipes.NewMaster(
		"tcp",
		"127.0.0.1:9051",
		"127.0.0.1:9050",
		nil,
		nil,
		http2.NewHTTP(nil, nil, nil),
	)
	go masterPipe.Serve()
	slavePipe = pipes.NewSlave(
		"tcp",
		"127.0.0.1:9051",
		nil,
	)
	go slavePipe.Serve()
}

func TestNoAuthMasterSlaveHTTPRequest(t *testing.T) {
	if getRequest(testUrl, "", "") != Success {
		t.Fatal(testUrl)
	}
}

func TestCloseNoAuthMasterSlavePipe(t *testing.T) {
	closingError := masterPipe.ProxyListener.Close()
	if closingError != nil {
		t.Fatal(closingError)
	}
	closingError = masterPipe.MasterConnection.Close()
	if closingError != nil {
		t.Fatal(closingError)
	}
	closingError = slavePipe.MasterConnection.Close()
	if closingError != nil {
		t.Fatal(closingError)
	}
}

//// Test Auth

func TestUsernamePasswordMasterSlaveInitialization(t *testing.T) {
	masterPipe = pipes.NewMaster(
		"tcp",
		"127.0.0.1:9051",
		"127.0.0.1:9050",
		nil,
		nil,
		http2.NewHTTP(basicAuthFunc, nil, nil),
	)
	go masterPipe.Serve()
	slavePipe = pipes.NewSlave(
		"tcp",
		"127.0.0.1:9051",
		nil,
	)
	go slavePipe.Serve()
}

func TestUsernamePasswordMasterSlaveHTTPRequest(t *testing.T) {
	if getRequest(testUrl, "sulcud", "password") != Success {
		t.Fatal(testUrl)
	}
}

func TestCloseUsernamePasswordAuthMasterSlavePipe(t *testing.T) {
	closingError := masterPipe.ProxyListener.Close()
	if closingError != nil {
		t.Fatal(closingError)
	}
	closingError = masterPipe.MasterConnection.Close()
	if closingError != nil {
		t.Fatal(closingError)
	}
	closingError = slavePipe.MasterConnection.Close()
	if closingError != nil {
		t.Fatal(closingError)
	}
}

//// Test inbound rules

func TestInboundMasterSlaveInitialization(t *testing.T) {
	masterPipe = pipes.NewMaster(
		"tcp",
		"127.0.0.1:9051",
		"127.0.0.1:9050",
		nil,
		basicInboundRule,
		http2.NewHTTP(basicAuthFunc, nil, nil),
	)
	go masterPipe.Serve()
	slavePipe = pipes.NewSlave(
		"tcp",
		"127.0.0.1:9051",
		nil,
	)
	go slavePipe.Serve()
}

func TestInboundMasterSlaveHTTPRequest(t *testing.T) {
	if getRequest(testUrl, "sulcud", "password") != FailedRequest {
		t.Fatal("Bypassed inbound rule")
	}
}

func TestCloseInboundMasterSlavePipe(t *testing.T) {
	closingError := masterPipe.ProxyListener.Close()
	if closingError != nil {
		t.Fatal(closingError)
	}
	closingError = masterPipe.MasterConnection.Close()
	if closingError != nil {
		t.Fatal(closingError)
	}
	closingError = slavePipe.MasterConnection.Close()
	if closingError != nil {
		t.Fatal(closingError)
	}
}

//// Test outbound rules

func TestOutboundMasterSlaveInitialization(t *testing.T) {
	masterPipe = pipes.NewMaster(
		"tcp",
		"127.0.0.1:9051",
		"127.0.0.1:9050",
		nil,
		nil,
		http2.NewHTTP(basicAuthFunc, nil, basicOutboundRule),
	)
	go masterPipe.Serve()
	slavePipe = pipes.NewSlave(
		"tcp",
		"127.0.0.1:9051",
		nil,
	)
	go slavePipe.Serve()
}

func TestInvalidOutboundMasterSlaveHTTPRequest(t *testing.T) {
	if getRequest("https://google.com", "sulcud", "password") != FailedRequest {
		t.Fatal("Bypassed outbound rule")
	}
}

func TestOutboundMasterSlaveHTTPRequest(t *testing.T) {
	if getRequest(testUrl, "sulcud", "password") != Success {
		t.Fatal(testUrl)
	}
}

func TestCloseOutboundMasterSlavePipe(t *testing.T) {
	closingError := masterPipe.ProxyListener.Close()
	if closingError != nil {
		t.Fatal(closingError)
	}
	closingError = masterPipe.MasterConnection.Close()
	if closingError != nil {
		t.Fatal(closingError)
	}
	closingError = slavePipe.MasterConnection.Close()
	if closingError != nil {
		t.Fatal(closingError)
	}
}

// Finally, close the http server

func TestClose(t *testing.T) {
	_ = httpListener.Close()
}
