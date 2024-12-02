package parse

import (
	"strconv"
	"strings"
)

// IntList takes a string with comma-separated integers (may be surrounded by whitespace) and returns the integers as a slice.
func IntList(s string) ([]int, error) {
	return IntListSep(s, ",")
}

// IntListWhitespace takes a string with whitespace-separated integers (may be surrounded by whitespace) and returns the integers as a slice.
func IntListWhitespace(s string) ([]int, error) {
	fields := strings.Fields(s)
	ns := make([]int, 0, len(fields))

	for _, field := range fields {
		n, err := strconv.Atoi(field)
		if err != nil {
			return ns, err
		}
		ns = append(ns, n)
	}
	return ns, nil
}

// IntListSep takes a string with sep-separated integers (may be surrounded by whitespace) and returns the integers as a slice.
func IntListSep(s string, sep string) ([]int, error) {
	fields := strings.Split(s, sep)
	ns := make([]int, 0, len(fields))

	for _, field := range fields {
		n, err := strconv.Atoi(strings.TrimSpace(field))
		if err != nil {
			return ns, err
		}
		ns = append(ns, n)
	}
	return ns, nil
}

// DigitList takes a string with digits integers and returns the digits as a slice.
func DigitList(s string) ([]int, error) {
	ns := make([]int, 0, len(s))

	for _, r := range s {
		n, err := strconv.Atoi(string(r))
		if err != nil {
			return ns, err
		}
		ns = append(ns, n)
	}
	return ns, nil
}
