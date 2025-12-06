package main

import (
	"bufio"
	"fmt"
	"github.com/Xjs/aoc/parse"
	"log"
	"os"
	"strconv"
	"strings"
)

type problem struct {
	op string
	ns []int
}

func (p problem) solve() int {
	var op func(a, b int) int
	var cur int

	switch p.op {
	case "+":
		op = func(a, b int) int { return a + b }
		cur = 0
	case "*":
		op = func(a, b int) int { return a * b }
		cur = 1
	}

	for _, n := range p.ns {
		cur = op(cur, n)
	}

	return cur
}

func main() {
	sc := bufio.NewScanner(os.Stdin)

	var numbers [][]int
	var symbols []string

	var rawLines []string

	length := 0
	for sc.Scan() {
		line := sc.Text()
		rawLines = append(rawLines, line)

		if ns, err := parse.IntListWhitespace(line); err == nil {
			if length == 0 {
				length = len(ns)
			} else if length != len(ns) {
				log.Fatalf("encountered line of length %d, previous length was %d", len(ns), length)
			}
			numbers = append(numbers, ns)
		} else {
			symbols = strings.Fields(line)
			if len(symbols) != length {
				log.Fatalf("encountered %d symbols instead of %d", len(symbols), length)
			}
		}
	}

	grandTotal := 0
	for col := 0; col < length; col++ {
		prob := problem{}
		prob.op = symbols[col]
		for row := 0; row < len(numbers); row++ {
			prob.ns = append(prob.ns, numbers[row][col])
		}

		grandTotal += prob.solve()
	}

	fmt.Println("part1", grandTotal)

	problems := []problem{problem{}}
	cur := 0

	for _, rl := range transpose(rawLines) {
		rl = strings.TrimSpace(rl)
		if rl == "" {
			cur++
			problems = append(problems, problem{})
			continue
		}

		if problems[cur].op == "" {
			problems[cur].op = string(rl[len(rl)-1])
			rl = strings.TrimSpace(rl[:len(rl)-1])
		}

		n, err := strconv.Atoi(rl)
		if err != nil {
			log.Fatal(err)
		}
		problems[cur].ns = append(problems[cur].ns, n)
	}

	grandTotal2 := 0
	for _, problem := range problems {
		grandTotal2 += problem.solve()
	}

	log.Println("part2", grandTotal2)
}

func transpose(lines []string) []string {
	result := make([]string, 0)
	for row := 0; row < len(lines); row++ {
		for col := 0; col < len(lines[row]); col++ {
			for len(result) <= col {
				result = append(result, "")
			}
			result[col] += string(lines[row][col])
		}
	}
	return result
}
