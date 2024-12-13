package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Xjs/aoc/grid"
)

type clawMachine struct {
	δA grid.Point
	δB grid.Point
	P  grid.Point
}

var buttonRE = regexp.MustCompile(`Button ([AB]): X\+(\d+), Y\+(\d+)`)
var prizeRE = regexp.MustCompile(`Prize: X=(\d+), Y=(\d+)`)

func parseClawMachine(input string) (clawMachine, error) {
	sp := strings.Split(input, "\n")
	if len(sp) != 3 {
		return clawMachine{}, errors.New("need exactly 3 lines")
	}

	δA, err := parseButton(sp[0], "A")
	if err != nil {
		return clawMachine{}, err
	}

	δB, err := parseButton(sp[1], "B")
	if err != nil {
		return clawMachine{}, err
	}

	pr := prizeRE.FindAllStringSubmatch(sp[2], -1)
	if len(pr) != 1 || len(pr[0]) != 3 {
		return clawMachine{}, errors.New("syntax: Prize: X=px, Y=py")
	}

	px, err := strconv.Atoi(pr[0][1])
	if err != nil {
		return clawMachine{}, err
	}
	py, err := strconv.Atoi(pr[0][2])
	if err != nil {
		return clawMachine{}, err
	}
	prize := grid.P(grid.Coordinate(px), grid.Coordinate(py))

	return clawMachine{
		δA: δA,
		δB: δB,
		P:  prize,
	}, nil
}

func parseButton(input string, letter string) (grid.Point, error) {
	but := buttonRE.FindAllStringSubmatch(input, -1)
	if len(but) != 1 || len(but[0]) != 4 || but[0][1] != letter {
		return grid.Point{}, fmt.Errorf("syntax: Button %s: X+ax, Y+ay", letter)
	}

	x, err := strconv.Atoi(but[0][2])
	if err != nil {
		return grid.Point{}, err
	}
	y, err := strconv.Atoi(but[0][3])
	if err != nil {
		return grid.Point{}, err
	}
	return grid.P(grid.Coordinate(x), grid.Coordinate(y)), nil
}

// solve solves the claw machine's puzzle, and returns if there is an integer solution, and if yes
// the number of times buttons A and B need to be pushed.
func (cm clawMachine) solve() (bool, int, int) {
	PX := int64(cm.P.X)
	PY := int64(cm.P.Y)

	ax := int64(cm.δA.X)
	ay := int64(cm.δA.Y)

	bx := int64(cm.δB.X)
	by := int64(cm.δB.Y)

	// B = (PY - (ay/ax)*PX) / (by - (ay/ax)*bx)
	// A = PX/ax - B*bx/ax

	B := new(big.Rat).Quo(
		new(big.Rat).Sub(new(big.Rat).SetInt64(PY), big.NewRat(ay*PX, ax)),
		new(big.Rat).Sub(new(big.Rat).SetInt64(by), big.NewRat(ay*bx, ax)),
	)

	A := new(big.Rat).Sub(
		big.NewRat(PX, ax),
		new(big.Rat).Mul(B, big.NewRat(bx, ax)),
	)

	if !A.IsInt() || !B.IsInt() {
		return false, 0, 0
	}
	iA := int(A.Num().Int64())
	iB := int(B.Num().Int64())

	if iA*int(cm.δA.X)+iB*int(cm.δB.X) != int(cm.P.X) || iA*int(cm.δA.Y)+iB*int(cm.δB.Y) != int(cm.P.Y) {
		log.Printf("didn't work out")
		return false, 0, 0
	}

	return true, iA, iB
}

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	var clawMachines []clawMachine
	for i, rawMachine := range bytes.Split(input, []byte("\n\n")) {
		cm, err := parseClawMachine(strings.TrimSpace(string(rawMachine)))
		if err != nil {
			log.Fatalf("claw machine %d: %v (%q)", i, err, string(rawMachine))
		}
		clawMachines = append(clawMachines, cm)
	}

	tokens := 0
	tokens2 := 0
	for i, cm := range clawMachines {
		ok, A, B := cm.solve()
		if !ok {
			log.Printf("claw machine %d is unsolvable", i)
		} else {
			log.Printf("claw machine %d: A: %d (%d tokens), B: %d (%d tokens), total %d tokens", i, A, 3*A, B, B, 3*A+B)
		}

		tokens += 3*A + B

		cm2 := cm
		const offset = 10000000000000

		cm2.P.X += offset
		cm2.P.Y += offset

		ok2, A2, B2 := cm2.solve()
		if !ok2 {
			log.Printf("part2: claw machine %d is unsolvable", i)
		} else {
			log.Printf("part2: claw machine %d: A: %d (%d tokens), B: %d (%d tokens), total %d tokens", i, A2, 3*A2, B2, B2, 3*A2+B2)
		}

		tokens2 += 3*A2 + B2
	}
	log.Printf("part1: %d", tokens)
	log.Printf("part2: %d", tokens2)
}
