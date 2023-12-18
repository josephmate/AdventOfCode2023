package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type DigStep struct {
	Direction string
	Distance  int
	RGBCode   string
}

func parseDigPlan(input string) []DigStep {
	lines := strings.Split(input, "\n")
	var steps []DigStep

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)

		direction := parts[0]
		distance, _ := strconv.Atoi(parts[1])
		rgbCode := parts[2]

		step := DigStep{
			Direction: direction,
			Distance:  distance,
			RGBCode:   rgbCode,
		}

		steps = append(steps, step)
	}

	return steps
}

func printDigMap(digMap map[[2]int]bool) {
	var minR = 0
	var minC = 0
	var maxR = 0
	var maxC = 0
	for key := range digMap {
		minR = Min(minR, key[0])
		maxR = Max(maxR, key[0])
		minC = Min(minC, key[1])
		maxC = Max(maxC, key[1])
	}

	fmt.Println("printDigMap", "minR", minR, "maxR", maxR, "minC", minC, "maxC", maxC)

	for r := minR; r <= maxR; r++ {
		for c := minC; c <= maxC; c++ {
			if digMap[[2]int{r, c}] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func calcAreaOfDigPlan(digPlan []DigStep) int {
	// fill in the dig map
	digMap := map[[2]int]bool{}
	var currentRow = 0
	var currentCol = 0
	digMap[[2]int{0, 0}] = true
	for _, digStep := range digPlan {
		if DEBUG {
			fmt.Println("calcAreaOfDigPlan digStep=", digStep)
		}
		var rDelta = 0
		var cDelta = 0
		if digStep.Direction == "R" {
			cDelta = 1
		} else if digStep.Direction == "L" {
			cDelta = -1
		} else if digStep.Direction == "U" {
			rDelta = -1
		} else if digStep.Direction == "D" {
			rDelta = 1
		} else {
			fmt.Println("calcAreaOfDigPlan unrecognized direction", digStep.Direction)
			panic("calcAreaOfDigPlan unrecognized direction")
		}

		for i := 0; i < digStep.Distance; i++ {
			currentRow += rDelta
			currentCol += cDelta
			digMap[[2]int{currentRow, currentCol}] = true
		}
		if DEBUG {
			fmt.Println("calcAreaOfDigPlan currentRow=", currentRow, "currentCol=", currentCol)
		}
	}
	if DEBUG {
		fmt.Println("calcAreaOfDigPlan len(digMap)=", len(digMap))
	}
	if DEBUG {
		fmt.Println("calcAreaOfDigPlan")
		printDigMap(digMap)
	}

	// starting from the top left
	var topLeft = [2]int{0, 0}
	for key := range digMap {
		if key[0] < topLeft[0] {
			topLeft = key
		} else if key[0] == topLeft[0] && key[1] < topLeft[1] {
			topLeft = key
		}
	}

	// calculate the area using BFS
	visited := map[[2]int]bool{}
	var queue [][2]int
	startPosn := [2]int{topLeft[0] + 1, topLeft[1] + 1}
	queue = append(queue, startPosn)
	for len(queue) > 0 {
		currentPosn := queue[0]
		queue = queue[1:]

		if visited[currentPosn] || digMap[currentPosn] {
			continue
		}

		visited[currentPosn] = true
		queue = append(queue, [2]int{currentPosn[0] + 1, currentPosn[1]})
		queue = append(queue, [2]int{currentPosn[0] - 1, currentPosn[1]})
		queue = append(queue, [2]int{currentPosn[0], currentPosn[1] + 1})
		queue = append(queue, [2]int{currentPosn[0], currentPosn[1] - 1})
	}
	return len(visited) + len(digMap)
}

func Day18() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: aoc 18 <input>")
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
	digSteps := parseDigPlan(text)
	if DEBUG {
		fmt.Println(digSteps)
	}
	fmt.Println("Part 1:")
	fmt.Println(calcAreaOfDigPlan(digSteps))
	fmt.Println("Part 2:")
}
