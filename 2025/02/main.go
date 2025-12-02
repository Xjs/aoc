package main

import (
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type idRange struct {
	start, end int
}

func (r idRange) Foreach(f func(int, string)) {
	for i := r.start; i <= r.end; i++ {
		f(i, strconv.Itoa(i))
	}
}

func parseRange(s string) idRange {
	sp := strings.Split(strings.TrimSpace(s), "-")
	if len(sp) != 2 {
		panic(s)
	}
	start, err := strconv.Atoi(sp[0])
	if err != nil {
		panic(err)
	}
	end, err := strconv.Atoi(sp[1])
	if err != nil {
		panic(err)
	}
	return idRange{start, end}
}

func main() {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	input := strings.Split(strings.TrimSpace(string(b)), ",")
	var ranges []idRange
	for _, rg := range input {
		ranges = append(ranges, parseRange(rg))
	}

	sum := 0
	sum2 := 0

	invalidIDs := make(map[int]struct{})

	for _, rg := range ranges {
		log.Printf("processing range %v", rg)
		rg.Foreach(func(i int, s string) {
			_, alreadyFound := invalidIDs[i]
			if alreadyFound {
				return
			}

			l := len(s)

			for _, partLen := range divisors(l) {
				ref := s[:partLen]
				found := true
				parts := l / partLen
				if parts == 1 {
					continue
				}
				for i := 1; i < parts; i++ {
					if s[i*partLen:(i+1)*partLen] != ref {
						found = false
						break
					}
				}
				if found {
					invalidIDs[i] = struct{}{}
					if parts == 2 {
						sum += i
					}
					if !alreadyFound {
						log.Printf("found %v for %d parts", i, parts)
						sum2 += i
						alreadyFound = true
					}
				}
			}
		})
	}

	log.Println(sum)
	log.Println(sum2)
}

var divisorsMap = make(map[int][]int)

func divisors(x int) []int {
	if res, ok := divisorsMap[x]; ok {
		return res
	}

	if x == 1 {
		return []int{1}
	}

	resSet := make(map[int]struct{})
	for i := x / 2; i >= 1; i-- {
		if x%i == 0 {
			resSet[i] = struct{}{}
			if i > 1 {
				for _, div := range divisors(x / i) {
					resSet[div] = struct{}{}
				}
			}
		}
	}

	var res []int
	for div := range resSet {
		res = append(res, div)
	}

	sort.Ints(res)
	divisorsMap[x] = res
	return res
}
