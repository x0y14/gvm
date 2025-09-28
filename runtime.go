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

func NewRuntime(program Program, config *Config) *Runtime {
	regs := map[Register]Operand{
		// Specials
		PC: ProgramAddress(0),
		BP: BasePointer(0),
		SP: StackPointer(config.StackSize - 1),
		HP: HeapAddress(0),
		// GeneralPurposes
		R1:   nil,
		R2:   nil,
		R3:   nil,
		ACM1: nil,
		ACM2: nil,
		// Flags
		ZF: Bool(false),
	}
	return &Runtime{
		program:   program,
		registers: regs,
		stack:     make([]Operand, config.StackSize),
		heap:      make([]Operand, config.HeapSize),
	}
}
