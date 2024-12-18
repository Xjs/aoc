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

// Dijkstra finds the shortest paths from source to target and returns it as well as its weight.
func Dijkstras[T any, W EdgeWeight](g Graph[T, W], source, target int, limit int) ([][]int, W) {
	q := newPQ()
	// unvisited := make(map[int]struct{})
	// distances := make(map[int]float64)
	previous := make(map[int][]int)

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

			if change, equal := q.decrease(id(nb), alt); change {
				if equal {
					previous[nb] = append(previous[nb], int(current))
				} else {
					previous[nb] = []int{int(current)}
				}
			}
		}

		if current == id(target) {
			break
		}
	}

	if _, ok := previous[target]; !ok {
		return nil, -1
	}

	gs := graphs(source, target, previous, limit)
	sum := W(0)
	cur := source
	for i := 1; i < len(gs[0]); i++ {
		next := gs[0][i]
		sum += g.Edges[cur][next]
		cur = next
	}

	return gs, sum
}

func graphs(source, target int, previous map[int][]int, limit int) [][]int {
	if target == source {
		return [][]int{{target}}
	}

	var result [][]int
	for _, i := range previous[target] {
		for _, path := range graphs(source, i, previous, limit) {
			result = append(result, append(path, target))
			if len(result) >= limit {
				return result
			}
		}
	}
	return result
}
