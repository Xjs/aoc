package main

import (
	"testing"
)

func Test_complexity(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{"029A", 68 * 29},
		{"980A", 60 * 980},
		{"179A", 68 * 179},
		{"456A", 64 * 456},
		{"379A", 64 * 379},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := complexity(tt.name); got != tt.want {
				t.Errorf("complexity() = %v, want %v", got, tt.want)
			}
		})
	}
}
