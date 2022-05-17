package main

import (
	"errors"
	"github.com/shoriwe/gplasma"
	"github.com/shoriwe/gplasma/pkg/vm"
	"os"
	"sync"
)

var (
	authenticationFailed     = errors.New("failed authentication")
	permissionDeniedInbound  = errors.New("permission denied for inbound connection")
	permissionDeniedOutbound = errors.New("permission denied for outbound connection")
	permissionDeniedListen   = errors.New("permission denied for listen address")
	permissionDeniedAccept   = errors.New("permission denied for accepted client")
)

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
	mutex    *sync.Mutex
	auth     func(string, string) error
	inbound  func(string) error
	outbound func(string) error
	listen   func(string) error
	accept   func(string) error
}

func (d *Driver) Auth(username, password []byte) error {
	if d.auth == nil {
		return nil
	}
	return d.auth(string(username), string(password))
}

func (d *Driver) Inbound(address string) error {
	if d.inbound == nil {
		return nil
	}
	return d.inbound(address)
}

func (d *Driver) Outbound(address string) error {
	if d.outbound == nil {
		return nil
	}
	return d.outbound(address)
}

func (d *Driver) Listen(address string) error {
	if d.listen == nil {
		return nil
	}
	return d.listen(address)
}

func (d *Driver) Accept(address string) error {
	if d.accept == nil {
		return nil
	}
	return d.accept(address)
}

func interpretAsString(v *vm.Value, ctx *vm.Context, p *vm.Plasma) string {
	if v.IsTypeById(vm.StringId) {
		return v.String
	}
	toString, getToStringException := v.Get(p, ctx, "ToString")
	if getToStringException != nil {
		return "Failed to get ToString of object"
	}
	result, succeed := p.CallFunction(ctx, toString)
	if !succeed {
		return "Failed to call ToString"
	}
	return result.String
}

func (d *Driver) feature() vm.Feature {
	return vm.Feature{
		"SetAuth": func(context *vm.Context, plasma *vm.Plasma) *vm.Value {
			return plasma.NewFunction(
				context, true, context.PeekSymbolTable(),
				vm.NewBuiltInFunction(1,
					func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
						d.auth = func(username string, password string) error {
							d.mutex.Lock()
							defer d.mutex.Unlock()
							result, succeed := plasma.CallFunction(context, arguments[0],
								plasma.NewString(context, false, username),
								plasma.NewString(context, false, password),
							)
							if !succeed {
								return errors.New(interpretAsString(result, context, plasma))
							}
							if result.Bool {
								return nil
							}
							return authenticationFailed
						}
						return plasma.GetNone(), true
					},
				),
			)
		},
		"SetInbound": func(context *vm.Context, plasma *vm.Plasma) *vm.Value {
			return plasma.NewFunction(
				context, true, context.PeekSymbolTable(),
				vm.NewBuiltInFunction(1,
					func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
						d.inbound = func(address string) error {
							d.mutex.Lock()
							defer d.mutex.Unlock()
							result, succeed := plasma.CallFunction(context, arguments[0],
								plasma.NewString(context, false, address),
							)
							if !succeed {
								return errors.New(interpretAsString(result, context, plasma))
							}
							if result.Bool {
								return nil
							}
							return permissionDeniedInbound
						}
						return plasma.GetNone(), true
					},
				),
			)
		},
		"SetOutbound": func(context *vm.Context, plasma *vm.Plasma) *vm.Value {
			return plasma.NewFunction(
				context, true, context.PeekSymbolTable(),
				vm.NewBuiltInFunction(1,
					func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
						d.outbound = func(address string) error {
							d.mutex.Lock()
							defer d.mutex.Unlock()
							result, succeed := plasma.CallFunction(context, arguments[0],
								plasma.NewString(context, false, address),
							)
							if !succeed {
								return errors.New(interpretAsString(result, context, plasma))
							}
							if result.Bool {
								return nil
							}
							return permissionDeniedOutbound
						}
						return plasma.GetNone(), true
					},
				),
			)
		},
		"SetListen": func(context *vm.Context, plasma *vm.Plasma) *vm.Value {
			return plasma.NewFunction(
				context, true, context.PeekSymbolTable(),
				vm.NewBuiltInFunction(1,
					func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
						d.listen = func(address string) error {
							d.mutex.Lock()
							defer d.mutex.Unlock()
							result, succeed := plasma.CallFunction(context, arguments[0],
								plasma.NewString(context, false, address),
							)
							if !succeed {
								return errors.New(interpretAsString(result, context, plasma))
							}
							if result.Bool {
								return nil
							}
							return permissionDeniedListen
						}
						return plasma.GetNone(), true
					},
				),
			)
		},
		"SetAccept": func(context *vm.Context, plasma *vm.Plasma) *vm.Value {
			return plasma.NewFunction(
				context, true, context.PeekSymbolTable(),
				vm.NewBuiltInFunction(1,
					func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
						d.accept = func(address string) error {
							d.mutex.Lock()
							defer d.mutex.Unlock()
							result, succeed := plasma.CallFunction(context, arguments[0],
								plasma.NewString(context, false, address),
							)
							if !succeed {
								return errors.New(interpretAsString(result, context, plasma))
							}
							if result.Bool {
								return nil
							}
							return permissionDeniedAccept
						}
						return plasma.GetNone(), true
					},
				),
			)
		},
	}
}

func loadDriver(script string) (*Driver, error) {
	d := &Driver{
		mutex: &sync.Mutex{},
	}
	v := gplasma.NewVirtualMachine()
	v.LoadFeature(d.feature())

	scriptContents, readError := os.ReadFile(script)
	if readError != nil {
		return nil, readError
	}
	_, succeed := v.ExecuteMain(string(scriptContents))
	if !succeed {
		return nil, errors.New("driver script execution error")
	}
	return d, nil
}