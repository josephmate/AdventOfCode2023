package main

import (
	"fmt"
	"os"
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

func generateGardenMovesInfinite(posn [2]int, gardenMap [][]byte) [][2]int {
	row := posn[0]
	col := posn[1]

	var result [][2]int

	var upRow = row - 1
	if upRow < 0 {
		upRow = len(gardenMap) - 1
	}
	if gardenMap[upRow][col] == '.' {
		result = append(result, [2]int{upRow, col})
	}

	var downRow = row + 1
	if downRow >= len(gardenMap) {
		downRow = 0
	}
	if gardenMap[downRow][col] == '.' {
		result = append(result, [2]int{downRow, col})
	}

	var leftCol = col - 1
	if leftCol < 0 {
		leftCol = len(gardenMap[row]) - 1
	}
	if gardenMap[row][leftCol] == '.' {
		result = append(result, [2]int{row, leftCol})
	}

	var rightCol = col + 1
	if rightCol >= len(gardenMap[row]) {
		rightCol = 0
	}
	if gardenMap[row][rightCol] == '.' {
		result = append(result, [2]int{row, rightCol})
	}

	return result
}

func runOneStepInfinite(prevVisited map[[2]int]bool, gardenMap [][]byte) map[[2]int]bool {
	nextVisited := map[[2]int]bool{}

	for key := range prevVisited {
		for _, move := range generateGardenMovesInfinite(key, gardenMap) {
			nextVisited[move] = true
		}
	}

	return nextVisited
}

func countFlowersInfinite(steps int, gardenMap [][]byte) int {
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
		currentVisited = runOneStepInfinite(currentVisited, gardenMap)
	}

	return len(currentVisited)
}

func Day21() {

	if len(os.Args) < 4 {
		fmt.Println("Usage: aoc 21 <steps> <input part 1> <steps> <input part 2>")
		os.Exit(1)
	}

	part1Steps := ParseIntOrExit(os.Args[2])
	part1Text := ReadFileOrExit(os.Args[3])
	part2Steps := ParseIntOrExit(os.Args[4])
	part2Text := ReadFileOrExit(os.Args[5])

	fmt.Println("Part 1:")
	if DEBUG {
		fmt.Println(part1Steps)
		fmt.Println(part1Text)
	}
	part1GardenMap := parseGardenMap(part1Text)
	if DEBUG {
		fmt.Println(part1Steps)
		fmt.Println(part1GardenMap)
	}
	fmt.Println(countFlowers(part1Steps, part1GardenMap))

	fmt.Println("Part 2:")
	if DEBUG {
		fmt.Println(part2Steps)
		fmt.Println(part2Text)
	}
	part2GardenMap := parseGardenMap(part2Text)
	if DEBUG {
		fmt.Println(part2Steps)
		fmt.Println(part2GardenMap)
	}
	fmt.Println(countFlowersInfinite(part2Steps, part2GardenMap))
}
