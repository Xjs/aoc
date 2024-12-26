package main

import (
	"bufio"
	"log"
	"os"
	"slices"
	"strings"
)

type network struct {
	computers map[string]map[string]struct{}
}

func (n *network) insert(c1, c2 string) {
	if n.computers == nil {
		n.computers = make(map[string]map[string]struct{})
	}
	if n.computers[c1] == nil {
		n.computers[c1] = make(map[string]struct{})
	}
	if n.computers[c2] == nil {
		n.computers[c2] = make(map[string]struct{})
	}
	n.computers[c1][c2] = struct{}{}
	n.computers[c2][c1] = struct{}{}
}

func copySet(s map[string]struct{}) map[string]struct{} {
	res := make(map[string]struct{})
	for k := range s {
		res[k] = struct{}{}
	}
	return res
}

func intersect(m1, m2 map[string]struct{}) map[string]struct{} {
	res := make(map[string]struct{})
	for k1 := range m1 {
		if _, ok := m2[k1]; ok {
			res[k1] = struct{}{}
		}
	}
	return res
}

func setToSlice(s map[string]struct{}) []string {
	keys := make([]string, 0, len(s))
	for k := range s {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

func fac(x int) int {
	if x <= 1 {
		return 1
	}
	return x * fac(x-1)
}

func comb(a, b int) int {
	return fac(a) / (fac(b) * fac(a-b))
}

func main() {
	n := new(network)

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		t := s.Text()
		names := strings.Split(t, "-")
		if len(names) != 2 {
			log.Fatalf("%q should be of syntax name1-name2", t)
		}

		// everybody is connected to themselves, makes processing easier
		n.insert(names[0], names[0])
		n.insert(names[0], names[1])
	}

	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	sets := make(map[string]struct{})
	triplets := make(map[string]struct{})
	for c, net := range n.computers {
		for neighbour := range net {
			for n2 := range net {
				if n2 == neighbour || neighbour == c || n2 == c {
					continue
				}

				if _, ok := n.computers[neighbour][n2]; ok {
					triplet := strings.Join(setToSlice(map[string]struct{}{c: {}, neighbour: {}, n2: {}}), ",")
					triplets[triplet] = struct{}{}
				}
			}

			conns := intersect(net, n.computers[neighbour])
			for neighbour2 := range conns {
				conns = intersect(conns, n.computers[neighbour2])
			}

			sl := setToSlice(conns)
			if len(sl) < 3 {
				continue
			}

			set := strings.Join(sl, ",")
			sets[set] = struct{}{}
		}
	}

	count := 0
	for triplet := range triplets {
		for _, k := range strings.Split(triplet, ",") {
			if k[0] == 't' {
				count++
				break
			}
		}
	}

	var largest string
	for set := range sets {
		if len(set) > len(largest) {
			largest = set
		}
	}

	log.Printf("part1: %d", count)
	log.Printf("part2: %s", largest)
}
