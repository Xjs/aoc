package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"strings"
	"time"
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

	t := time.Now()
	possibleT := 0
	allT := 0
	for _, t := range tasks {
		c := buildAll(t, patterns, maxLen)
		if c > 0 {
			possibleT++
		}
		allT += c
	}
	log.Printf("part1: %d", possibleT)
	log.Printf("part2: %d", allT)

	log.Printf("%s", time.Since(t))
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
