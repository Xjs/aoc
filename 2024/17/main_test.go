package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_machine_exec(t *testing.T) {
	tests := []struct {
		name string
		prog []int
		m    machine
		want machine
	}{
		{"ex1", []int{2, 6}, machine{C: 9}, machine{B: 1, C: 9}},
		{"ex2", []int{5, 0, 5, 1, 5, 4}, machine{A: 10}, machine{A: 10, output: []int{0, 1, 2}}},
		{"ex3", []int{0, 1, 5, 4, 3, 0}, machine{A: 2024}, machine{output: []int{4, 2, 5, 6, 7, 7, 7, 7, 3, 1, 0}}},
		{"ex4", []int{1, 7}, machine{B: 29}, machine{B: 26}},
		{"ex5", []int{4, 0}, machine{B: 2024, C: 43690}, machine{B: 44354, C: 43690}},
		{"ex6", []int{0, 1, 5, 4, 3, 0}, machine{A: 729}, machine{output: []int{4, 6, 3, 5, 6, 3, 5, 2, 1, 0}}},
		// According to the website, this is wrong
		{"prod", []int{2, 4, 1, 1, 7, 5, 1, 5, 4, 2, 5, 5, 0, 3, 3, 0}, machine{A: 28422061}, machine{B: 5, output: []int{3, 6, 7, 0, 5, 7, 3, 1, 5}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.program = tt.prog
			tt.want.program = tt.prog
			tt.want.i = len(tt.prog)
			for (&tt.m).exec() {
			}
			if got := tt.m; !cmp.Equal(got, tt.want, cmp.AllowUnexported(machine{})) {
				t.Errorf("machine.exec() = %v, want %v", got, tt.want)
			}
		})
	}
}
