package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Xjs/aoc/graph"
	"github.com/Xjs/aoc/grid"
	"github.com/Xjs/aoc/parse"
)

func main() {
	w := grid.Coordinate(71)
	h := grid.Coordinate(71)
	limit := 1024

	if len(os.Args) == 4 {
		ww, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		w = grid.Coordinate(ww)

		hh, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		h = grid.Coordinate(hh)

		ll, err := strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatal(err)
		}
		limit = ll
	}

	log.Printf("w: %d, h: %d, limit: %d", w, h, limit)

	g := grid.NewGrid[rune](w, h)
	g.Foreach(func(p grid.Point) { g.Set(p, '.') })

	done := 0
	var bs []grid.Point

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		is, err := parse.IntList(s.Text())
		if err != nil {
			log.Fatal(err)
		}
		if len(is) != 2 {
			log.Fatal("length of coordinate must be 2")
		}
		done++

		bs = append(bs, grid.P(grid.Coordinate(is[0]), grid.Coordinate(is[1])))
	}

	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	for _, p := range bs[:limit] {
		g.Set(p, '#')
	}

	fmt.Println(grid.StringCharGrid(g))

	l := pathLength(&g)
	log.Printf("part1: %d", l)

	for _, p := range bs[limit:] {
		g.Set(p, '#')
		l := pathLength(&g)
		if l < 0 {
			log.Fatal("part2: %v", p)
		}
	}
}

func pathLength(g *grid.Grid[rune]) int {
	var gr graph.Graph[grid.Point, int]
	gr.Edges = make(map[int]map[int]int)
	id := 0
	ids := make(map[grid.Point]int)
	g.Foreach(func(p grid.Point) {
		if g.MustAt(p) == '#' {
			return
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

	nEdges := 0
	for _, edge := range gr.Edges {
		nEdges += len(edge)
	}

	start := grid.P(0, 0)
	end := grid.P(g.Width()-1, g.Height()-1)

	_, l := graph.Dijkstras(gr, ids[start], ids[end], 1)

	return l
}
