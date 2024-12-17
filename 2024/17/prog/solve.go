package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	input, _ := strconv.Atoi(os.Args[1])
	offset, _ := strconv.Atoi(os.Args[2])
	limit, _ := strconv.Atoi(os.Args[3])
	for start := input + offset; start < input+offset+limit; start++ {
		A := start
		B := 0
		C := 0

		var output []int

		for A != 0 {
			// Take the lower 3 bits of A and flip the lowest one
			B = (A % 8) ^ 1
			// Take all bits except the lower B ones
			C = A / (1 << B)
			// XOR B with 5 (i. e. flip high and low)
			B ^= 5
			// XOR the bits from B
			// (which are now effectively the lower 3 bits from A with the high one flipped)
			// with the bits from C,
			// which are the bits from A where the last ones were flipped (see above)
			B ^= C
			// Output the lower three bits from B
			output = append(output, B%8)
			// Drop the lowest three bits
			A /= (1 << 3)
		}

		fmt.Printf("%b: %v\n", start, output)
	}
}
