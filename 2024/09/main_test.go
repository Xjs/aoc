package main

import (
	"reflect"
	"testing"
)

func Test_parseFS(t *testing.T) {
	type args struct {
		diskMap []byte
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{"example", args{[]byte("2333133121414131402")}, []int{0, 0, -1, -1, -1, 1, 1, 1, -1, -1, -1, 2, -1, -1, -1, 3, 3, 3, -1, 4, 4, -1, 5, 5, 5, 5, -1, 6, 6, 6, 6, -1, 7, 7, 7, -1, 8, 8, 8, 8, 9, 9}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseFS(tt.args.diskMap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseFS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compact(t *testing.T) {
	type args struct {
		fs []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{"example", args{parseFS([]byte("2333133121414131402"))}, []int{0, 0, 9, 9, 8, 1, 1, 1, 8, 8, 8, 2, 7, 7, 7, 3, 3, 3, 6, 4, 4, 6, 5, 5, 5, 5, 6, 6, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compact(tt.args.fs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compact() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checksum(t *testing.T) {
	type args struct {
		fs []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"example", args{compact(parseFS([]byte("2333133121414131402")))}, 1928},
		{"example", args{compact2(parseFS([]byte("2333133121414131402")))}, 2858},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checksum(tt.args.fs); got != tt.want {
				t.Errorf("checksum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compact2(t *testing.T) {
	type args struct {
		fs []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{"example", args{parseFS([]byte("2333133121414131402"))}, []int{0, 0, 9, 9, 2, 1, 1, 1, 7, 7, 7, -1, 4, 4, -1, 3, 3, 3, -1, -1, -1, -1, 5, 5, 5, 5, -1, 6, 6, 6, 6, -1, -1, -1, -1, -1, 8, 8, 8, 8, -1, -1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compact2(tt.args.fs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compact2() = %v, want %v", got, tt.want)
			}
		})
	}
}
