package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func parseRollingBolders(text string) [][]byte {
	var rollingBolders [][]byte
	for _, row := range strings.Split(text, "\n") {
		trimmedPart := strings.TrimSpace(row)
		if trimmedPart != "" {
			rollingBolders = append(rollingBolders, []byte(trimmedPart))
		}
	}

	return rollingBolders
}

func findLastEmptyNorthSpace(rollingBolders [][]byte, startR int, startC int) int {
	var r = startR
	for r >= 0 && rollingBolders[r][startC] != 'O' && rollingBolders[r][startC] != '#' {
		r--
	}
	return r + 1
}

func moveEverythingNorth(rollingBolders [][]byte) int {
	for r, row := range rollingBolders {
		for c, col := range row {
			if col == 'O' {
				rollingBolders[r][c] = '.'
				newRow := findLastEmptyNorthSpace(rollingBolders, r, c)
				fmt.Println("moveEverythingNorth", "r=", r, "c=", c, "newRow=", newRow)
				rollingBolders[newRow][c] = 'O'
			}
		}
	}

	if DEBUG {
		fmt.Println("moveEverythingNorth")
		for _, row := range rollingBolders {
			fmt.Println(string(row))
		}
	}

	var sum = 0
	for r, row := range rollingBolders {
		for _, col := range row {
			if col == 'O' {
				sum += len(rollingBolders) - r
			}
		}
	}

	return sum
}

func Day14() {

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
	rollingBolders := parseRollingBolders(text)
	if DEBUG {
		fmt.Println(rollingBolders)
	}
	fmt.Println("Part 1:")
	fmt.Println(moveEverythingNorth(rollingBolders))

	fmt.Println("Part 2:")
}
