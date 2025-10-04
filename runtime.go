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

func (r *Runtime) solve(offset word.Offset) (word.Storable, error) {
	switch offset.Target {
	case word.SP:
		return r.stack[r.sp()+word.Address(offset.Diff)], nil
	case word.BP:
		return r.stack[r.bp()+word.Address(offset.Diff)], nil
	default:
		return nil, fmt.Errorf("solve: unsupported offset: %s", offset.String())
	}
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

func (r *Runtime) alloc(size int) {
	baseAddr := r.hp()
	if len(r.heap) <= int(r.hp())+size {
		panic("heap: out of memory")
	}
	r.setSpecial(word.HP, word.Address(int(r.hp())+size))
	r.push(baseAddr)
}
func (r *Runtime) store(dst word.Address, src word.Storable) {
	if 0 <= int(dst) && int(dst) < len(r.heap) {
		r.heap[dst] = src
		return
	}
	panic("heap: out of bounds") // 不法侵入
}
func (r *Runtime) load(dst word.Register, src word.Address) {
	if 0 <= int(src) && int(src) < len(r.heap) {
		r.registers[dst] = r.heap[src]
		return
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

func (r *Runtime) mov(dst word.Register, src word.Storable) error {
	switch dst.(type) {
	case word.SpecialRegister, word.GeneralPurposeRegister:
		// src がレジスタの場合：レジスタの値をコピー
		if srcReg, ok := src.(word.Register); ok {
			val := r.registers[srcReg]
			if val == nil {
				return fmt.Errorf("mov: src register %s is nil", srcReg.String())
			}
			// 特殊レジスタなら Address に適合させる
			if _, isSpecial := dst.(word.SpecialRegister); isSpecial {
				switch v := val.(type) {
				case word.Address:
					r.registers[dst] = v
					return nil
				case word.Integer:
					r.registers[dst] = word.Address(v.Value())
					return nil
				default:
					return fmt.Errorf("mov: cannot mov %T to special register %s", v, dst.String())
				}
			}
			// 汎用レジスタへはそのままコピー
			r.registers[dst] = val
			return nil
		}

		// src が即値の場合：型に応じて格納（特殊レジスタは Address 必須、Integer は Address に変換可）
		switch v := src.(type) {
		case word.Address:
			r.registers[dst] = v
			return nil
		case word.Integer:
			if _, isSpecial := dst.(word.SpecialRegister); isSpecial {
				r.registers[dst] = word.Address(v.Value())
			} else {
				r.registers[dst] = v
			}
			return nil
		case word.Bool:
			if _, isSpecial := dst.(word.SpecialRegister); isSpecial {
				return fmt.Errorf("mov: cannot mov Bool to special register %s", dst.String())
			}
			r.registers[dst] = v
			return nil
		default:
			return fmt.Errorf("mov: unsupported src type %T", v)
		}
	default:
		return fmt.Errorf("mov: unsupported dst: %T", dst)
	}
}

func (r *Runtime) exec() error {
	switch w := r.program[r.pc()].(type) {
	case word.Opcode:
		switch w {
		case word.NOP:
			defer func() { r.setSpecial(word.PC, word.Address(int(r.pc())+1+w.NumOperands())) }()
			return nil
		case word.MOV:
			defer func() { r.setSpecial(word.PC, word.Address(int(r.pc())+1+w.NumOperands())) }()
			dst := r.program[r.pc()+1]
			src := r.program[r.pc()+2]
			if _, ok := dst.(word.Register); !ok {
				return fmt.Errorf("mov: unsupported dst: %s", dst.String())
			}
			switch src.(type) {
			case word.Register:
				return r.mov(dst.(word.Register), r.registers[src.(word.Register)])
			case word.Offset:
				v, err := r.solve(dst.(word.Offset))
				if err != nil {
					return err
				}
				return r.mov(dst.(word.Register), v)
			case word.Immediate:
				return r.mov(dst.(word.Register), src.(word.Immediate))
			default:
				return fmt.Errorf("mov: unsupported src: %s", src.String())
			}
		case word.PUSH:
			defer func() { r.setSpecial(word.PC, word.Address(int(r.pc())+1+w.NumOperands())) }()
			src := r.program[r.pc()+1]
			switch src.(type) {
			case word.Register:
				r.push(r.registers[src.(word.Register)])
				return nil
			case word.Offset:
				v, err := r.solve(src.(word.Offset))
				if err != nil {
					return err
				}
				r.push(v)
				return nil
			case word.Immediate:
				r.push(src.(word.Immediate))
				return nil
			default:
				return fmt.Errorf("push: unsuported src: %s", src.String())
			}
		case word.POP:
			defer func() { r.setSpecial(word.PC, word.Address(int(r.pc())+1+w.NumOperands())) }()
			dst := r.program[r.pc()+1]
			switch dst.(type) {
			case word.Register:
				r.registers[dst.(word.Register)] = r.pop()
				return nil
			default:
				return fmt.Errorf("pop: unsupported dst: %s", dst.String())
			}
		case word.ALLOC:
			defer func() { r.setSpecial(word.PC, word.Address(int(r.pc())+1+w.NumOperands())) }()
			var size word.Integer
			op := r.program[r.pc()+1]
			switch op.(type) {
			case word.Register:
				size = r.registers[op.(word.Register)].(word.Integer)
			case word.Offset:
				s, err := r.solve(op.(word.Offset))
				if err != nil {
					return err
				}
				size = s.(word.Integer)
			case word.Integer:
				size = op.(word.Integer)
			default:
				return fmt.Errorf("unsupported alloc size: %s", op.String())
			}
			r.alloc(size.Value())
			return nil
		case word.STORE:
			defer func() { r.setSpecial(word.PC, word.Address(int(r.pc())+1+w.NumOperands())) }()
			dst := r.program[r.pc()+1]
			src := r.program[r.pc()+2]
			if _, ok := dst.(word.Address); !ok {
				return fmt.Errorf("unsupported store dst: %s", dst.String())
			}
			if _, ok := src.(word.Storable); !ok {
				return fmt.Errorf("unsupported store src: %s", src.String())
			}
			r.store(dst.(word.Address), src.(word.Storable))
			return nil
		case word.LOAD:
			defer func() { r.setSpecial(word.PC, word.Address(int(r.pc())+1+w.NumOperands())) }()
			dst := r.program[r.pc()+1]
			src := r.program[r.pc()+2]
			if _, ok := dst.(word.Address); !ok {
				return fmt.Errorf("unsupported load dst: %s", dst.String())
			}
			if _, ok := src.(word.Address); !ok {
				return fmt.Errorf("unsupported load src: %s", src.String())
			}
			r.load(dst.(word.Register), src.(word.Address))
			return nil
		case word.CALL:
		case word.RET:
		case word.ADD:
			defer func() { r.setSpecial(word.PC, word.Address(int(r.pc())+1+w.NumOperands())) }()
		case word.SUB:
			defer func() { r.setSpecial(word.PC, word.Address(int(r.pc())+1+w.NumOperands())) }()
		case word.JMP:
		case word.JE:
		case word.JNE:
		case word.EQ:
			defer func() { r.setSpecial(word.PC, word.Address(int(r.pc())+1+w.NumOperands())) }()
		case word.NE:
			defer func() { r.setSpecial(word.PC, word.Address(int(r.pc())+1+w.NumOperands())) }()
		case word.LT:
			defer func() { r.setSpecial(word.PC, word.Address(int(r.pc())+1+w.NumOperands())) }()
		case word.LE:
			defer func() { r.setSpecial(word.PC, word.Address(int(r.pc())+1+w.NumOperands())) }()
		default:
			return fmt.Errorf("unsupported opcode: %s", w.String())
		}
		return nil
	default:
		return fmt.Errorf("unsupported word: %s", w.String())
	}
}
