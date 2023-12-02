package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Cube struct {
	Count  int
	Colour string
}

type GameSet struct {
	Cubes []Cube
}

type Game struct {
	Id       int
	GameSets []GameSet
}

func parseGame(line string) (Game, error) {
	game := Game{}
	gameHighLevelEntry := strings.Split(line, ":")
	gameIdStr := strings.Split(gameHighLevelEntry[0], " ")[1]
	gameId, err := strconv.Atoi(gameIdStr)
	if err != nil {
		return game, err
	}
	game.Id = gameId

	gameSetsStr := strings.Split(gameHighLevelEntry[1], ";")
	for _, gameSetStr := range gameSetsStr {
		gameSet := GameSet{}

		cubesStrs := strings.Split(gameSetStr, ",")
		for _, cubeStr := range cubesStrs {
			cube := Cube{}
			cubeStrSplit := strings.Split(cubeStr, " ")
			count, err := strconv.Atoi(cubeStrSplit[1])
			if err != nil {
				return game, err
			}
			cube.Count = count
			cube.Colour = cubeStrSplit[2]
			gameSet.Cubes = append(gameSet.Cubes, cube)
		}

		game.GameSets = append(game.GameSets, gameSet)
	}

	return game, nil
}

func isPossible(game Game) bool {
	// 12 red cubes, 13 green cubes, and 14 blue cubes
	for _, gameSet := range game.GameSets {
		for _, cube := range gameSet.Cubes {
			if cube.Colour == "red" && cube.Count > 12 {
				return false
			}
			if cube.Colour == "green" && cube.Count > 13 {
				return false
			}
			if cube.Colour == "blue" && cube.Count > 14 {
				return false
			}
		}
	}
	return true
}

func Day02() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: aoc 2 <part 1 input> <part 2 input>")
		os.Exit(1)
	}

	filenamePart1 := os.Args[2]
	// Open the file
	file, err := os.Open(filenamePart1)
	if err != nil {
		log.Fatalf("Error opening file: %s", err)
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	var sum = 0
	// Read line by line
	for scanner.Scan() {
		line := scanner.Text()
		game, err := parseGame(line)
		if err != nil {
			fmt.Println("Could not parse input", err)
			os.Exit(1)
		}

		if isPossible(game) {
			sum += game.Id
		}
	}
	// Check for any scanner errors
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %s", err)
	}
	file.Close()
	fmt.Println("Part 1:")
	fmt.Println(sum)

}
