package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/Xjs/aoc/graph"
	"github.com/Xjs/aoc/grid"
)

var numPos = make(map[rune]grid.Point)
var arrowPos = make(map[rune]grid.Point)
var shortestPathsNum map[[2]rune][][]grid.Delta
var shortestPathsArrow map[[2]rune][][]grid.Delta

var layer2phrases = make(map[string][]string)
var shortestPathsNum2 = make(map[[2]rune]string)

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

	phrases := make(map[string][]string)
	for _, sp := range shortestPathsNum {
		for _, p := range sp {
			phrases[toRunes(p)+"A"] = nil
		}
	}

	for phrase := range phrases {
		phrases[phrase] = resolve(arrowPos, shortestPathsArrow, phrase)
	}

	for n, sp := range shortestPathsNum {
		var path string
		for _, p := range sp {
			pp := filterShortest(phrases[toRunes(p)+"A"])[0]
			if path == "" || len(path) > len(pp) {
				path = toRunes(p)
			}
		}
		shortestPathsNum2[n] = path
	}

	for _, sp := range shortestPathsArrow {
		for _, p := range sp {
			cur := toRunes(p)
			layer2phrases[cur] = filterShortest(resolve(arrowPos, shortestPathsArrow, "A"+cur+"A"))
		}
	}
	layer2phrases[""] = []string{"A"}
	// for _, paths := range phrases {
	// 	for _, phrase := range paths {
	// 		for _, sp := range strings.Split(phrase, "A") {
	// 			layer2phrases[sp] = filterShortest(resolve(arrowPos, shortestPathsArrow, sp+"A"))
	// 		}
	// 	}
	// }

	// log.Print(layer2phrases)
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
	iter1 := 2
	iter2 := 25
	sum := 0
	sum2 := 0
	sum22 := 0
	for s.Scan() {
		t := s.Text()
		c := complexity(t, iter1)
		sum += c
		log.Printf("%s: %d", t, c)

		c2 := complexity(t, iter2)
		sum2 += c2
		log.Printf("%s (%d): %d", t, iter2, c2)

		// c3 := complexity2(t, iter2)
		// sum22 += c3
		// log.Printf("%s (%d, correct): %d", t, iter2, c3)

		// if c2 != c3 {
		// 	log.Fatalf("%s not correct", t)
		// }
	}
	log.Printf("part1: %d", sum)
	log.Printf("part2: %d", sum2)
	log.Printf("part2_alt: %d", sum22)
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

func complexity(code string, iter int) int {
	n, err := strconv.Atoi(code[:len(code)-1])
	if err != nil {
		// meh, can't be bothered
		panic(err)
	}

	pathsNum := resolve(numPos, shortestPathsNum, "A"+code)
	var pathsArrow []string = pathsNum
	for i := 0; i < 1; i++ {
		pathsArrow = resolveArrows(pathsArrow)
	}

	lim := iter
	if lim > depth+1 {
		lim = depth + 1
	}

	for i := 1; i < lim; i++ {
		pathsArrow = resolveArrows2(pathsArrow)
	}

	shortest := math.MaxInt
	for _, pp := range pathsArrow {
		r := resolveArrows4(pp, iter-depth-1)
		if r < shortest {
			shortest = r
		}
	}

	return n * shortest
}

func complexity2(code string, iter int) int {
	n, err := strconv.Atoi(code[:len(code)-1])
	if err != nil {
		// meh, can't be bothered
		panic(err)
	}

	pathsNum := resolve(numPos, shortestPathsNum, "A"+code)
	var pathsArrow []string = pathsNum
	for i := 0; i < 1; i++ {
		pathsArrow = resolveArrows(pathsArrow)
	}

	for i := 1; i < iter; i++ {
		pathsArrow = resolveArrows2(pathsArrow)
	}
	pathsArrow = filterShortest(pathsArrow)

	return n * len(pathsArrow[0])
}

type mem struct {
	path string
	iter int
}

var memo = make(map[mem]int)

func resolveArrows3(path string, iter int) int {
	if iter <= 0 {
		return len(path)
	}
	if i, ok := memo[mem{path, iter}]; ok {
		return i
	}

	sps := splitA(path)
	if len(sps) == 1 {
		res := resolveArrows3(resolveArrows2one(path), iter-1)
		memo[mem{path, iter}] = res
		return res
	}

	l := 0
	for _, sp := range sps {
		l += resolveArrows3(sp+"A", iter)
	}

	return l
}

const depth = 4

func resolveArrows4(path string, iter int) int {

	if iter <= 0 {
		return len(path)
	}
	if i, ok := memo[mem{path, iter}]; ok {
		return i
	}
	if iter%depth != 0 {
		res := resolveArrows4(resolveArrows2one(path), iter-1)
		memo[mem{path, iter}] = res
		return res
	}

	sps := splitA(path)
	if len(sps) == 1 {
		np := filterShortest(resolveArrows2depth([]string{path}, depth))[0]
		res := resolveArrows4(np, iter-depth)
		memo[mem{path, iter}] = res
		return res
	}

	l := 0
	for _, sp := range sps {
		l += resolveArrows3(sp+"A", iter)
	}

	return l
}

func resolveArrows2(ps []string) []string {
	result := make([]string, len(ps))
	for i, p := range ps {
		result[i] = resolveArrows2one(p)
	}
	return result
}

func splitA(s string) []string {
	var splices []string
	for len(s) > 0 {
		var i int
		for i = 0; s[i] != 'A'; i++ {
		}
		splices = append(splices, s[:i])
		s = s[i+1:]
	}
	return splices
}

func resolveArrows2one(p string) string {
	// log.Printf("input: %v (%q)", p, splitA(p))
	newP := new(strings.Builder)
	for _, sp := range splitA(p) {
		pp := layer2phrases[sp]
		if pp == nil {
			log.Fatal(sp)
		}
		if pp[0] == "" {
			log.Fatal(sp)
		}
		newP.WriteString(layer2phrases[sp][0])
	}
	np := newP.String()
	// log.Printf("result: %s", np)
	return np
}

func resolveArrows2depth(ps []string, depth int) []string {
	// log.Printf("input: %v (%q)", p, splitA(p))
	if depth < 1 {
		return ps
	}

	nps := make([][][]string, len(ps))
	for i, p := range ps {
		splits := splitA(p)
		nps[i] = make([][]string, len(splits))
		for j, sp := range splits {
			nps[i][j] = make([]string, len(layer2phrases[sp]))
			copy(nps[i][j], layer2phrases[sp])
		}
	}
	return flatten(nps)
}

func flatten(nps [][][]string) []string {
	var result []string
	for _, nnps := range nps {
		result = append(result, flatten2(nnps)...)
	}
	return result
}

func flatten2(nnps [][]string) []string {
	if len(nnps) == 0 {
		return []string{""}
	}
	start, nnps := nnps[0], nnps[1:]
	var result []string
	for _, st := range start {
		for _, end := range flatten2(nnps) {
			result = append(result, st+end)
		}
	}
	return result
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
		log.Fatal(d)
		return 'X'
	}
}

func toRunes(ds []grid.Delta) string {
	res := make([]rune, len(ds))
	for i, d := range ds {
		res[i] = toRune(d)
	}
	return string(res)
}
