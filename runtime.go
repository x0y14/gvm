package gvm

import (
	"fmt"
)

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

func (r *Runtime) Run() error {
	for _, word := range r.program {
		switch word.(type) {
		case Opcode:
			return r.do()
		default:
			return fmt.Errorf("unsupported word: %s", word.String())
		}
	}
	return nil
}

func (r *Runtime) set(reg Register, operand Operand) {
	r.registers[reg] = operand
}

func (r *Runtime) pc() ProgramAddress {
	return r.registers[PC].(ProgramAddress)
}
func (r *Runtime) bp() BasePointer {
	return r.registers[BP].(BasePointer)
}
func (r *Runtime) sp() StackPointer {
	return r.registers[SP].(StackPointer)
}
func (r *Runtime) hp() HeapAddress {
	return r.registers[HP].(HeapAddress)
}

func (r *Runtime) push(operand Operand) {
	r.set(SP, r.sp()-1)
	if r.sp() < 0 {
		panic("stack overflow")
	}
	r.stack[r.sp()] = operand
}
func (r *Runtime) pop() Operand {
	v := r.stack[r.sp()]
	r.stack[r.sp()] = nil
	r.set(SP, r.sp()+1)
	return v
}

func (r *Runtime) do() error {
	switch word := r.program[r.pc()].(type) {
	case Opcode:
		switch word {
		case PUSH:
			defer func() { r.set(PC, r.pc()+1+ProgramAddress(word.NumOperands())) }()
			src, ok := r.program[r.pc()+1].(Operand)
			if !ok {
				return fmt.Errorf("invalid push src: want=operand, got=%s", word.String())
			}
			switch src.(type) {
			case Register:
			case Offset:
			case Immediate:
				r.push(src)
				return nil
			default:
				return fmt.Errorf("unsupported push src: %s", src.String())
			}
		case POP:
			defer func() { r.set(PC, r.pc()+1+ProgramAddress(word.NumOperands())) }()
			dst, ok := r.program[r.pc()+1].(Operand)
			if !ok {
				return fmt.Errorf("invalid pop dst: want=operand, got=%s", word.String())
			}
			switch dst := dst.(type) {
			case Register:
				r.registers[dst] = r.pop()
				return nil
			default:
				return fmt.Errorf("invalid pop dst: %s", word.String())
			}
		default:
			return fmt.Errorf("unsupported opcode: %s", word.String())
		}
	default:
		return fmt.Errorf("unsupported word: %s", word.String())
	}
	return nil
}
