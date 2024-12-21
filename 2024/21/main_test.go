package main

import (
	"testing"
)

func Test_complexity(t *testing.T) {
	tests := []struct {
		name string
		iter int
		want int
	}{
		{"029A", 2, 68 * 29},
		{"980A", 2, 60 * 980},
		{"179A", 2, 68 * 179},
		{"456A", 2, 64 * 456},
		{"379A", 2, 64 * 379},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := complexity(tt.name, tt.iter); got != tt.want {
				t.Errorf("complexity() = %v, want %v", got, tt.want)
			}
		})
	}
}
