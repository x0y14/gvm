package program

import (
	"fmt"
	"strconv"
)

// Address 絶対位置表現
type Address int

func (a Address) String() string {
	s := strconv.Itoa(int(a))
	if int(a) >= 0 {
		s = "+" + s
	}
	return fmt.Sprintf("@%s", s)
}
func (a Address) isOperand()  {}
func (a Address) isStorable() {}
