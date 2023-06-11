package http

import (
	"testing"

	"github.com/shoriwe/fullproxy/v3/utils/network"
)

func TestNewMux(t *testing.T) {
	l := network.ListenAny()
	defer l.Close()
	checkCh := make(chan struct{}, 1)
	defer close(checkCh)
	go func() {
		NewMux(l)
		checkCh <- struct{}{}
	}()
	<-checkCh
}
