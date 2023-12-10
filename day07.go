package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Hand struct {
	Cards string
	Bid   int
}

func parseHands(input string) []Hand {
	lines := strings.Split(input, "\n")
	hands := make([]Hand, len(lines))

	for i, line := range lines {
		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Println("Invalid input format")
			return nil
		}

		cards := parts[0]
		bid, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Println("Error converting bid to int:", err)
			return nil
		}

		hands[i] = Hand{Cards: cards, Bid: bid}
	}
	return hands
}

func compareCards(i string, j string) bool {
	return false // TODO
}

func calcTotalWinnings(hands []Hand) int {
	return 0 // TODO
}

func Day07() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: aoc 7 <part 1 input>")
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
	hands := parseHands(text)
	fmt.Println("Part 1:")
	if DEBUG {
		fmt.Println(hands)
	}
	fmt.Println("Part 2:")
}
