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
				stack: []Stockable{nil, nil},
				heap:  []Stockable{nil, nil},
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
				stack: []Stockable{nil, nil},
				heap:  []Stockable{},
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
				stack: []Stockable{nil, nil},
				heap:  []Stockable{},
			},
		},
		{
			"alloc, store, load",
			[]Word{
				ALLOC, Integer(1), // ヒープに1つ分確保、base addrをpush
				POP, R1, // base addrをR1にpop
				PUSH, Integer(42), // 42をpush
				POP, R2, // R2にpop
				STORE, R1, R2, // R1(addr)にR2の値(42)をstore
				LOAD, R3, R1, // R1(addr)からloadしてR3へ
			},
			&Config{4, 4},
			&Runtime{
				program: nil,
				registers: map[Register]Operand{
					PC:   ProgramAddress(14),
					BP:   BasePointer(0),
					SP:   StackPointer(3),
					HP:   HeapAddress(1),
					R1:   HeapAddress(0),
					R2:   Integer(42),
					R3:   Integer(42),
					ACM1: nil,
					ACM2: nil,
					ZF:   Bool(false),
				},
				stack: []Stockable{nil, nil, nil, nil},
				heap:  []Stockable{Integer(42), nil, nil, nil},
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
