package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

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

	// Read line by line
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line) // Process each line as needed
	}

	// Check for any scanner errors
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %s", err)
	}
}
