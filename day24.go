package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type HailRecord struct {
	Position [3]int64
	Velocity [3]int64
}

func parseHailRecords(input string) []HailRecord {
	lines := strings.Split(input, "\n")
	records := make([]HailRecord, len(lines))

	for i, line := range lines {
		fields := strings.Split(line, "@")
		posFields := strings.Split(strings.TrimSpace(fields[0]), ",")
		velFields := strings.Split(strings.TrimSpace(fields[1]), ",")

		record := HailRecord{}

		for j := 0; j < 3; j++ {
			posVal, _ := strconv.ParseInt(strings.TrimSpace(posFields[j]), 10, 64)
			velVal, _ := strconv.ParseInt(strings.TrimSpace(velFields[j]), 10, 64)

			record.Position[j] = posVal
			record.Velocity[j] = velVal
		}

		records[i] = record
	}

	return records
}

func Day24() {

	if len(os.Args) < 5 {
		fmt.Println("Usage: aoc 24 lowerbound upperbound <input>")
		os.Exit(1)
	}

	lowerbound := ParseIntOrExit(os.Args[2])
	upperbound := ParseIntOrExit(os.Args[3])
	text := ReadFileOrExit(os.Args[4])

	fmt.Println("Part 1:")
	if DEBUG {
		fmt.Println(text)
		fmt.Println(lowerbound)
		fmt.Println(upperbound)
	}
	hailRecords := parseHailRecords(text)
	if DEBUG {
		fmt.Println(hailRecords)
		fmt.Println(lowerbound)
		fmt.Println(upperbound)
	}

	fmt.Println("Part 2:")
}
