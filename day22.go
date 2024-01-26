package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseBricks(input string) ([][2][3]int, error) {
	lines := strings.Split(input, "\n")
	result := make([][2][3]int, len(lines))

	for i, line := range lines {
		parts := strings.Split(line, "~")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid input format at line %d", i+1)
		}

		segment1 := strings.Split(parts[0], ",")
		segment2 := strings.Split(parts[1], ",")

		if len(segment1) != 3 || len(segment2) != 3 {
			return nil, fmt.Errorf("invalid number of segments at line %d", i+1)
		}

		var entry [2][3]int

		for j, numStr := range segment1 {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				return nil, fmt.Errorf("failed to convert string to int at line %d, segment %d", i+1, j+1)
			}
			entry[0][j] = num
		}

		for j, numStr := range segment2 {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				return nil, fmt.Errorf("failed to convert string to int at line %d, segment %d", i+1, j+1)
			}
			entry[1][j] = num
		}

		result[i] = entry
	}

	return result, nil
}

func canMoveDown(posnToBrickId map[[3]int]int, brick [][3]int, brickId int) bool {
	for _, brickColumn := range brick {
		downOne := [3]int{brickColumn[0],brickColumn[1],brickColumn[2]-1}

		if downOne[2] <= 0 || (posnToBrickId[downOne] != 0 && posnToBrickId[downOne] != brickId) {
			return false
		}
	}

	return true
}

func simulateFallingBricks(brickIdToPosn map[int][][3]int, posnToBrickId map[[3]int]int) {
	var somethingChanged = true

	for somethingChanged {
		somethingChanged = false

		for brickId := range brickIdToPosn {
			brickPosns := brickIdToPosn[brickId]
			
		}
	}
}

func countRemovableBricks(bricks [][2][3]int) int {

	brickIdToPosn := map[int][][3]int{}
	posnToBrickId := map[[3]int]int{}
	
	for idx, brick := range bricks {
		for direction := 0; direction < 3; direction++ {
			if brick[0][direction] != brick[1][direction] {
				lowerBrick := min(brick[0][direction], brick[1][direction])
				higherBrick := max(brick[0][direction], brick[1][direction])
				for brickColumn := lowerBrick; brickColumn <= higherBrick; brickColumn++ {
					brickColumnPosn := [3]int{brick[0][0], brick[0][1], brick[0][2]}
					brickColumnPosn[direction] = brickColumn
					
					// idx + 1 so that 0 can mean not found
					// default value of map when key is no present is 0
					posnToBrickId[brickColumnPosn] = idx+1
					brickIdToPosn[idx+1] = append(brickIdToPosn[idx+1], brickColumnPosn)
				}
			}
		}
	}

	if DEBUG {
		fmt.Println("countRemovableBricks", "brickIdToPosn", brickIdToPosn)
		fmt.Println("countRemovableBricks", "posnToBrickId", posnToBrickId)
	}

	return 0
}

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

	bricks, err := parseBricks(part1Text)
	if err != nil {
		fmt.Println("file did not parse correctly", err)
		os.Exit(-1)
	}
	if DEBUG {
		fmt.Println(bricks)
	}
	fmt.Println(countRemovableBricks(bricks))


	fmt.Println("Part 2:")
}
