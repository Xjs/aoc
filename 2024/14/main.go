package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Xjs/aoc/grid"
	"github.com/Xjs/aoc/parse"
)

func main() {
	var (
		rw = grid.Coordinate(101)
		rh = grid.Coordinate(103)
	)

	if len(os.Args) == 3 {
		w, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		rw = grid.Coordinate(w)

		h, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		rh = grid.Coordinate(h)
	}

	s := bufio.NewScanner(os.Stdin)
	robots := make(map[id]robot)
	currentID := id(0)
	g := grid.NewGrid[map[id]struct{}](rw, rh)
	for s.Scan() {
		pos, robot, err := parseRobot(s.Text())
		if err != nil {
			log.Fatal(err)
		}
		robots[currentID] = robot
		addRobot(&g, pos, currentID)
		currentID++
	}

	orig := g.Copy()

	// for i := 0; i < 100; i++ {
	// 	iter(&g, robots)
	// }

	product := 1
	for _, pair := range [][2]grid.Point{
		{grid.P(0, 0), grid.P(rw/2-1, rh/2-1)},
		{grid.P(rw/2+1, 0), grid.P(rw-1, rh/2-1)},
		{grid.P(0, rh/2+1), grid.P(rw/2-1, rh-1)},
		{grid.P(rw/2+1, rh/2+1), grid.P(rw-1, rh-1)},
	} {
		log.Print(pair[0], pair[1])
		product *= count(&g, pair[0], pair[1])
	}

	log.Printf("part1: %d", product)

	i := 0
	for {
		score := 0
		orig.Foreach(func(p grid.Point) {
			m := orig.MustAt(p)
			if len(m) == 0 {
				return
			}
			locscore := 0
			for _, pp := range orig.Environment8(p) {
				if len(orig.MustAt(pp)) > 0 {
					locscore++
				}
			}
			if locscore > 2 {
				score++
			}
		})
		if score > 200 {
			log.Printf("part2: %d seconds: %d\n", i, score)
			printCounts(&orig)
			break
		}

		iter(&orig, robots)
		i++
	}
}

func addRobot(g *grid.Grid[map[id]struct{}], pos grid.Point, currentID id) {
	m := g.MustAt(pos)
	if m == nil {
		m = make(map[id]struct{})
	}
	m[currentID] = struct{}{}
	g.Set(pos, m)
}

func removeRobot(g *grid.Grid[map[id]struct{}], pos grid.Point, currentID id) {
	m := g.MustAt(pos)
	if m == nil {
		return
	}
	delete(m, currentID)
	g.Set(pos, m)
}

func printCounts(g *grid.Grid[map[id]struct{}]) {
	countGrid := grid.NewGrid[int](g.Width(), g.Height())
	g.Foreach(func(p grid.Point) {
		countGrid.Set(p, len(g.MustAt(p)))
	})
	fmt.Print(strings.ReplaceAll(grid.StringIntGrid(countGrid), "0", "."))
}

func iter(g *grid.Grid[map[id]struct{}], robots map[id]robot) {
	newGrid := grid.NewGrid[map[id]struct{}](g.Width(), g.Height())
	g.Foreach(func(p grid.Point) {
		for theID := range g.MustAt(p) {
			vel := robots[theID].velocity
			newP := g.DeltaWrap(p, vel)
			addRobot(&newGrid, newP, theID)
		}
	})
	*g = newGrid
}

func count(g *grid.Grid[map[id]struct{}], min, max grid.Point) int {
	sum := 0
	for x := min.X; x <= max.X; x++ {
		for y := min.Y; y <= max.Y; y++ {
			sum += len(g.MustAt(grid.P(x, y)))
		}
	}
	return sum
}

type id int

type robot struct {
	velocity grid.Delta
}

func parseRobot(input string) (grid.Point, robot, error) {
	sp := strings.Fields(input)
	if len(sp) != 2 {
		return grid.Point{}, robot{}, errors.New("need exactly 2 fields")
	}

	pos, err := parse.IntList(strings.TrimPrefix(sp[0], "p="))
	if err != nil {
		return grid.Point{}, robot{}, err
	}

	if len(pos) != 2 {
		return grid.Point{}, robot{}, errors.New("need exactly 2 coordinates")
	}

	vel, err := parse.IntList(strings.TrimPrefix(sp[1], "v="))
	if err != nil {
		return grid.Point{}, robot{}, err
	}

	if len(vel) != 2 {
		return grid.Point{}, robot{}, errors.New("need exactly 2 velocity components")
	}

	return grid.P(grid.Coordinate(pos[0]), grid.Coordinate(pos[1])), robot{velocity: grid.D(vel[0], vel[1])}, nil
}
