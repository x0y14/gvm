package gvm

import (
	"fmt"
	"strconv"
)

type Operand interface {
	Word
	Value() int
	isOperand()
}

type Register interface {
	Operand
	isRegister()
}

type SpecialRegister int

const (
	_ SpecialRegister = iota
	PC
	BP
	SP
	HP
)

func (s SpecialRegister) String() string {
	return []string{
		PC: "pc",
		BP: "bp",
		SP: "sp",
		HP: "hp",
	}[s]
}
func (s SpecialRegister) Value() int  { return int(s) }
func (s SpecialRegister) isOperand()  {}
func (s SpecialRegister) isRegister() {}

type GeneralPurposeRegister int

const (
	_ GeneralPurposeRegister = iota
	R1
	R2
	R3
	ACM1
	ACM2
)

func (g GeneralPurposeRegister) String() string {
	return []string{
		R1:   "r1",
		R2:   "r2",
		R3:   "r3",
		ACM1: "acm1",
		ACM2: "acm2",
	}[g]
}
func (g GeneralPurposeRegister) Value() int  { return int(g) }
func (g GeneralPurposeRegister) isOperand()  {}
func (g GeneralPurposeRegister) isRegister() {}

type FlagRegister int

const (
	_ FlagRegister = iota
	ZF
)

func (f FlagRegister) String() string {
	return []string{
		ZF: "zf",
	}[f]
}
func (f FlagRegister) Value() int  { return int(f) }
func (f FlagRegister) isOperand()  {}
func (f FlagRegister) isRegister() {}

type PrimitiveType int

const (
	_ PrimitiveType = iota
	TInteger
	TChar
	TBool
)

type Immediate interface {
	Operand
	Type() PrimitiveType
}

type Integer int

func (i Integer) String() string {
	return fmt.Sprintf("%d", i)
}
func (i Integer) Value() int          { return int(i) }
func (i Integer) isOperand()          {}
func (i Integer) Type() PrimitiveType { return TInteger }

type Char int

func (c Char) String() string {
	return strconv.QuoteRune(rune(c))
}
func (c Char) Value() int          { return int(c) }
func (c Char) isOperand()          {}
func (c Char) Type() PrimitiveType { return TChar }

type Bool bool

func (b Bool) String() string {
	if b {
		return "true"
	}
	return "false"
}
func (b Bool) Value() int {
	if b {
		return 1
	}
	return 0
}
func (b Bool) isOperand()          {}
func (b Bool) Type() PrimitiveType { return TBool }

type Location interface {
	Word
	isOperand()
	isLocation()
}

// Offset Relative Location
type Offset interface {
	Location
	isOffset()
}

type BpOffset int

func (b BpOffset) String() string {
	var op = ""
	if b >= 0 {
		op = "+"
	}
	return fmt.Sprintf("[bp%s%d]", op, b)
}
func (b BpOffset) Value() int  { return int(b) }
func (b BpOffset) isOperand()  {}
func (b BpOffset) isLocation() {}
func (b BpOffset) isOffset()   {}

type SpOffset int

func (s SpOffset) String() string {
	var op = ""
	if s >= 0 {
		op = "+"
	}
	return fmt.Sprintf("[sp%s%d]", op, s)
}
func (s SpOffset) Value() int  { return int(s) }
func (s SpOffset) isOperand()  {}
func (s SpOffset) isLocation() {}
func (s SpOffset) isOffset()   {}

// Address Abstract Location
type Address interface {
	Location
	isAddress()
}

type ProgramAddress int

func (p ProgramAddress) String() string {
	return fmt.Sprintf("@%d", p)
}
func (p ProgramAddress) Value() int  { return int(p) }
func (p ProgramAddress) isOperand()  {}
func (p ProgramAddress) isLocation() {}
func (p ProgramAddress) isAddress()  {}

type HeapAddress int

func (h HeapAddress) String() string {
	return fmt.Sprintf("@%d", h)
}
func (h HeapAddress) Value() int  { return int(h) }
func (h HeapAddress) isOperand()  {}
func (h HeapAddress) isLocation() {}
func (h HeapAddress) isAddress()  {}

type Pointer interface {
	Location
	isPointer()
}
type BasePointer int

func (b BasePointer) String() string { return fmt.Sprintf("@%d", b) }
func (b BasePointer) Value() int     { return int(b) }
func (b BasePointer) isOperand()     {}
func (b BasePointer) isLocation()    {}
func (b BasePointer) isPointer()     {}

type StackPointer int

func (s StackPointer) String() string { return fmt.Sprintf("@%d", s) }
func (s StackPointer) Value() int     { return int(s) }
func (s StackPointer) isOperand()     {}
func (s StackPointer) isLocation()    {}
func (s StackPointer) isPointer()     {}
