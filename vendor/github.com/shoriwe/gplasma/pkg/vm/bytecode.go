package vm

type Instruction struct {
	OpCode uint8
	Line   int
}

func NewInstruction(opCode uint8) Instruction {
	return Instruction{
		OpCode: opCode,
	}
}

type Code struct {
	Instruction Instruction
	Value       interface{}
	Line        int
}

func NewCode(opCode uint8, line int, value interface{}) *Code {
	return &Code{
		Instruction: NewInstruction(opCode),
		Value:       value,
		Line:        line,
	}
}

type Bytecode struct {
	instructions []*Code
	length       int
	index        int
}

func (bytecode *Bytecode) HasNext() bool {
	return bytecode.index < bytecode.length
}

func (bytecode *Bytecode) Peek() *Code {
	return bytecode.instructions[bytecode.index]
}

func (bytecode *Bytecode) Next() *Code {
	result := bytecode.instructions[bytecode.index]
	bytecode.index++
	return result
}

func (bytecode *Bytecode) NextN(n int) []*Code {
	result := bytecode.instructions[bytecode.index : bytecode.index+n]
	bytecode.index += n
	return result
}

func (bytecode *Bytecode) Jump(offset int) {
	bytecode.index += offset
}

func NewBytecodeFromArray(codes []*Code) *Bytecode {
	return &Bytecode{
		instructions: codes,
		length:       len(codes),
		index:        0,
	}
}
