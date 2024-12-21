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

func Test_resolveArrows2one(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name string
		want string
	}{
		{"v<<A>>^A<A>AvA<^AA>A<vAAA>^A", "<vA<AA>>^AvAA<^A>A<v<A>>^AvA^A<vA>^A<v<A>^A>AAvA^A<v<A>A>^AAAvA<^A>A"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveArrows2one(tt.name); got != tt.want {
				t.Errorf("resolveArrows2one() = %v, want %v", got, tt.want)
			}
		})
	}
}
