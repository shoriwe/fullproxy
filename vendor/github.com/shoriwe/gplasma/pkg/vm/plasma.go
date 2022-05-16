package vm

import (
	"bufio"
	"crypto/rand"
	"hash"
	"hash/crc32"
	"io"
	"math/big"
)

const (
	polySize = 0xffffffff
)

type (
	ObjectLoader func(*Context, *Plasma) *Value
	Feature      map[string]ObjectLoader
)

type Plasma struct {
	currentId      int64
	BuiltInContext *Context
	Crc32Hash      hash.Hash32
	seed           uint64
	stdinScanner   *bufio.Scanner
	Stdin          io.Reader
	Stdout         io.Writer
	Stderr         io.Writer
}

func (p *Plasma) HashString(s string) int64 {
	_, hashingError := p.Crc32Hash.Write([]byte(s))
	if hashingError != nil {
		panic(hashingError)
	}
	hashValue := p.Crc32Hash.Sum32()
	p.Crc32Hash.Reset()
	return int64(hashValue)
}

func (p *Plasma) HashBytes(s []byte) int64 {
	_, hashingError := p.Crc32Hash.Write(s)
	if hashingError != nil {
		panic(hashingError)
	}
	hashValue := p.Crc32Hash.Sum32()
	p.Crc32Hash.Reset()
	return int64(hashValue)
}

func (p *Plasma) Seed() uint64 {
	return p.seed
}

func (p *Plasma) LoadFeature(symbolMap Feature) {
	for symbol, loader := range symbolMap {
		p.BuiltInContext.PeekSymbolTable().Set(symbol, loader(p.BuiltInContext, p))
	}
}

/*
	InitializeBytecode
	Loads the bytecode and clears the stack
*/

func (p *Plasma) StdInScanner() *bufio.Scanner {
	return p.stdinScanner
}

func (p *Plasma) StdIn() io.Reader {
	return p.Stdin
}

func (p *Plasma) StdOut() io.Writer {
	return p.Stdout
}

func (p *Plasma) StdErr() io.Writer {
	return p.Stderr
}

func (p *Plasma) BuiltInSymbols() *SymbolTable {
	return p.BuiltInContext.PeekSymbolTable()
}

func (p *Plasma) NextId() int64 {
	result := p.currentId
	p.currentId++
	return result
}

func newContext() *Context {
	result := &Context{
		ObjectStack: NewObjectStack(),
		SymbolStack: NewSymbolStack(),
		LastState:   NoState,
	}
	return result
}

func (p *Plasma) NewContext() *Context {
	result := newContext()
	symbols := NewSymbolTable(p.BuiltInContext.PeekSymbolTable())
	result.SymbolStack.Push(symbols)

	builtIn := &Value{
		id:              p.NextId(),
		typeName:        ValueName,
		class:           nil,
		subClasses:      nil,
		onDemandSymbols: map[string]OnDemandLoader{},
		IsBuiltIn:       true,
		symbols:         p.BuiltInContext.PeekSymbolTable(),
	}
	p.ObjectInitialize(true)(result, builtIn)
	symbols.Set("__built_in__", builtIn)

	self := &Value{
		id:              p.NextId(),
		typeName:        ValueName,
		class:           nil,
		onDemandSymbols: map[string]OnDemandLoader{},
		subClasses:      nil,
		symbols:         symbols,
	}
	p.ObjectInitialize(true)(result, self)
	symbols.Set(Self, self)
	return result
}

func NewPlasmaVM(stdin io.Reader, stdout io.Writer, stderr io.Writer) *Plasma {
	number, randError := rand.Int(rand.Reader, big.NewInt(polySize))
	if randError != nil {
		panic(randError)
	}
	vm := &Plasma{
		currentId:      1,
		BuiltInContext: newContext(),
		Crc32Hash:      crc32.New(crc32.MakeTable(polySize)),
		seed:           number.Uint64(),
		stdinScanner:   bufio.NewScanner(stdin),
		Stdin:          stdin,
		Stdout:         stdout,
		Stderr:         stderr,
	}
	vm.BuiltInContext.PushSymbolTable(NewSymbolTable(nil))
	vm.InitializeBuiltIn()
	return vm
}
