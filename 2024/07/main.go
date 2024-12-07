package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Xjs/aoc/parse"
)

func main() {
	s := bufio.NewScanner(os.Stdin)

	sum := 0
	sum2 := 0
	for s.Scan() {
		sp := strings.Split(s.Text(), ":")
		if len(sp) != 2 {
			log.Fatal("syntax: 'ref: numbers numbers numbers'")
		}
		ref, err := strconv.Atoi(sp[0])
		if err != nil {
			log.Fatal(err)
		}
		numbers, err := parse.IntListWhitespace(sp[1])
		if err != nil {
			log.Fatal(err)
		}

		opts, err := possibleResults(numbers, []rune("+*"))
		if err != nil {
			log.Fatal(err)
		}

		opts2, err := possibleResults(numbers, []rune("+*|"))
		if err != nil {
			log.Fatal(err)
		}

		if _, ok := opts[ref]; ok {
			sum += ref
		}

		if _, ok := opts2[ref]; ok {
			sum2 += ref
		}
	}

	log.Printf("part1: %d", sum)
	log.Printf("part2: %d", sum2)
}

// possibleResults returns a map containing all possible left-to-right evaluations
// with all operators in ops (currently implemented: + *).
func possibleResults(numbers []int, ops []rune) (map[int]struct{}, error) {
	var options [][]int = [][]int{numbers}
	results := make(map[int]struct{})

	for len(options) > 0 {
		var opt []int
		opt, options = options[0], options[1:]
		if len(opt) == 1 {
			results[opt[0]] = struct{}{}
			continue
		}

		for _, op := range ops {
			res, err := eval(op, opt[0], opt[1])
			if err != nil {
				return nil, err
			}
			options = append(options, append([]int{res}, opt[2:]...))
		}
	}

	return results, nil
}

func eval(op rune, a, b int) (int, error) {
	switch op {
	case '+':
		return a + b, nil
	case '*':
		return a * b, nil
	case '|':
		// too lazy to do this algebraically 🙈
		return strconv.Atoi(fmt.Sprintf("%d%d", a, b))
	default:
		return 0, fmt.Errorf("not implemented: %c", op)
	}
}
