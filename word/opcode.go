package word

type Opcode int

const (
	NOP Opcode = iota

	MOV

	PUSH
	POP

	ALLOC
	STORE
	LOAD

	CALL
	RET

	ADD
	SUB

	JMP
	JE
	JNE

	EQ
	NE
	LT
	LE
)

func (op Opcode) String() string {
	return []string{
		NOP:   "nop",
		MOV:   "mov",
		PUSH:  "push",
		POP:   "pop",
		ALLOC: "alloc",
		STORE: "store",
		LOAD:  "load",
		CALL:  "call",
		RET:   "ret",
		ADD:   "add",
		SUB:   "sub",
		JMP:   "jmp",
		JE:    "je",
		JNE:   "jne",
		EQ:    "eq",
		NE:    "ne",
		LT:    "lt",
		LE:    "le",
	}[op]
}

func (op Opcode) NumOperands() int {
	return []int{
		NOP:   0,
		MOV:   2,
		PUSH:  1,
		POP:   1,
		ALLOC: 1,
		STORE: 2,
		LOAD:  2,
		CALL:  1,
		RET:   0,
		ADD:   2,
		SUB:   2,
		JMP:   1,
		JE:    1,
		JNE:   1,
		EQ:    2,
		NE:    2,
		LT:    2,
		LE:    2,
	}[op]
}
