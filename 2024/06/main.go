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

	var next *grid.Point
	for {
		next = step(g, next)
		if next == nil {
			break
		}
	}

	sum := 0
	g.Foreach(func(p grid.Point) {
		if g.MustAt(p) == 'X' {
			sum++
		}
	})

	log.Printf("part1: %d", sum)
}

// Make the guard perform one action:
// * turn 90° right if something is in the way
// or
// * move in the direction its arrow faces,
// mark visited places with an X (in-place),
// and return nil if the guard leaves the map,
// or the guard's new position otherwise.
//
// If start is nil, will search for the guard
// (indicated by ^, <, v, or >), and start at the last guard found,
// or return nil if no guard has been found.
func step(g *grid.Grid[rune], start *grid.Point) *grid.Point {
	if start == nil {
		g.Foreach(func(p grid.Point) {
			switch g.MustAt(p) {
			case '<', 'v', '>', '^':
				foundPoint := grid.P(p.X, p.Y)
				start = &foundPoint
			}
		})
	}

	if start == nil {
		return nil
	}

	p := *start
	guard := g.MustAt(p)

	nextPoint, err := g.Delta(p, steps[guard])
	if err != nil {
		// leaving the grid
		g.Set(p, 'X')
		return nil
	}

	if g.MustAt(nextPoint) == '#' {
		g.Set(p, rotations[guard])
		return &p
	}

	g.Set(p, 'X')
	g.Set(nextPoint, guard)
	return &nextPoint
}

var steps = map[rune]grid.Delta{
	'<': grid.D(-1, 0),
	'^': grid.D(0, -1),
	'>': grid.D(1, 0),
	'v': grid.D(0, 1),
}

var rotations = map[rune]rune{
	'<': '^',
	'^': '>',
	'>': 'v',
	'v': '<',
}
