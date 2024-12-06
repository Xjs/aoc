package main

import (
	"log"
	"os"

	"github.com/Xjs/aoc/grid"
)

func main() {
	g, err := grid.ReadRuneGrid(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("part1: %d", part1(g))
	log.Printf("part2: %d", part2(g))
}

func part1(g *grid.Grid[rune]) int {
	var search = []rune("XMAS")

	foundInstances := 0

	g.Foreach(func(p grid.Point) {
		if g.MustAt(p) == search[0] {
			for _, dir := range directions(p, len(search)) {
				found := true
				for i, np := range dir {
					r, err := g.At(np)
					if err != nil {
						found = false
						break
					}
					if r != search[i] {
						found = false
						break
					}
				}
				if found {
					foundInstances++
				}
			}
		}
	})

	return foundInstances
}

func part2(g *grid.Grid[rune]) int {
	type dp struct{ dx, dy int }

	var search = []rune("MAS")

	foundInstances := 0

	g.Foreach(func(p grid.Point) {
		if g.MustAt(p) == search[1] {
			points := make(map[dp]rune)
			for _, dx := range []int{-1, 1} {
				for _, dy := range []int{-1, 1} {
					x := int(p.X) + dx
					y := int(p.Y) + dy
					if x < 0 || y < 0 {
						return
					}

					pp, err := g.At(grid.P(grid.Coordinate(x), grid.Coordinate(y)))
					if err != nil {
						return
					}
					points[dp{dx, dy}] = pp
				}
			}

			// top left
			tl := points[dp{-1, -1}]
			// top right
			tr := points[dp{-1, 1}]
			// bottom left
			bl := points[dp{1, -1}]
			// bottom right
			br := points[dp{1, 1}]

			tl2br := false
			if tl == search[0] && br == search[2] || tl == search[2] && br == search[0] {
				tl2br = true
			}

			tr2bl := false
			if tr == search[0] && bl == search[2] || tr == search[2] && bl == search[0] {
				tr2bl = true
			}

			if tl2br && tr2bl {
				foundInstances++
			}
		}
	})

	return foundInstances
}

func directions(c grid.Point, length int) [][]grid.Point {
	var result [][]grid.Point

	for _, dx := range []int{-1, 0, 1} {
		for _, dy := range []int{-1, 0, 1} {
			if dx == 0 && dy == 0 {
				continue
			}

			var dir []grid.Point
			for i := 0; i < length; i++ {
				x := int(c.X) + i*dx
				y := int(c.Y) + i*dy

				if x < 0 || y < 0 {
					break
				}

				dir = append(dir, grid.P(grid.Coordinate(x), grid.Coordinate(y)))
			}

			if len(dir) < length {
				// encountered a negative here, can't complete the sequence
				continue
			}

			result = append(result, dir)
		}
	}

	return result
}
