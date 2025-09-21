package program

import (
	"fmt"
	"strconv"
)

type Operand interface {
	Word
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
func (g GeneralPurposeRegister) isOperand()  {}
func (g GeneralPurposeRegister) isRegister() {}

type FlagRegister int

const (
	_ FlagRegister = iota
	ZF
)

func (f FlagRegister) String() string {
	return []string{}[f]
}
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
func (i Integer) isOperand()          {}
func (i Integer) Type() PrimitiveType { return TInteger }

type Char int

func (c Char) String() string {
	return strconv.QuoteRune(rune(c))
}
func (c Char) isOperand()          {}
func (c Char) Type() PrimitiveType { return TChar }

type Bool bool

func (b Bool) String() string {
	if b {
		return "true"
	}
	return "false"
}
func (b Bool) isOperand()          {}
func (b Bool) Type() PrimitiveType { return TBool }
