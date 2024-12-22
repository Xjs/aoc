package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"sync"
)

type rng struct {
	seed int
}

func (r *rng) next() int {
	r.mix(r.seed * 64)
	r.prune()

	r.mix(r.seed / 32)
	r.prune()

	r.mix(r.seed * 2048)
	r.prune()

	return r.seed
}

func (r *rng) mix(x int) int {
	r.seed ^= x
	return r.seed
}

func (r *rng) prune() int {
	r.seed %= (1 << 24)
	return r.seed
}

const count = 2000

func (r *rng) prices() ([]int, []int) {
	seq := make([]int, count+1)
	diffs := make([]int, count)
	seq[0] = r.seed
	for i := 0; i < count; i++ {
		r.next()
		seq[i+1] = (r.seed % 10)
		diffs[i] = seq[i+1] - seq[i]
	}

	return seq, diffs
}

func search(prices []int, diffs []int, seq [4]int) int {
	for i := 0; i < count-4; i++ {
		if diffs[i] == seq[0] &&
			diffs[i+1] == seq[1] &&
			diffs[i+2] == seq[2] &&
			diffs[i+3] == seq[3] {
			return prices[i+4]
		}
	}

	return 0
}

func main() {
	sum := 0
	s := bufio.NewScanner(os.Stdin)
	r := new(rng)

	var allPrices [][]int
	var allDiffs [][]int

	for s.Scan() {
		seed, err := strconv.Atoi(s.Text())
		if err != nil {
			log.Print(err)
			continue
		}

		r.seed = seed
		prices, diffs := r.prices()

		allPrices = append(allPrices, prices)
		allDiffs = append(allDiffs, diffs)

		sum += r.seed
	}

	log.Printf("part1: %d", sum)

	highest := make(map[[4]int][]int)
	hm := new(sync.RWMutex)

	var wg sync.WaitGroup
	for i := -9; i < 10; i++ {
		for j := -9; j < 10; j++ {
			for k := -9; k < 10; k++ {
				for l := -9; l < 10; l++ {
					wg.Add(1)
					seq := [4]int{i, j, k, l}

					hm.Lock()
					highest[seq] = make([]int, len(allPrices))
					hm.Unlock()

					go func(seq [4]int) {
						for m := 0; m < len(allPrices); m++ {
							p := search(
								allPrices[m],
								allDiffs[m],
								seq,
							)

							hm.Lock()
							highest[seq][m] = p
							hm.Unlock()

						}
						wg.Done()
					}(seq)
				}
			}
		}
	}

	wg.Wait()

	highestSum := 0
	var highestSeq [4]int
	sums := make(map[[4]int]int)
	for seq, highestPrice := range highest {
		sum := 0
		for _, i := range highestPrice {
			sum += i
		}
		sums[seq] = sum
		if sum > highestSum {
			highestSum = sum
			highestSeq = seq
		}
	}

	log.Printf("part2: %d (with %v)", highestSum, highestSeq)
}
