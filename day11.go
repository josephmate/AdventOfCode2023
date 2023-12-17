package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func parseDay11(input string) []string {
	lines := strings.Split(input, "\n")

	var result []string

	for _, line := range lines {
		trimmedStr := strings.TrimSpace(line)
		if len(trimmedStr) > 0 {
			result = append(result, trimmedStr)
		}
	}

	return result
}

func isEmptyRow(spaceImage []string, r int) bool {
	for _, col := range spaceImage[r] {
		if col != '.' {
			return false
		}
	}
	return true
}
func isEmptyCol(spaceImage []string, c int) bool {
	for _, row := range spaceImage {
		if row[c] != '.' {
			return false
		}
	}
	return true
}

func expandSpace(spaceImage []string) []string {
	emptyRows := map[int]bool{}
	emptyCols := map[int]bool{}
	for r, _ := range spaceImage {
		if isEmptyRow(spaceImage, r) {
			emptyRows[r] = true
		}
	}
	for c, _ := range spaceImage[0] {
		if isEmptyCol(spaceImage, c) {
			emptyCols[c] = true
		}
	}

	var expandedSpace []string
	for r, row := range spaceImage {
		var builder strings.Builder
		for c, col := range row {
			builder.WriteRune(col)
			if emptyCols[c] {
				builder.WriteRune('.')
			}
		}
		expandedRow := builder.String()
		if emptyRows[r] {
			expandedSpace = append(expandedSpace, strings.Repeat(".", len(expandedRow)))
		}
		expandedSpace = append(expandedSpace, expandedRow)
	}

	if DEBUG {
		fmt.Println("expandSpace")
		fmt.Println(strings.Join(expandedSpace, "\n"))
	}
	return expandedSpace
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func allPairShortestPathSum(spaceImage []string) int {
	expandedSpace := expandSpace(spaceImage)

	var galaxies [][2]int
	for r, row := range expandedSpace {
		for c, col := range row {
			if col == '#' {
				galaxies = append(galaxies, [2]int{r, c})
			}
		}
	}

	var sum = 0
	for i := 0; i < len(galaxies); i++ {
		for j := i + 1; j < len(galaxies); j++ {
			iPos := galaxies[i]
			jPos := galaxies[j]
			sum += ManhattanDistance(iPos, jPos)
		}
	}

	return sum
}

func Day11() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: aoc 11 <part 1 input>")
		os.Exit(1)
	}

	filenamePart1 := os.Args[2]
	// Open the file
	file, err := os.Open(filenamePart1)
	if err != nil {
		log.Fatalf("Error opening file: %s", err)
	}
	defer file.Close()
	reader := io.Reader(file)
	textBytes, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	text := string(textBytes)
	if DEBUG {
		fmt.Println(text)
	}
	spaceImage := parseDay11(text)
	if DEBUG {
		fmt.Println(spaceImage)
	}
	fmt.Println("Part 1:")
	fmt.Println(allPairShortestPathSum(spaceImage))
	fmt.Println("Part 2:")
}
