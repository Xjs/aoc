package main

import (
	"log"
	"os"

	"github.com/Xjs/aoc/grid"
)

func main() {
	input, err := grid.ReadDigitGrid(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	sum := 0
	sum2 := 0
	input.Foreach(func(p grid.Point) {
		if input.MustAt(p) == 0 {
			sum += trailheadScore1(input, p)
			sum2 += trailheadScore2(input, p)
		}
	})
	log.Printf("part1: %d", sum)
	log.Printf("part2: %d", sum2)
}

func trailheadScore2(g *grid.Grid[int], start grid.Point) int {
	if g.MustAt(start) != 0 {
		return 0
	}

	score := 0

	points := []grid.Point{start}
	for len(points) > 0 {
		var newPoints []grid.Point
		for _, point := range points {
			for _, pp := range hike(g, point) {
				if g.MustAt(pp) == 9 {
					score++
				} else {
					newPoints = append(newPoints, pp)
				}
			}
		}
		points = newPoints
	}

	return score
}

func trailheadScore1(g *grid.Grid[int], start grid.Point) int {
	if g.MustAt(start) != 0 {
		return 0
	}

	points := make(map[grid.Point]struct{})
	scorePoints := make(map[grid.Point]struct{})
	points[start] = struct{}{}
	for len(points) > 0 {
		newPoints := make(map[grid.Point]struct{})
		for point := range points {
			for _, pp := range hike(g, point) {
				if g.MustAt(pp) == 9 {
					scorePoints[pp] = struct{}{}
				} else {
					newPoints[pp] = struct{}{}
				}
			}
		}
		points = newPoints
	}

	return len(scorePoints)
}

func hike(g *grid.Grid[int], p grid.Point) []grid.Point {
	vp := g.MustAt(p)

	var result []grid.Point
	for _, env := range g.Environment4(p) {
		v := g.MustAt(env)
		if v-vp == 1 {
			result = append(result, env)
		}
	}

	return result
}
