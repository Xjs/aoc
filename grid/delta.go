package grid

import "errors"

// Delta is a delta that can be applied to a point.
type Delta struct {
	Dx, Dy int
}

// Diff calculates the difference between two points p1-p2,
// such that p2 + δ = p1
func Diff(p1, p2 Point) Delta {
	return D(int(p1.X)-int(p2.X), int(p1.Y)-int(p2.Y))

}

// D constructs a delta from two integers.
func D(dx, dy int) Delta { return Delta{dx, dy} }

// Delta applies the delta to the given point or returns an out-of-bounds error.
func (g *Grid[T]) Delta(p Point, d Delta) (Point, error) {
	x := int(p.X) + d.Dx
	y := int(p.Y) + d.Dy

	if x < 0 || y < 0 || Coordinate(x) >= g.Width() || Coordinate(y) >= g.Height() {
		return Point{}, errors.New("out of bounds")
	}

	return P(Coordinate(x), Coordinate(y)), nil
}

// DeltaWrap applies the delta to the given point, wrapping around if it would be out of bounds.
func (g *Grid[T]) DeltaWrap(p Point, d Delta) Point {
	x := (int(p.X) + d.Dx) % int(g.width)
	y := (int(p.Y) + d.Dy) % int(g.height)

	for x < 0 {
		x += int(g.width)
	}
	for y < 0 {
		y += int(g.height)
	}

	return P(Coordinate(x), Coordinate(y))
}

// GeneralDirections contains the four cardinal Directions.
var GeneralDirections = map[rune]Delta{
	'<': D(-1, 0),
	'^': D(0, -1),
	'>': D(1, 0),
	'v': D(0, 1),
}
