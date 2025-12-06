package main

import (
	"bufio"
	"github.com/Xjs/aoc/integer"
	"log"
	"os"
	"sort"
	"strconv"
)

func main() {
	sc := bufio.NewScanner(os.Stdin)

	var ranges []integer.Range
	idSet := make(map[int]struct{})

	rangeMode := true
	for sc.Scan() {
		l := sc.Text()
		if l == "" {
			rangeMode = !rangeMode
			continue
		}

		if rangeMode {
			r, err := integer.ParseRange(l)
			if err != nil {
				log.Fatal(err)
			}
			ranges = append(ranges, r)
		} else {
			id, err := strconv.Atoi(l)
			if err != nil {
				log.Fatal(err)
			}
			idSet[id] = struct{}{}
		}
	}

	var ids []int
	for id := range idSet {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	ranges = mergeRanges(ranges)

	goodIDs := 0

	currentRange := 0
outer:
	for _, id := range ids {
		for id > ranges[currentRange].Max {
			currentRange++
			if currentRange >= len(ranges) {
				break outer
			}
		}

		if ranges[currentRange].Contains(id) {
			goodIDs++
		}
	}

	log.Printf("part1: %d", goodIDs)

	total := 0
	for _, r := range ranges {
		total += r.Length()
	}

	log.Printf("part2: %d", total)
}

func mergeRanges(ranges []integer.Range) []integer.Range {
	// Sort ranges by lower bound
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].Min < ranges[j].Min
	})

	removeRanges := make(map[int]struct{}, len(ranges))
	for i, r := range ranges {
		if i == 0 {
			continue
		}
		// Find ranges that overlap (after sorting)
		if r.Min <= ranges[i-1].Max {
			removeRanges[i-1] = struct{}{}
			// Merge ranges into the second one
			ranges[i].Min = ranges[i-1].Min
			// Account for complete inclusion
			if ranges[i-1].Max > ranges[i].Max {
				ranges[i].Max = ranges[i-1].Max
			}
		}
	}

	var mergedRanges []integer.Range
	for i, r := range ranges {
		if _, ok := removeRanges[i]; ok {
			continue
		}
		mergedRanges = append(mergedRanges, r)
	}
	return mergedRanges
}
