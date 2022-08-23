package main

import (
	"errors"
	"github.com/shoriwe/gplasma/pkg/ast3"
	"github.com/shoriwe/gplasma/pkg/bytecode/assembler"
	"github.com/shoriwe/gplasma/pkg/vm"
	"os"
)

var (
	authenticationFailed     = errors.New("failed authentication")
	permissionDeniedInbound  = errors.New("permission denied for inbound connection")
	permissionDeniedOutbound = errors.New("permission denied for outbound connection")
	permissionDeniedListen   = errors.New("permission denied for listen address")
	permissionDeniedAccept   = errors.New("permission denied for accepted client")
)

func authScript(username, password []byte) []byte {
	bytecode, _ := assembler.Assemble(
		ast3.Program{
			&ast3.Call{
				Function: &ast3.Selector{
					X: &ast3.Identifier{
						Symbol: "-driver",
					},
					Identifier: &ast3.Identifier{
						Symbol: "-auth",
					},
				},
				Arguments: []ast3.Expression{
					&ast3.String{
						Contents: username,
					},
					&ast3.String{
						Contents: password,
					},
				},
			},
		},
	)
	return bytecode
}

func inboundScript(address []byte) []byte {
	bytecode, _ := assembler.Assemble(
		ast3.Program{
			&ast3.Call{
				Function: &ast3.Selector{
					X: &ast3.Identifier{
						Symbol: "-driver",
					},
					Identifier: &ast3.Identifier{
						Symbol: "-inbound",
					},
				},
				Arguments: []ast3.Expression{
					&ast3.String{
						Contents: address,
					},
				},
			},
		},
	)
	return bytecode
}

func outboundScript(address []byte) []byte {
	bytecode, _ := assembler.Assemble(
		ast3.Program{
			&ast3.Call{
				Function: &ast3.Selector{
					X: &ast3.Identifier{
						Symbol: "-driver",
					},
					Identifier: &ast3.Identifier{
						Symbol: "-outbound",
					},
				},
				Arguments: []ast3.Expression{
					&ast3.String{
						Contents: address,
					},
				},
			},
		},
	)
	return bytecode
}

func listenScript(address []byte) []byte {
	bytecode, _ := assembler.Assemble(
		ast3.Program{
			&ast3.Call{
				Function: &ast3.Selector{
					X: &ast3.Identifier{
						Symbol: "-driver",
					},
					Identifier: &ast3.Identifier{
						Symbol: "-listen",
					},
				},
				Arguments: []ast3.Expression{
					&ast3.String{
						Contents: address,
					},
				},
			},
		},
	)
	return bytecode
}

func acceptScript(address []byte) []byte {
	bytecode, _ := assembler.Assemble(
		ast3.Program{
			&ast3.Call{
				Function: &ast3.Selector{
					X: &ast3.Identifier{
						Symbol: "-driver",
					},
					Identifier: &ast3.Identifier{
						Symbol: "-accept",
					},
				},
				Arguments: []ast3.Expression{
					&ast3.String{
						Contents: address,
					},
				},
			},
		},
	)
	return bytecode
}

type filter struct {
	inbound  func(string) error
	outbound func(string) error
	listen   func(string) error
	accept   func(string) error
}

func (f *filter) Inbound(address string) error {
	if f.inbound == nil {
		return nil
	}
	return f.inbound(address)
}

func (f *filter) Outbound(address string) error {
	if f.outbound == nil {
		return nil
	}
	return f.outbound(address)
}

func (f *filter) Listen(address string) error {
	if f.listen == nil {
		return nil
	}
	return f.listen(address)
}

func (f *filter) Accept(address string) error {
	if f.accept == nil {
		return nil
	}
	return f.accept(address)
}

type Driver struct {
	plasma *vm.Plasma
}

func (d *Driver) Auth(username, password []byte) error {
	bytecode := authScript(username, password)
	resultChannel, errorChannel, _ := d.plasma.Execute(bytecode)
	if err := <-errorChannel; err != nil {
		return err
	}
	result := <-resultChannel
	if !result.Bool() {
		return authenticationFailed
	}
	return nil
}

func (d *Driver) Inbound(address string) error {
	bytecode := inboundScript([]byte(address))
	resultChannel, errorChannel, _ := d.plasma.Execute(bytecode)
	if err := <-errorChannel; err != nil {
		return err
	}
	result := <-resultChannel
	if !result.Bool() {
		return permissionDeniedInbound
	}
	return nil
}

func (d *Driver) Outbound(address string) error {
	bytecode := outboundScript([]byte(address))
	resultChannel, errorChannel, _ := d.plasma.Execute(bytecode)
	if err := <-errorChannel; err != nil {
		return err
	}
	result := <-resultChannel
	if !result.Bool() {
		return permissionDeniedOutbound
	}
	return nil
}

func (d *Driver) Listen(address string) error {
	bytecode := listenScript([]byte(address))
	resultChannel, errorChannel, _ := d.plasma.Execute(bytecode)
	if err := <-errorChannel; err != nil {
		return err
	}
	result := <-resultChannel
	if !result.Bool() {
		return permissionDeniedListen
	}
	return nil
}

func (d *Driver) Accept(address string) error {
	bytecode := acceptScript([]byte(address))
	resultChannel, errorChannel, _ := d.plasma.Execute(bytecode)
	if err := <-errorChannel; err != nil {
		return err
	}
	result := <-resultChannel
	if !result.Bool() {
		return permissionDeniedAccept
	}
	return nil
}

func (r *runner) loadDriver(scriptPath string) (*Driver, error) {
	scriptContents, readError := os.ReadFile(scriptPath)
	if readError != nil {
		return nil, readError
	}
	d := &Driver{
		plasma: initializeVM(),
	}
	// Setup the names used by the script
	_, errorChannel, _ := d.plasma.ExecuteString(string(scriptContents))
	//

	if err := <-errorChannel; err != nil {
		return nil, err
	}
	return d, nil
}

func initializeVM() *vm.Plasma {
	plasma := vm.NewVM(os.Stdin, os.Stdout, os.Stderr)
	plasma.Load("-driver", func(_ *vm.Plasma) *vm.Value { return plasma.Value() })
	plasma.Load("set_auth",
		func(p *vm.Plasma) *vm.Value {
			return p.NewBuiltInFunction(
				p.Symbols(),
				func(argument ...*vm.Value) (*vm.Value, error) {
					driver, _ := p.Symbols().Get("-driver")
					driver.Set("-auth", argument[0])
					return p.None(), nil
				},
			)
		},
	)
	plasma.Load("set_inbound",
		func(p *vm.Plasma) *vm.Value {
			return p.NewBuiltInFunction(
				p.Symbols(),
				func(argument ...*vm.Value) (*vm.Value, error) {
					driver, _ := p.Symbols().Get("-driver")
					driver.Set("-inbound", argument[0])
					return p.None(), nil
				},
			)
		},
	)
	plasma.Load("set_outbound",
		func(p *vm.Plasma) *vm.Value {
			return p.NewBuiltInFunction(
				p.Symbols(),
				func(argument ...*vm.Value) (*vm.Value, error) {
					driver, _ := p.Symbols().Get("-driver")
					driver.Set("-outbound", argument[0])
					return p.None(), nil
				},
			)
		},
	)
	plasma.Load("set_listen",
		func(p *vm.Plasma) *vm.Value {
			return p.NewBuiltInFunction(
				p.Symbols(),
				func(argument ...*vm.Value) (*vm.Value, error) {
					driver, _ := p.Symbols().Get("-driver")
					driver.Set("-listen", argument[0])
					return p.None(), nil
				},
			)
		},
	)
	plasma.Load("set_accept",
		func(p *vm.Plasma) *vm.Value {
			return p.NewBuiltInFunction(
				p.Symbols(),
				func(argument ...*vm.Value) (*vm.Value, error) {
					driver, _ := p.Symbols().Get("-driver")
					driver.Set("-accept", argument[0])
					return p.None(), nil
				},
			)
		},
	)
	return plasma
}
