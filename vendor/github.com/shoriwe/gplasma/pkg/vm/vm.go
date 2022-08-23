package vm

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/compiler"
	"io"
)

type (
	Loader func(plasma *Plasma) *Value
	Plasma struct {
		Stdin             io.Reader
		Stdout, Stderr    io.Writer
		rootSymbols       *Symbols
		onDemand          map[string]func(self *Value) *Value
		true, false, none *Value
		value             *Value
		string            *Value
		bytes             *Value
		bool              *Value
		noneType          *Value
		int               *Value
		float             *Value
		array             *Value
		tuple             *Value
		hash              *Value
		function          *Value
		class             *Value
	}
)

func (plasma *Plasma) Symbols() *Symbols {
	return plasma.rootSymbols
}

func (plasma *Plasma) True() *Value {
	return plasma.true
}

func (plasma *Plasma) False() *Value {
	return plasma.false
}

func (plasma *Plasma) None() *Value {
	return plasma.none
}

func (plasma *Plasma) Value() *Value {
	return plasma.value
}

func (plasma *Plasma) String() *Value {
	return plasma.string
}

func (plasma *Plasma) Bytes() *Value {
	return plasma.bytes
}

func (plasma *Plasma) Bool() *Value {
	return plasma.bool
}

func (plasma *Plasma) NoneType() *Value {
	return plasma.noneType
}

func (plasma *Plasma) Int() *Value {
	return plasma.int
}

func (plasma *Plasma) Float() *Value {
	return plasma.float
}

func (plasma *Plasma) Array() *Value {
	return plasma.array
}

func (plasma *Plasma) Tuple() *Value {
	return plasma.tuple
}

func (plasma *Plasma) Hash() *Value {
	return plasma.hash
}

func (plasma *Plasma) Function() *Value {
	return plasma.function
}

func (plasma *Plasma) Class() *Value {
	return plasma.class
}

func (plasma *Plasma) executeCtx(ctx *context) {
	defer func() {
		err := recover()
		if err != nil {
			ctx.err <- fmt.Errorf("execution error: %v", err)
		} else {
			ctx.err <- nil
		}
		ctx.result <- ctx.register
		return
	}()
	for ctx.hasNext() {
		select {
		case <-ctx.stop:
			return
		default:
			plasma.do(ctx)
		}
	}
}

func (plasma *Plasma) Load(symbol string, loader Loader) {
	plasma.rootSymbols.Set(symbol, loader(plasma))
}

func (plasma *Plasma) Execute(bytecode []byte) (result chan *Value, err chan error, stop chan struct{}) {
	// Create new context
	ctx := plasma.newContext(bytecode)
	ctx.result = make(chan *Value, 1)
	ctx.err = make(chan error, 1)
	ctx.stop = make(chan struct{}, 1)
	// Execute bytecode with context
	go plasma.executeCtx(ctx)
	return ctx.result, ctx.err, ctx.stop
}

func (plasma *Plasma) ExecuteString(scriptCode string) (result chan *Value, err chan error, stop chan struct{}) {
	bytecode, compileError := compiler.Compile(scriptCode)
	// Create new context
	ctx := plasma.newContext(bytecode)
	ctx.result = make(chan *Value, 1)
	ctx.err = make(chan error, 1)
	ctx.stop = make(chan struct{}, 1)
	if compileError != nil {
		ctx.result <- nil
		ctx.err <- compileError
	} else {
		// Execute bytecode with context
		go plasma.executeCtx(ctx)
	}
	return ctx.result, ctx.err, ctx.stop
}

func NewVM(stdin io.Reader, stdout, stderr io.Writer) *Plasma {
	plasma := &Plasma{
		Stdin:       stdin,
		Stdout:      stdout,
		Stderr:      stderr,
		rootSymbols: NewSymbols(nil),
	}
	plasma.init()
	return plasma
}
