package main

import (
	"crypto/tls"
	_ "embed"
	"errors"
	"fmt"
	"github.com/shoriwe/fullproxy/v3/internal/listeners"
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers"
	http2 "github.com/shoriwe/fullproxy/v3/internal/proxy/servers/http"
	http_hosts "github.com/shoriwe/fullproxy/v3/internal/proxy/servers/http-hosts"
	port_forward "github.com/shoriwe/fullproxy/v3/internal/proxy/servers/port-forward"
	reverse2 "github.com/shoriwe/fullproxy/v3/internal/proxy/servers/reverse"
	socks52 "github.com/shoriwe/fullproxy/v3/internal/proxy/servers/socks5"
	pf_to_socks5 "github.com/shoriwe/fullproxy/v3/internal/proxy/servers/translation/pf-to-socks5"
	"gopkg.in/yaml.v3"
	http3 "net/http"
	"os"
	"strings"
)

type runner struct {
	drivers map[string]*Driver
}

func (r *runner) serveListener(
	c struct {
		Config   ListenerConfig `yaml:"config"`
		Protocol ProtocolConfig `yaml:"protocol"`
	},
	errorChan chan error,
) {
	if c.Config.Type == "slave" {
		slaveListener, newSlaveError := listeners.NewSlave(
			c.Config.MasterNetwork,
			c.Config.MasterAddress,
			&tls.Config{
				InsecureSkipVerify: c.Config.SlaveTrust,
			},
		)
		if newSlaveError != nil {
			errorChan <- newSlaveError
			return
		}
		errorChan <- slaveListener.Serve()
		return
	}

	// Prepare listener
	var (
		l                listeners.Listener
		listenError      error
		protocol         servers.Protocol
		newProtocolError error
		tlsConfig        *tls.Config = nil
		masterTLSConfig  *tls.Config = nil
	)
	if c.Config.TLS != nil {
		tlsConfig = &tls.Config{}
		for _, keyPem := range c.Config.TLS {
			split := strings.Split(keyPem, ":")
			if len(split) != 2 {
				errorChan <- errors.New("expected key:pem in tls list")
				return
			}
			cert, certError := tls.LoadX509KeyPair(split[0], split[1])
			if certError != nil {
				errorChan <- certError
				return
			}
			tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
		}
	}
	if c.Config.MasterTLS != nil {
		masterTLSConfig = &tls.Config{}
		for _, keyPem := range c.Config.MasterTLS {
			split := strings.Split(keyPem, ":")
			if len(split) != 2 {
				errorChan <- errors.New("expected key:pem in tls list")
				return
			}
			cert, certError := tls.LoadX509KeyPair(split[0], split[1])
			if certError != nil {
				errorChan <- certError
				return
			}
			masterTLSConfig.Certificates = append(masterTLSConfig.Certificates, cert)
		}
	}
	switch c.Config.Type {
	case "basic":
		l, listenError = listeners.NewBindListener(c.Config.Network, c.Config.Address, tlsConfig)
	case "master":
		l, listenError = listeners.NewMaster(
			c.Config.Network,
			c.Config.Address,
			tlsConfig,
			c.Config.MasterNetwork,
			c.Config.MasterAddress,
			masterTLSConfig,
		)
	}
	if listenError != nil {
		errorChan <- listenError
		return
	}
	var (
		httpReverseHosts map[string]*reverse2.Target
		rawReverseHosts  []*reverse2.Host
	)
	for _, h := range c.Protocol.RawHosts {
		rawReverseHosts = append(rawReverseHosts, &h)
	}
	for hostname, t := range c.Protocol.HTTPHosts {
		tt := &reverse2.Target{
			RequestHeader:  http3.Header{},
			ResponseHeader: http3.Header{},
			Path:           t.Path,
			CurrentTarget:  0,
			Targets:        nil,
		}
		for key, value := range t.RequestHeaders {
			tt.RequestHeader[key] = []string{value}
		}
		for key, value := range t.ResponseHeaders {
			tt.ResponseHeader[key] = []string{value}
		}
		for _, h := range t.Pool {
			tt.Targets = append(tt.Targets, &h)
		}
		httpReverseHosts[hostname] = tt
	}
	// Prepare protocol
	switch c.Protocol.Type {
	case "socks5":
		driver, found := r.drivers[c.Protocol.Authentication]
		if !found {
			errorChan <- fmt.Errorf("unknown driver %s", driver)
		}
		protocol = socks52.NewSocks5(driver.Auth)
	case "http":
		driver, found := r.drivers[c.Protocol.Authentication]
		if !found {
			errorChan <- fmt.Errorf("unknown driver %s", driver)
		}
		protocol = http2.NewHTTP(driver.Auth)
	case "forward":
		protocol = port_forward.NewForward(
			c.Protocol.TargetNetwork,
			c.Protocol.TargetAddress,
		)
	case "translate":
		protocol, newProtocolError = pf_to_socks5.NewForwardToSocks5(
			c.Protocol.ProxyNetwork,
			c.Protocol.ProxyNetwork,
			nil, // TODO: Fix this, let the yaml config this field
			c.Protocol.TargetNetwork,
			c.Protocol.TargetAddress,
		)
	case "http-hosts":
		protocol = http_hosts.NewHosts()
	case "reverse-http":
		protocol = reverse2.NewHTTP(httpReverseHosts)
	case "reverse-raw":
		protocol = reverse2.NewRaw(rawReverseHosts)
	}
	if newProtocolError != nil {
		errorChan <- newProtocolError
		return
	}
	listenerFilter := &filter{}
	if f, found := r.drivers[c.Config.Filters.Inbound]; found {
		listenerFilter.inbound = f.inbound
	}
	if f, found := r.drivers[c.Config.Filters.Outbound]; found {
		listenerFilter.outbound = f.outbound
	}
	if f, found := r.drivers[c.Config.Filters.Listen]; found {
		listenerFilter.listen = f.listen
	}
	if f, found := r.drivers[c.Config.Filters.Accept]; found {
		listenerFilter.accept = f.accept
	}
	l.SetFilters(listenerFilter)
	switch protocol.(type) {
	case servers.HTTPHandler:
		errorChan <- listeners.ServeHTTPHandler(l, protocol.(servers.HTTPHandler), nil) // TODO: Set logging file
	default:
		errorChan <- listeners.Serve(l, protocol, nil) // TODO: Set logging file
	}
}

func (r *runner) startConfig(c ConfigFile) {
	var (
		err error
	)
	r.drivers = map[string]*Driver{}
	for name, script := range c.Drivers {
		r.drivers[name], err = loadDriver(script)
		if err != nil {
			printAndExit(err.Error(), 1)
		}
	}
	serveError := make(chan error, 1)
	for _, listener := range c.Listeners {
		go r.serveListener(listener, serveError)
	}
	printAndExit((<-serveError).Error(), 1)
}

func config() {
	if len(os.Args) != 3 {
		printAndExit(fmt.Sprintf("Usage: fullproxy config YAML_CONFIG\n\n%s\n\n%s\n", yamlMarkdown, scriptingMarkdown), 0)
	}
	configContents, readError := os.ReadFile(os.Args[2])
	if readError != nil {
		printAndExit(readError.Error(), 1)
	}
	var c ConfigFile
	unmarshalError := yaml.Unmarshal(configContents, &c)
	if unmarshalError != nil {
		printAndExit(unmarshalError.Error(), 1)
	}
	r := runner{}
	r.startConfig(c)
}
