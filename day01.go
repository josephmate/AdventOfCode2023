package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func calcNumber(str string) int {
	var firstLast []int
	for i := 0; i < len(str); i++ {
		char := str[i]
		if char >= '0' && char <= '9' {
			firstLast = append(firstLast, int(char-'0'))
		}
	}
	return firstLast[0]*10 + firstLast[len(firstLast)-1]
}

func Day1() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: aoc 1 <inputfile>")
		os.Exit(1)
	}

	filename := os.Args[2]
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening file: %s", err)
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	var sum = 0
	// Read line by line
	for scanner.Scan() {
		line := scanner.Text()
		sum += calcNumber(line)
		calcNumber(line)
	}

	fmt.Println("Part 1:")
	fmt.Println(sum)

	// Check for any scanner errors
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %s", err)
	}
}
