package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	left, right, err := read(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	p1, err := part1(left, right)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("part1: %d", p1)
	log.Printf("part2: %d", part2(left, right))
}

func read(r io.Reader) (left, right []int, err error) {
	s := bufio.NewScanner(r)
	line := 0
	for s.Scan() {
		line++
		for i, field := range strings.Fields(s.Text()) {
			x, err := strconv.Atoi(field)
			if err != nil {
				return nil, nil, err
			}
			if i == 0 {
				left = append(left, x)
			} else if i == 1 {
				right = append(right, x)
			} else {
				return nil, nil, fmt.Errorf("too many columns in line %d", line)
			}
		}
	}

	return left, right, s.Err()
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

// part1 sums the distances between pairs of numbers ordered by size in the two input lists
func part1(left, right []int) (int, error) {
	if len(left) != len(right) {
		return 0, fmt.Errorf("lists are of different lengths: %d vs. %d", len(left), len(right))
	}

	sort.Ints(left)
	sort.Ints(right)

	sum := 0

	for i := 0; i < len(left); i++ {
		sum += abs(left[i] - right[i])
	}

	return sum, nil
}

// part2 calculates a "similarity score" by multiplying each number from the left list by the number of its occurrences in the right list
// and summing these products.
func part2(left, right []int) int {
	occurrences := make(map[int]int)
	for _, n := range right {
		occurrences[n] += 1
	}

	similarity := 0

	for _, n := range left {
		similarity += n * occurrences[n]
	}

	return similarity
}
