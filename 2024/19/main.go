package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"strings"
)

func main() {
	s := bufio.NewScanner(os.Stdin)

	patterns := make(map[string]struct{})
	var tasks []string
	maxLen := 0
	minLen := math.MaxInt

	for s.Scan() {
		t := strings.TrimSpace(s.Text())
		if sp := strings.Split(t, ", "); len(sp) > 1 {
			for _, pattern := range sp {
				if len(pattern) > maxLen {
					maxLen = len(pattern)
				}
				if len(pattern) < minLen {
					minLen = len(pattern)
				}
				patterns[pattern] = struct{}{}
			}
			continue
		}

		if t == "" {
			continue
		}

		tasks = append(tasks, t)
	}

	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	possible := 0
	all := 0
	for i, t := range tasks {
		if len(build(t, nil, patterns, maxLen)) != 0 {
			log.Printf("%d: %q (%d), %d, %d", i, t, len(t), minLen, maxLen)
			possible++
			all += buildAll(t, patterns, maxLen)
		}
	}

	log.Printf("part1: %d", possible)
	log.Printf("part2: %d", all)
}

func build(pattern string, state []string, patterns map[string]struct{}, maxLen int) []string {
	if pattern == "" {
		return state
	}

	for i := maxLen; i > 0; i-- {
		if len(pattern) < i {
			continue
		}
		pat := pattern[:i]
		_, ok := patterns[pat]
		if !ok {
			continue
		}

		st := build(strings.TrimPrefix(pattern, pat), append(state, pat), patterns, maxLen)
		if len(st) > 0 {
			return st
		}
	}

	return nil
}

var lenCache = make(map[string]int)

func buildAll(pattern string, patterns map[string]struct{}, maxLen int) int {
	if l, ok := lenCache[pattern]; ok {
		return l
	}

	if len(pattern) == 0 {
		return 0
	}

	var sum int

	for i := 1; i <= maxLen; i++ {
		if i > len(pattern) {
			continue
		}
		pat := pattern[:i]
		if _, ok := patterns[pat]; !ok {
			continue
		}
		current, ok := strings.CutPrefix(pattern, pat)
		if !ok {
			continue
		}
		if current == "" {
			sum += 1
			continue
		}
		sum += buildAll(current, patterns, maxLen)
	}

	lenCache[pattern] = sum
	return sum
}
