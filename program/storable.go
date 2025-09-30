package program

type Storable interface {
	Operand
	isStorable()
}
