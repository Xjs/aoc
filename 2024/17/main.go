package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Xjs/aoc/parse"
)

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	m, err := parseMachine(string(input))
	if err != nil {
		log.Fatal(err)
	}
	for m.exec() {
	}

	log.Printf("part1: %v", concatInts(m.output))
}

func concatInts(is []int) string {
	var b strings.Builder
	for i, x := range is {
		if i != 0 {
			b.WriteRune(',')
		}
		b.WriteString(strconv.Itoa(x))
	}
	return b.String()
}

func parseMachine(input string) (machine, error) {
	m := machine{}

	for i, line := range strings.Split(input, "\n") {
		if reg, ok := strings.CutPrefix(line, "Register "); ok {
			if reg[1] != ':' {
				return machine{}, fmt.Errorf("line %d: invalid register syntax: %q", i, line)
			}

			n, err := strconv.Atoi(strings.TrimSpace(reg[2:]))
			if err != nil {
				return machine{}, fmt.Errorf("line %d: %w", i, err)
			}

			switch reg[0] {
			case 'A':
				m.A = n
			case 'B':
				m.B = n
			case 'C':
				m.C = n
			}
		} else if prog, ok := strings.CutPrefix(line, "Program: "); ok {
			p, err := parse.IntList(prog)
			if err != nil {
				return machine{}, fmt.Errorf("line %d: %w", i, err)
			}
			m.program = p
		} else if strings.TrimSpace(line) != "" {
			return machine{}, fmt.Errorf("line %d: extra input: %q", i, line)
		}
	}
	return m, nil
}

type machine struct {
	A, B, C int
	// instruction pointer
	i int

	program []int
	output  []int
}

func (m *machine) instruction(i int) instruction {
	var inst instruction
	switch i {
	case 0:
		inst = m.adv
	case 1:
		inst = m.bxl
	case 2:
		inst = m.bst
	case 3:
		inst = m.jnz
	case 4:
		inst = m.bxc
	case 5:
		inst = m.out
	case 6:
		inst = m.bdv
	case 7:
		inst = m.cdv
	}
	return inst
}

// exec returns false if the machine halts
func (m *machine) exec() bool {
	if m.i >= len(m.program)-1 {
		return false
	}

	p, op := m.program[m.i], m.program[m.i+1]
	if jump := m.instruction(p)(op); !jump {
		m.i += 2
	}

	log.Print(p, op)
	log.Print(m)

	return true
}

func (m *machine) combo(op int) int {
	switch op {
	case 0, 1, 2, 3:
		return op
	case 4:
		return m.A
	case 5:
		return m.B
	case 6:
		return m.C
	case 7:
		panic(op)
	default:
		panic(op)
	}
}

// instructions take an argument, and return whether they jumped.
// Unless they jumped, execution should continue by increasing the program counter by 2.
type instruction func(int) bool

func pow(a, b int) int {
	product := 1
	for ; b > 0; b-- {
		product *= a
	}
	return product
}

func (m *machine) dv(arg int) int {
	return m.A / pow(2, m.combo(arg))
}

func (m *machine) adv(arg int) bool {
	m.A = m.dv(arg)
	return false
}

func (m *machine) bxl(arg int) bool {
	m.B ^= arg
	return false
}

func (m *machine) bst(arg int) bool {
	m.B = m.combo(arg) % 8
	return false
}

func (m *machine) jnz(arg int) bool {
	if m.A == 0 {
		return false
	}
	m.i = arg
	return true
}

func (m *machine) bxc(arg int) bool {
	m.B ^= m.C
	return false
}

func (m *machine) out(arg int) bool {
	m.output = append(m.output, m.combo(arg)%8)
	return false
}

func (m *machine) bdv(arg int) bool {
	m.B = m.dv(arg)
	return false
}

func (m *machine) cdv(arg int) bool {
	m.C = m.dv(arg)
	return false
}
