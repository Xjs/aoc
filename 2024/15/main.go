package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"unicode"

	"github.com/Xjs/aoc/grid"
)

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	sp := bytes.Split(input, []byte("\n\n"))
	if len(sp) < 2 {
		log.Fatal("need: grid with boxes, robot and walls; list of movements (<^>v)")
	}

	g, err := grid.ReadRuneGrid(bytes.NewReader(bytes.TrimSpace(sp[0])))
	if err != nil {
		log.Fatal(err)
	}

	g2 := part2grid(g)

	log.Printf("w = %d, h = %d", g.Width(), g.Height())
	fmt.Println(grid.StringCharGrid(*g))

	log.Printf("part2: w = %d, h = %d", g2.Width(), g2.Height())
	fmt.Println(grid.StringCharGrid(*g2))

	var robot grid.Point
	g.Foreach(func(p grid.Point) {
		if g.MustAt(p) == '@' {
			robot = p
		}
	})

	for _, dir := range string(bytes.TrimSpace(sp[1])) {
		if unicode.IsSpace(dir) {
			// Apparently the include newlines somewhere inside the list…
			continue
		}
		d, ok := grid.GeneralDirections[dir]
		if !ok {
			log.Fatalf("illegal direction %q", dir)
		}
		robot, err = move(g, robot, d)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(grid.StringCharGrid(*g))

	log.Printf("part1: %d", score(g))

	var robot2 grid.Point
	g2.Foreach(func(p grid.Point) {
		if g2.MustAt(p) == '@' {
			robot2 = p
		}
	})

	for _, dir := range string(bytes.TrimSpace(sp[1])) {
		if unicode.IsSpace(dir) {
			// Apparently the include newlines somewhere inside the list…
			continue
		}
		d, ok := grid.GeneralDirections[dir]
		if !ok {
			log.Fatalf("illegal direction %q", dir)
		}
		robot2, err = move2(g2, robot2, d)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(grid.StringCharGrid(*g2))

	log.Printf("part2: %d", score(g2))
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func move(g *grid.Grid[rune], p grid.Point, dir grid.Delta) (grid.Point, error) {
	if abs(dir.Dx)+abs(dir.Dy) != 1 {
		return grid.Point{}, errors.New("only cardinal movements of size 1 allowed")
	}

	next, err := g.Delta(p, dir)
	if err != nil {
		// we don't move
		return p, nil
	}

	switch tt := g.MustAt(next); tt {
	case '#':
		return p, nil
	case '0', 'O':
		if _, err := move(g, next, dir); err != nil {
			return p, err
		}
	case '.':
	default:
		return grid.Point{}, fmt.Errorf("cannot move onto unknown tile %q", tt)
	}

	// Maybe something has moved next to us…
	if g.MustAt(next) != '.' {
		return p, nil
	}

	us := g.MustAt(p)
	g.Set(next, us)
	g.Set(p, '.')
	return next, nil
}

func score(g *grid.Grid[rune]) int {
	score := 0
	g.Foreach(func(p grid.Point) {
		switch g.MustAt(p) {
		case '0', 'O', '[':
			score += int(100*p.Y + p.X)
		}
	})
	return score
}
