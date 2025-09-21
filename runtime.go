package gvm

type Config struct {
	StackSize int
	HeapSize  int
}

type Runtime struct {
	program   Program
	registers map[Register]Operand
	stack     []Operand
	heap      []Operand
}

func NewRuntime(config *Config) *Runtime {
	regs := map[Register]Operand{
		// Specials
		PC: Integer(0),
		BP: Integer(0),
		SP: Integer(0),
		HP: Integer(0),
		// GeneralPurposes
		R1:   nil,
		R2:   nil,
		R3:   nil,
		ACM1: nil,
		ACM2: nil,
		// Flags
		ZF: Integer(0),
	}
	return &Runtime{
		program:   nil,
		registers: regs,
		stack:     make([]Operand, config.StackSize),
		heap:      make([]Operand, config.HeapSize),
	}
}
