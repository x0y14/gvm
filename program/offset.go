package program

import (
	"fmt"
	"strconv"
)

// Offset 相対位置表現
type Offset struct {
	Target Register
	Diff   int
}

func (o Offset) String() string {
	switch o.Target {
	case PC, BP, SP:
		s := strconv.Itoa(o.Diff)
		if o.Diff >= 0 {
			s = "+" + s
		}
		return fmt.Sprintf("[%s%s]", o.Target.String(), s)
	default:
		panic("unsupported :" + o.Target.String())
	}
}

func (o Offset) isOperand() {}
