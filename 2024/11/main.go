package main

import (
	"io"
	"log"
	"os"
	"strconv"

	"github.com/Xjs/aoc/parse"
)

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	start, err := parse.IntListWhitespace(string(input))
	if err != nil {
		log.Fatal(err)
	}

	stones := NewStones(start)

	for i := 0; i < 25; i++ {
		stones.Blink()
	}

	log.Printf("part1: %d", len(stones.IntSlice()))

	sum := 0
	for _, stone := range start {
		sum += blinkSingleStone(stone, 75)
	}
	log.Printf("part2: %d", sum)
}

type Stone struct {
	number int
	next   *Stone
}

func NewStones(stones []int) *Stone {
	if len(stones) == 0 {
		return nil
	}

	current := &Stone{number: stones[0]}
	start := current

	for i := 1; i < len(stones); i++ {
		current.next = &Stone{number: stones[i]}
		current = current.next
	}

	return start
}

func (s *Stone) IntSlice() []int {
	var sl []int
	for stone := s; stone != nil; stone = stone.next {
		sl = append(sl, stone.number)
	}
	return sl
}

func (s *Stone) Blink() {
	for stone := s; stone != nil; stone = stone.next {
		if stone.number == 0 {
			stone.number = 1
		} else if ss := strconv.Itoa(stone.number); len(ss)%2 == 0 {
			left, right := ss[:len(ss)/2], ss[len(ss)/2:]
			stone.number = atoi(left)
			oldnext := stone.next
			stone.next = &Stone{number: atoi(right), next: oldnext}
			stone = stone.next
		} else {
			stone.number *= 2024
		}
	}
}

func atoi(s string) int {
	x, _ := strconv.Atoi(s)
	return x
}

// recursionMemo stores the answer to the question "how many stones do I have
// if I use the production rule `times` times on the stone with the starting number `number`"
type recursionMemo struct {
	number int
	times  int
}

// Learned this from Axel today, thanks :)
var memo = make(map[recursionMemo]int)

// blinkSingleStone returns the number of stones after blinking on the stone with the given inscribed number the
// given number of times
func blinkSingleStone(number int, times int) (result int) {
	mem := recursionMemo{number, times}
	if m, ok := memo[mem]; ok {
		return m
	}
	defer func() { memo[mem] = result }()

	if times == 0 {
		return 1
	}
	if number == 0 {
		return blinkSingleStone(1, times-1)
	}
	if ss := strconv.Itoa(number); len(ss)%2 == 0 {
		left, right := ss[:len(ss)/2], ss[len(ss)/2:]

		return blinkSingleStone(atoi(left), times-1) + blinkSingleStone(atoi(right), times-1)
	}

	return blinkSingleStone(2024*number, times-1)
}
