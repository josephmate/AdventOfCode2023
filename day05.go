package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type RangeDesc struct {
	DestRangeStart   int
	SourceRangeStart int
	RangeLength      int
}

type BySourceRange []RangeDesc

func (a BySourceRange) Len() int           { return len(a) }
func (a BySourceRange) Less(i, j int) bool { return a[i].SourceRangeStart < a[j].SourceRangeStart }
func (a BySourceRange) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

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

func parseMap(section string, startLineNum int) ([]RangeDesc, int, error) {
	var linNum = startLineNum
	var data []RangeDesc
	lines := strings.Split(section, "\n")

	// Skip the first line ("seed-to-soil map:")
	if len(lines) > 0 {
		lines = lines[1:]
	}
	linNum++

	// Read and parse lines into the map
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		fields := strings.Fields(line)
		linNum++

		if len(fields) != 3 {
			return data, linNum, fmt.Errorf("invalid input format lineNum=%d. Expected 3 fields got %d. Line is '%s'", linNum, len(fields), line)
		}

		col1, err := strconv.Atoi(fields[0])
		if err != nil {
			return data, linNum, fmt.Errorf("error converting key to int: %s", err)
		}

		col2, err := strconv.Atoi(fields[1])
		if err != nil {
			return data, linNum, fmt.Errorf("error converting first value to int: %s", err)
		}

		col3, err := strconv.Atoi(fields[2])
		if err != nil {
			return data, linNum, fmt.Errorf("error converting second value to int: %s", err)
		}

		data = append(data, RangeDesc{
			DestRangeStart:   col1,
			SourceRangeStart: col2,
			RangeLength:      col3,
		})
	}

	return data, linNum, nil
}

func parseAlmanac(input string) (Almanac, error) {
	var lineNum = 0
	result := Almanac{}

	splitByDoubleNewline := regexp.MustCompile(`\n\n`)
	sections := splitByDoubleNewline.Split(input, -1)

	seeds, err := parseSeeds(sections[0])
	if err != nil {
		return result, err
	}
	result.Seeds = seeds
	lineNum++
	lineNum++

	for i := 0; i < 7; i++ {
		maps, nextLineNum, err := parseMap(sections[i+1], lineNum)
		if err != nil {
			return result, fmt.Errorf("error converting section %d to map starting at line %d. %s", i+1, lineNum, err)
		}
		lineNum = nextLineNum
		result.ConvertMaps = append(result.ConvertMaps, maps)
	}

	return result, nil
}

func getLowestNumInRange(convertMap []RangeDesc, id int) (int, bool) {
	sort.Sort(BySourceRange(convertMap))

	for _, rangeDesc := range convertMap {
		if rangeDesc.SourceRangeStart <= id && id < (rangeDesc.SourceRangeStart+rangeDesc.RangeLength) {
			return rangeDesc.DestRangeStart + (id - rangeDesc.SourceRangeStart), true
		}
	}

	return 0, false
}

func calcLocationNumber(almanac Almanac, seed int) int {
	var currentId = seed
	for i := 0; i < len(almanac.ConvertMaps); i++ {
		nextId, hasIt := getLowestNumInRange(almanac.ConvertMaps[i], currentId)
		if hasIt {
			currentId = nextId
		}
		if DEBUG {
			fmt.Println("calcLocationNumber inner", "seed=", seed, "level=", i, "currentId=", currentId)
		}
	}
	if DEBUG {
		fmt.Println("calcLocationNumber", "seed=", seed, "location=", currentId)
	}
	return currentId
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

// TODO: need to think more about all the edge cases
//
//	draw a diagram
//
// assumes convertMap is sorted
func mapToNext(convertMap []RangeDesc, seed int, seedLen int) (int, int) {
	// for _, rangeDesc := range convertMap {
	// 	// fully fits into the range
	// 	if rangeDesc.SourceRangeStart <= id && id < (rangeDesc.SourceRangeStart+rangeDesc.RangeLength) {
	// 		return rangeDesc.DestRangeStart + (id - rangeDesc.SourceRangeStart), rangeDesc.RangeLength
	// 	}
	// 	// part of the range is satified
	// 	if rangeDesc.SourceRangeStart > id && id < (rangeDesc.SourceRangeStart+rangeDesc.RangeLength) {
	// 		// if nothing is found, it just maps to itself
	// 		return rangeDesc.SourceRangeStart, id - rangeDesc.SourceRangeStart
	// 	}
	// 	if rangeDesc.SourceRangeStart <= id && id < (rangeDesc.SourceRangeStart+rangeDesc.RangeLength) {
	// 		return rangeDesc.DestRangeStart + (id - rangeDesc.SourceRangeStart), rangeDesc.RangeLength
	// 	}
	// }
	// // fully does not map to anything
	// return rangeDesc.SourceRangeStart, seedLen
	return 0, 0
}

func calcNextLevel(convertMap []RangeDesc, currentlyProcessing [][2]int) [][2]int {
	sort.Sort(BySourceRange(convertMap))
	var nextLevel = [][2]int{}

	for _, startRangeSlice := range currentlyProcessing {
		var seed = startRangeSlice[0]
		var seedLen = startRangeSlice[1]
		for seedLen > 0 {
			var destination, destinationLen = mapToNext(convertMap, seed, seedLen)
			nextLevel = append(nextLevel, [2]int{destination, destinationLen})

			if destinationLen > seedLen {
				destinationLen = seedLen
			}
			seed += destinationLen
			seedLen -= destinationLen
		}
	}

	return nextLevel
}

func calcLowestLocationNumberRanged(almanac Almanac) int {
	var currentlyProcessing = [][2]int{}
	for i := 0; i < len(almanac.Seeds); i += 2 {
		seedRangePair := [2]int{almanac.Seeds[i], almanac.Seeds[i+1]}
		currentlyProcessing = append(currentlyProcessing, seedRangePair)
	}
	for i := 0; i < len(almanac.ConvertMaps); i++ {
		currentlyProcessing = calcNextLevel(almanac.ConvertMaps[i], currentlyProcessing)
	}

	var min = currentlyProcessing[0][0]
	for i := 1; i < len(currentlyProcessing); i++ {
		curr := currentlyProcessing[0][0]
		if curr < min {
			min = curr
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
	fmt.Println(calcLowestLocationNumberRanged(almanac))
}
