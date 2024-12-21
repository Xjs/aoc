package main

import (
	"math"
)

var costs = make(map[[2]rune]map[int]int)

func cost(r1, r2 rune, layer int) int {
	if c, ok := costs[[2]rune{r1, r2}][layer]; ok {
		return c
	}

	if r1 == r2 {
		return 1
	}

	paths := shortestPathsArrow[[2]rune{r1, r2}]
	lowestCost := math.MaxInt
	for _, path := range paths {
		pairs := make([][2]rune, len(path)+1)
		for i, d := range path {
			if i == 0 {
				pairs[i] = [2]rune{'A', toRune(d)}
				continue
			}
			pairs[i] = [2]rune{toRune(path[i-1]), toRune(d)}
		}
		pairs[len(path)] = [2]rune{toRune(path[len(path)-1]), 'A'}

		total := 0
		for _, pair := range pairs {
			pairCost := 1
			if layer > 0 {
				pairCost = cost(pair[0], pair[1], layer-1)
			}
			total += pairCost
		}

		if total < lowestCost {
			lowestCost = total
		}
	}

	if costs[[2]rune{r1, r2}] == nil {
		costs[[2]rune{r1, r2}] = make(map[int]int)
	}
	costs[[2]rune{r1, r2}][layer] = lowestCost
	return lowestCost
}
