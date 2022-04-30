package test

import (
	"bytes"
	"errors"
	"fmt"
	haochensocks5 "github.com/haochen233/socks5"
	"github.com/shoriwe/fullproxy/v3/internal/pipes"
	proxy2 "github.com/shoriwe/fullproxy/v3/internal/proxy"
	"github.com/shoriwe/fullproxy/v3/internal/proxy/socks5"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"
)

var (
	SampleMessage = []byte("HELLO")
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

func basicOutboundRule(addr net.Addr) error {
	tcpAddr, resolveError := net.ResolveTCPAddr("tcp", "google.com:443")
	if resolveError != nil {
		return resolveError
	}
	if addr.String() == tcpAddr.String() {
		return errors.New("host denied")
	}
	return nil
}

func basicInboundRule(addr net.Addr) error {
	if addr.(*net.TCPAddr).IP.String() == "127.0.0.1" {
		return errors.New("host denied")
	}
	return nil
}

func NewBindPipe(protocol proxy2.Protocol, inboundFilter, outboundFilter pipes.IOFilter) net.Listener {
	bindPipe := pipes.NewBindPipe(
		networkType, proxyAddress,
		protocol,
		nil,
		inboundFilter,
		outboundFilter,
	)
	go bindPipe.Serve()
	time.Sleep(2 * time.Second)
	return bindPipe.(*pipes.Bind).Server
}

func NewMasterSlave(protocol proxy2.Protocol, inboundFilter, outboundFilter pipes.IOFilter) (net.Listener, net.Listener) {
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
		outboundFilter,
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
	return masterPipe.(*pipes.Master).ProxyListener, masterPipe.(*pipes.Master).C2Listener
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

func Socks5BindTest(
	proxyAddress string,
	authMethod proxy2.AuthenticationMethod,
	auth map[haochensocks5.METHOD]haochensocks5.Authenticator,
	t *testing.T,
) {
	time.Sleep(10 * time.Second)
	proxyServer := NewBindPipe(
		socks5.NewSocks5(authMethod),
		nil,
		basicOutboundRule,
	)
	defer proxyServer.Close()
	socksClient := haochensocks5.Client{
		ProxyAddr:      proxyAddress,
		Auth:           auth,
		Timeout:        time.Minute,
		DisableSocks4A: true,
	}
	var (
		address         *haochensocks5.Address
		server          net.Conn
		connEstablished = new(sync.WaitGroup)
	)
	connEstablished.Add(1)
	go func() {
		var (
			secondError <-chan error
			err         error
		)
		address, secondError, server, err = socksClient.Bind(haochensocks5.Version5, "127.0.0.1:9999")
		if err != nil {
			t.Fatal(err)
		}
		connEstablished.Done()
		err = <-secondError
		if err != nil {
			t.Fatal(err)
		}
	}()
	connEstablished.Wait()
	client, connError := net.Dial("tcp", address.String())
	if connError != nil {
		t.Fatal(connError)
	}
	defer server.Close()
	defer client.Close()
	_, writeError := server.Write(SampleMessage)
	if writeError != nil {
		t.Fatal(writeError)
	}
	message := make([]byte, len(SampleMessage))
	_, readError := client.Read(message[:])
	if readError != nil {
		t.Fatal(readError)
	}
	if !bytes.Equal(message[:], SampleMessage) {
		t.Fatal(string(message[:]))
	}
}
