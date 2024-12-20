package grid

import (
	"github.com/Xjs/aoc/graph"
)

// Maze sets up g as a aoc-style maze, making everything that contains a '#' uncrossable
// and everything else reachable with edge weight 1.
// To quickly identify the point IDs, a map is returned.
// If there are S and E points, it will return them as start and end
func Maze(g *Grid[rune]) (gr *graph.Graph[Point, int], ids map[Point]int, start, end Point) {
	gr = new(graph.Graph[Point, int])
	gr.Edges = make(map[int]map[int]int)
	id := 0
	ids = make(map[Point]int)

	g.Foreach(func(p Point) {
		r := g.MustAt(p)
		if r == '#' {
			return
		}
		if r == 'S' {
			start = p
		}
		if r == 'E' {
			end = p
		}

		gr.Points = append(gr.Points, p)
		ids[p] = id
		id++
	})

	for i, p := range gr.Points {
		for _, nb := range g.Environment4(p) {
			if g.MustAt(nb) == '#' {
				continue
			}
			if gr.Edges[i] == nil {
				gr.Edges[i] = make(map[int]int)
			}
			gr.Edges[i][ids[nb]] = 1
		}
	}

	return
}
