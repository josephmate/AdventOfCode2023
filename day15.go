package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func parseInitSeq(input string) []string {
	var result []string
	for _, column := range strings.Split(input, ",") {
		trimmed := strings.TrimSpace(column)
		if column != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func hashStep(step string) int {
	ascii := []byte(step)

	var currentHash = 0
	for _, num := range ascii {
		currentHash += int(num)
		currentHash *= 17
		currentHash = currentHash % 256
	}

	return currentHash
}

func hashInitSeq(initSeq []string) int {
	var sum = 0
	for _, step := range initSeq {
		sum += hashStep(step)
	}

	return sum
}

func Day15() {

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
	initSeq := parseInitSeq(text)
	if DEBUG {
		fmt.Println(initSeq)
	}
	fmt.Println("Part 1:")
	fmt.Println(hashInitSeq(initSeq))
	fmt.Println("Part 2:")
}
