package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type rotation int

func (r rotation) String() string {
	if r < 0 {
		return fmt.Sprintf("L%d", -r)
	}
	return fmt.Sprintf("R%d", r)
}

func parseRotation(s string) rotation {

	num, err := strconv.ParseInt(s[1:], 10, 64)
	if err != nil {
		panic(err)
	}

	switch s[0] {
	case 'L':
		return rotation(-num)
	case 'R':
		return rotation(num)
	default:
		panic(s[0])
	}
}

func main() {
	current := rotation(50)
	zeroes := 0
	zeroes2 := 0

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		t := parseRotation(s.Text())

		dir := rotation(1)
		lim := int(t)
		if t < 0 {
			dir = rotation(-1)
			lim = -lim
		}
		for i := 0; i < lim; i++ {
			current += dir

			if current == -1 {
				current = 99
			}
			if current == 100 {
				current = 0
			}

			if current == 0 {
				zeroes2++
			}
		}

		if current == 0 {
			zeroes++
		}
	}

	log.Printf("zeroes: %d", zeroes)
	log.Printf("zeroes2: %d", zeroes2)
}
