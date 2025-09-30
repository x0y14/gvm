package word

type Operand interface {
	Word
	isOperand()
}
