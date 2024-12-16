package graph

import (
	"container/heap"
	"math"
)

type EdgeWeight interface {
	int | float64
}

// A Graph is a network of points of type T and weighted, directed edges between them.
// Edges are represented by pairs of indexes into the Points slice.
// An undirected edge is represented as a symmetric pair of edges.
type Graph[T any, W EdgeWeight] struct {
	Points []T
	Edges  map[int]map[int]W
}

func NewUndirectedGraph[T any, W EdgeWeight](points []T, edges map[[2]int]W) Graph[T, W] {
	directedEdges := make(map[int]map[int]W)
	for edge, weight := range edges {
		x, y := edge[0], edge[1]
		if directedEdges[x] == nil {
			directedEdges[x] = make(map[int]W)
		}
		directedEdges[x][y] = weight

		if directedEdges[y] == nil {
			directedEdges[y] = make(map[int]W)
		}
		directedEdges[y][x] = weight
	}
	return Graph[T, W]{
		Points: points,
		Edges:  directedEdges,
	}
}

// Dijkstra finds the shortest path from source to target and returns it as well as its weight.
func Dijkstra[T any, W EdgeWeight](g Graph[T, W], source, target int) ([]int, W) {
	q := newPQ()
	// unvisited := make(map[int]struct{})
	// distances := make(map[int]float64)
	previous := make(map[int]int)

	for i := range g.Points {
		prio := math.Inf(1)
		if i == source {
			prio = 0
		}
		heap.Push(q, id(i))
		q.update(id(i), prio)
	}

	heap.Init(q)

	for {
		if q.Len() == 0 {
			break
		}

		current, curDist := q.pop()

		if math.IsInf(curDist, 1) {
			break
		}

		for nb, w := range g.Edges[int(current)] {
			alt := curDist + float64(w)

			if q.decrease(id(nb), alt) {
				previous[nb] = int(current)
			}
		}

		if current == id(target) {
			break
		}
	}

	if _, ok := previous[target]; !ok {
		return nil, -1
	}

	hops := 0
	for current := target; current != source; current = previous[current] {
		hops++
	}
	hops++

	path := make([]int, hops)
	current := target
	for i := hops - 1; i >= 0; i-- {
		path[i] = current
		current = previous[current]
	}

	sum := W(0)
	cur := source
	for i := 1; i < len(path); i++ {
		sum += g.Edges[cur][path[i]]
		cur = path[i]
	}

	return path, sum
}
