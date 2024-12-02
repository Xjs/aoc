package grid

import (
	"reflect"
	"testing"
)

func TestGrid_Environment4(t *testing.T) {
	type fields struct {
		width  Coordinate
		height Coordinate
		values [][]int
	}
	type args struct {
		p Point
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Point
	}{
		{"regular", fields{42, 42, nil}, args{P(5, 5)}, []Point{P(4, 5), P(6, 5), P(5, 4), P(5, 6)}},
		{"lower bound", fields{42, 42, nil}, args{P(0, 0)}, []Point{P(1, 0), P(0, 1)}},
		{"upper bound", fields{42, 42, nil}, args{P(41, 41)}, []Point{P(40, 41), P(41, 40)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Grid[int]{
				width:  tt.fields.width,
				height: tt.fields.height,
				values: tt.fields.values,
			}
			if got := g.Environment4(tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Grid.Environment4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrid_Environment8(t *testing.T) {
	type fields struct {
		width  Coordinate
		height Coordinate
		values [][]int
	}
	type args struct {
		p Point
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Point
	}{
		{"regular", fields{42, 42, nil}, args{P(5, 5)}, []Point{
			P(4, 5), P(6, 5), P(5, 4), P(5, 6),
			P(4, 4), P(6, 6), P(4, 6), P(6, 4)},
		},
		{"regular", fields{42, 42, nil}, args{P(10, 5)}, []Point{
			P(9, 5), P(11, 5), P(10, 4), P(10, 6),
			P(9, 4), P(11, 6), P(9, 6), P(11, 4)},
		},
		{"lower bound", fields{42, 42, nil}, args{P(0, 0)}, []Point{P(1, 0), P(0, 1), P(1, 1)}},
		{"upper bound", fields{42, 42, nil}, args{P(41, 41)}, []Point{P(40, 41), P(41, 40), P(40, 40)}},
		{"edge top", fields{42, 42, nil}, args{P(5, 0)}, []Point{
			P(4, 0), P(6, 0), P(5, 1),
			P(6, 1), P(4, 1)},
		},
		{"edge bottom", fields{42, 42, nil}, args{P(5, 41)}, []Point{
			P(4, 41), P(6, 41), P(5, 40),
			P(4, 40), P(6, 40)},
		},
		{"edge left", fields{42, 42, nil}, args{P(0, 5)}, []Point{
			P(1, 5), P(0, 4), P(0, 6),
			P(1, 6), P(1, 4)},
		},
		{"edge right", fields{42, 42, nil}, args{P(41, 5)}, []Point{
			P(40, 5), P(41, 4), P(41, 6),
			P(40, 4), P(40, 6)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Grid[int]{
				width:  tt.fields.width,
				height: tt.fields.height,
				values: tt.fields.values,
			}
			if got := g.Environment8(tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Grid.Environment8() = %v, want %v", got, tt.want)
			}
		})
	}
}
