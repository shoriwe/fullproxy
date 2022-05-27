package main

type (
	YAML struct {
		InitOrder []string           `yaml:"init-order"`
		Drivers   map[string]string  `yaml:"drivers"`
		Services  map[string]Service `yaml:"services"`
	}

	Service struct {
		Log   string `yaml:"log"`
		Sniff struct {
			Incoming string `yaml:"incoming"`
			Outgoing string `yaml:"outgoing"`
		} `yaml:"sniff"`
		Listener Listener `yaml:"listener"`
		Protocol Protocol `yaml:"protocol"`
	}

	Listener struct {
		Type          string   `yaml:"type"`
		Network       string   `yaml:"network"`
		Address       string   `yaml:"address"`
		TLS           []string `yaml:"tls"`
		Filters       Filters  `yaml:"filters"`
		MasterNetwork string   `yaml:"master-network"`
		MasterAddress string   `yaml:"master-address"`
		MasterTLS     []string `yaml:"master-tls"`
		SlaveTrust    bool     `yaml:"slave-trust"`
	}

	Filters struct {
		Inbound  string `yaml:"inbound"`
		Outbound string `yaml:"outbound"`
		Listen   string `yaml:"listen"`
		Accept   string `yaml:"accept"`
	}

	Protocol struct {
		Type           string           `yaml:"type"`
		Authentication string           `yaml:"authentication"`
		DialTLS        *DialTLS         `yaml:"dial-tls"`
		TargetNetwork  string           `yaml:"target-network"`
		TargetAddress  string           `yaml:"target-address"`
		ProxyNetwork   string           `yaml:"proxy-network"`
		ProxyAddress   string           `yaml:"proxy-address"`
		Translation    string           `yaml:"translation"`
		Credentials    string           `yaml:"credentials"`
		RawHosts       map[string]*Host `yaml:"raw-hosts"`
		HTTPHosts      map[string]struct {
			URI             string            `yaml:"uri"`
			ResponseHeaders map[string]string `yaml:"response-headers"`
			RequestHeaders  map[string]string `yaml:"request-headers"`
			Pool            map[string]*Host  `yaml:"pool"`
		} `yaml:"http-hosts"`
	}

	DialTLS struct {
		Trust       bool   `yaml:"trust"`
		Certificate string `yaml:"certificate"`
	}

	Host struct {
		WebsocketReadBufferSize  int            `yaml:"websocket-read-buffer-size"`
		WebsocketWriteBufferSize int            `yaml:"websocket-write-buffer-size"`
		Scheme                   string         `yaml:"scheme"`
		URI                      string         `yaml:"uri"`
		Network                  string         `yaml:"network"`
		Address                  string         `yaml:"address"`
		TLSConfig                *HostTLSConfig `yaml:"tls-config"`
	}

	HostTLSConfig struct {
		Trust        bool     `yaml:"trust"`
		Certificates []string `yaml:"certificates"`
	}
)
