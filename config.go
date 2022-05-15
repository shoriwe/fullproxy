package main

import (
	_ "embed"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

var (
	//go:embed docs/yaml.md
	yamlMarkdown string
	//go:embed docs/scripting.md
	scriptingMarkdown string
)

type ConfigFile struct {
	Drivers   map[string]string `yaml:"drivers"`
	Listeners map[string]struct {
		Config struct {
			Type          string   `yaml:"type"`
			Network       string   `yaml:"network"`
			Address       string   `yaml:"address"`
			TLS           []string `yaml:"tls"`
			MasterNetwork string   `yaml:"master-network"`
			MasterAddress string   `yaml:"master-address"`
			MasterTLS     []string `yaml:"master-tls"`
			SlaveTrust    bool     `yaml:"slave-trust"`
		} `yaml:"config"`
		Protocol struct {
			Type    string `yaml:"type"`
			Filters struct {
				Inbound  string `yaml:"inbound"`
				Outbound string `yaml:"outbound"`
				Listen   string `yaml:"listen"`
				Accept   string `yaml:"accept"`
			} `yaml:"filters"`
			Authentication string `yaml:"authentication"`
			TargetNetwork  string `yaml:"target-network"`
			TargetAddress  string `yaml:"target-address"`
			ProxyNetwork   string `yaml:"proxy-network"`
			ProxyAddress   string `yaml:"proxy-address"`
			Translation    string `yaml:"translation"`
		} `yaml:"protocol"`
	} `yaml:"listeners"`
}

func startConfig(c ConfigFile) {
	fmt.Printf("%+v\n", c)
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
	startConfig(c)
}
