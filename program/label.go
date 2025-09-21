package program

import "fmt"

type Label int

func (l Label) String() string {
	return fmt.Sprintf("label(%d)", l)
}
