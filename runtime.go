package gvm

import (
	"fmt"

	"github.com/x0y14/gvm/word"
)

type Config struct {
	StackSize int
	HeapSize  int
}

type Runtime struct {
	program   []word.Word
	registers map[word.Register]word.Storable
	stack     []word.Storable
	heap      []word.Storable
}

func NewRuntime(program []word.Word, config *Config) *Runtime {
	regs := map[word.Register]word.Storable{
		// Specials
		word.PC: word.Address(0),
		word.BP: word.Address(0),
		word.SP: word.Address(config.StackSize - 1),
		word.HP: word.Address(0),
		// GeneralPurposes
		word.R1:   nil,
		word.R2:   nil,
		word.R3:   nil,
		word.ACM1: nil,
		word.ACM2: nil,
		// Flags
		word.ZF: word.Bool(false),
	}
	return &Runtime{
		program:   program,
		registers: regs,
		stack:     make([]word.Storable, config.StackSize),
		heap:      make([]word.Storable, config.HeapSize),
	}
}

func (r *Runtime) pc() word.Address {
	return r.registers[word.PC].(word.Address)
}
func (r *Runtime) bp() word.Address {
	return r.registers[word.BP].(word.Address)
}
func (r *Runtime) sp() word.Address {
	return r.registers[word.SP].(word.Address)
}
func (r *Runtime) hp() word.Address {
	return r.registers[word.HP].(word.Address)
}
func (r *Runtime) setSpecial(register word.SpecialRegister, addr word.Address) {
	r.registers[register] = addr
}

func (r *Runtime) push(storable word.Storable) {
	r.setSpecial(word.SP, r.sp()-1)
	if r.sp() < 0 {
		panic("stack overflow")
	}
	r.stack[r.sp()] = storable
}
func (r *Runtime) pop() word.Storable {
	v := r.stack[r.sp()]
	r.stack[r.sp()] = nil
	r.setSpecial(word.SP, r.sp()+1)
	return v
}

func (r *Runtime) alloc(size int) word.Address {
	baseAddr := r.hp()
	if len(r.heap) <= int(r.hp())+size {
		panic("heap: out of memory")
	}
	r.setSpecial(word.HP, word.Address(int(r.hp())+size))
	return baseAddr
}
func (r *Runtime) store(addr word.Address, storable word.Storable) {
	if 0 <= int(addr) && int(addr) < len(r.heap) {
		r.heap[addr] = storable
		return
	}
	panic("heap: out of bounds") // 不法侵入
}
func (r *Runtime) load(addr word.Address) word.Storable {
	if 0 <= int(addr) && int(addr) < len(r.heap) {
		return r.heap[addr]
	}
	panic("heap: out of bounds")
}

func (r *Runtime) add(dst word.Register, src word.Immediate) error {
	switch dst.(type) {
	case word.SpecialRegister, word.GeneralPurposeRegister:
		// srcがintであることを確認
		imm, ok := src.(word.Integer)
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
		case word.Address:
			sum = int(r.registers[dst].(word.Address)) + imm.Value()
			r.registers[dst] = word.Address(sum)
		case word.Integer:
			sum = int(r.registers[dst].(word.Integer)) + imm.Value()
			r.registers[dst] = word.Integer(sum)
		default:
			return fmt.Errorf("add: dst register %s dose not contain Integer/Address (got %T)", dst.String(), curt)
		}
		// zf
		r.registers[word.ZF] = word.Bool(int(sum) == 0)
		return nil
	default:
		return fmt.Errorf("add: unsupported dst: %T", dst)
	}
}
func (r *Runtime) sub(dst word.Register, src word.Immediate) error {
	switch dst.(type) {
	case word.SpecialRegister, word.GeneralPurposeRegister:
		// srcがintであることを確認
		imm, ok := src.(word.Integer)
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
		case word.Address:
			sum = int(r.registers[dst].(word.Address)) - imm.Value()
			r.registers[dst] = word.Address(sum)
		case word.Integer:
			sum = int(r.registers[dst].(word.Integer)) - imm.Value()
			r.registers[dst] = word.Integer(sum)
		default:
			return fmt.Errorf("sub: dst register %s dose not contain Integer/Address (got %T)", dst.String(), curt)
		}
		// zf
		r.registers[word.ZF] = word.Bool(int(sum) == 0)
		return nil
	default:
		return fmt.Errorf("sub: unsupported dst: %T", dst)
	}
}
func (r *Runtime) Run() error {
	return nil
}
