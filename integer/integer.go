package integer

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Sum returns the sum of the list of integers
func Sum(ns []int) int {
	var sum int
	for _, n := range ns {
		sum += n
	}
	return sum
}

// Abs returns the absolute value of x.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Max(xs []int) int {
	max := math.MinInt
	for _, x := range xs {
		if x > max {
			max = x
		}
	}
	return max
}

// A Range represents a range of integers (inclusive of the edges)
type Range struct {
	Min, Max int
}

func (r Range) String() string {
	return fmt.Sprintf("%d-%d", r.Min, r.Max)
}

func ParseRange(s string) (r Range, err error) {
	sp := strings.Split(s, "-")
	if len(sp) != 2 {
		return r, fmt.Errorf("range %q is not of format X-Y", s)
	}

	min, err := strconv.Atoi(sp[0])
	if err != nil {
		return r, fmt.Errorf("error parsing lower bound: %w", err)
	}

	max, err := strconv.Atoi(sp[1])
	if err != nil {
		return r, fmt.Errorf("error parsing upper bound: %w", err)
	}

	return Range{min, max}, nil
}

func (r Range) Contains(x int) bool {
	return x >= r.Min && x <= r.Max
}

func (r Range) Length() int {
	return r.Max - r.Min + 1
}
