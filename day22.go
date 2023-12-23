package main

import (
	"fmt"
	"os"
)

func Day22() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: aoc 22 <input>")
		os.Exit(1)
	}

	part1Text := ReadFileOrExit(os.Args[2])

	fmt.Println("Part 1:")
	if DEBUG {
		fmt.Println(part1Text)
	}

	fmt.Println("Part 2:")
}
