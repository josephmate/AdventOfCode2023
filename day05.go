package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type RangeDesc struct {
	DestRangeStart   int
	SourceRangeStart int
	RangeLength      int
}

type Almanac struct {
	Seeds       []int
	ConvertMaps [][]RangeDesc
}

func parseSeeds(line string) ([]int, error) {
	numStrs := strings.Split(line, " ")

	var seeds []int

	for _, seedStr := range numStrs[1:] {
		seed, err := strconv.Atoi(seedStr)
		if err != nil {
			return seeds, err
		}
		seeds = append(seeds, seed)
	}
	return seeds, nil
}

func parseMap(section string) ([]RangeDesc, error) {
	var data []RangeDesc
	lines := strings.Split(section, "\n")

	// Skip the first line ("seed-to-soil map:")
	if len(lines) > 0 {
		lines = lines[1:]
	}

	// Read and parse lines into the map
	for _, line := range lines {
		fields := strings.Fields(line)

		if len(fields) != 3 {
			return data, errors.New("invalid input format")
		}

		col1, err := strconv.Atoi(fields[0])
		if err != nil {
			return data, fmt.Errorf("error converting key to int: %s", err)
		}

		col2, err := strconv.Atoi(fields[1])
		if err != nil {
			return data, fmt.Errorf("error converting first value to int: %s", err)
		}

		col3, err := strconv.Atoi(fields[2])
		if err != nil {
			return data, fmt.Errorf("error converting second value to int: %s", err)
		}

		data = append(data, RangeDesc{
			DestRangeStart:   col1,
			SourceRangeStart: col2,
			RangeLength:      col3,
		})
	}

	return data, nil
}

func parseAlmanac(input string) (Almanac, error) {
	result := Almanac{}

	splitByDoubleNewline := regexp.MustCompile(`\n\n`)
	sections := splitByDoubleNewline.Split(input, -1)

	seeds, err := parseSeeds(sections[0])
	if err != nil {
		return result, err
	}
	result.Seeds = seeds

	for i := 0; i < 7; i++ {
		maps, err := parseMap(sections[i+1])
		if err != nil {
			return result, fmt.Errorf("error converting section %d to map. %s", i+1, err)
		}
		result.ConvertMaps = append(result.ConvertMaps, maps)
	}

	return result, nil
}

func calcLocationNumber(almanac Almanac, seed int) int {

}

func calcLowestLocationNumber(almanac Almanac) int {
	var min = calcLocationNumber(almanac, almanac.Seeds[0])
	for i := 1; i < len(almanac.Seeds); i++ {
		locNum := calcLocationNumber(almanac, almanac.Seeds[i])
		if locNum < min {
			min = locNum
		}
	}
	return min
}

func Day05() {

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

	// Convert the byte slice to a string
	text := string(textBytes)
	almanac, err := parseAlmanac(text)
	if err != nil {
		fmt.Println("Could not parse", filenamePart1, "got error", err)
		os.Exit(1)
	}
	fmt.Println(almanac)

	fmt.Println("Part 1:")
	fmt.Println(calcLowestLocationNumber(almanac))
	fmt.Println("Part 2:")
}
