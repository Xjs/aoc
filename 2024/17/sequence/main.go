package main

import (
	"cmp"
	"fmt"
	"log"
	"slices"
)

func main() {
	ref := []int{2, 4, 1, 1, 7, 5, 1, 5, 4, 2, 5, 5, 0, 3, 3, 0}
	lookup := buildLookup()

	numbers := buildDFS([]pattern{{}}, nil, ref, lookup)
	for _, num := range numbers {
		log.Printf("part2: %d", num.number)
	}

	/*
		candidates2 := build(pattern{}, []int{}, 2, lookup)
		log.Print(candidates2)
		for _, cand := range candidates2 {
			log.Print(program(cand.number)[:1])
			candidates24 := build(cand, []int{2}, 4, lookup)
			log.Print(candidates24)
			for _, cand := range candidates24 {
				log.Print(program(cand.number)[:2])
				candidates241 := build(cand, []int{2, 4}, 1, lookup)
				log.Print(candidates241)
				for _, cand := range candidates241 {
					log.Print(program(cand.number)[:3])
				}
			}
		}
	*/
	return

	// numbers := buildAll(ref)
	// sort.IntSlice(numbers).Sort()

}

// pattern represents a pattern of bits that leads to a specific output
type pattern struct {
	// number is the start number for A that leads to the given output
	number int
	// mask is the mask of bits that are significant for this pattern to work
	// (i. e. all bits that are not in number & mask don't matter and can have any value)
	mask int
}

func (p pattern) String() string {
	return fmt.Sprintf("%b (%b)", p.number, p.mask)
}

func buildAll(wantOutput []int) []int {
	var current []pattern
	lookup := buildLookup()

	for i := 0; i < len(wantOutput)-1; i++ {
		wo, add := wantOutput[:i], wantOutput[i]
		current = buildMulti(current, wo, add, lookup)
		log.Print(i, len(current))
	}

	numbers := make([]int, len(current))
	for i, c := range current {
		numbers[i] = c.number
	}

	return numbers
}

func buildDFS(start []pattern, wantOutput []int, ref []int, lookup map[int]map[pattern]struct{}) []pattern {
	if len(wantOutput) == len(ref) {
		return start
	}

	var result []pattern
	for _, pat := range start {
		next := ref[len(wantOutput)]
		candidates := build(pat, wantOutput, next, lookup)
		slices.SortFunc(candidates, func(i, j pattern) int {
			return cmp.Compare(i.number, j.number)
		})
		result = append(result, buildDFS(candidates, append(wantOutput, next), ref, lookup)...)

		if len(result) > 20 {
			return result
		}
	}
	return result
}

func buildMulti(start []pattern, wantOutput []int, add int, lookup map[int]map[pattern]struct{}) []pattern {
	var result []pattern

	if len(start) == 0 {
		start = []pattern{{}}
	}

	for _, pat := range start {
		result = append(result,
			build(pat, wantOutput, add, lookup)...,
		)
	}
	return result
}

// build returns all numbers that can be built using start, with the numbers from wantOutput already fixed,
// adding add to the output.
// startMask indicates all bits that must not be changed in start (can be more, in which case it indicates the leading zeros are important)
// It is possible to start from 0, 0, meaning everything is open.
func build(start pattern, wantOutput []int, add int, lookup map[int]map[pattern]struct{}) []pattern {
	// We can ignore the lower 3*len(wantOutput) bits as they will not influence the next number anymore
	shift := 3 * len(wantOutput)
	cur := start
	cur.number >>= shift
	cur.mask >>= shift

	// Collect all matching candidates
	var candidates []pattern
	for candidate := range lookup[add] {
		// the bits that are significant in both cases need to match
		mask := candidate.mask & cur.mask

		if cur.number&mask == candidate.number&mask {
			candidates = append(candidates, candidate)
		}
	}

	var candidateNumbers []pattern
	for _, cand := range candidates {
		candidate := ((cand.number) << shift) | start.number
		candidateMask := (cand.mask << shift) | start.mask

		output := program(candidate)

		works := true
		for i, x := range append(wantOutput, add) {
			if output[i] != x {
				works = false
				break
			}
		}
		if !works {
			continue
		}

		candidateNumbers = append(candidateNumbers, pattern{number: candidate, mask: candidateMask})
	}

	return candidateNumbers
}

// buildLookup builds a lookup table from output numbers to patterns the full bit pattern needs to observer
// to output these numbers.
// Since based no the algorithm in [program] all first output numbers are only influenced by
//  1. the rightmost three bits
//  2. the three bits remaining if the input is shifted to the right by (input % 8) ^ 1 (i. e. at most 7)
//
// it cannot be influenced by more than 10 bits, so it goes through all 10-bit numbers,
// analytically building up all patterns based on the first number of the output
func buildLookup() map[int]map[pattern]struct{} {
	lookup := make(map[int]map[pattern]struct{})
	for i := 0; i < 8; i++ {
		lookup[i] = make(map[pattern]struct{})
	}
	for i := 0; i < 1<<10; i++ {
		output := program(i)
		if len(output) == 0 {
			continue
		}

		A := i
		B := ((A % 8) ^ 1)

		// The three bits shifted by B in the formula of the program are significant (they become C)
		mask := (0b111) << B

		// the lowest three bits are always significant (but that doesn't really concern us)
		mask |= 0b111

		pat := pattern{number: i, mask: mask}
		log.Print(output, pat)
		lookup[output[0]][pat] = struct{}{}
	}
	return lookup
}

// program is the program that is coded in my personal puzzle input.
func program(A int) []int {
	B := 0
	C := 0
	var output []int
	for A != 0 {
		// Take the lower 3 bits of A and flip the lowest one
		B = (A % 8) ^ 1
		// Take all bits except the lower B ones
		C = A >> B
		// XOR B with 5 (i. e. flip high and low)
		B ^= 5
		// XOR the bits from B
		// (which are now effectively the lower 3 bits from A with the high one flipped)
		// with the bits from C,
		// which are the bits from A where the last ones were flipped (see above)
		B ^= C
		// Output the lower three bits from B
		output = append(output, B%8)
		// Drop the lowest three bits
		A >>= 3
	}

	return output
}
