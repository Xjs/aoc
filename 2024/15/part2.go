package main

import (
	"errors"

	"github.com/Xjs/aoc/grid"
)

func part2grid(g *grid.Grid[rune]) *grid.Grid[rune] {
	g2 := grid.NewGrid[rune](g.Width()*2, g.Height())

	g.Foreach(func(p grid.Point) {
		p2 := grid.P(2*p.X, p.Y)
		pp := grid.P(2*p.X+1, p.Y)
		switch g.MustAt(p) {
		case '#':
			g2.Set(p2, '#')
			g2.Set(pp, '#')
		case 'O':
			g2.Set(p2, '[')
			g2.Set(pp, ']')
		case '.':
			g2.Set(p2, '.')
			g2.Set(pp, '.')
		case '@':
			g2.Set(p2, '@')
			g2.Set(pp, '.')
		}
	})

	return &g2
}

func getNext(g *grid.Grid[rune], p grid.Point, dir grid.Delta) (grid.Point, *grid.Point, error) {
	next, err := g.Delta(p, dir)
	if err != nil {
		return p, nil, err
	}

	us := g.MustAt(p)

	var next2 *grid.Point
	if dir.Dy != 0 {
		var boxDir grid.Delta
		switch us {
		case '[':
			boxDir = grid.GeneralDirections['>']
		case ']':
			boxDir = grid.GeneralDirections['<']
		}
		switch us {
		case '[', ']':
			p2, err := g.Delta(p, boxDir)
			if err != nil {
				return p, nil, err
			}

			n2, err := g.Delta(p2, dir)
			if err != nil {
				return p, nil, err
			}

			next2 = &n2
		}
	}

	return next, next2, nil
}

func canMove(g *grid.Grid[rune], p grid.Point, dir grid.Delta) (bool, error) {
	if abs(dir.Dx)+abs(dir.Dy) != 1 {
		return false, errors.New("only cardinal movements of size 1 allowed")
	}

	next, next2, err := getNext(g, p, dir)
	if err != nil {
		// we can't move because there is no next point.
		return false, nil
	}

	nextTile := g.MustAt(next)
	var nextTile2 rune
	if next2 != nil {
		nextTile2 = g.MustAt(*next2)
	}

	if nextTile == '#' {
		return false, nil
	}

	if nextTile2 == '#' {
		return false, nil
	}

	var canMoveTile1, canMoveTile2 bool
	if isBox2(nextTile2) {
		cm, err := canMove(g, *next2, dir)
		if err != nil {
			return false, err
		}
		canMoveTile2 = cm
	}

	if isBox2(nextTile) {
		cm, err := canMove(g, next, dir)
		if err != nil {
			return false, err
		}
		canMoveTile1 = cm
	}

	if nextTile == '.' {
		canMoveTile1 = true
	}

	if nextTile2 == '.' {
		canMoveTile2 = true
	}

	if next2 == nil {
		canMoveTile2 = true
	}

	return canMoveTile1 && canMoveTile2, nil
}

func isBox2(r rune) bool {
	return r == '[' || r == ']'
}

func move2(g *grid.Grid[rune], p grid.Point, dir grid.Delta) (grid.Point, error) {
	if abs(dir.Dx)+abs(dir.Dy) != 1 {
		return grid.Point{}, errors.New("only cardinal movements of size 1 allowed")
	}

	cm, err := canMove(g, p, dir)
	if err != nil {
		return grid.Point{}, err
	}
	if !cm {
		// we don't move
		return p, nil
	}

	next, next2, err := getNext(g, p, dir)
	if err != nil {
		// we can't move because there is no next point.
		return p, nil
	}

	if isBox2(g.MustAt(next)) {
		if _, err := move2(g, next, dir); err != nil {
			return p, err
		}
	}
	if next2 != nil && isBox2(g.MustAt(*next2)) {
		if _, err := move2(g, *next2, dir); err != nil {
			return p, err
		}
	}

	us := g.MustAt(p)
	g.Set(next, us)
	g.Set(p, '.')

	if next2 != nil {
		d := grid.Diff(*next2, next)
		p2, err := g.Delta(p, d)
		if err != nil {
			return p, err
		}

		us2 := g.MustAt(p2)
		g.Set(*next2, us2)
		g.Set(p2, '.')
	}

	return next, nil
}
