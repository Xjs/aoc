package main

import (
	"testing"
)

func Test_safe(t *testing.T) {
	tests := []struct {
		name string
		args []int
		want bool
	}{
		{"l1", []int{7, 6, 4, 2, 1}, true},
		{"l2", []int{1, 2, 7, 8, 9}, false},
		{"l3", []int{9, 7, 6, 2, 1}, false},
		{"l4", []int{1, 3, 2, 4, 5}, false},
		{"l5", []int{8, 6, 4, 4, 1}, false},
		{"l6", []int{1, 3, 6, 7, 9}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := safe(tt.args); got != tt.want {
				t.Errorf("safe() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dampenersafe(t *testing.T) {
	tests := []struct {
		name string
		args []int
		want bool
	}{
		{"l1", []int{7, 6, 4, 2, 1}, true},
		{"l2", []int{1, 2, 7, 8, 9}, false},
		{"l3", []int{9, 7, 6, 2, 1}, false},
		{"l4", []int{1, 3, 2, 4, 5}, true},
		{"l5", []int{8, 6, 4, 4, 1}, true},
		{"l6", []int{1, 3, 6, 7, 9}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dampenersafe(tt.args); got != tt.want {
				t.Errorf("dampenersafe() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_safereports(t *testing.T) {
	type args struct {
		ls       [][]int
		dampener bool
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "part1", args: args{ls: [][]int{
			{7, 6, 4, 2, 1},
			{1, 2, 7, 8, 9},
			{9, 7, 6, 2, 1},
			{1, 3, 2, 4, 5},
			{8, 6, 4, 4, 1},
			{1, 3, 6, 7, 9}}, dampener: false,
		}, want: 2},
		{name: "part2", args: args{ls: [][]int{
			{7, 6, 4, 2, 1},
			{1, 2, 7, 8, 9},
			{9, 7, 6, 2, 1},
			{1, 3, 2, 4, 5},
			{8, 6, 4, 4, 1},
			{1, 3, 6, 7, 9}}, dampener: true,
		}, want: 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := safereports(tt.args.ls, tt.args.dampener); got != tt.want {
				t.Errorf("safereports() = %v, want %v", got, tt.want)
			}
		})
	}
}
