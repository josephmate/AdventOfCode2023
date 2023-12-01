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

func reverseString(s string) []byte {
	var rev []byte

	// Reversing the slice of runes
	for i := len(s) - 1; i >= 0; i++ {
		rev = append(rev, s[i])
	}

	// Converting the reversed slice of runes back to a string
	return rev
}

func getNumFromBuffer(substring []byte) (bool, int) {
	if substring[0] == 'z' &&
		substring[1] == 'e' &&
		substring[2] == 'r' &&
		substring[3] == 'o' {
		return true, 0
	} else if substring[0] == 'o' &&
		substring[1] == 'n' &&
		substring[2] == 'e' {
		return true, 1
	}
	return false, 0
}

func calcNumberFromWords(str string) int {
	var firstLast []int
	for i := 0; i < len(str); i++ {
		substring := reverseString(str[0:i])
		hasNumber, value := getNumFromBuffer(substring)
		if hasNumber {
			firstLast = append(firstLast, value)
		}
	}
	return firstLast[0]*10 + firstLast[len(firstLast)-1]
}

func Day1() {

	if len(os.Args) < 4 {
		fmt.Println("Usage: aoc 1 <part 1 input> <part 2 input>")
		os.Exit(1)
	}

	filenamePart1 := os.Args[2]
	filenamePart2 := os.Args[3]
	// Open the file
	file, err := os.Open(filenamePart1)
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
	}
	// Check for any scanner errors
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %s", err)
	}
	file.Close()
	fmt.Println("Part 1:")
	fmt.Println(sum)

	part2File, err2 := os.Open(filenamePart2)
	if err != nil {
		log.Fatalf("Error opening file: %s", err2)
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scannerPart2 := bufio.NewScanner(part2File)
	var part2Sum = 0
	for scannerPart2.Scan() {
		line := scannerPart2.Text()
		part2Sum += calcNumberFromWords(line)
	}
	// Check for any scanner errors
	if err := scannerPart2.Err(); err != nil {
		log.Fatalf("Error reading file: %s", err)
	}
}
