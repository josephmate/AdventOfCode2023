package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

func parseDay08(text string) (string, map[string][2]string) {
	splitByTwoSections := regexp.MustCompile(`\r?\n\r?\n`)
	sections := splitByTwoSections.Split(text, 2)
	firstList := strings.TrimSpace(sections[0])
	charMappings := strings.TrimSpace(sections[1])
	lines := strings.Split(charMappings, "\n")
	m := make(map[string][2]string)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, " = ")
		key := parts[0]
		values := strings.Trim(parts[1], "()")
		valueSlice := strings.Split(values, ", ")
		if len(valueSlice) >= 2 {
			m[key] = [2]string{valueSlice[0], valueSlice[1]}
		}
	}

	return firstList, m
}

func runUntilZZZ(directions string, mapping map[string][2]string) int {
	var current = "AAA"
	var steps = 0

	for current != "ZZZ" {
		nextStepDirection := directions[steps%len(directions)]
		if nextStepDirection == 'R' {
			current = mapping[current][1]
		} else {
			current = mapping[current][0]
		}

		steps++
	}

	return steps
}

func Day08() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: aoc 5 <part 1 input>")
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

	fmt.Println("Part 1:")
	directions, mapping := parseDay08(text)
	if DEBUG {
		fmt.Println(directions)
		fmt.Println(mapping)
	}
	fmt.Println(runUntilZZZ(directions, mapping))
	fmt.Println("Part 2:")
}
