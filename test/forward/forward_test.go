package socks5

import (
	"bytes"
	"github.com/shoriwe/FullProxy/pkg/Pipes"
	"github.com/shoriwe/FullProxy/pkg/Pipes/Reverse/Master"
	"github.com/shoriwe/FullProxy/pkg/Pipes/Reverse/Slave"
	"github.com/shoriwe/FullProxy/pkg/Proxies/PortForward"
	"log"
	"net"
	"net/http"
	"testing"
)

const (
	testUrl     = "http://127.0.0.1:8080/big.txt"
	networkType = "tcp"
)

var (
	httpListener net.Listener
	bindPipe     *Pipes.Bind
	masterPipe   *Master.Master
	slavePipe    *Slave.Slave
)

func basicInboundRule(host string) bool {
	if host == "127.0.0.1" {
		return false
	}
	return true
}

const (
	Success = iota
	FailedProxySetup
	FailedRequest
	proxyAddress = "127.0.0.1:8081"
)

func getRequest(url string) uint8 {
	_, getError := http.Get(url)
	if getError != nil {
		return FailedRequest
	}
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

func TestBindInitialization(t *testing.T) {
	bindPipe = Pipes.NewBindPipe(
		networkType,
		proxyAddress,
		PortForward.NewForward(networkType, "127.0.0.1:8080", nil),
		nil,
		nil,
	)
	go bindPipe.Serve()
}

func TestBindRequest(t *testing.T) {
	if getRequest(testUrl) != Success {
		t.Fatal(testUrl)
	}
}

func TestCloseBindPipe(t *testing.T) {
	closingError := bindPipe.Server.Close()
	if closingError != nil {
		t.Fatal(closingError)
	}
}

// Test Master/slave

func TestMasterSlaveInitialization(t *testing.T) {
	masterPipe = Master.NewMaster(
		networkType,
		"127.0.0.1:8082",
		proxyAddress,
		nil,
		nil,
		PortForward.NewForward(networkType, "127.0.0.1:8080", nil),
	)
	go masterPipe.Serve()

	slavePipe = Slave.NewSlave(
		networkType,
		"127.0.0.1:8082",
		nil,
	)
	go slavePipe.Serve()
}

func TestMasterSlaveRequest(t *testing.T) {
	if getRequest(testUrl) != Success {
		t.Fatal(testUrl)
	}
}

func TestCloseMasterSlavePipe(t *testing.T) {
	closingError := masterPipe.C2Listener.Close()
	if closingError != nil {
		t.Fatal(closingError)
	}

	closingError = masterPipe.ProxyListener.Close()
	if closingError != nil {
		t.Fatal(closingError)
	}

	closingError = slavePipe.MasterConnection.Close()
	if closingError != nil {
		t.Fatal(closingError)
	}
}

// Finally, close the HTTP server

func TestClose(t *testing.T) {
	_ = httpListener.Close()
}