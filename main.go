package main

import (
	"fmt"
	"os"
)

func main() {

	// Check if enough arguments are provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: aoc <day> <arguments>")
		os.Exit(1)
	}

	day := os.Args[1]

	if len(day) == 0 {
		fmt.Println("Usage: aoc <day> <arguments>")
		fmt.Println("day cannot be empty")
		return
	}

	if day == "1" {
		Day1()
	} else {
		fmt.Println(day, "is not recognized.")
		os.Exit(1)
	}
}