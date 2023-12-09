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

func allEndsWithZ(currents []string) bool {
	for _, current := range currents {
		if !strings.HasSuffix(current, "Z") {
			return false
		}
	}
	return true
}

// TODO: takes too long to
//
//	next try to detect repeats from each of the parallel paths and figure out when all parallel
//	paths align.
func runUntilAllEndsWithZ(directions string, mapping map[string][2]string) int {
	var currents = []string{}
	for key := range mapping {
		if strings.HasSuffix(key, "A") {
			currents = append(currents, key)
		}
	}
	var steps = 0

	if DEBUG {
		fmt.Println("runUntilAllEndsWithZ", currents)
	}
	for !allEndsWithZ(currents) {
		var nextCurrents = []string{}
		for _, current := range currents {
			nextStepDirection := directions[steps%len(directions)]
			if nextStepDirection == 'R' {
				current = mapping[current][1]
			} else {
				current = mapping[current][0]
			}
			nextCurrents = append(nextCurrents, current)
		}
		currents = nextCurrents
		steps++
		if DEBUG {
			fmt.Println("runUntilAllEndsWithZ", currents)
		}
	}

	return steps
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

	if len(os.Args) < 4 {
		fmt.Println("Usage: aoc 5 <part 1 input>")
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

	file2, err := os.Open(filenamePart2)
	if err != nil {
		log.Fatalf("Error opening file: %s", err)
	}
	defer file2.Close()
	reader2 := io.Reader(file2)
	textBytes2, err := io.ReadAll(reader2)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	text2 := string(textBytes2)
	directions2, mapping2 := parseDay08(text2)
	fmt.Println(runUntilAllEndsWithZ(directions2, mapping2))

}
