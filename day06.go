package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func parseTimetable(input string) [2][]int {
	timetable := [2][]int{}

	lines := strings.Split(input, "\n")
	// only two lines
	// ingnore empty line at end if its there
	for i, line := range lines[0:2] {
		fields := strings.Fields(line)
		values := make([]int, len(fields)-1)

		// 0th field is the column header
		for j, field := range fields[1:] {
			value, err := strconv.Atoi(field)
			if err != nil {
				panic(err)
			}
			values[j] = value
		}
		timetable[i] = values
	}

	return timetable
}

func calcNumberOfWaysToWin(numMsecAvailable int, distanceMilliM int) int {
	var numOfWays = 0
	for speed := 1; speed < numMsecAvailable; speed++ {
		if (numMsecAvailable-speed)*speed > distanceMilliM {
			if DEBUG {
				fmt.Println("calcNumberOfWaysToWin", "numMsecAvailable=", numMsecAvailable, "speed=", speed, "distanceMilliM", distanceMilliM, "distanceTravelled=", (numMsecAvailable-speed)*speed)
			}
			numOfWays++
		}
	}
	return numOfWays
}

func calcMarginOfRaceError(timetable [2][]int) int {
	var marginOfError = 1

	for c := 0; c < len(timetable[0]); c++ {
		numOfWaysToWin := calcNumberOfWaysToWin(timetable[0][c], timetable[1][c])
		if DEBUG {
			fmt.Println("calcMarginOfRaceError", "numMsecAvailable=", timetable[0][c], "distanceMilliM=", timetable[1][c], "numOfWaysToWin", numOfWaysToWin)
		}
		marginOfError = marginOfError * numOfWaysToWin
	}

	return marginOfError
}

func parseTimetableIgnoreSpaces(input string) (int, int) {
	result := [2]int{}
	lines := strings.Split(input, "\n")
	// only two lines
	// ingnore empty line at end if its there
	for i, line := range lines[0:2] {
		fields := strings.Fields(line)
		combinedNumStr := strings.Join(fields[1:], "")
		value, err := strconv.Atoi(combinedNumStr)
		if err != nil {
			panic(err)
		}
		result[i] = value
	}
	return result[0], result[1]
}

func Day06() {

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
	timetable := parseTimetable(text)

	fmt.Println("Part 1:")
	fmt.Println(calcMarginOfRaceError(timetable))
	fmt.Println("Part 2:")
	time, dist := parseTimetableIgnoreSpaces(text)
	if DEBUG {
		fmt.Println(time, dist)
	}
	fmt.Println(calcNumberOfWaysToWin(time, dist))
}
