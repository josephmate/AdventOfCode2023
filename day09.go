package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func parseDay09(input string) [][]int {
	lines := strings.Split(input, "\n")
	result := make([][]int, len(lines))

	for i, line := range lines {
		values := strings.Fields(line)
		intValues := make([]int, len(values))

		for j, val := range values {
			num, err := strconv.Atoi(val)
			if err != nil {
				fmt.Println("Error converting string to int:", err)
				return nil
			}
			intValues[j] = num
		}
		result[i] = intValues
	}
	return result
}

func allZero(deltas []int) bool {
	for _, delta := range deltas {
		if delta != 0 {
			return false
		}
	}
	return true
}

func predictNextValue(history []int) int {
	var lastEntries []int
	var deltas [][]int
	deltas = append(deltas, history)
	lastEntries = append(lastEntries, history[len(history)-1])
	for !allZero(deltas[len(deltas)-1]) {
		var nextRow []int
		prevRow := deltas[len(deltas)-1]
		var prevVal = prevRow[0]

		for _, val := range prevRow[1:] {
			nextRow = append(nextRow, val-prevVal)
			prevVal = val
		}

		deltas = append(deltas, nextRow)
		lastEntries = append(lastEntries, nextRow[len(nextRow)-1])
	}

	var result = 0
	for i := len(lastEntries) - 1; i >= 0; i-- {
		result += lastEntries[i]
	}
	return result
}

func predictNextValues(histories [][]int) int {
	var sum = 0
	for _, history := range histories {
		sum += predictNextValue(history)
	}
	return sum
}

func predictPrevValue(history []int) int {
	var firstEntries []int
	var deltas [][]int
	deltas = append(deltas, history)
	firstEntries = append(firstEntries, history[0])
	for !allZero(deltas[len(deltas)-1]) {
		var nextRow []int
		prevRow := deltas[len(deltas)-1]
		var prevVal = prevRow[0]

		for _, val := range prevRow[1:] {
			nextRow = append(nextRow, val-prevVal)
			prevVal = val
		}

		deltas = append(deltas, nextRow)
		firstEntries = append(firstEntries, nextRow[0])
	}

	var result = 0
	for i := len(firstEntries) - 1; i >= 0; i-- {
		if DEBUG {
			fmt.Println("predictPrevValue", firstEntries[i]-result, "-", firstEntries[i], "-", result)
		}
		result = firstEntries[i] - result
	}
	return result
}

func predictPrevValues(histories [][]int) int {
	var sum = 0
	for _, history := range histories {
		result := predictPrevValue(history)
		if DEBUG {
			fmt.Println("predictPrevValues", result, history)
		}
		sum += result
	}
	return sum
}

func Day09() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: aoc 9 <part 1 input>")
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
	histories := parseDay09(text)
	fmt.Println("Part 1:")
	if DEBUG {
		fmt.Println(histories)
	}
	fmt.Println(predictNextValues(histories))
	fmt.Println("Part 2:")
	fmt.Println(predictPrevValues(histories))
}
