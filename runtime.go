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

func (r *Runtime) pc() Integer {
	return r.registers[PC].(Integer)
}
func (r *Runtime) bp() Integer {
	return r.registers[BP].(Integer)
}
func (r *Runtime) sp() Integer {
	return r.registers[SP].(Integer)
}
func (r *Runtime) hp() Integer {
	return r.registers[HP].(Integer)
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
			defer func() { r.set(PC, r.pc()+1+Integer(word.NumOperands())) }()
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
			defer func() { r.set(PC, r.pc()+1+Integer(word.NumOperands())) }()
			dst, ok := r.program[r.pc()+1].(Operand)
			if !ok {
				return fmt.Errorf("invalid pop dst: want=operand, got=%s", word.String())
			}
			switch dst.(type) {
			case Register:
				r.registers[dst.(Register)] = r.pop()
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
