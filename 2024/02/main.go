package main

import (
	"bufio"
	"io"
	"log"
	"os"

	"github.com/Xjs/aoc/parse"
)

func main() {
	slices, err := parseInput(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("part1: %d", safereports(slices, false))
	log.Printf("part2: %d", safereports(slices, true))
}

func parseInput(r io.Reader) ([][]int, error) {
	s := bufio.NewScanner(r)
	var result [][]int
	for s.Scan() {
		res, err := parse.IntListWhitespace(s.Text())
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}
	return result, s.Err()
}

func safe(s []int) bool {
	dir := 0
	for i := 0; i < len(s)-1; i++ {
		a, b := s[i], s[i+1]
		ldir := 0
		if a < b {
			ldir = 1
		} else if a > b {
			ldir = -1
		}
		if ldir == 0 {
			// unsafe, cannot have equal numbers
			return false
		}

		if dir == 0 {
			dir = ldir
		}

		if ldir != dir {
			// unsafe, direction changes
			return false
		}

		diff := ldir * (b - a)
		if diff < 1 || diff > 3 {
			// unsafe, too high a difference
			return false
		}
	}

	return true
}

func cut(s []int, i int) []int {
	ns := make([]int, len(s)-1)
	copy(ns[:i], s[:i])
	copy(ns[i:], s[i+1:])
	return ns
}

func dampenersafe(l []int) bool {
	for i := 0; i < len(l); i++ {
		nl := cut(l, i)
		if safe(nl) {
			return true
		}
	}
	return false
}

func safereports(ls [][]int, dampener bool) int {
	sum := 0
	for _, l := range ls {
		if safe(l) {
			sum++
		} else if dampener {
			if dampenersafe(l) {
				sum++
				continue
			}
		}
	}
	return sum
}
