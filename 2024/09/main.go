package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("part1: %d", checksum(compact(parseFS(b))))
	log.Printf("part2: %d", checksum(compact2(parseFS(b))))
}

// parseFS reads a "disk map" and returns the whole filesystem represented by it
func parseFS(diskMap []byte) []int {
	var result []int

	isFile := true
	id := 0
	for _, b := range diskMap {
		l := int(b - '0')
		content := -1
		if isFile {
			content = id
		}

		for i := 0; i < l; i++ {
			result = append(result, content)
		}

		if isFile {
			id++
		}
		isFile = !isFile
	}

	return result
}

func compact(fs []int) []int {
	gap := 0
	for i := len(fs) - 1; i >= 0; i-- {
		b := fs[i]
		if b == -1 {
			continue
		}

		gap = nextGap(fs, gap)
		if gap == -1 {
			break
		}

		if gap < i {
			fs[gap] = b
			fs[i] = -1
		}
	}
	return fs
}

func compact2(fs []int) []int {
	type file struct {
		position int
		length   int
	}
	files := make(map[int]file)
	gaps := make(map[int]file)
	var fileIDs, gapIDs []int

	gapID := 0
	for i := 0; i < len(fs); i++ {
		if fs[i] == -1 {
			g, ok := gaps[gapID]
			if ok {
				g.length++
			} else {
				g.position = i
				g.length = 1
			}
			gaps[gapID] = g
		} else {
			f, ok := files[fs[i]]
			if ok {
				f.length++
			} else {
				f.position = i
				f.length = 1
				fileIDs = append(fileIDs, fs[i])
				gapIDs = append(gapIDs, gapID)
				gapID++
			}
			files[fs[i]] = f
		}
	}

	for i := len(fileIDs) - 1; i >= 0; i-- {
		f := files[fileIDs[i]]
		for _, gapID := range gapIDs {
			gap := gaps[gapID]
			if gap.position > f.position {
				continue
			}
			if gap.length >= f.length {
				gp := gap.position
				for k := 0; k < f.length; k++ {
					fs[gp+k] = fs[f.position+k]
					fs[f.position+k] = -1
					gap.length--
					gap.position++
				}
				gaps[gapID] = gap
				break
			}
		}
	}

	return fs
}

func printFS(fs []int) {
	for _, b := range fs {
		if b == -1 {
			fmt.Print(".")
		} else {
			fmt.Print(b)
		}
	}
	fmt.Println()
}

func nextGap(fs []int, gap int) int {
	for ; fs[gap] != -1 && gap < len(fs); gap++ {
	}
	if fs[gap] != -1 {
		return -1
	}
	return gap
}

func checksum(fs []int) int {
	sum := 0
	for i, b := range fs {
		if b == -1 {
			continue
		}
		sum += i * b
	}
	return sum
}
