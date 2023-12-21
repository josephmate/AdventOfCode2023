package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func parseGardenMap(input string) [][]byte {
	lines := strings.Split(input, "\n")

	var byteSlices [][]byte
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if len(trimmedLine) > 0 {
			byteSlices = append(byteSlices, []byte(trimmedLine))
		}
	}

	return byteSlices
}

func generateGardenMoves(posn [2]int, gardenMap [][]byte) [][2]int {
	row := posn[0]
	col := posn[1]

	var result [][2]int
	if row-1 >= 0 && gardenMap[row-1][col] == '.' {
		result = append(result, [2]int{row - 1, col})
	}
	if row+1 < len(gardenMap) && gardenMap[row+1][col] == '.' {
		result = append(result, [2]int{row + 1, col})
	}
	if col-1 >= 0 && gardenMap[row][col-1] == '.' {
		result = append(result, [2]int{row, col - 1})
	}
	if col+1 >= 0 && gardenMap[row][col+1] == '.' {
		result = append(result, [2]int{row, col + 1})
	}

	return result
}

func runOneStep(prevVisited map[[2]int]bool, gardenMap [][]byte) map[[2]int]bool {
	nextVisited := map[[2]int]bool{}

	for key := range prevVisited {
		for _, move := range generateGardenMoves(key, gardenMap) {
			nextVisited[move] = true
		}
	}

	return nextVisited
}

func countFlowers(steps int, gardenMap [][]byte) int {
	// find start location
	startRow, startCol, found := FindChar('S', gardenMap)
	if !found {
		fmt.Println("Did not find S in the garden map!")
		return -1
	}
	if DEBUG {
		fmt.Println("countFlowers startRow", startRow, "startCol", startCol)
	}
	gardenMap[startRow][startCol] = '.'

	var currentVisited = map[[2]int]bool{}
	currentVisited[[2]int{startRow, startCol}] = true
	for step := 0; step < steps; step++ {
		currentVisited = runOneStep(currentVisited, gardenMap)
	}

	return len(currentVisited)
}

func Day21() {

	if len(os.Args) < 4 {
		fmt.Println("Usage: aoc 21 <steps> <input>")
		os.Exit(1)
	}

	steps, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Error occurred parsing steps:", os.Args[2], err)
		return
	}
	filenamePart1 := os.Args[3]
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
		fmt.Println(steps)
		fmt.Println(text)
	}
	gardenMap := parseGardenMap(text)
	if DEBUG {
		fmt.Println(steps)
		fmt.Println(gardenMap)
	}
	fmt.Println("Part 1:")
	fmt.Println(countFlowers(steps, gardenMap))
	fmt.Println("Part 2:")
}
