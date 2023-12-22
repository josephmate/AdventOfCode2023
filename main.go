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
		Day01()
	} else if day == "2" {
		Day02()
	} else if day == "3" {
		Day03()
	} else if day == "4" {
		Day04()
	} else if day == "5" {
		Day05()
	} else if day == "6" {
		Day06()
	} else if day == "7" {
		Day07()
	} else if day == "8" {
		Day08()
	} else if day == "9" {
		Day09()
	} else if day == "10" {
		Day10()
	} else if day == "11" {
		Day11()
	} else if day == "13" {
		Day13()
	} else if day == "14" {
		Day14()
	} else if day == "15" {
		Day15()
	} else if day == "16" {
		Day16()
	} else if day == "17" {
		Day17()
	} else if day == "18" {
		Day18()
	} else if day == "20" {
		Day20()
	} else if day == "21" {
		Day21()
	} else if day == "22" {
		Day22()
	} else {
		fmt.Println(day, "is not recognized.")
		os.Exit(1)
	}
}
