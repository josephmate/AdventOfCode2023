package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type CommModRecord struct {
	IsFlipFlop    bool
	IsConjunction bool
	Name          string
	Destinations  []string
}

func parseCommModRecords(input string) []CommModRecord {
	lines := strings.Split(input, "\n")
	records := make([]CommModRecord, 0)

	for _, line := range lines {
		parts := strings.Split(line, " -> ")

		name := parts[0]
		destinations := strings.Split(parts[1], ", ")

		var isFlipFlop, isConjunction bool
		if strings.HasPrefix(name, "%") {
			isFlipFlop = true
			name = strings.TrimPrefix(name, "%")
		} else if strings.HasPrefix(name, "&") {
			isConjunction = true
			name = strings.TrimPrefix(name, "&")
		}

		record := CommModRecord{
			IsFlipFlop:    isFlipFlop,
			IsConjunction: isConjunction,
			Name:          name,
			Destinations:  destinations,
		}

		records = append(records, record)
	}

	return records
}

type CommMod struct {
	IsFlipFlop    bool
	IsConjunction bool
	Name          string
	Outputs       []CommMod
	Inputs        []CommMod
}

func simulate1000(commModRecords []CommModRecord) int {
	// preprocessing for efficient lookup
	lookUpMap := map[string]CommMod{}
	for _, commModRec := range commModRecords {
		lookUpMap[commModRec.Name] = CommMod{
			IsFlipFlop:    commModRec.IsFlipFlop,
			IsConjunction: commModRec.IsConjunction,
			Name:          commModRec.Name,
		}
	}

	// hook up everything as a graph

	return 0
}

func Day20() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: aoc 20 <input>")
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
	commModRecords := parseCommModRecords(text)
	if DEBUG {
		fmt.Println(commModRecords)
	}
	fmt.Println("Part 1:")
	fmt.Println(simulate1000(commModRecords))
	fmt.Println("Part 2:")
}
