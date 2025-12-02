package main

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestDivisors(t *testing.T) {
	if !cmp.Equal(divisors(10), []int{1, 2, 5}) {
		t.Errorf("divisors of 10 != [1, 2, 5], got: %v", divisors(10))
	}
}
