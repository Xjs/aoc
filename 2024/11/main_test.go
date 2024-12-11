package main

import (
	"reflect"
	"testing"
)

func TestStones_Blink(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		want []int
	}{
		{"example1", []int{0, 1, 10, 99, 999}, []int{1, 2024, 1, 0, 9, 9, 2021976}},
		{"example2-1", []int{125, 17}, []int{253000, 1, 7}},
		{"example2-2", []int{253000, 1, 7}, []int{253, 0, 2024, 14168}},
		{"example2-3", []int{253, 0, 2024, 14168}, []int{512072, 1, 20, 24, 28676032}},
		{"example2-4", []int{512072, 1, 20, 24, 28676032}, []int{512, 72, 2024, 2, 0, 2, 4, 2867, 6032}},
		{"example2-5", []int{512, 72, 2024, 2, 0, 2, 4, 2867, 6032}, []int{1036288, 7, 2, 20, 24, 4048, 1, 4048, 8096, 28, 67, 60, 32}},
		{"example2-6", []int{1036288, 7, 2, 20, 24, 4048, 1, 4048, 8096, 28, 67, 60, 32}, []int{2097446912, 14168, 4048, 2, 0, 2, 4, 40, 48, 2024, 40, 48, 80, 96, 2, 8, 6, 7, 6, 0, 3, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := NewStones(tt.s)
			st.Blink()
			if got := st.IntSlice(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stones.Blink() = %v, want %v", got, tt.want)
			}
		})
	}
}
