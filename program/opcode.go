package program

type Opcode int

const (
	NOP Opcode = iota
)

func (op Opcode) String() string {
	return []string{
		NOP: "nop",
	}[op]
}

func (op Opcode) NumOperands() int {
	return []int{
		NOP: 0,
	}[op]
}
