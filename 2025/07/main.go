package main

import (
	"errors"
	"github.com/Xjs/aoc/grid"
	"log"
	"os"
)

func main() {
	g, err := grid.ReadRuneGrid(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	g0 := g.Copy()

	total := 0

	for {
		gg, hits, lowest := propagate(g)
		g = gg
		if lowest {
			break
		}
		total += hits
	}

	log.Println("part1", total)
	log.Println("part2", analyse(&g0))
}

// propagate applies the transformation rules:
// Empty space is represented by '.'
// A beam is represented by '|'.
// A splitter is represented by '^'.
// If an empty space has a beam directly above it, it will propagate downwards into it.
// If a splitter has a beam directly above it, it will propagate the beam to its left and right, given there is empty
// space.
// The number of splitter hits will be returned. The last return value will be true if only the grid has been updated
// this round.
func propagate(g *grid.Grid[rune]) (*grid.Grid[rune], int, bool) {
	hits := 0
	static := true

	gg := g.Copy()

	g.Foreach(func(p grid.Point) {
		c := g.MustAt(p)
		above, err := g.At(grid.P(p.X, p.Y-1))
		if errors.Is(err, grid.ErrOutOfBounds) {
			// out of bounds
			return
		} else if err != nil {
			panic(err)
		}

		if above == 'S' {
			above = '|'
		}

		if above != '|' {
			return
		}

		switch c {
		case '.':
			gg.Set(p, '|')
		case '^':
			left := grid.P(p.X-1, p.Y)
			right := grid.P(p.X+1, p.Y)

			update := false
			if v, err := g.At(left); err == nil && v == '.' {
				gg.Set(left, '|')
				update = true
			}
			if v, err := g.At(right); err == nil && v == '.' {
				gg.Set(right, '|')
				update = true
			}
			if !update {
				return
			}
			hits++
		default:
			return
		}

		static = false
	})

	return &gg, hits, static
}

// analyse analyses the properties of the grid line by line. It will return a number of possible beams.
func analyse(g *grid.Grid[rune]) int {
	w, h := g.Width(), g.Height()
	gg := grid.NewGrid[int](w, h)

	startLines := make(map[int]bool)
	g.Foreach(func(p grid.Point) {
		if g.MustAt(p) == 'S' {
			gg.Set(p, 1)
			startLines[int(p.Y)] = true
		}
	})

	if len(startLines) != 1 || !startLines[0] {
		panic("starting elsewhere than in line 0 is not supported")
	}

	last := func() int {
		sum := 0
		for x := grid.Coordinate(0); x < w; x++ {
			sum += gg.MustAt(grid.P(x, h-1))
		}
		return sum
	}

	for y := grid.Coordinate(1); y < h; y++ {
		for x := grid.Coordinate(0); x < w; x++ {
			p := grid.P(x, y)
			above := grid.P(x, y-1)
			left := grid.P(x-1, y)
			right := grid.P(x+1, y)
			amt := gg.MustAt(above)
			if amt == 0 {
				continue
			}
			if g.MustAt(grid.P(x, y)) == '^' {
				ll, _ := gg.At(left)
				gg.Set(left, ll+amt)

				rr, _ := gg.At(right)
				gg.Set(right, rr+amt)
			} else {
				gg.Set(p, gg.MustAt(p)+amt)
			}
		}
	}

	return last()
}
