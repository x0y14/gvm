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
	for r.pc().Value() < len(r.program) {
		switch word := r.program[r.pc()]; word.(type) {
		case Opcode:
			err := r.do()
			if err != nil {
				return err
			}
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

func (r *Runtime) calcOffset(offset Offset) int {
	switch offset.(type) {
	case BpOffset:
		return r.bp().Value() + offset.(BpOffset).Value()
	case SpOffset:
		return r.sp().Value() + offset.(SpOffset).Value()
	default:
		panic("unknown offset")
	}
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
		case NOP:
			defer func() { r.set(PC, r.pc()+1+ProgramAddress(word.NumOperands())) }()
			return nil
		case MOV:
			defer func() { r.set(PC, r.pc()+1+ProgramAddress(word.NumOperands())) }()
			dst := r.program[r.pc()+1].(Operand)
			src := r.program[r.pc()+2].(Operand)
			switch dst.(type) {
			case Register: // ex) mov r1, ??
				switch src.(type) {
				case Register:
					r.registers[dst.(Register)] = r.registers[src.(Register)]
					return nil
				case Offset:
					r.registers[dst.(Register)] = r.stack[r.calcOffset(src.(Offset))]
					return nil
				case Immediate:
					r.registers[dst.(Register)] = src
					return nil
				default:
					return fmt.Errorf("unsupported mov src: %s", word.String())
				}
			case Offset:
				switch src.(type) {
				case Register:
					r.stack[r.calcOffset(dst.(Offset))] = r.registers[src.(Register)]
					return nil
				case Offset:
					r.stack[r.calcOffset(dst.(Offset))] = r.stack[r.calcOffset(src.(Offset))]
					return nil
				case Immediate:
					r.stack[r.calcOffset(dst.(Offset))] = src
					return nil
				default:
					return fmt.Errorf("unsupported mov src: %s", word.String())
				}
			default:
				return fmt.Errorf("unsupported mov dst: %s", word.String())
			}
		case PUSH:
			defer func() { r.set(PC, r.pc()+1+ProgramAddress(word.NumOperands())) }()
			src, ok := r.program[r.pc()+1].(Operand)
			if !ok {
				return fmt.Errorf("invalid push src: want=operand, got=%s", word.String())
			}
			switch src.(type) {
			case Register:
				r.push(r.registers[src.(Register)])
				return nil
			case Offset:
				r.push(r.stack[r.calcOffset(src.(Offset))])
				return nil
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
}
