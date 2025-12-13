package main

import (
	"container/heap"
	"os"
	"testing"
)

func equalsPair(p *pair, ref1, ref2 box) bool {
	return (p.b1 == ref1 && p.b2 == ref2) || (p.b2 == ref1 && p.b1 == ref2)
}

func TestParseAndOrder(t *testing.T) {
	f, err := os.Open("testdata/example.txt")
	if err != nil {
		t.Fatalf("Open(testdata/example.txt) error = %v", err)
	}
	t.Cleanup(func() { f.Close() })
	coll, err := readCollection(f)
	p := heap.Pop(coll.distHeap).(*pair)

	ref1 := box{162, 817, 812}
	ref2 := box{425, 690, 689}
	if !equalsPair(p, ref1, ref2) {
		t.Errorf("wrong first pair")
	}
}

func TestPart1(t *testing.T) {
	f, err := os.Open("testdata/example.txt")
	if err != nil {
		t.Fatalf("Open(testdata/example.txt) error = %v", err)
	}
	t.Cleanup(func() { f.Close() })
	coll, err := readCollection(f)

	connectStraightLines(coll, 10)
	sizes := coll.sortedCircuitSizes()
	product := 1
	for i := 0; i < 3; i++ {
		product *= sizes[i]
	}
	if product != 40 {
		t.Errorf("part1 example: got %d (%v), want 40", product, sizes)
	}

	circuits := 0
	for _, size := range sizes {
		if size > 0 {
			circuits++
		}
	}
	if circuits != 11 {
		t.Errorf("part1 example: got %d circuits, want 11", circuits)
	}
}
