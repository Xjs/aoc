package main

import (
	"github.com/Xjs/aoc/grid"
	"log"
	"os"
)

func main() {
	g, err := grid.ReadDigitGrid(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	sum := 0
	sum2 := 0

	for _, line := range g.Lines() {
		sum += findMaxJoltage(line, 2)
		sum2 += findMaxJoltage(line, 12)
	}

	log.Println("part1:", sum)
	log.Println("part2:", sum2)
}

func findMaxJoltage(line []int, length int) int {
	maxVal := make([]int, length)
	maxIdx := make([]int, length)

	sum := 0

	for l := 0; l < length; l++ {
		remainder := length - l - 1

		prevMaxIdx := -1
		if l != 0 {
			prevMaxIdx = maxIdx[l-1]
		}
		for i := prevMaxIdx + 1; i < len(line)-remainder; i++ {
			if line[i] > maxVal[l] {
				maxVal[l] = line[i]
				maxIdx[l] = i
			}
		}

		sum += maxVal[l] * pow10(uint(remainder))
	}
	return sum
}

func pow10(n uint) int {
	if n == 0 {
		return 1
	}
	return 10 * pow10(n-1)
}
