package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Scratchcard struct {
	Id      int
	Winners []int
	Plays   []int
}

func parseNumbers(substring string) ([]int, error) {
	splitByContiguousSpaces := regexp.MustCompile(` +`)
	numberStrs := splitByContiguousSpaces.Split(substring, -1)

	var nums []int
	for _, numStr := range numberStrs {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return nums, err
		}
		nums = append(nums, num)
	}

	return nums, nil
}

func parseScratchcard(line string) (Scratchcard, error) {
	scratch := Scratchcard{}
	splitByContiguousSpaces := regexp.MustCompile(` +`)

	upperEntrySplit := strings.Split(line, ":")
	gameIdStr := splitByContiguousSpaces.Split(upperEntrySplit[0], -1)[1]
	gameId, err := strconv.Atoi(gameIdStr)
	if err != nil {
		return scratch, err
	}
	scratch.Id = gameId

	innerEntrySplit := strings.Split(upperEntrySplit[1], "|")
	winnersStr := strings.TrimSpace(innerEntrySplit[0])
	winners, err := parseNumbers(winnersStr)
	if err != nil {
		return scratch, err
	}
	scratch.Winners = winners
	playsStr := strings.TrimSpace(innerEntrySplit[1])
	plays, err := parseNumbers(playsStr)
	if err != nil {
		return scratch, err
	}
	scratch.Plays = plays

	return scratch, nil
}

func scratchScore(matches int) int {
	if matches == 0 {
		return 0
	}

	var score = 1
	for i := 1; i < matches; i++ {
		score = score * 2
	}

	if DEBUG {
		fmt.Println("scratchScore", matches, score)
	}
	return score
}

func calcScore(scratch Scratchcard) int {
	var matches = 0

	// all though O(N^2), the input is bounded by a small constant 10*26
	// which executes quickly
	for _, num := range scratch.Plays {
		if Contains(scratch.Winners, num) {
			if DEBUG {
				fmt.Println("calcScore found matches", scratch.Winners, scratch.Plays, num)
			}
			matches++
		}
	}

	score := scratchScore(matches)
	if DEBUG {
		fmt.Println("calcScore", matches, score)
	}
	return score
}

func Day04() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: aoc 2 <part 1 input>")
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

	// Read line by line
	var sum = 0
	for scanner.Scan() {
		line := scanner.Text()
		scratch, err := parseScratchcard(line)
		if err != nil {
			fmt.Println(scratch, err)
			continue
		}
		sum += calcScore(scratch)
		if err != nil {
			fmt.Println("Could not parse input", err)
			os.Exit(1)
		}
	}

	fmt.Println("Part 1:")
	fmt.Println(sum)

	// Check for any scanner errors
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

}
