package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Xjs/aoc/graph"
	"github.com/Xjs/aoc/grid"
)

type node struct {
	p   grid.Point
	dir grid.Delta
}

const debug = false

func main() {
	g, err := grid.ReadRuneGrid(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	var nodes []node
	nodeIDs := make(map[grid.Point]map[grid.Delta]int)
	edges := make(map[int]map[int]int)

	id := 0

	var source, target grid.Point
	g.Foreach(func(p grid.Point) {
		switch tt := g.MustAt(p); tt {
		case '.', 'S', 'E':
			if tt == 'S' {
				source = p
			}
			if tt == 'E' {
				target = p
			}
			for _, dir := range grid.GeneralDirections {
				nodes = append(nodes, node{
					p:   p,
					dir: dir,
				})
				if nodeIDs[p] == nil {
					nodeIDs[p] = make(map[grid.Delta]int)
				}
				nodeIDs[p][dir] = id
				id++
			}
		}
	})

	for p := range nodeIDs {
		for _, p2 := range g.Environment4(p) {
			switch g.MustAt(p2) {
			default:
				continue
			case '.', 'S', 'E':
			}

			dir := grid.Diff(p2, p)
			nodeID := nodeIDs[p][dir]
			nodeID2 := nodeIDs[p2][dir]

			if edges[nodeID] == nil {
				edges[nodeID] = make(map[int]int)
			}
			edges[nodeID][nodeID2] = 1
		}
	}

	for p := range nodeIDs {
		for glyph, dir := range grid.GeneralDirections {
			nextDir := grid.GeneralDirections[rotations[glyph]]
			nodeID := nodeIDs[p][dir]
			nextNodeID := nodeIDs[p][nextDir]
			if edges[nodeID] == nil {
				edges[nodeID] = make(map[int]int)
			}
			edges[nodeID][nextNodeID] = 1000
			if edges[nextNodeID] == nil {
				edges[nextNodeID] = make(map[int]int)
			}
			edges[nextNodeID][nodeID] = 1000
		}
	}

	gr := graph.Graph[node, int]{
		Points: nodes,
		Edges:  edges,
	}

	startID := nodeIDs[source][grid.GeneralDirections['>']]
	var targetIDs []int
	for _, dir := range grid.GeneralDirections {
		targetIDs = append(targetIDs, nodeIDs[target][dir])
	}

	var shortestPath []int
	shortestPathLength := -1
	for _, targetID := range targetIDs {
		path, length := graph.Dijkstra(gr, startID, targetID)
		if shortestPathLength == -1 || length < shortestPathLength {
			shortestPathLength = length
			shortestPath = path
		}
	}

	log.Printf("part1: %d", shortestPathLength)

	if debug {
		printPath(nodes, shortestPath)
	}
}

var rotations = map[rune]rune{
	'<': '^',
	'^': '>',
	'>': 'v',
	'v': '<',
}

func printPath(nodes []node, path []int) {
	for _, hop := range path {
		fmt.Println(nodes[hop])
	}
}
