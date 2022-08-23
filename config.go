package main

import (
	"crypto/tls"
	_ "embed"
	"fmt"
	"github.com/shoriwe/fullproxy/v3/internal/listeners"
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers"
	http2 "github.com/shoriwe/fullproxy/v3/internal/proxy/servers/http"
	port_forward "github.com/shoriwe/fullproxy/v3/internal/proxy/servers/port-forward"
	reverse2 "github.com/shoriwe/fullproxy/v3/internal/proxy/servers/reverse"
	socks52 "github.com/shoriwe/fullproxy/v3/internal/proxy/servers/socks5"
	pf_to_socks5 "github.com/shoriwe/fullproxy/v3/internal/proxy/servers/translation/pf-to-socks5"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	http3 "net/http"
	"net/url"
	"os"
	"strings"
)

type runner struct {
	drivers map[string]*Driver
}

func (r *runner) listenerConfig(listener Listener) (l listeners.Listener, err error) {
	var (
		tlsConfig, masterTLSConfig *tls.Config
	)
	if listener.TLS != nil {
		tlsConfig = &tls.Config{}
		for _, keyPem := range listener.TLS {
			split := strings.Split(keyPem, ":")
			if len(split) != 2 {
				return nil, KeyPemError
			}
			cert, certError := tls.LoadX509KeyPair(split[0], split[1])
			if certError != nil {
				return nil, certError
			}
			tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
		}
	}
	switch listener.Type {
	case "basic":
		return listeners.NewBindListener(listener.Network, listener.Address, tlsConfig)
	case "master":
		if listener.MasterTLS != nil {
			masterTLSConfig = &tls.Config{}
			for _, keyPem := range listener.MasterTLS {
				split := strings.Split(keyPem, ":")
				if len(split) != 2 {
					return nil, KeyPemError
				}
				cert, certError := tls.LoadX509KeyPair(split[0], split[1])
				if certError != nil {
					return nil, certError
				}
				masterTLSConfig.Certificates = append(masterTLSConfig.Certificates, cert)
			}
		}
		return listeners.NewMaster(
			listener.Network,
			listener.Address,
			tlsConfig,
			listener.MasterNetwork,
			listener.MasterAddress,
			masterTLSConfig,
		)
	case "slave":
		return listeners.NewSlave(
			listener.MasterNetwork,
			listener.MasterAddress,
			&tls.Config{
				InsecureSkipVerify: listener.SlaveTrust,
			},
		)
	default:
		return nil, fmt.Errorf(UnknownListenerTypeError, listener.Type)
	}
}

func (r *runner) socks5Config(p Protocol) (servers.Protocol, error) {
	var auth servers.AuthenticationMethod
	if p.Authentication != "" {
		driver, found := r.drivers[p.Authentication]
		if !found {
			return nil, fmt.Errorf(UnknownDriverError, p.Authentication)
		}
		auth = driver.Auth
	}
	return socks52.NewSocks5(auth), nil
}

func (r *runner) httpConfig(p Protocol) (servers.Protocol, error) {
	var auth servers.AuthenticationMethod
	if p.Authentication != "" {
		driver, found := r.drivers[p.Authentication]
		if !found {
			return nil, fmt.Errorf(UnknownDriverError, p.Authentication)
		}
		auth = driver.Auth
	}
	return http2.NewHTTP(auth), nil
}

func (r *runner) forwardConfig(p Protocol) (servers.Protocol, error) {
	var dialTlsConfig *tls.Config
	if p.DialTLS != nil {
		var certificates []tls.Certificate
		split := strings.Split(p.DialTLS.Certificate, ":")
		if len(split) != 2 {
			return nil, KeyPemError
		}
		cert, loadError := tls.LoadX509KeyPair(split[0], split[1])
		if loadError != nil {
			return nil, loadError
		}
		certificates = append(certificates, cert)
		dialTlsConfig = &tls.Config{
			InsecureSkipVerify: p.DialTLS.Trust,
			Certificates:       certificates,
		}
	}
	return port_forward.NewForward(
		p.TargetNetwork,
		p.TargetAddress,
		dialTlsConfig,
	), nil
}

func (r *runner) translateConfig(p Protocol) (servers.Protocol, error) {
	var userInfo *url.Userinfo = nil
	split := strings.Split(p.Credentials, ":")
	if len(split) == 2 {
		userInfo = url.UserPassword(split[0], split[1])
	}
	return pf_to_socks5.NewForwardToSocks5(
		p.ProxyNetwork,
		p.ProxyAddress,
		userInfo,
		p.TargetNetwork,
		p.TargetAddress,
	)
}

func (r *runner) parseHost(h *Host) (*reverse2.Host, error) {
	hh := &reverse2.Host{
		WebsocketReadBufferSize:  h.WebsocketReadBufferSize,
		WebsocketWriteBufferSize: h.WebsocketWriteBufferSize,
		Scheme:                   h.Scheme,
		URI:                      h.URI,
		Network:                  h.Network,
		Address:                  h.Address,
		TLSConfig:                nil,
	}
	if h.TLSConfig != nil {
		var certificates []tls.Certificate
		for _, pair := range h.TLSConfig.Certificates {
			split := strings.Split(pair, ":")
			if len(split) != 2 {
				return nil, KeyPemError
			}
			cert, loadError := tls.LoadX509KeyPair(split[0], split[1])
			if loadError != nil {
				return nil, loadError
			}
			certificates = append(certificates, cert)
		}
		hh.TLSConfig = &tls.Config{
			InsecureSkipVerify: h.TLSConfig.Trust,
			Certificates:       certificates,
		}
	}
	return hh, nil
}

func (r *runner) loadReverseHTTPHosts(p Protocol) (map[string]*reverse2.Target, error) {
	httpReverseHosts := map[string]*reverse2.Target{}
	for hostname, rawTarget := range p.HTTPHosts {
		target := &reverse2.Target{
			RequestHeader:  http3.Header{},
			ResponseHeader: http3.Header{},
			URI:            rawTarget.URI,
			CurrentHost:    0,
			Hosts:          nil,
		}
		for key, value := range rawTarget.RequestHeaders {
			target.RequestHeader[key] = []string{value}
		}
		for key, value := range rawTarget.ResponseHeaders {
			target.ResponseHeader[key] = []string{value}
		}
		for _, rawHost := range rawTarget.Pool {
			parsedHost, parseError := r.parseHost(rawHost)
			if parseError != nil {
				return nil, parseError
			}
			target.Hosts = append(target.Hosts, parsedHost)
		}
		httpReverseHosts[hostname] = target
	}
	return httpReverseHosts, nil
}

func (r *runner) loadReverseRawHosts(p Protocol) (rawReverseHosts []*reverse2.Host, err error) {
	for _, rawHost := range p.RawHosts {
		parsedHost, parseError := r.parseHost(rawHost)
		if parseError != nil {
			return nil, parseError
		}
		rawReverseHosts = append(rawReverseHosts, parsedHost)
	}
	return rawReverseHosts, nil
}

func (r *runner) loadSniffers(incoming, outgoing string) (i io.WriteCloser, o io.WriteCloser, err error) {
	if incoming != "" {
		i, err = os.Create(incoming)
		if err != nil {
			return nil, nil, err
		}
	}
	if outgoing != "" {
		o, err = os.Create(outgoing)
		if err != nil {
			return nil, nil, err
		}
	}
	return i, o, nil
}

func (r *runner) loadFilters(filters Filters) listeners.Filters {
	listenerFilter := &filter{}
	if f, found := r.drivers[filters.Inbound]; found {
		listenerFilter.inbound = f.Inbound
	}
	if f, found := r.drivers[filters.Outbound]; found {
		listenerFilter.outbound = f.Outbound
	}
	if f, found := r.drivers[filters.Listen]; found {
		listenerFilter.listen = f.Listen
	}
	if f, found := r.drivers[filters.Accept]; found {
		listenerFilter.accept = f.Accept
	}
	return listenerFilter
}

func (r *runner) serveListener(
	listenerName string,
	service Service,
	errorChan chan error,
) {
	logger := &log.Logger{}
	logger.SetOutput(os.Stderr)

	// Prepare listener
	listener, listenerConfigError := r.listenerConfig(service.Listener)
	if listenerConfigError != nil {
		errorChan <- listenerConfigError
		return
	}
	if slaveListener, ok := listener.(*listeners.Slave); ok {
		errorChan <- slaveListener.Serve()
		return
	}

	var (
		protocol            servers.Protocol
		protocolConfigError error
	)
	// Prepare protocol
	switch service.Protocol.Type {
	case "socks5":
		protocol, protocolConfigError = r.socks5Config(service.Protocol)
	case "http":
		protocol, protocolConfigError = r.httpConfig(service.Protocol)
	case "forward":
		protocol, protocolConfigError = r.forwardConfig(service.Protocol)
	case "translate":
		protocol, protocolConfigError = r.translateConfig(service.Protocol)
	case "reverse-http":
		reverseHTTPHosts, reverseHTTPHostsError := r.loadReverseHTTPHosts(service.Protocol)
		if reverseHTTPHostsError != nil {
			errorChan <- reverseHTTPHostsError
			return
		}
		protocol = reverse2.NewHTTP(reverseHTTPHosts)
	case "reverse-raw":
		reverseRawHosts, reverseRawHostsError := r.loadReverseRawHosts(service.Protocol)
		if reverseRawHostsError != nil {
			errorChan <- reverseRawHostsError
			return
		}
		protocol = reverse2.NewRaw(reverseRawHosts)
	default:
		protocolConfigError = fmt.Errorf(UnknownProtocolError, service.Protocol.Type)
	}
	if protocolConfigError != nil {
		errorChan <- protocolConfigError
		return
	}
	listener.SetFilters(r.loadFilters(service.Listener.Filters))
	var logFunc listeners.LogFunc = nil
	if service.Log != "" {
		f, createError := os.Create(service.Log)
		if createError != nil {
			errorChan <- createError
			return
		}
		defer f.Close()
		logger.SetOutput(f)
		logFunc = func(args ...interface{}) {
			logger.Print(args...)
		}
	}
	incoming, outgoing, snifferConfigError := r.loadSniffers(service.Sniff.Incoming, service.Sniff.Outgoing)
	if snifferConfigError != nil {
		errorChan <- snifferConfigError
		return
	}
	protocol.SetSniffers(incoming, outgoing)

	log.Println(listenerName, "Started")
	switch protocol.(type) {
	case servers.HTTPHandler:
		errorChan <- listeners.ServeHTTPHandler(listener, protocol.(servers.HTTPHandler), logFunc)
	default:
		errorChan <- listeners.Serve(listener, protocol, logFunc)
	}
}

func (r *runner) startConfig(c YAML) {
	var (
		err error
	)
	log.Print("Loading drivers")
	r.drivers = map[string]*Driver{}
	for name, scriptCode := range c.Drivers {
		r.drivers[name], err = r.loadDriver(scriptCode)
		if err != nil {
			printAndExit(err.Error(), 1)
		}
	}
	log.Print("Drivers loaded")
	log.Println("Starting listeners")
	serveError := make(chan error, 1)
	if c.InitOrder != nil {
		for _, listenerName := range c.InitOrder {
			listener, found := c.Services[listenerName]
			if !found {
				printAndExit(fmt.Sprintf("listener with name %s never set", listenerName), 1)
			}
			log.Println("-", listenerName)
			go r.serveListener(listenerName, listener, serveError)
		}
	} else {
		for listenerName, listener := range c.Services {
			log.Println("-", listenerName)
			go r.serveListener(listenerName, listener, serveError)
		}
	}

	log.Println("Listeners started")
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
	var c YAML
	unmarshalError := yaml.Unmarshal(configContents, &c)
	if unmarshalError != nil {
		printAndExit(unmarshalError.Error(), 1)
	}
	r := runner{}
	r.startConfig(c)
}
