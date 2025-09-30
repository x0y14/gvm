package word

type Storable interface {
	Operand
	isStorable()
}
