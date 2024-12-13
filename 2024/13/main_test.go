package main

import (
	"reflect"
	"testing"

	"github.com/Xjs/aoc/grid"
)

const (
	ex1 = `Button A: X+94, Y+34
Button B: X+22, Y+67
Prize: X=8400, Y=5400`
	ex2 = `Button A: X+26, Y+66
Button B: X+67, Y+21
Prize: X=12748, Y=12176`
	ex3 = `Button A: X+17, Y+86
Button B: X+84, Y+37
Prize: X=7870, Y=6450`
	ex4 = `Button A: X+69, Y+23
Button B: X+27, Y+71
Prize: X=18641, Y=10279`
)

var (
	cm1 = clawMachine{grid.P(94, 34), grid.P(22, 67), grid.P(8400, 5400)}
	cm2 = clawMachine{grid.P(26, 66), grid.P(67, 21), grid.P(12748, 12176)}
	cm3 = clawMachine{grid.P(17, 86), grid.P(84, 37), grid.P(7870, 6450)}
	cm4 = clawMachine{grid.P(69, 23), grid.P(27, 71), grid.P(18641, 10279)}
)

func Test_parseClawMachine(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    clawMachine
		wantErr bool
	}{
		{"ex1", args{ex1}, cm1, false},
		{"ex2", args{ex2}, cm2, false},
		{"ex3", args{ex3}, cm3, false},
		{"ex4", args{ex4}, cm4, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseClawMachine(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseClawMachine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseClawMachine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_clawMachine_solve(t *testing.T) {
	tests := []struct {
		name  string
		cm    clawMachine
		want  bool
		wantA int
		wantB int
	}{
		{"ex1", cm1, true, 80, 40},
		{"ex2", cm2, false, 0, 0},
		{"ex3", cm3, true, 38, 86},
		{"ex4", cm4, false, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := tt.cm.solve()
			if got != tt.want {
				t.Errorf("clawMachine.solve() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.wantA {
				t.Errorf("clawMachine.solve() got1 = %v, want %v", got1, tt.wantA)
			}
			if got2 != tt.wantB {
				t.Errorf("clawMachine.solve() got2 = %v, want %v", got2, tt.wantB)
			}
		})
	}
}
