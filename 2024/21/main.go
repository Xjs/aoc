package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/Xjs/aoc/graph"
	"github.com/Xjs/aoc/grid"
)

var numPos = make(map[rune]grid.Point)
var arrowPos = make(map[rune]grid.Point)
var shortestPathsNum map[[2]rune][][]grid.Delta
var shortestPathsArrow map[[2]rune][][]grid.Delta

func init() {
	numpad, err := grid.GridFrom[rune]([][]rune{
		{'7', '8', '9'},
		{'4', '5', '6'},
		{'1', '2', '3'},
		{'#', '0', 'A'},
	})
	if err != nil {
		log.Fatal(err)
	}

	numPos = make(map[rune]grid.Point)
	numpad.Foreach(func(p grid.Point) {
		numPos[numpad.MustAt(p)] = p
	})

	arrowPad, err := grid.GridFrom[rune]([][]rune{
		{'#', '^', 'A'},
		{'<', 'v', '>'},
	})
	if err != nil {
		log.Fatal(err)
	}

	arrowPos = make(map[rune]grid.Point)
	arrowPad.Foreach(func(p grid.Point) {
		arrowPos[arrowPad.MustAt(p)] = p
	})

	shortestPathsNum = makeShortestPaths(numpad, numPos)
	shortestPathsArrow = makeShortestPaths(arrowPad, arrowPos)
}

func makeShortestPaths(pad grid.Grid[rune], pos map[rune]grid.Point) map[[2]rune][][]grid.Delta {
	var runes []rune
	pad.Foreach(func(p grid.Point) {
		r := pad.MustAt(p)
		if r == '#' {
			return
		}
		runes = append(runes, r)
	})

	res := make(map[[2]rune][][]grid.Delta)
	maze, ids, _, _ := grid.Maze(&pad)
	for _, r1 := range runes {
		for _, r2 := range runes {
			ps, _ := graph.Dijkstras(*maze, ids[pos[r1]], ids[pos[r2]], math.MaxInt)
			for _, path := range ps {
				var pp []grid.Delta
				for i := 0; i < len(path)-1; i++ {
					p1 := maze.Points[path[i]]
					p2 := maze.Points[path[i+1]]

					d := grid.Diff(p2, p1)
					pp = append(pp, d)
				}
				res[[2]rune{r1, r2}] = append(res[[2]rune{r1, r2}], pp)
			}
		}
	}
	return res
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	sum := 0
	for s.Scan() {
		t := s.Text()
		c := complexity(t)
		log.Printf("%s: %d", t, c)
		sum += c
	}
	log.Printf("part1: %d", sum)
}

func filterShortest(ps []string) []string {
	shortest := math.MaxInt
	for _, p := range ps {
		if len(p) < shortest {
			shortest = len(p)
		}
	}

	var result []string
	for _, p := range ps {
		if len(p) == shortest {
			result = append(result, p)
		}
	}

	return result
}

func complexity(code string) int {
	n, err := strconv.Atoi(code[:len(code)-1])
	if err != nil {
		// meh, can't be bothered
		panic(err)
	}

	pathsNum := resolve(numPos, shortestPathsNum, "A"+code)
	var pathsArrow []string = pathsNum
	for i := 0; i < 2; i++ {
		pathsArrow = resolveArrows(pathsArrow)
	}
	pathsArrow = filterShortest(pathsArrow)

	return len(pathsArrow[0]) * n
}

func resolveArrows(ps []string) []string {
	var paths []string
	for _, p := range ps {
		paths = append(paths, resolve(arrowPos, shortestPathsArrow, "A"+p)...)
	}
	return paths
}

func resolve(pos map[rune]grid.Point, s map[[2]rune][][]grid.Delta, code string) []string {
	if len(code) < 2 {
		return []string{""}
	}

	from := rune(code[0])
	to := rune(code[1])

	var rr [][]rune
	for _, p := range s[[2]rune{from, to}] {
		var onePath []rune
		for _, d := range p {
			onePath = append(onePath, toRune(d))
		}
		onePath = append(onePath, 'A')
		rr = append(rr, onePath)
	}

	if from == to {
		rr = append(rr, []rune{'A'})
	}

	var result []string
	for _, p := range rr {
		ps := resolve(pos, s, code[1:])
		for _, follow := range ps {
			complete := append(p, []rune(follow)...)
			result = append(result, string(complete))
		}
	}

	return result
}

func path(pos map[rune]grid.Point, from, to rune) []grid.Delta {
	pFrom := pos[from]
	pTo := pos[to]

	var result []grid.Delta
	// When going left, first go to the correct row, then to the correct column
	if pFrom.X > pTo.X {
		result = append(result, goToRow(pFrom, pTo)...)
		result = append(result, goToCol(pFrom, pTo)...)
	} else {
		result = append(result, goToCol(pFrom, pTo)...)
		result = append(result, goToRow(pFrom, pTo)...)
	}
	return result
}

func goToCol(p1, p2 grid.Point) []grid.Delta {
	dir := 1
	if p1.X > p2.X {
		dir = -1
	}
	var result []grid.Delta
	for i := int(p1.X); i != int(p2.X); i += dir {
		result = append(result, grid.D(dir, 0))
	}
	return result
}

func goToRow(p1, p2 grid.Point) []grid.Delta {
	dir := 1
	if p1.Y > p2.Y {
		dir = -1
	}
	var result []grid.Delta
	for i := int(p1.Y); i != int(p2.Y); i += dir {
		result = append(result, grid.D(0, dir))
	}
	return result

}

func toRune(d grid.Delta) rune {
	switch d {
	case grid.GeneralDirections['<']:
		return '<'
	case grid.GeneralDirections['>']:
		return '>'
	case grid.GeneralDirections['^']:
		return '^'
	case grid.GeneralDirections['v']:
		return 'v'
	default:
		log.Print(d)
		return 'X'
	}
}
