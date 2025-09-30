package program

import (
	"fmt"
	"strconv"
)

type Immediate interface {
	Operand
}

type Integer int

func (i Integer) String() string {
	return fmt.Sprintf("%d", i)
}
func (i Integer) Value() int  { return int(i) }
func (i Integer) isOperand()  {}
func (i Integer) isStorable() {}

type Char int

func (c Char) String() string {
	return strconv.QuoteRune(rune(c))
}
func (c Char) Value() int  { return int(c) }
func (c Char) isOperand()  {}
func (c Char) isStorable() {}

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
func (b Bool) isOperand()  {}
func (b Bool) isStorable() {}
