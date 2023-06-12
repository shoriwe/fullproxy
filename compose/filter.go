package compose

import (
	"net"
	"regexp"

	"github.com/shoriwe/fullproxy/v4/filter"
	"github.com/shoriwe/fullproxy/v4/utils/network"
)

type PortRange struct {
	From *int `yaml:"from,omitempty" json:"from,omitempty"`
	To   *int `yaml:"to,omitempty" json:"to,omitempty"`
}

type Match struct {
	Host  *string    `yaml:"host,omitempty" json:"host,omitempty"`
	Port  *int       `yaml:"port,omitempty" json:"port,omitempty"`
	Range *PortRange `yaml:"portRange,omitempty" json:"portRange,omitempty"`
}

func (m *Match) Compile() (filter.Match, error) {
	var (
		err  error
		host *regexp.Regexp
		port int = -1
		from int = -1
		to   int = -1
	)
	if m.Host != nil {
		host, err = regexp.Compile(*m.Host)
	}
	if m.Port != nil {
		port = *m.Port
	}
	if m.Range != nil {
		if m.Range.From != nil {
			from = *m.Range.From
		}
		if m.Range.To != nil {
			to = *m.Range.To
		}
	}
	f := filter.Match{
		Host:      host,
		Port:      port,
		PortRange: [2]int{from, to},
	}
	return f, err
}

type Filter struct {
	Whitelist []Match `yaml:"whitelist,omitempty" json:"whitelist,omitempty"`
	Blacklist []Match `yaml:"blacklist,omitempty" json:"blacklist,omitempty"`
}

func (f *Filter) Listener(l net.Listener) (net.Listener, error) {
	var whitelist, blacklist []filter.Match
	for _, white := range f.Whitelist {
		compiled, err := white.Compile()
		if err != nil {
			return nil, err
		}
		whitelist = append(whitelist, compiled)
	}
	for _, black := range f.Blacklist {
		compiled, err := black.Compile()
		if err != nil {
			return nil, err
		}
		blacklist = append(blacklist, compiled)
	}
	ll := &filter.Listener{
		Listener:  l,
		Whitelist: whitelist,
		Blacklist: blacklist,
	}
	return ll, nil
}

func (f *Filter) DialFunc(dialFunc network.DialFunc) (*filter.DialFunc, error) {
	var whitelist, blacklist []filter.Match
	for _, white := range f.Whitelist {
		compiled, err := white.Compile()
		if err != nil {
			return nil, err
		}
		whitelist = append(whitelist, compiled)
	}
	for _, black := range f.Blacklist {
		compiled, err := black.Compile()
		if err != nil {
			return nil, err
		}
		blacklist = append(blacklist, compiled)
	}
	df := &filter.DialFunc{
		DialFunc:  dialFunc,
		Whitelist: whitelist,
		Blacklist: blacklist,
	}
	return df, nil
}
