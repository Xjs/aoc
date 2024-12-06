package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/Xjs/aoc/parse"
)

type rule [2]int

func main() {

	// rule is a page number indicating that this page needs to be printed after the index page in the rulebook
	rulebook := make(map[rule]struct{})

	s := bufio.NewScanner(os.Stdin)
	line := 0

	var updates [][]int
	for s.Scan() {
		line++
		text := s.Text()
		if strings.Contains(text, "|") {
			ns, err := parse.IntListSep(text, "|")
			if err != nil {
				log.Fatalf("line %d: %v", line, err)
			}

			if len(ns) != 2 {
				log.Fatalf("line %d: rules must be of length 2", line)
			}

			rulebook[rule(ns)] = struct{}{}
		}

		if strings.Contains(text, ",") {
			update, err := parse.IntList(text)
			if err != nil {
				log.Fatalf("line %d: %v", line, err)
			}
			updates = append(updates, update)
		}
	}

	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	var validUpdates [][]int
	var invalidUpdates [][]int
	for _, update := range updates {
		if valid(update, rulebook) {
			validUpdates = append(validUpdates, update)
		} else {
			invalidUpdates = append(invalidUpdates, update)
		}
	}

	sum := 0
	for _, vu := range validUpdates {
		sum += vu[len(vu)/2]
	}

	log.Printf("part1: %d", sum)

	var fixedInvalids [][]int
	for _, iu := range invalidUpdates {
		valid := false
		for !valid {
			// I hope this terminates
			iu, valid = swapFirstInvalid(iu, rulebook)
		}
		fixedInvalids = append(fixedInvalids, iu)
	}

	sum2 := 0
	for _, fu := range fixedInvalids {
		sum2 += fu[len(fu)/2]
	}

	log.Printf("part2: %d", sum2)

}

func valid(update []int, rulebook map[rule]struct{}) bool {
	for i, n := range update {
		for _, seen := range update[:i] {
			if _, have := rulebook[rule{n, seen}]; have {
				return false
			}
		}
	}

	return true
}

func swapFirstInvalid(update []int, rulebook map[rule]struct{}) ([]int, bool) {
	for i, n := range update {
		for k, seen := range update[:i] {
			if _, have := rulebook[rule{n, seen}]; have {
				u2 := make([]int, len(update))
				copy(u2, update)
				u2[i] = update[k]
				u2[k] = update[i]

				return u2, false
			}
		}
	}

	return update, true
}
