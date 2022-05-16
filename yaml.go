package main

import reverse2 "github.com/shoriwe/fullproxy/v3/internal/proxy/servers/reverse"

type ListenerConfig struct {
	Type    string   `yaml:"type"`
	Network string   `yaml:"network"`
	Address string   `yaml:"address"`
	TLS     []string `yaml:"tls"`
	Filters struct {
		Inbound  string `yaml:"inbound"`
		Outbound string `yaml:"outbound"`
		Listen   string `yaml:"listen"`
		Accept   string `yaml:"accept"`
	} `yaml:"filters"`
	MasterNetwork string   `yaml:"master-network"`
	MasterAddress string   `yaml:"master-address"`
	MasterTLS     []string `yaml:"master-tls"`
	SlaveTrust    bool     `yaml:"slave-trust"`
}

type ProtocolConfig struct {
	Type           string                   `yaml:"type"`
	Authentication string                   `yaml:"authentication"`
	TargetNetwork  string                   `yaml:"target-network"`
	TargetAddress  string                   `yaml:"target-address"`
	ProxyNetwork   string                   `yaml:"proxy-network"`
	ProxyAddress   string                   `yaml:"proxy-address"`
	Translation    string                   `yaml:"translation"`
	RawHosts       map[string]reverse2.Host `yaml:"raw-hosts"'`
	HTTPHosts      map[string]struct {
		Path            string                   `yaml:"path"`
		ResponseHeaders map[string]string        `yaml:"response-headers"`
		RequestHeaders  map[string]string        `yaml:"request-headers"`
		Pool            map[string]reverse2.Host `yaml:"pool"`
	} `yaml:"http-hosts"`
}

type ConfigFile struct {
	Drivers   map[string]string `yaml:"drivers"`
	Listeners map[string]struct {
		Config   ListenerConfig `yaml:"config"`
		Protocol ProtocolConfig `yaml:"protocol"`
	} `yaml:"listeners"`
}
