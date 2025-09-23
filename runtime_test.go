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
					PC:   Integer(0),
					BP:   Integer(0),
					SP:   Integer(0),
					HP:   Integer(0),
					R1:   nil,
					R2:   nil,
					R3:   nil,
					ACM1: nil,
					ACM2: nil,
					ZF:   Integer(0),
				},
				stack: []Operand{nil, nil},
				heap:  []Operand{nil, nil},
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
