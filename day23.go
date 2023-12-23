package main

import (
	"fmt"
	"os"
)

type HikingPath struct {
	CurrentPosition [2]int
	Visited         map[[2]int]bool
}

func cloneHikingPath(hikingPath HikingPath) HikingPath {
	cloned := HikingPath{}

	cloned.Visited = make(map[[2]int]bool)
	for key, value := range hikingPath.Visited {
		cloned.Visited[key] = value
	}

	return cloned
}

func printHikingPath(hikingMap [][]byte, longestSoFarPath HikingPath) {
	if !DEBUG {
		return
	}

	for r, row := range hikingMap {
		for c, col := range row {
			if longestSoFarPath.Visited[[2]int{r, c}] {
				fmt.Print("O")
			} else {
				fmt.Print(string(col))
			}
		}
		fmt.Println()
	}
}

func findLongestPath(hikingMap [][]byte) int {
	startPosn := [2]int{0, 1}
	numRows := len(hikingMap)
	numCols := len(hikingMap[0])
	endPosn := [2]int{numRows - 1, numCols - 2}

	var queue []HikingPath
	startPath := HikingPath{}
	startPath.CurrentPosition = startPosn
	startPath.Visited = map[[2]int]bool{}
	startPath.Visited[startPosn] = true
	queue = append(queue, startPath)

	longestSoFar := map[[2]int]int{}
	longestSoFarPath := map[[2]int]HikingPath{}

	for len(queue) > 0 {
		currentPath := queue[0]
		queue = queue[1:]

		if DEBUG {
			fmt.Println(currentPath)
		}

		// we already have a longer path thru this point
		if longestSoFar[currentPath.CurrentPosition] >= len(currentPath.Visited) {
			continue
		}
		longestSoFar[currentPath.CurrentPosition] = len(currentPath.Visited)
		longestSoFarPath[currentPath.CurrentPosition] = currentPath

		// out of bounds
		if currentPath.CurrentPosition[0] < 0 ||
			currentPath.CurrentPosition[0] >= len(hikingMap) ||
			currentPath.CurrentPosition[1] < 0 ||
			currentPath.CurrentPosition[1] >= len(hikingMap[0]) {
			continue
		}

		currentTerrain := hikingMap[currentPath.CurrentPosition[0]][currentPath.CurrentPosition[1]]
		if currentTerrain == '#' {
			continue // cannot travel thru the forest
		}

		if currentTerrain == '>' || currentTerrain == '.' {
			nextPosn := [2]int{currentPath.CurrentPosition[0], currentPath.CurrentPosition[1] + 1}
			if !currentPath.Visited[nextPosn] &&
				nextPosn[0] >= 0 &&
				nextPosn[0] < len(hikingMap) &&
				nextPosn[1] >= 0 &&
				nextPosn[1] < len(hikingMap[0]) &&
				hikingMap[nextPosn[0]][nextPosn[1]] != '#' {
				nextPath := cloneHikingPath(currentPath)
				nextPath.CurrentPosition = nextPosn
				nextPath.Visited[nextPosn] = true
				queue = append(queue, nextPath)
			}
		}
		if currentTerrain == 'v' || currentTerrain == '.' {
			nextPosn := [2]int{currentPath.CurrentPosition[0] + 1, currentPath.CurrentPosition[1]}
			if !currentPath.Visited[nextPosn] &&
				nextPosn[0] >= 0 &&
				nextPosn[0] < len(hikingMap) &&
				nextPosn[1] >= 0 &&
				nextPosn[1] < len(hikingMap[0]) &&
				hikingMap[nextPosn[0]][nextPosn[1]] != '#' {
				nextPath := cloneHikingPath(currentPath)
				nextPath.CurrentPosition = nextPosn
				nextPath.Visited[nextPosn] = true
				queue = append(queue, nextPath)
			}
		}
		if currentTerrain == '<' || currentTerrain == '.' {
			nextPosn := [2]int{currentPath.CurrentPosition[0], currentPath.CurrentPosition[1] - 1}
			if !currentPath.Visited[nextPosn] &&
				nextPosn[0] >= 0 &&
				nextPosn[0] < len(hikingMap) &&
				nextPosn[1] >= 0 &&
				nextPosn[1] < len(hikingMap[0]) &&
				hikingMap[nextPosn[0]][nextPosn[1]] != '#' {
				nextPath := cloneHikingPath(currentPath)
				nextPath.CurrentPosition = nextPosn
				nextPath.Visited[nextPosn] = true
				queue = append(queue, nextPath)
			}
		}
		if currentTerrain == '^' || currentTerrain == '.' {
			nextPosn := [2]int{currentPath.CurrentPosition[0] - 1, currentPath.CurrentPosition[1]}
			if !currentPath.Visited[nextPosn] &&
				nextPosn[0] >= 0 &&
				nextPosn[0] < len(hikingMap) &&
				nextPosn[1] >= 0 &&
				nextPosn[1] < len(hikingMap[0]) &&
				hikingMap[nextPosn[0]][nextPosn[1]] != '#' {
				nextPath := cloneHikingPath(currentPath)
				nextPath.CurrentPosition = nextPosn
				nextPath.Visited[nextPosn] = true
				queue = append(queue, nextPath)
			}
		}
	}

	printHikingPath(hikingMap, longestSoFarPath[endPosn])

	return longestSoFar[endPosn] - 1
}

func Day23() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: aoc 23 <input>")
		os.Exit(1)
	}

	text := ReadFileOrExit(os.Args[2])

	fmt.Println("Part 1:")
	if DEBUG {
		fmt.Println(text)
	}
	hikingMap := ParseAs2DMatrix(text)
	if DEBUG {
		fmt.Println(hikingMap)
	}
	fmt.Println(findLongestPath(hikingMap))

	fmt.Println("Part 2:")
}
