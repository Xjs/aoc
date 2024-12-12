package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Xjs/aoc/grid"
)

func main() {
	input, err := grid.ReadRuneGrid(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	sum := 0
	sum2 := 0
	for _, region := range regions(input) {
		ar := region.area()
		pr := region.perimeter()
		pr2 := region.perimeter2()
		log.Printf("region: %d %d %d", ar, pr, pr2)
		sum += ar * pr
		sum2 += ar * pr2
	}
	log.Printf("part1: %d", sum)
	log.Printf("part2: %d", sum2)
}

type region struct {
	letter rune
	edges  map[grid.Point]map[direction]struct{}
}

func (r region) area() int {
	return len(r.edges)
}

func (r region) perimeter() int {
	sum := 0

	for _, edges := range r.edges {
		sum += len(edges)
	}

	return sum
}

func (r region) perimeter2() int {
	edgeIDs := make(map[grid.Point]map[direction]int)
	edgesSeen := make(map[grid.Point]map[direction]bool)
	var points []grid.Point

	for p, edges := range r.edges {
		if len(edges) == 0 {
			continue
		}
		edgeIDs[p] = make(map[direction]int)
		for dir := range edges {
			edgeIDs[p][dir] = 0
		}
		points = append(points, p)

		edgesSeen[p] = make(map[direction]bool)
	}

	id := 1

	var current grid.Point
	for len(points) > 0 {
		current, points = points[0], points[1:]
		edges, ok := edgeIDs[current]
		if !ok {
			continue
		}

		if edgesSeen[current] == nil {
			edgesSeen[current] = make(map[direction]bool)
		}

		for _, dir := range []direction{top, right, bottom, left} {
			if _, has := edges[dir]; has {
				traceEdge(current, dir, id, edgeIDs, edgesSeen)

				id++
			}
		}
	}

	edgeIDset := make(map[int]struct{})

	for _, ids := range edgeIDs {
		for _, id := range ids {
			edgeIDset[id] = struct{}{}
		}
	}

	return len(edgeIDset)
}

func traceEdge(curr grid.Point, checking direction, id int, edgeIDs map[grid.Point]map[direction]int, edgesSeen map[grid.Point]map[direction]bool) {
	if edgesSeen[curr][checking] {
		return
	}

	var mvmtA, mvmtB direction
	switch checking {
	case top:
		mvmtA = left
		mvmtB = right
	case left:
		mvmtA = bottom
		mvmtB = top
	case bottom:
		mvmtA = right
		mvmtB = left
	case right:
		mvmtA = top
		mvmtB = bottom
	}

	traceEdgeOnesided(curr, checking, mvmtA, id, edgeIDs, edgesSeen)
	edgesSeen[curr][checking] = false
	traceEdgeOnesided(curr, checking, mvmtB, id, edgeIDs, edgesSeen)
}

func traceEdgeOnesided(curr grid.Point, checking, movement direction, id int, edgeIDs map[grid.Point]map[direction]int, edgesSeen map[grid.Point]map[direction]bool) {
	if edgesSeen[curr][checking] {
		return
	}

	for {
		_, ok := edgeIDs[curr]
		if !ok {
			break
		}

		_, has := edgeIDs[curr][checking]
		if !has {
			break
		}

		if existingID := edgeIDs[curr][checking]; existingID != 0 && existingID != id {
			panic(fmt.Sprintf("edgeID seen double: %d/%d", existingID, id))
		}

		edgeIDs[curr][checking] = id
		edgesSeen[curr][checking] = true

		if _, ok := edgeIDs[curr][movement]; ok {
			break
		}

		curr = move(curr, movement)
	}
}

func regions(g *grid.Grid[rune]) []region {
	seen := grid.NewGrid[bool](g.Width(), g.Height())

	var result []region

	g.Foreach(func(p grid.Point) {
		if seen.MustAt(p) {
			return
		}

		result = append(result, traceRegion(p, g, &seen))
	})

	return result
}

func traceRegion(p grid.Point, g *grid.Grid[rune], seen *grid.Grid[bool]) region {
	r := region{
		letter: g.MustAt(p),
		edges:  make(map[grid.Point]map[direction]struct{}),
	}

	points := []grid.Point{p}
	current := p
	for len(points) > 0 {
		current, points = points[0], points[1:]

		if seen.MustAt(current) {
			continue
		}

		r.edges[current] = make(map[direction]struct{})
		if current.X == 0 {
			r.edges[current][left] = struct{}{}
		}
		if current.X == g.Width()-1 {
			r.edges[current][right] = struct{}{}
		}
		if current.Y == 0 {
			r.edges[current][top] = struct{}{}
		}
		if current.Y == g.Height()-1 {
			r.edges[current][bottom] = struct{}{}
		}

		seen.Set(current, true)

		neighbours := g.Environment4(current)

		for _, nb := range neighbours {
			if g.MustAt(nb) != r.letter {
				d, _ := dir(current, nb)
				r.edges[current][d] = struct{}{}

				continue
			}

			if seen.MustAt(nb) {
				continue
			}

			points = append(points, nb)
		}
	}

	return r
}

type direction rune

const (
	noDirection direction = iota
	top
	right
	bottom
	left
)

func dir(a, b grid.Point) (direction, error) {
	if a.X == b.X {
		switch int(a.Y) - int(b.Y) {
		case -1:
			return bottom, nil
		case 1:
			return top, nil
		}
	}
	if a.Y == b.Y {
		switch int(a.X) - int(b.X) {
		case -1:
			return right, nil
		case 1:
			return left, nil
		}
	}

	return noDirection, errors.New("points are not directly adjacent")
}

func move(p grid.Point, dir direction) grid.Point {
	switch dir {
	case left:
		p.X--
	case right:
		p.X++
	case top:
		p.Y--
	case bottom:
		p.Y++
	}
	return p
}
