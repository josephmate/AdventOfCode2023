package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func convertToNumber(line string, startCol int, endCol int) int {
	var sum = 0
	for c := startCol; c < endCol; c++ {
		sum = sum * 10
		sum += int(line[c] - '0')
	}
	return sum
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func isSymbol(lines []string, r int, c int) bool {

	if r < 0 || r >= len(lines) {
		return false
	}

	if c < 0 || c >= len(lines[r]) {
		return false
	}

	char := lines[r][c]

	return !isDigit(char) && char != '.'
}

func isDigitAdjacentToSymbol(lines []string, r int, c int) bool {
	for rDelta := -1; rDelta <= 1; rDelta++ {
		for cDelta := -1; cDelta <= 1; cDelta++ {
			if isSymbol(lines, r+rDelta, c+cDelta) {
				return true
			}
		}
	}

	return false
}

func isNumberAdjacentToSymbol(lines []string, r int, startCol int, endCol int) bool {

	for c := startCol; c < endCol; c++ {
		if isDigitAdjacentToSymbol(lines, r, c) {
			return true
		}
	}

	return false
}

func calcAdjacentNumbers(lines []string) int {
	rows := len(lines)
	cols := len(lines[0])
	var sum = 0
	for r := 0; r < rows; r++ {
		var c = 0
		for c < cols {
			startChar := lines[r][c]
			if startChar >= '0' && startChar <= '9' {
				startCol := c
				for c < cols && lines[r][c] >= '0' && lines[r][c] <= '9' {
					c++
				}
				endCol := c
				number := convertToNumber(lines[r], startCol, endCol)
				if isNumberAdjacentToSymbol(lines, r, startCol, endCol) {
					sum += number
				}
			} else {
				c++
			}
		}
	}

	return sum
}

func getAdjacentNumber(lines []string, r int, c int) (bool, int) {

	if r < 0 || r >= len(lines) {
		if DEBUG {
			fmt.Println("getAdjacentNumber", r, c, "r out of bounds")
		}
		return false, 0
	}

	if c < 0 || c >= len(lines[r]) {
		if DEBUG {
			fmt.Println("getAdjacentNumber", r, c, "c out of bounds")
		}
		return false, 0
	}

	var char = lines[r][c]
	if char < '0' || char > '9' {
		if DEBUG {
			fmt.Println("getAdjacentNumber", r, c, "not number")
		}
		return false, 0
	}

	var startCol = c
	for startCol >= 0 && lines[r][startCol] >= '0' && lines[r][startCol] <= '9' {
		if DEBUG {
			fmt.Println("getAdjacentNumber startCol", r, startCol, string(lines[r][startCol]))
		}
		startCol--
	}
	startCol++
	var endCol = c
	for endCol < len(lines[r]) && lines[r][endCol] >= '0' && lines[r][endCol] <= '9' {
		if DEBUG {
			fmt.Println("getAdjacentNumber endCol", r, endCol, string(lines[r][endCol]))
		}
		endCol++
	}

	number := convertToNumber(lines[r], startCol, endCol)
	if DEBUG {
		fmt.Println("getAdjacentNumber", r, c, "got number", number, "start", startCol, "end", endCol)
	}

	return true, number
}

func getExactlyTwoAdjacent(lines []string, r int, c int) (bool, int, int) {
	var adjacentNums []int

	for rDelta := -1; rDelta <= 1; rDelta++ {
		for cDelta := -1; cDelta <= 1; cDelta++ {
			hasNumber, number := getAdjacentNumber(lines, r+rDelta, c+cDelta)
			if hasNumber && !Contains(adjacentNums, number) {
				adjacentNums = append(adjacentNums, number)
			}
		}
	}

	if DEBUG {
		fmt.Println("getExactlyTwoAdjacent", r, c, "len", len(adjacentNums))
	}
	if len(adjacentNums) == 2 {
		return true, adjacentNums[0], adjacentNums[1]
	}
	return false, 0, 0
}

func calcAdjacentGears(lines []string) int {
	rows := len(lines)
	cols := len(lines[0])
	var sum = 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if lines[r][c] == '*' {
				hasTwo, g1, g2 := getExactlyTwoAdjacent(lines, r, c)
				if DEBUG {
					fmt.Println(hasTwo, g1, g2)
				}
				if hasTwo {
					sum += g1 * g2
				}
			}
		}
	}

	return sum
}

func Day03() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: aoc 3 <part 1 input>")
		os.Exit(1)
	}

	filenamePart1 := os.Args[2]
	// Open the file
	file, err := os.Open(filenamePart1)
	if err != nil {
		log.Fatalf("Error opening file: %s", err)
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	var lines []string
	// Read line by line
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		if err != nil {
			fmt.Println("Could not parse input", err)
			os.Exit(1)
		}
	}
	// Check for any scanner errors
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	fmt.Println("Part 1:")
	fmt.Println(calcAdjacentNumbers(lines))
	fmt.Println("Part 2:")
	fmt.Println(calcAdjacentGears(lines))
}
