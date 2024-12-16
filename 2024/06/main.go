package main

import (
	"log"
	"os"

	"github.com/Xjs/aoc/grid"
)

type visit map[rune]struct{}

func main() {
	g, err := grid.ReadRuneGrid(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	gbak := g.Copy()

	loop := walk(g)
	if loop {
		log.Fatal("encountered loop in part1")
	}

	sum := 0
	var visitedPoints []grid.Point
	g.Foreach(func(p grid.Point) {
		if g.MustAt(p) == 'X' {
			sum++
			visitedPoints = append(visitedPoints, p)
		}
	})

	log.Printf("part1: %d", sum)

	loops := 0
	for _, vp := range visitedPoints {
		gNew := gbak.Copy()
		gNew.Set(vp, '#')
		if loop := walk(&gNew); loop {
			loops++
		}
	}

	log.Printf("part2: %d", loops)
}

// walk walks the grid until the guard leaves the grid or a loop is encountered.
// It leaves the visited places marked.
func walk(g *grid.Grid[rune]) bool {
	var next *grid.Point
	// visit in the given rotation
	visits := grid.NewGrid[visit](g.Width(), g.Height())
	loop := false
	for !loop {
		next, loop = step(g, next, &visits)
		if next == nil {
			break
		}
	}

	return loop
}

// Make the guard perform one action:
// * turn 90° right if something is in the way
// or
// * move in the direction its arrow faces,
// mark visited places with an X (in-place),
// and return nil if the guard leaves the map,
// or the guard's new position otherwise.
// The second return value indicates if we have run into a loop,
// i. e. we have been at this position with the same rotation before.
//
// If start is nil, will search for the guard
// (indicated by ^, <, v, or >), and start at the last guard found,
// or return nil if no guard has been found.
func step(g *grid.Grid[rune], start *grid.Point, visits *grid.Grid[visit]) (*grid.Point, bool) {
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
		return nil, false
	}

	p := *start
	guard := g.MustAt(p)

	nextPoint, err := g.Delta(p, grid.GeneralDirections[guard])
	if err != nil {
		// leaving the grid
		g.Set(p, 'X')
		return nil, false
	}

	if g.MustAt(nextPoint) == '#' {
		g.Set(p, rotations[guard])
		return &p, false
	}

	vs := visits.MustAt(p)
	if vs == nil {
		vs = make(visit)
	}
	if _, ok := vs[guard]; ok {
		return &p, true
	}

	vs[guard] = struct{}{}
	visits.Set(p, vs)

	g.Set(p, 'X')
	g.Set(nextPoint, guard)
	return &nextPoint, false
}

var rotations = map[rune]rune{
	'<': '^',
	'^': '>',
	'>': 'v',
	'v': '<',
}
