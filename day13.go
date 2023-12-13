package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func parseDay13(text string) [][][]byte {
	parts := strings.Split(text, "\n\n")

	var result [][][]byte

	for _, part := range parts {
		trimmedPart := strings.TrimSpace(part)
		if trimmedPart != "" {
			var stoneMap [][]byte
			for _, row := range strings.Split(trimmedPart, "\n") {
				stoneMap = append(stoneMap, []byte(row))
			}
			result = append(result, stoneMap)
		}
	}

	return result
}

func areColumnsTheSame(stoneMap [][]byte, a int, b int) bool {
	for r, _ := range stoneMap {
		if stoneMap[r][a] != stoneMap[r][b] {
			return false
		}
	}
	return true
}

func isReflectedVertically(stoneMap [][]byte, startCol int) bool {

	// 0 1 2
	//   ^ ^
	var otherCol = startCol + 1
	for c := startCol; c >= 0 && otherCol < len(stoneMap[0]); c-- {
		if !areColumnsTheSame(stoneMap, c, otherCol) {
			return false
		}

		otherCol++
	}

	// 0 1 2 3
	//   ^ ^
	// ^     ^

	return true
}

func findVerticalReflection(stoneMap [][]byte) (bool, int) {
	for c := 0; c < len(stoneMap[0])-1; c++ {
		if isReflectedVertically(stoneMap, c) {
			return true, c + 1
		}
	}
	return false, 0
}

func areRowsTheSame(stoneMap [][]byte, a int, b int) bool {
	for c, _ := range stoneMap[0] {
		if stoneMap[a][c] != stoneMap[b][c] {
			return false
		}
	}
	return true
}

func isReflectedHorizontally(stoneMap [][]byte, startRow int) bool {

	var otherRow = startRow + 1
	for r := startRow; r >= 0 && otherRow < len(stoneMap); r-- {
		if !areRowsTheSame(stoneMap, r, otherRow) {
			return false
		}

		otherRow++
	}

	return true
}

func findHorizontalReflection(stoneMap [][]byte) (bool, int) {
	for r := 0; r < len(stoneMap)-1; r++ {
		if isReflectedHorizontally(stoneMap, r) {
			return true, r + 1
		}
	}
	return false, 0
}

func calcReflections(maps [][][]byte) int {
	var sum = 0
	for _, stoneMap := range maps {
		hasVerticalReflection, column := findVerticalReflection(stoneMap)
		if hasVerticalReflection {
			sum += column
		}
		hasHorizontalReflection, column := findHorizontalReflection(stoneMap)
		if hasHorizontalReflection {
			sum += 100 * column
		}
	}

	return sum
}

func calcSmudgeReflections(maps [][][]byte) int {
	var sum = 0
	// for _, stoneMap := range maps {
	// 	for r, row := range stoneMap {
	// 		for c, col := range row {
	// 			hasVerticalReflection, column := findVerticalReflection(stoneMap)
	// 			if hasVerticalReflection {
	// 				sum += column
	// 			}
	// 			hasHorizontalReflection, column := findHorizontalReflection(stoneMap)
	// 			if hasHorizontalReflection {
	// 				sum += 100 * column
	// 			}
	// 		}
	// 	}
	// }

	return sum
}

func Day13() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: aoc 13 <input>")
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
	maps := parseDay13(text)
	if DEBUG {
		fmt.Println(maps)
	}
	fmt.Println("Part 1:")
	fmt.Println(calcReflections(maps))
	fmt.Println("Part 2:")
	fmt.Println(calcSmudgeReflections(maps))
}
