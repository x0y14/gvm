package gvm

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRuntime_Run(t *testing.T) {
	tests := []struct {
		name   string
		prog   Program
		config *Config
		result *Runtime
	}{
		{
			"init",
			[]Word{},
			&Config{2, 2},
			&Runtime{
				program: nil,
				registers: map[Register]Operand{
					PC:   ProgramAddress(0),
					BP:   BasePointer(0),
					SP:   StackPointer(2 - 1),
					HP:   HeapAddress(0),
					R1:   nil,
					R2:   nil,
					R3:   nil,
					ACM1: nil,
					ACM2: nil,
					ZF:   Bool(false),
				},
				stack: []Operand{nil, nil},
				heap:  []Operand{nil, nil},
			},
		},
		{
			"push, pop",
			[]Word{
				PUSH, Integer(99),
				POP, R3,
			},
			&Config{2, 0},
			&Runtime{
				program: nil,
				registers: map[Register]Operand{
					PC:   ProgramAddress(4),
					BP:   BasePointer(0),
					SP:   StackPointer(1),
					HP:   HeapAddress(0),
					R1:   nil,
					R2:   nil,
					R3:   Integer(99),
					ACM1: nil,
					ACM2: nil,
					ZF:   Bool(false),
				},
				stack: []Operand{nil, nil},
				heap:  []Operand{},
			},
		},
		{
			"mov",
			[]Word{
				PUSH, Integer(99),
				POP, R3,
				MOV, R1, R3,
			},
			&Config{2, 0},
			&Runtime{
				program: nil,
				registers: map[Register]Operand{
					PC:   ProgramAddress(7),
					BP:   BasePointer(0),
					SP:   StackPointer(1),
					HP:   HeapAddress(0),
					R1:   Integer(99),
					R2:   nil,
					R3:   Integer(99),
					ACM1: nil,
					ACM2: nil,
					ZF:   Bool(false),
				},
				stack: []Operand{nil, nil},
				heap:  []Operand{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRuntime(tt.prog, tt.config)
			err := r.Run()
			if err != nil {
				t.Error(err)
			}

			tt.result.program = tt.prog
			if diff := cmp.Diff(tt.result, r, cmp.AllowUnexported(Runtime{})); diff != "" {
				t.Errorf("diff: %s", diff)
			}
		})
	}
}
