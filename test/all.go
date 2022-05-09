package test

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/shoriwe/fullproxy/v3/internal/listeners"
	socks52 "github.com/shoriwe/fullproxy/v3/internal/proxy/clients/socks5"
	proxy2 "github.com/shoriwe/fullproxy/v3/internal/proxy/servers"
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers/socks5"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
)

var (
	SampleMessage = []byte("HELLO")
	SampleAddress = &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 9999,
	}
)

const (
	testUrl          = "http://127.0.0.1:8080/big.txt"
	testUrlIPv6      = "http://[::1]:8080/big.txt"
	networkType      = "tcp"
	httpIPv6Address  = "[::1]:8080"
	httpAddress      = "127.0.0.1:8080"
	proxyAddress     = "127.0.0.1:9050"
	ipv6ProxyAddress = "[::1]:9050"
	c2Address        = "127.0.0.1:9051"
	Success          = iota
	FailedProxySetup
	FailedRequest
)

func basicAuthFunc(username []byte, password []byte) error {
	if bytes.Equal(username, []byte("sulcud")) &&
		bytes.Equal(password, []byte("password")) {
		return nil
	}
	return errors.New("auth failed")
}

type BasicInbound struct{}

func (b *BasicInbound) Outbound(address string) error {
	return nil
}

func (b *BasicInbound) Listen(address string) error {
	return nil
}

func (b *BasicInbound) Accept(address string) error {
	return nil
}

func (b *BasicInbound) Inbound(address string) error {
	host, _, splitError := net.SplitHostPort(address)
	if splitError != nil {
		return splitError
	}
	if host == "127.0.0.1" {
		return errors.New("host denied")
	}
	return nil
}

type BasicOutbound struct{}

func (b *BasicOutbound) Outbound(address string) error {
	host, _, splitError := net.SplitHostPort(address)
	if splitError != nil {
		return splitError
	}
	if host == "google.com" {
		return errors.New("host denied")
	}
	return nil
}

func (b *BasicOutbound) Listen(address string) error {
	return nil
}

func (b *BasicOutbound) Accept(address string) error {
	return nil
}

func (b *BasicOutbound) Inbound(address string) error {
	return nil
}

type BasicListen struct{}

func (b *BasicListen) Outbound(address string) error {
	return nil
}

func (b *BasicListen) Listen(address string) error {
	host, _, splitError := net.SplitHostPort(address)
	if splitError != nil {
		return splitError
	}
	if host == "127.0.0.1" {
		return errors.New("listen in localhost denied")
	}
	return nil
}

func (b *BasicListen) Accept(address string) error {
	return nil
}

func (b *BasicListen) Inbound(address string) error {
	return nil
}

type BasicAccept struct{}

func (b *BasicAccept) Outbound(address string) error {
	return nil
}

func (b *BasicAccept) Listen(address string) error {
	return nil
}

func (b *BasicAccept) Accept(address string) error {
	host, _, splitError := net.SplitHostPort(address)
	if splitError != nil {
		return splitError
	}
	if host == "127.0.0.1" {
		return errors.New("connections from localhost denied")
	}
	return nil
}

func (b *BasicAccept) Inbound(address string) error {
	return nil
}

func NewBindHandler(handler proxy2.HTTPHandler, filters listeners.Filters) net.Listener {
	if filters == nil {
		filters = &listeners.NoFilter{}
	}
	listener, listenError := listeners.NewBindListener(networkType, proxyAddress, nil)
	if listenError != nil {
		panic(listenError)
	}
	listener.SetFilters(filters)
	go listeners.ServeHTTPHandler(listener, handler, nil)
	return listener
}

func NewMasterSlaveHandler(handler proxy2.HTTPHandler, filters listeners.Filters) (net.Listener, net.Listener) {
	if filters == nil {
		filters = &listeners.NoFilter{}
	}
	master, listenError := listeners.NewMaster(
		networkType,
		proxyAddress,
		nil,
		networkType,
		c2Address,
		nil,
	)
	if listenError != nil {
		panic(listenError)
	}
	master.SetFilters(filters)
	go listeners.ServeHTTPHandler(master, handler, nil)
	time.Sleep(1 * time.Second)
	slave, listenError := listeners.NewSlave(networkType, c2Address, nil)
	if listenError != nil {
		panic(listenError)
	}
	go slave.Serve()
	time.Sleep(1 * time.Second)
	return master, slave
}

func NewBindPipe(protocol proxy2.Protocol, filters listeners.Filters) net.Listener {
	if filters == nil {
		filters = &listeners.NoFilter{}
	}
	bindPipe, listenError := listeners.NewBindListener(
		networkType, proxyAddress,
		nil,
	)
	if listenError != nil {
		panic(listenError)
	}
	bindPipe.SetFilters(filters)
	go listeners.Serve(bindPipe, protocol, nil)
	return bindPipe
}

func NewMasterSlave(protocol proxy2.Protocol, filters listeners.Filters) (net.Listener, net.Listener) {
	if filters == nil {
		filters = &listeners.NoFilter{}
	}
	masterPipe, listenError := listeners.NewMaster(
		networkType,
		proxyAddress,
		nil,
		networkType,
		c2Address,
		nil,
	)
	if listenError != nil {
		panic(listenError)
	}
	masterPipe.SetFilters(filters)
	go listeners.Serve(masterPipe, protocol, nil)
	time.Sleep(1 * time.Second)
	slavePipe, listenError := listeners.NewSlave(
		networkType,
		c2Address,
		nil,
	)
	if listenError != nil {
		panic(listenError)
	}
	go slavePipe.Serve()
	time.Sleep(1 * time.Second)
	return masterPipe.(*listeners.Master).ProxyListener, masterPipe.(*listeners.Master).C2Listener
}

func StartIPv4HTTPServer(t *testing.T) net.Listener {
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

func StartIPv6HTTPServer(t *testing.T) net.Listener {
	httpListener, listenError := net.Listen(networkType, httpIPv6Address)
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

func Socks6BindSucceed(
	proxyAddress string,
	authMethod proxy2.AuthenticationMethod,
	filters listeners.Filters,
	username, password string,
	t *testing.T,
) {
	time.Sleep(10 * time.Second)
	proxyServer := NewBindPipe(
		socks5.NewSocks5(authMethod),
		filters,
	)
	defer proxyServer.Close()
	socksClient := socks52.Socks5{
		Network:  "tcp",
		Address:  proxyAddress,
		Username: username,
		Password: password,
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
			t.Fatal(writeError)
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
	if readError != nil {
		t.Fatal(readError)
	}
	if !bytes.Equal(response, SampleMessage) {
		t.Fatal(string(response))
	}
}
