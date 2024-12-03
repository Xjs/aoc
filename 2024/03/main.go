package main

import (
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("part1: %d", parse(input, false))
	log.Printf("part2: %d", parse(input, true))
}

func parse(input []byte, parseDonts bool) int {
	var re *regexp.Regexp
	mulRE := `mul\([0-9]{1,3},[0-9]{1,3}\)`
	if parseDonts {
		re = regexp.MustCompile(mulRE + `|do\(\)|don't\(\)`)
	} else {
		re = regexp.MustCompile(mulRE)
	}

	results := re.FindAll(input, -1)
	sum := 0

	enabled := true

	for _, result := range results {
		result := string(result)
		if parseDonts {
			switch result {
			case "do()":
				enabled = true
				continue
			case "don't()":
				enabled = false
				continue
			}
		}

		if !enabled {
			continue
		}

		product := 1
		for _, num := range strings.Split(strings.TrimSuffix(strings.TrimPrefix(result, "mul("), ")"), ",") {
			n, err := strconv.Atoi(num)
			if err != nil {
				log.Fatal(err)
			}
			product *= n
		}
		sum += product
	}

	return sum
}
