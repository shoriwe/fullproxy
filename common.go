package main

import (
	_ "embed"
	"errors"
	"github.com/shoriwe/fullproxy/v3/internal/listeners"
	reverse2 "github.com/shoriwe/fullproxy/v3/internal/proxy/servers/reverse"
	"net/url"
)

var (
	//go:embed docs/yaml.md
	yamlMarkdown string
	//go:embed docs/scripting.md
	scriptingMarkdown string
)

type hostsSlice struct {
	contents []*reverse2.Host
}

func (h *hostsSlice) String() string {
	return "{}"
}

func (h *hostsSlice) Set(ss string) error {
	u, err := url.Parse(ss)
	if err != nil {
		return err
	}
	h.contents = append(h.contents, &reverse2.Host{
		Network: u.Scheme,
		Address: u.Host,
	})
	return nil
}

func createListener(listen, master string) (listeners.Listener, error) {
	if listen == "" {
		return nil, errors.New("no listen address provided")
	}
	listenURL, parseError := url.Parse(listen)
	if parseError != nil {
		return nil, parseError
	}
	if master == "" {
		return listeners.NewBindListener(listenURL.Scheme, listenURL.Host, nil)
	}
	masterURL, parseMasterError := url.Parse(master)
	if parseMasterError != nil {
		printAndExit(parseMasterError.Error(), 1)
	}
	return listeners.NewMaster(listenURL.Scheme, listenURL.Host, nil, masterURL.Scheme, masterURL.Host, nil)
}
