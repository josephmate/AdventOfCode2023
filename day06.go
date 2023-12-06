package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

// func parseTimetable(input string) [2][]int {
// 	matrix := [2][]int{}

// 	lines := strings.Split(input, "\n")
// 	for i, line := range lines {
// 		fields := strings.Fields(line)
// 		values []int, len(fields)-1)

// 		for j, field := range fields[1:] {
// 			value, err := strconv.Atoi(field)
// 			if err != nil {
// 				panic(err)
// 			}
// 			values[j] = value
// 		}
// 		matrix[i] = values
// 	}

// 	return matrix
// }

func Day06() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: aoc 5 <part 1 input>")
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
	fmt.Println(text)

	fmt.Println("Part 1:")
	fmt.Println("Part 2:")
}
