package http

import (
	"net"
	"net/http"
)

const (
	EchoRoute = "/echo"
	EchoMsg   = "ECHO"
)

func NewMux(l net.Listener) {
	mux := http.NewServeMux()
	mux.HandleFunc(EchoRoute, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(EchoMsg))
	})
	go http.Serve(l, mux)
}
