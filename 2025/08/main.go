package main

import (
	"bufio"
	"container/heap"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

// A box is a junction box in 3D space.
type box struct {
	x, y, z int
}

func (b box) distance(b2 box) float64 {
	return math.Sqrt(math.Pow(float64(b2.x-b.x), 2) +
		math.Pow(float64(b2.y-b.y), 2) +
		math.Pow(float64(b2.z-b.z), 2))
}

// A collection is a collection of boxes
type collection struct {
	boxes []box
	// connections maps connections from one box to another.
	connections map[int]map[int]struct{}
	// boxCircuit maps boxes to a circuit ID
	boxCircuit map[int]int
	// Map of existant ciruits and their sizes
	circuitSizes map[int]int
	// Distances between boxes
	distances map[int]map[int]float64

	distHeap *distanceHeap
}

type distanceHeap []*pair

func (h distanceHeap) Len() int {
	return len(h)
}
func (h distanceHeap) Less(i, j int) bool {
	return h[i].distance < h[j].distance
}
func (h distanceHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h *distanceHeap) Push(x interface{}) {
	*h = append(*h, x.(*pair))
}
func (h *distanceHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// A pair is a pair of boxes and their distance.
type pair struct {
	b1, b2   box
	id1, id2 int
	distance float64
}

func (c *collection) setDistance(id1, id2 int, distance float64) {
	if id2 < id1 {
		id1, id2 = id2, id1
	}
	if c.distances[id1] == nil {
		c.distances[id1] = make(map[int]float64)
	}
	c.distances[id1][id2] = distance
}

func (c *collection) distance(id1, id2 int) float64 {
	if id2 < id1 {
		id1, id2 = id2, id1
	}
	return c.distances[id1][id2]
}

func (c *collection) add(boxNew box) {
	if c.connections == nil {
		c.connections = make(map[int]map[int]struct{})
	}
	if c.boxCircuit == nil {
		c.boxCircuit = make(map[int]int)
	}
	if c.circuitSizes == nil {
		c.circuitSizes = make(map[int]int)
	}
	if c.distances == nil {
		c.distances = make(map[int]map[int]float64)
	}
	if c.distHeap == nil {
		c.distHeap = new(distanceHeap)
		heap.Init(c.distHeap)
	}

	// idNew is our new ID
	idNew := len(c.boxes)
	c.boxes = append(c.boxes, boxNew)
	c.connections[idNew] = make(map[int]struct{})
	// A new boxNew belongs to its own circuit
	c.boxCircuit[idNew] = idNew
	c.circuitSizes[idNew] = 1

	// Cache distances to all existing boxes
	for idOld, boxOld := range c.boxes {
		if idOld == idNew {
			continue
		}

		dist := boxNew.distance(boxOld)
		c.setDistance(idNew, idOld, dist)

		lower, higher := idNew, idOld
		lowerBox, higherBox := boxNew, boxOld
		if lower > higher {
			lower, higher = higher, lower
			lowerBox, higherBox = higherBox, lowerBox
		}

		p := pair{
			b1: lowerBox, b2: higherBox,
			id1: lower, id2: higher,
			distance: dist,
		}

		heap.Push(c.distHeap, &p)
	}
}

// connect connects the two boxes. If they are part of the same circuit, the function returns false (no connection).
func (c *collection) connect(id1, id2 int) bool {
	cid1 := c.boxCircuit[id1]
	cid2 := c.boxCircuit[id2]

	if cid1 == cid2 {
		return false
	}

	newID, oldID := cid1, cid2
	if cid2 < cid1 {
		newID, oldID = cid2, cid1
	}
	c.boxCircuit[id1] = newID
	c.boxCircuit[id2] = newID
	c.circuitSizes[newID] += c.circuitSizes[oldID]
	c.circuitSizes[oldID] = 0

	return true
}

func parseBox(coords string) (box, error) {
	parts := strings.Split(coords, ",")
	if len(parts) != 3 {
		return box{}, errors.New("need exactly three coordinates")
	}
	var result box
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return box{}, fmt.Errorf("error parsing %q (coordinate %d) as int: %w", p, i, err)
		}
		switch i {
		case 0:
			result.x = n
		case 1:
			result.y = n
		case 2:
			result.z = n
		}
	}
	return result, nil
}

func readCollection(r io.Reader) (*collection, error) {
	c := new(collection)

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}

		box, err := parseBox(line)
		if err != nil {
			return nil, fmt.Errorf("error parsing %q as box: %w", line, err)
		}
		c.add(box)
	}
	return c, nil
}

func connectStraightLines(c *collection, nConnections int) {

	for connections := 0; connections < nConnections; connections++ {
		p := heap.Pop(c.distHeap).(*pair)

		c.connect(p.id1, p.id2)
	}
}

// sortedCircuitSizes returns the sizes of the circuits, sorted by size (descending)
func (c *collection) sortedCircuitSizes() []int {
	sizes := make([]int, 0, len(c.circuitSizes))
	for _, size := range c.circuitSizes {
		sizes = append(sizes, size)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))
	return sizes
}

func main() {
	coll, err := readCollection(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("parsed collection with %d boxes.", len(coll.boxes))

	connectStraightLines(coll, 1000)
	sizes := coll.sortedCircuitSizes()
	product := 1
	for i := 0; i < 3; i++ {
		product *= sizes[i]
	}
	log.Printf("part1: %d", product)
	log.Println(coll.sortedCircuitSizes())
}
