package gvm

import (
	"fmt"

	"github.com/x0y14/gvm/program"
)

type Config struct {
	StackSize int
	HeapSize  int
}

type Runtime struct {
	program   []program.Word
	registers map[program.Register]program.Storable
	stack     []program.Storable
	heap      []program.Storable
}

func NewRuntime(program Program, config *Config) *Runtime {
	regs := map[program.Register]program.Storable{
		// Specials
		program.PC: program.Address(0),
		program.BP: program.Address(0),
		program.SP: program.Address(config.StackSize - 1),
		program.HP: program.Address(0),
		// GeneralPurposes
		program.R1:   nil,
		program.R2:   nil,
		program.R3:   nil,
		program.ACM1: nil,
		program.ACM2: nil,
		// Flags
		program.ZF: program.Bool(false),
	}
	return &Runtime{
		program:   program,
		registers: regs,
		stack:     make([]program.Storable, config.StackSize),
		heap:      make([]program.Storable, config.HeapSize),
	}
}

func (r *Runtime) pc() program.Address {
	return r.registers[program.PC].(program.Address)
}
func (r *Runtime) bp() program.Address {
	return r.registers[program.BP].(program.Address)
}
func (r *Runtime) sp() program.Address {
	return r.registers[program.SP].(program.Address)
}
func (r *Runtime) hp() program.Address {
	return r.registers[program.HP].(program.Address)
}
func (r *Runtime) setSpecial(register program.SpecialRegister, addr program.Address) {
	r.registers[register] = addr
}

func (r *Runtime) push(storable program.Storable) {
	r.setSpecial(program.SP, r.sp()-1)
	if r.sp() < 0 {
		panic("stack overflow")
	}
	r.stack[r.sp()] = storable
}
func (r *Runtime) pop() program.Storable {
	v := r.stack[r.sp()]
	r.stack[r.sp()] = nil
	r.setSpecial(program.SP, r.sp()+1)
	return v
}

func (r *Runtime) alloc(size int) program.Address {
	baseAddr := r.hp()
	if len(r.heap) <= int(r.hp())+size {
		panic("heap: out of memory")
	}
	r.setSpecial(program.HP, program.Address(int(r.hp())+size))
	return baseAddr
}
func (r *Runtime) store(addr program.Address, storable program.Storable) {
	if 0 <= int(addr) && int(addr) < len(r.heap) {
		r.heap[addr] = storable
		return
	}
	panic("heap: out of bounds") // 不法侵入
}
func (r *Runtime) load(addr program.Address) program.Storable {
	if 0 <= int(addr) && int(addr) < len(r.heap) {
		return r.heap[addr]
	}
	panic("heap: out of bounds")
}

func (r *Runtime) add(dst program.Register, src program.Immediate) error {
	switch dst.(type) {
	case program.SpecialRegister, program.GeneralPurposeRegister:
		// srcがintであることを確認
		imm, ok := src.(program.Integer)
		if !ok {
			return fmt.Errorf("add: src must be integer, got %T", src)
		}
		// dst
		curt := r.registers[dst]
		if curt == nil {
			return fmt.Errorf("add: dst register %s is nil", dst.String())
		}
		var sum int
		switch curt.(type) {
		case program.Address:
			sum = int(r.registers[dst].(program.Address)) + imm.Value()
			r.registers[dst] = program.Address(sum)
		case program.Integer:
			sum = int(r.registers[dst].(program.Integer)) + imm.Value()
			r.registers[dst] = program.Integer(sum)
		default:
			return fmt.Errorf("add: dst register %s dose not contain Integer/Address (got %T)", dst.String(), curt)
		}
		// zf
		r.registers[program.ZF] = program.Bool(int(sum) == 0)
		return nil
	default:
		return fmt.Errorf("add: unsupported dst: %T", dst)
	}
}
func (r *Runtime) sub(dst program.Register, src program.Immediate) error {
	switch dst.(type) {
	case program.SpecialRegister, program.GeneralPurposeRegister:
		// srcがintであることを確認
		imm, ok := src.(program.Integer)
		if !ok {
			return fmt.Errorf("sub: src must be integer, got %T", src)
		}
		// dst
		curt := r.registers[dst]
		if curt == nil {
			return fmt.Errorf("sub: dst register %s is nil", dst.String())
		}
		var sum int
		switch curt.(type) {
		case program.Address:
			sum = int(r.registers[dst].(program.Address)) - imm.Value()
			r.registers[dst] = program.Address(sum)
		case program.Integer:
			sum = int(r.registers[dst].(program.Integer)) - imm.Value()
			r.registers[dst] = program.Integer(sum)
		default:
			return fmt.Errorf("sub: dst register %s dose not contain Integer/Address (got %T)", dst.String(), curt)
		}
		// zf
		r.registers[program.ZF] = program.Bool(int(sum) == 0)
		return nil
	default:
		return fmt.Errorf("sub: unsupported dst: %T", dst)
	}
}
func (r *Runtime) Run() error {
	return nil
}
