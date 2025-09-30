package program

type Operand interface {
	Word
	isOperand()
}
