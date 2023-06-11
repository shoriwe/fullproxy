package http

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/shoriwe/fullproxy/v3/utils/network"
)

func TestNewMux(t *testing.T) {
	l := network.ListenAny()
	defer l.Close()
	NewMux(l)
	expect := httpexpect.Default(t, "http://"+l.Addr().String())
	expect.GET(EchoRoute).Expect().Status(http.StatusOK)
}
