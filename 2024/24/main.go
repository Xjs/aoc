package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

type wiring struct {
	wires map[string]func() bool
	cache map[string]bool

	graph map[string][]string
}

func newWiring() *wiring {
	w := new(wiring)
	w.wires = make(map[string]func() bool)
	w.cache = make(map[string]bool)
	w.graph = map[string][]string{}
	return w
}

func (w *wiring) addConst(s string, v bool) {
	w.wires[s] = func() bool {
		return v
	}
}

func (w *wiring) eval(s string) bool {
	if v, ok := w.cache[s]; ok {
		return v
	}

	v := w.wires[s]()

	w.cache[s] = v

	return v
}

func and(a, b bool) bool { return a && b }
func or(a, b bool) bool  { return a || b }
func xor(a, b bool) bool { return a != b }

func (w *wiring) addLogic(s string, left, right string, op string) error {
	var f func(a, b bool) bool
	switch op {
	case "AND":
		f = and
	case "OR":
		f = or
	case "XOR":
		f = xor
	default:
		return fmt.Errorf("%q not known", op)
	}
	w.wires[s] = func() bool {
		return f(w.wires[left](), w.wires[right]())
	}
	w.graph[s] = []string{left, right}
	return nil
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	w := newWiring()
	for s.Scan() {
		t := s.Text()
		if sp := strings.Split(t, ": "); len(sp) == 2 {
			v, err := strconv.Atoi(sp[1])
			if err != nil {
				log.Fatal(err)
			}
			var vv bool
			switch v {
			case 0:
				vv = false
			case 1:
				vv = true
			default:
				log.Fatal(fmt.Errorf("%d not valid", v))
			}

			w.addConst(sp[0], vv)
		}

		if sp := strings.Split(t, " -> "); len(sp) == 2 {
			abc := strings.Fields(sp[0])
			if len(abc) != 3 {
				log.Fatalf("%q invalid", t)
			}
			if err := w.addLogic(sp[1], abc[0], abc[2], abc[1]); err != nil {
				log.Fatal(err)
			}
		}
	}

	x, _ := w.getNum('x')
	y, _ := w.getNum('y')
	z, bits := w.getNum('z')

	log.Printf("part1: x: %d", x)
	log.Printf("part1: y: %d", y)
	log.Printf("part1: z: %d (%t, should be %d)", z, z == x+y, x+y)

	log.Printf("bits that need fixing: %b", (x+y)^z)

	for i := 0; i < bits; i++ {
		if ((x+y)^z)&(1<<i) > 0 {
			// bit i needs fixing
			log.Printf("# bit %d, candidates", i)
			cm := make(map[string]struct{})
			for _, cand := range traverse(w.graph, fmt.Sprintf("z%2d", i)) {
				cm[cand] = struct{}{}
			}
			var candidates []string
			for c := range cm {
				candidates = append(candidates, c)
			}
			slices.Sort(candidates)
			log.Print(candidates)
		}
	}
}

func traverse(graph map[string][]string, start string) []string {
	st := graph[start]
	if len(st) == 0 {
		return nil
	}
	if start == "" {
		panic(start)
	}

	left := traverse(graph, st[0])
	right := traverse(graph, st[1])
	var result []string
	result = append(result, st[:]...)
	result = append(result, left...)
	result = append(result, right...)
	return result
}

func (w *wiring) getNum(prefix rune) (int, int) {
	var zs []string
	for cable := range w.wires {
		if []rune(cable)[0] == prefix {
			zs = append(zs, cable)
		}
	}

	slices.Sort(zs)

	result := 0
	for i, z := range zs {
		result += (toint(w.eval(z)) << i)
	}
	return result, len(zs)
}

func toint(x bool) int {
	if x {
		return 1
	}
	return 0
}
