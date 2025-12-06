package main

import (
	"fmt"
	"github.com/Xjs/aoc/grid"
	"log"
	"math"
	"os"
)

func main() {
	g, err := grid.ReadRuneGrid(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	count := 0

	gg := g.Copy()

	g.Foreach(func(p grid.Point) {
		if g.MustAt(p) != '@' {
			return
		}

		rolls := 0
		for _, pp := range g.Environment8(p) {
			if g.MustAt(pp) == '@' {
				rolls++
			}
		}
		if rolls < 4 {
			gg.Set(p, 'x')
			count++
		}
	})

	fmt.Println(grid.StringCharGrid(gg))

	log.Println("step1", count)

	removables := math.MaxInt
	totalRemoved := 0
	for removables > 0 {
		removables = 0

		gg := g.Copy()

		g.Foreach(func(p grid.Point) {
			if g.MustAt(p) != '@' {
				return
			}

			rolls := 0
			for _, pp := range g.Environment8(p) {
				if g.MustAt(pp) == '@' {
					rolls++
				}
			}
			if rolls < 4 {
				gg.Set(p, 'x')
				removables++
			}
		})

		log.Printf("Removing %d rolls", removables)

		gg.Foreach(func(p grid.Point) {
			if gg.MustAt(p) == 'x' {
				g.Set(p, '.')
				totalRemoved++
			}
		})
	}

	log.Println("step2", totalRemoved)
}
