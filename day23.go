package main

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/bits-and-blooms/bloom/v3"
)

type HikingPath struct {
	CurrentPosition [2]int
	Visited         map[[2]int]bool
	HashOfPath      uint
	VisitedBloom    *bloom.BloomFilter
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

/*
https://stackoverflow.com/questions/28326965/good-hash-function-for-list-of-integers-where-order-doesnt-change-value
*/
func hash(hikingMap [][]byte, posn [2]int) uint {
	var x = uint(posn[0]*len(hikingMap[0]) + posn[1])
	x ^= x >> 17;
	x *= uint(0xed5ad4bb);
	x ^= x >> 11;
	x *= uint(0xac4c1b51);
	x ^= x >> 15;
	x *= uint(0x31848bab);
	x ^= x >> 14;
	return x
}

func addToBloom(hikingMap [][]byte, hikingPath HikingPath, posn [2]int) {
	i := uint32(posn[0] * len(hikingMap[0]) + posn[1])
	n1 := make([]byte, 4)
	binary.BigEndian.PutUint32(n1, i)
	hikingPath.VisitedBloom.Add(n1)
}

func containsBloom(hikingMap [][]byte, hikingPath HikingPath, posn [2]int) bool {
	i := uint32(posn[0] * len(hikingMap[0]) + posn[1])
	n1 := make([]byte, 4)
	binary.BigEndian.PutUint32(n1, i)
	return hikingPath.VisitedBloom.Test(n1)
}

func findLongestPathIgnoreSlopes(hikingMap [][]byte) int {
	startPosn := [2]int{0, 1}
	numRows := len(hikingMap)
	numCols := len(hikingMap[0])
	endPosn := [2]int{numRows - 1, numCols - 2}
	var queue []HikingPath
	startPath := HikingPath{}
	startPath.CurrentPosition = startPosn
	startPath.VisitedBloom = bloom.NewWithEstimates(uint(len(hikingMap))*uint(len(hikingMap[0])), 0.0001)
	m, k := bloom.EstimateParameters(uint(len(hikingMap))*uint(len(hikingMap[0])), 0.0001)
	addToBloom(hikingMap, startPath, startPosn)
	fmt.Println("estimate params m=", m, "k=", k)
	startPath.HashOfPath = hash(hikingMap, startPosn)
	queue = append(queue, startPath)

	longestSoFar := map[[2]int]int{}
	longestSoFarPath := map[[2]int]HikingPath{}
	// try to deduplicate paths using the hash of the path so far
	dedupePaths := map[uint]bool{}

	for len(queue) > 0 {
		currentPath := queue[0]
		queue = queue[1:]

		if DEBUG {
			fmt.Println(currentPath)
		}

		// we already have a longer path thru this point
		// if longestSoFar[currentPath.CurrentPosition] >= len(currentPath.Visited) {
		// 	continue
		// }
		// longestSoFar[currentPath.CurrentPosition] = len(currentPath.Visited)
		// longestSoFarPath[currentPath.CurrentPosition] = currentPath
		if len(currentPath.Visited) > longestSoFar[currentPath.CurrentPosition] {
			longestSoFar[currentPath.CurrentPosition] = len(currentPath.Visited)
			longestSoFarPath[currentPath.CurrentPosition] = currentPath
		}

		dedupePaths[currentPath.HashOfPath] = true

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

		var nextPosn = [2]int{currentPath.CurrentPosition[0], currentPath.CurrentPosition[1] + 1}
		var nextHash = currentPath.HashOfPath + hash(hikingMap, nextPosn)
		if !containsBloom(hikingMap, currentPath, nextPosn) &&
			!dedupePaths[nextHash] &&
			nextPosn[0] >= 0 &&
			nextPosn[0] < len(hikingMap) &&
			nextPosn[1] >= 0 &&
			nextPosn[1] < len(hikingMap[0]) &&
			hikingMap[nextPosn[0]][nextPosn[1]] != '#' {
			nextPath := HikingPath{}
			nextPath.VisitedBloom = currentPath.VisitedBloom.Copy()
			addToBloom(hikingMap, nextPath, nextPosn)
			nextPath.CurrentPosition = nextPosn
			nextPath.HashOfPath = nextHash
			queue = append(queue, nextPath)
		}

		nextPosn = [2]int{currentPath.CurrentPosition[0] + 1, currentPath.CurrentPosition[1]}
		nextHash = currentPath.HashOfPath + hash(hikingMap, nextPosn)
		if !containsBloom(hikingMap, currentPath, nextPosn) &&
			!dedupePaths[nextHash] &&
			nextPosn[0] >= 0 &&
			nextPosn[0] < len(hikingMap) &&
			nextPosn[1] >= 0 &&
			nextPosn[1] < len(hikingMap[0]) &&
			hikingMap[nextPosn[0]][nextPosn[1]] != '#' {
			nextPath := HikingPath{}
			nextPath.VisitedBloom = currentPath.VisitedBloom.Copy()
			addToBloom(hikingMap, nextPath, nextPosn)
			nextPath.CurrentPosition = nextPosn
			nextPath.HashOfPath = nextHash
			queue = append(queue, nextPath)
		}

		nextPosn = [2]int{currentPath.CurrentPosition[0], currentPath.CurrentPosition[1] - 1}
		nextHash = currentPath.HashOfPath + hash(hikingMap, nextPosn)
		if !containsBloom(hikingMap, currentPath, nextPosn) &&
			!dedupePaths[nextHash] &&
			nextPosn[0] >= 0 &&
			nextPosn[0] < len(hikingMap) &&
			nextPosn[1] >= 0 &&
			nextPosn[1] < len(hikingMap[0]) &&
			hikingMap[nextPosn[0]][nextPosn[1]] != '#' {
			nextPath := HikingPath{}
			nextPath.VisitedBloom = currentPath.VisitedBloom.Copy()
			addToBloom(hikingMap, nextPath, nextPosn)
			nextPath.CurrentPosition = nextPosn
			nextPath.HashOfPath = nextHash
			queue = append(queue, nextPath)
		}

		nextPosn = [2]int{currentPath.CurrentPosition[0] - 1, currentPath.CurrentPosition[1]}
		nextHash = currentPath.HashOfPath + hash(hikingMap, nextPosn)
		if !containsBloom(hikingMap, currentPath, nextPosn) &&
			!dedupePaths[nextHash] &&
			nextPosn[0] >= 0 &&
			nextPosn[0] < len(hikingMap) &&
			nextPosn[1] >= 0 &&
			nextPosn[1] < len(hikingMap[0]) &&
			hikingMap[nextPosn[0]][nextPosn[1]] != '#' {
			nextPath := HikingPath{}
			nextPath.VisitedBloom = currentPath.VisitedBloom.Copy()
			addToBloom(hikingMap, nextPath, nextPosn)
			nextPath.CurrentPosition = nextPosn
			nextPath.HashOfPath = nextHash
			queue = append(queue, nextPath)
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
	fmt.Println(findLongestPathIgnoreSlopes(hikingMap))
}
