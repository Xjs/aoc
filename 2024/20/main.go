package main

import (
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/Xjs/aoc/graph"
	"github.com/Xjs/aoc/grid"
)

var limit = 100

func main() {
	if len(os.Args) == 2 {
		l, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		limit = l
	}

	const cheatRadius = 20
	const debug = false

	g, err := grid.ReadRuneGrid(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	maze, ids, start, end := grid.Maze(g)

	log.Printf("start: %v, end: %v; points: %d, edges: %d", start, end, len(maze.Points), maze.NEdges())

	tstart := time.Now()
	reg, regularLength := graph.Dijkstras(*maze, ids[start], ids[end], 1)
	tstop := time.Now()

	log.Printf("regular length: %d (found in %v)", regularLength, tstop.Sub(tstart))

	// cheats := findCheats(g, 20)
	cheats := findCheatsFromSet(maze.Points, reg[0], cheatRadius)
	// cheats = checkCheats(cheats, g)

	cheatsPart1 := make(map[[2]grid.Point]int)
	for ps, cheat := range cheats {
		if d(ps[0], ps[1]) <= 2 {
			cheatsPart1[ps] = cheat
		}
	}

	countPart1 := 0
	times1, ns1 := histogram(cheatsPart1)
	for i, saving := range times1 {
		if debug {
			log.Printf("part1: %d picoseconds: %d cheats", saving, ns1[i])
		}
		if saving >= limit {
			countPart1 += ns1[i]
		}
	}

	countPart2 := 0
	times, ns := histogram(cheats)
	for i, saving := range times {
		if debug {
			log.Printf("%d picoseconds: %d cheats", saving, ns[i])
		}
		if saving >= limit {
			countPart2 += ns[i]
		}
	}

	log.Printf("part1: %d", countPart1)
	log.Printf("part2: %d", countPart2)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func d(p1, p2 grid.Point) int {
	δ := grid.Diff(p1, p2)
	return abs(δ.Dx) + abs(δ.Dy)
}

func findCheatsFromSet(points []grid.Point, visited []int, cheatRadius int) map[[2]grid.Point]int {
	result := make(map[[2]grid.Point]int)
	for i, id := range visited {
		p1 := points[id]
		// doesn't make sense to jump anything that's farther than 4 picoseconds apart
		for j := i + 4; j < len(visited); j++ {
			id2 := visited[j]
			if id2 == id {
				continue
			}

			p2 := points[id2]
			// absolute difference between these two points
			d := d(p1, p2)

			// Not eligible
			if d > cheatRadius {
				continue
			}

			// how much time crossing these two points needs in the regular path
			distInPath := abs(i - j)

			// Time we could potentially save by jumping
			saving := distInPath - d
			if saving < limit {
				continue
			}

			result[[2]grid.Point{p1, p2}] = saving
		}
	}

	return result
}

func findCheats(g *grid.Grid[rune], cheatRadius int) map[[2]grid.Point]int {
	cheats := make(map[[2]grid.Point]int)
	g.Foreach(func(p grid.Point) {
		if g.MustAt(p) == '#' {
			return
		}
		g.Foreach(func(p2 grid.Point) {
			if g.MustAt(p2) == '#' {
				return
			}
			d := d(p, p2)
			if d < 2 {
				return
			}
			if d > cheatRadius {
				return
			}

			cheats[[2]grid.Point{p, p2}] = d
		})
	})
	return cheats
}

func checkCheats(cheats map[[2]grid.Point]int, g *grid.Grid[rune]) map[[2]grid.Point]int {
	log.Printf("%d cheats identified. Pre-checking them...", len(cheats))

	maze, ids, _, _ := grid.Maze(g)

	sensibleCheats := make(map[[2]grid.Point]int)
	checked := 0
	for cheat, length := range cheats {
		checked++
		if checked%200 == 0 {
			pct := float64(checked) / float64(len(cheats))
			pctSens := float64(len(sensibleCheats)) / float64(len(cheats))
			log.Printf("%f %% checked. %f %% sensible", pct*100, pctSens*100)
		}

		_, l := graph.Dijkstras(*maze, ids[cheat[0]], ids[cheat[1]], 1)
		if l-length >= 50 {
			sensibleCheats[cheat] = length
		}
	}

	log.Printf("%d sensible cheats identified. Checking them...", len(sensibleCheats))

	lengths := make(map[[2]grid.Point]int)
	checked = 0
	for cheat, length := range sensibleCheats {
		checked++

		if checked%200 == 0 {
			pct := float64(checked) / float64(len(cheats))
			log.Printf("%f %% checked.", pct*100)
		}

		ng := g.Copy()
		nm, ids, start, end := grid.Maze(&ng)

		idx := ids[cheat[0]]
		idy := ids[cheat[1]]

		if nm.Edges[idx] == nil {
			nm.Edges[idx] = make(map[int]int)
		}
		nm.Edges[idx][idy] = length
		_, l := graph.Dijkstras(*nm, ids[start], ids[end], 1)

		lengths[cheat] = l
	}

	return lengths
}

func histogram(cheats map[[2]grid.Point]int) ([]int, []int) {
	hist := make(map[int]int)
	for _, saving := range cheats {
		hist[saving]++
	}

	var savings, ns []int
	for saving := range hist {
		savings = append(savings, saving)
	}

	sort.IntSlice(savings).Sort()
	ns = make([]int, len(savings))

	for i, saving := range savings {
		ns[i] = hist[saving]
	}

	return savings, ns
}
