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

func getNumFromBuffer(substring string) (bool, int) {
	if len(substring) >= 4 &&
		substring[0] == 'z' &&
		substring[1] == 'e' &&
		substring[2] == 'r' &&
		substring[3] == 'o' {
		return true, 0
	} else if len(substring) >= 3 &&
		substring[0] == 'o' &&
		substring[1] == 'n' &&
		substring[2] == 'e' {
		return true, 1
	} else if len(substring) >= 3 &&
		substring[0] == 't' &&
		substring[1] == 'w' &&
		substring[2] == 'o' {
		return true, 2
	} else if len(substring) >= 5 &&
		substring[0] == 't' &&
		substring[1] == 'h' &&
		substring[2] == 'r' &&
		substring[3] == 'e' &&
		substring[4] == 'e' {
		return true, 3
	} else if len(substring) >= 4 &&
		substring[0] == 'f' &&
		substring[1] == 'o' &&
		substring[2] == 'u' &&
		substring[3] == 'r' {
		return true, 4
	} else if len(substring) >= 4 &&
		substring[0] == 'f' &&
		substring[1] == 'i' &&
		substring[2] == 'v' &&
		substring[3] == 'e' {
		return true, 5
	} else if len(substring) >= 3 &&
		substring[0] == 's' &&
		substring[1] == 'i' &&
		substring[2] == 'x' {
		return true, 6
	} else if len(substring) >= 5 &&
		substring[0] == 's' &&
		substring[1] == 'e' &&
		substring[2] == 'v' &&
		substring[3] == 'e' &&
		substring[4] == 'n' {
		return true, 7
	} else if len(substring) >= 5 &&
		substring[0] == 'e' &&
		substring[1] == 'i' &&
		substring[2] == 'g' &&
		substring[3] == 'h' &&
		substring[4] == 't' {
		return true, 8
	} else if len(substring) >= 4 &&
		substring[0] == 'n' &&
		substring[1] == 'i' &&
		substring[2] == 'n' &&
		substring[3] == 'e' {
		return true, 9
	}
	return false, 0
}

func calcNumberFromWords(str string) int {
	var firstLast []int
	for i := 0; i < len(str); i++ {
		char := str[i]
		if char >= '0' && char <= '9' {
			firstLast = append(firstLast, int(char-'0'))
		} else {
			substring := str[i:len(str)]
			hasNumber, value := getNumFromBuffer(substring)
			if hasNumber {
				firstLast = append(firstLast, value)
			}
		}
	}
	return firstLast[0]*10 + firstLast[len(firstLast)-1]
}

func Day01() {

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
		fmt.Println(line)
		num := calcNumberFromWords(line)
		fmt.Println(num)
		part2Sum += num
	}
	fmt.Println("Part 2:")
	fmt.Println(part2Sum)

	// Check for any scanner errors
	if err := scannerPart2.Err(); err != nil {
		log.Fatalf("Error reading file: %s", err)
	}
}
