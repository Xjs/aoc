package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"
)

type wiring struct {
	wires map[string]func([]string) (bool, error)
	// cache map[string]bool

	register map[rune]int

	graph map[string][]string
}

func newWiring() *wiring {
	w := new(wiring)
	w.wires = make(map[string]func([]string) (bool, error))
	w.register = make(map[rune]int)
	// w.cache = make(map[string]bool)
	w.graph = map[string][]string{}
	return w
}

func (w *wiring) addConst(s string, v bool) {
	if s[0] == 'x' || s[0] == 'y' {
		// fast access to input registers
		r := []rune(s)[0]
		bit, err := strconv.Atoi(s[1:])
		if v {
			w.register[r] |= (1 << bit)
		} else {
			w.register[r] &= ^(1 << bit)
		}

		// log.Printf("Adding %c%d = %t (register %c = %b)", r, bit, v, r, w.register[r])
		if err == nil {
			w.wires[s] = func([]string) (bool, error) {
				return (w.register[r] & (1 << bit)) > 0, nil
			}
			return
		}
	}
	w.wires[s] = func([]string) (bool, error) {
		return v, nil
	}
}

func (w *wiring) swap(s []string) {
	for i := 0; i < len(s); i += 2 {
		swap(w, s[i], s[i+1])
	}
}

func (w *wiring) eval(s string) (bool, error) {
	// if v, ok := w.cache[s]; ok {
	// 	return v
	// }

	return w.call(s, nil)

	// w.cache[s] = v
}

var errLoop = errors.New("loop")

func (w *wiring) call(s string, path []string) (bool, error) {
	for _, p := range path {
		if p == s {
			return false, errLoop
		}
	}

	return w.wires[s](append(path, s))
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
	w.wires[s] = func(path []string) (bool, error) {
		l, err := w.call(left, path)
		if err != nil {
			return false, err
		}
		r, err := w.call(right, path)
		if err != nil {
			return false, err
		}
		return f(l, r), nil
	}
	w.graph[s] = []string{left, right}
	return nil
}

func main() {
	init := os.Args[1:]
	offset := 1
	flag.IntVar(&offset, "offset", offset, "offset")
	flag.Parse()

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
				vv = true
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

	log.Printf("Doing initial swap: %v", init)
	w.swap(init)

	x, _, _ := w.getNum('x')
	y, _, _ := w.getNum('y')
	z, bits, err := w.getNum('z')
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("part1: x: %d ( %b)", x, x)
	log.Printf("part1: y: %d ( %b)", y, y)
	log.Printf("part1: z: %d (%b, %t, should be %d)", z, z, z == x+y, x+y)

	log.Printf("bits that need fixing:               %b", (x+y)^z)

	fix := math.MaxInt
	for i := 1; i < bits; i++ {
		// log.Printf("%d: coverage: %f", i, cov)
		if err := coverage(w, i, 1<<(i-1), 20); err != nil {
			log.Printf("coverage: %v, fix = %d", err, i)
			if strings.Count(strings.Split(err.Error(), "wrong bits: ")[1], "1") > 1 {
				os.Exit(4)
			}
			fix = i
			break
		}
	}

	var wires []string
	for wire := range w.wires {
		switch wire[0] {
		case 'x', 'y', 'z':
			continue
		}

		wires = append(wires, wire)
	}

	log.Printf("%d wires.", len(wires))

	swaps := make(map[[2]string]struct{})
	initMap := make(map[string]bool)
	for _, w := range init {
		initMap[w] = true
	}

	for _, wire1 := range wires {
		for _, wire2 := range wires {
			w1 := wire1
			w2 := wire2
			if w1 == w2 {
				continue
			}
			// don't include those that have been in init, no swap can be done multiple times
			if initMap[w1] || initMap[w2] {
				continue
			}

			if w2 < w1 {
				w2, w1 = w1, w2
			}
			swaps[[2]string{w1, w2}] = struct{}{}
		}
	}

	log.Printf("%d swaps identified.", len(swaps))

	candidates2 := make(map[string]struct{})

	for swap := range swaps {
		w.swap(swap[:])
		if err := coverage(w, fix, 1<<(fix-1), 100); err == nil {
			log.Printf("possibly %v", swap)
			for _, c := range swap {
				candidates2[c] = struct{}{}
			}
		}
		w.swap(swap[:])
	}

	cm := make(map[int]map[string]bool)
	bitsFixing := make(map[int]bool)

	for i := 0; i < bits; i++ {
		needsFixing := ((x+y)^z)&(1<<i) > 0
		if needsFixing {
			bitsFixing[i] = needsFixing
		}

		cm[i] = make(map[string]bool)

		for _, cand := range traverse(w.graph, fmt.Sprintf("z%2d", i)) {
			cm[i][cand] = needsFixing
		}
	}

	// start with the smallest bit that needs fixing
	smallest := math.MaxInt
	for bit := range bitsFixing {
		if bit < smallest {
			smallest = bit
		}
	}

	log.Printf("smallest = %d, fix = %d", smallest, fix)

	log.Printf("candidates before purging: %d", len(cm[smallest]))

	candidates := make(map[string]struct{})
	for bit := range bitsFixing {
		if bit != fix {
			continue
		}

		for i := bit - offset; i >= 0; i-- {
			for irrelevant := range cm[i] {
				delete(cm[bit], irrelevant)
			}
		}

		for cand := range cm[bit] {
			candidates[cand] = struct{}{}
		}
	}

	// for i := 0; i <= smallest+1; i++ {
	// 	if err := check(w, i); err != nil {
	// 		log.Printf("%d: %v", i, err)
	// 	}
	// }

	// I assume the x-s and y-s are not wrong, since they are directly the inputs.
	for cand := range candidates {
		if cand[0] == 'x' || cand[0] == 'y' {
			delete(candidates, cand)
		}
	}

	log.Printf("candidates after purging: %v (%d)", candidates, len(candidates))

	// if len(candidates) == 0 {
	// 	os.Exit(1)
	// }

	count := 0

	var combs [][]string

	possibleSwaps := make(map[string]map[[2]string]struct{})

	for cand := range candidates2 {
		possibleSwaps[cand] = make(map[[2]string]struct{})
		log.Printf("Trying out all swaps with %s", cand)
		for swap := range swaps {
			if swap[0] != cand && swap[1] != cand {
				continue
			}
			possibleSwaps[cand][swap] = struct{}{}
		}
	}

	// single := 0
	// for _, swaps := range possibleSwaps {
	// 	for swap := range swaps {
	// 		w.swap(swap[:])
	// 		if err := check(w, fix); err == nil {
	// 			if coverage(w, fix, 1<<(fix-1), 10) == nil {
	// 				log.Printf("found: %v", swap)
	// 				single++
	// 				// if coverage(w, bits, 1<<(fix-1), 100) == nil {
	// 				// }
	// 			}
	// 		}
	// 		w.swap(swap[:])
	// 	}
	// }

	// if single != 0 {
	// 	os.Exit(0)
	// }

	possibleDoubleSwaps := make(map[[4]string]struct{})

	for cand1, swaps1 := range possibleSwaps {
		for cand2, swaps2 := range possibleSwaps {
			if cand1 == cand2 {
				continue
			}

			for sw1 := range swaps1 {
				for sw2 := range swaps2 {
					if sw1[0] == sw2[0] || sw1[1] == sw2[0] || sw1[0] == sw2[1] || sw1[1] == sw2[1] {
						continue
					}

					possibleDoubleSwaps[[4]string{sw1[0], sw1[1], sw2[0], sw2[1]}] = struct{}{}
				}
			}
		}
	}

	log.Printf("possibleDoubleSwaps = %d", len(possibleDoubleSwaps))

	maybes := 0
	for ds := range possibleDoubleSwaps {
		w.swap(ds[:])
		if err := check(w, fix); err == nil {
			maybes++
		} else {
			delete(possibleDoubleSwaps, ds)
		}
		w.swap(ds[:])
	}

	log.Printf("%d maybes", maybes)
	if maybes == 0 {
		os.Exit(2)
	}

	for ds := range possibleDoubleSwaps {
		w.swap(ds[:])
		if coverage(w, fix, 1<<(fix-1), 10) == nil {
			if coverage(w, bits, 1<<(fix-1), 100) == nil {
				log.Printf("possibly %v", ds)
			}
		}
		w.swap(ds[:])
	}

	for c := 2; c <= len(candidates); c += 2 {
		log.Printf("trying out swaps of length %d", c)
		if err := combinations(candidates, nil, c, func(s []string) error {
			count++
			defer func() {
				if r := recover(); r == "loop" {
					// log.Printf("after %v: %v", s, r)
				}
			}()

			if count%1e5 == 0 {
				log.Printf("%d (%v)", count, s)
			}

			// log.Printf("swapping %v", s)
			w.swap(s)
			if err := check(w, fix); err == nil {
				// log.Printf("maybe %v", s)
				if err := check(w, math.MaxInt); err == nil {
					combs = append(combs, s)
				}
				return nil
			} else {
				// log.Printf("after %v: %v", s, err)
			}
			// swap back
			w.swap(s)
			return nil
		}); err != nil {
			log.Fatal(err)
		}
	}

	for _, comb := range combs {
		log.Print(comb)
	}
}

func combinations(candidates map[string]struct{}, current []string, l int, f func([]string) error) error {
	if len(current) == l {
		if err := f(current); err != nil {
			return err
		}
		return nil
	}

outer:
	for cand := range candidates {
		for _, cur := range current {
			if cur == cand {
				continue outer
			}
		}
		if err := combinations(candidates, append(current, cand), l, f); err != nil {
			return err
		}
	}

	return nil
}

func coverage(w *wiring, maxBits int, maxInput int, iter int) error {
	combinations := iter * iter
	success := 0
	for i := 0; i < iter; i++ {
		w.register['x'] = rand.Intn(maxInput)
		for j := 0; j < iter; j++ {
			w.register['y'] = rand.Intn(maxInput)
			if err := check(w, maxBits); err == nil {
				success++
			} else {
				return err
			}
		}
	}
	if success != combinations {
		return fmt.Errorf("cov: %d/%d", success, combinations)
	}
	return nil
}

func swap(w *wiring, a, b string) {
	w.wires[b], w.wires[a] = w.wires[a], w.wires[b]
	// empty the cache
	// TODO: Only empty everything downstream of these wires?
	// w.cache = make(map[string]bool)
}

func check(w *wiring, max int) error {
	x, _, _ := w.getNumMax('x', max-1)
	y, _, _ := w.getNumMax('y', max-1)
	z, _, err := w.getNumMax('z', max)
	if err != nil {
		return err
	}

	if z == x+y {
		return nil
	}

	wrongBits := z ^ (x + y)

	return fmt.Errorf("z = %d (should be %d = %d + %d, wrong bits: %b)", z, x+y, x, y, wrongBits)
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

func (w *wiring) getNum(prefix rune) (int, int, error) {
	return w.getNumMax(prefix, math.MaxInt)
}

func (w *wiring) getNumMax(prefix rune, max int) (int, int, error) {
	var zs []string
	for cable := range w.wires {
		if []rune(cable)[0] == prefix {
			zs = append(zs, cable)
		}
	}

	slices.Sort(zs)

	result := 0
	if max < 0 {
		max = 0
	}
	if max > len(zs) {
		max = len(zs)
	}
	zs = zs[:max]
	for i, z := range zs {
		v, err := w.eval(z)
		if err != nil {
			return 0, 0, err
		}

		result += (toint(v) << i)
	}
	return result, max, nil
}

func toint(x bool) int {
	if x {
		return 1
	}
	return 0
}
