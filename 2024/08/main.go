package main

import (
	"log"
	"os"
	"unicode"

	"github.com/Xjs/aoc/grid"
)

func main() {
	antennas, err := grid.ReadRuneGrid(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	antennaList := make(map[rune]map[grid.Point]struct{})

	antennas.Foreach(func(p grid.Point) {
		freq := antennas.MustAt(p)
		if unicode.IsDigit(freq) || unicode.IsLetter(freq) {
			if antennaList[freq] == nil {
				antennaList[freq] = make(map[grid.Point]struct{})
			}
			antennaList[freq][p] = struct{}{}
		}
	})

	antinodes := grid.NewGrid[rune](antennas.Width(), antennas.Height())
	antinodes2 := grid.NewGrid[rune](antennas.Width(), antennas.Height())

	for _, antennas := range antennaList {
		for antenna := range antennas {
			for antenna2 := range antennas {
				if antenna == antenna2 {
					continue
				}

				antinode := antenna
				antinodes2.Set(antinode, '#')

				d := grid.Diff(antenna, antenna2)

				pt1 := false
				for {
					var err error
					antinode, err = antinodes.Delta(antinode, d)
					if err != nil {
						// out of bounds
						break
					}
					if !pt1 {
						antinodes.Set(antinode, '#')
						pt1 = true
					}
					// Maybe this doesn't capture all of them, if so, calculate the greatest common divisor and divide the diff by it, also
					// sweep in positive and negative direction
					antinodes2.Set(antinode, '#')
				}
			}
		}
	}

	numAntinodes := 0
	antinodes.Foreach(func(p grid.Point) {
		if antinodes.MustAt(p) == '#' {
			numAntinodes++
		}
	})

	numAntinodes2 := 0
	antinodes2.Foreach(func(p grid.Point) {
		if antinodes2.MustAt(p) == '#' {
			numAntinodes2++
		}
	})

	log.Printf("part1: %d", numAntinodes)
	log.Printf("part2: %d", numAntinodes2)
}
