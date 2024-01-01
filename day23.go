package main

import (
	"fmt"
	"hash/fnv"
	"os"

	"github.com/bits-and-blooms/bitset"
)

type HikingPath struct {
	CurrentPosition [2]int
	HashOfPath      uint
	Visited    *bitset.BitSet
}

func mapPosnToId(hikingMap [][]byte) map[[2]int]uint {
	var id = uint(0)
	posnToId := map[[2]int]uint{}
	for r, row := range hikingMap {
		for c, col := range row {
			if col != '#' {
				posnToId[[2]int{r,c}] = id
				id++
			}
		}
	}
	return posnToId
}

func getId()

func findLongestPath(hikingMap [][]byte) uint {
	startPosn := [2]int{0, 1}
	numRows := len(hikingMap)
	numCols := len(hikingMap[0])
	endPosn := [2]int{numRows - 1, numCols - 2}

	posnToId := mapPosnToId(hikingMap)
	var queue []HikingPath
	startPath := HikingPath{}
	startPath.CurrentPosition = startPosn
	startPath.Visited = bitset.New(64)
	startPath.Visited.Set(posnToId[startPosn])
	queue = append(queue, startPath)

	longestSoFar := map[[2]int]uint{}

	for len(queue) > 0 {
		currentPath := queue[0]
		queue = queue[1:]

		if DEBUG {
			fmt.Println(currentPath)
		}

		// we already have a longer path thru this point
		if longestSoFar[currentPath.CurrentPosition] >= currentPath.Visited.Count() {
			continue
		}
		longestSoFar[currentPath.CurrentPosition] = currentPath.Visited.Count()

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
			nextHash := currentPath.HashOfPath + hash(hikingMap, nextPosn)
			if nextPosn[0] >= 0 &&
				nextPosn[0] < len(hikingMap) &&
				nextPosn[1] >= 0 &&
				nextPosn[1] < len(hikingMap[0]) &&
				hikingMap[nextPosn[0]][nextPosn[1]] != '#' &&
				!currentPath.Visited.Test(posnToId[nextPosn]) {
				newBitSet := currentPath.Visited.Clone()
				newBitSet.Set(posnToId[nextPosn])
				nextPath := HikingPath{
					CurrentPosition: nextPosn,
					HashOfPath: currentPath.HashOfPath + nextHash,
					Visited: newBitSet,
				}
				queue = append(queue, nextPath)
			}
		}
		if currentTerrain == 'v' || currentTerrain == '.' {
			nextPosn := [2]int{currentPath.CurrentPosition[0] + 1, currentPath.CurrentPosition[1]}
			nextHash := currentPath.HashOfPath + hash(hikingMap, nextPosn)
			if nextPosn[0] >= 0 &&
				nextPosn[0] < len(hikingMap) &&
				nextPosn[1] >= 0 &&
				nextPosn[1] < len(hikingMap[0]) &&
				hikingMap[nextPosn[0]][nextPosn[1]] != '#' &&
				!currentPath.Visited.Test(posnToId[nextPosn]) {
				newBitSet := currentPath.Visited.Clone()
				newBitSet.Set(posnToId[nextPosn])
				nextPath := HikingPath{
					CurrentPosition: nextPosn,
					HashOfPath: currentPath.HashOfPath + nextHash,
					Visited: newBitSet,
				}
				queue = append(queue, nextPath)
			}
		}
		if currentTerrain == '<' || currentTerrain == '.' {
			nextPosn := [2]int{currentPath.CurrentPosition[0], currentPath.CurrentPosition[1] - 1}
			nextHash := currentPath.HashOfPath + hash(hikingMap, nextPosn)
			if nextPosn[0] >= 0 &&
				nextPosn[0] < len(hikingMap) &&
				nextPosn[1] >= 0 &&
				nextPosn[1] < len(hikingMap[0]) &&
				hikingMap[nextPosn[0]][nextPosn[1]] != '#' &&
				!currentPath.Visited.Test(posnToId[nextPosn]){
				newBitSet := currentPath.Visited.Clone()
				newBitSet.Set(posnToId[nextPosn])
				nextPath := HikingPath{
					CurrentPosition: nextPosn,
					HashOfPath: currentPath.HashOfPath + nextHash,
					Visited: newBitSet,
				}
				queue = append(queue, nextPath)
			}
		}
		if currentTerrain == '^' || currentTerrain == '.' {
			nextPosn := [2]int{currentPath.CurrentPosition[0] - 1, currentPath.CurrentPosition[1]}
			nextHash := currentPath.HashOfPath + hash(hikingMap, nextPosn)
			if nextPosn[0] >= 0 &&
				nextPosn[0] < len(hikingMap) &&
				nextPosn[1] >= 0 &&
				nextPosn[1] < len(hikingMap[0]) &&
				hikingMap[nextPosn[0]][nextPosn[1]] != '#' &&
				!currentPath.Visited.Test(posnToId[nextPosn]){
				newBitSet := currentPath.Visited.Clone()
				newBitSet.Set(posnToId[nextPosn])
				nextPath := HikingPath{
					CurrentPosition: nextPosn,
					HashOfPath: currentPath.HashOfPath + nextHash,
					Visited: newBitSet,
				}
				queue = append(queue, nextPath)
			}
		}
	}

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

func countMovableSpaces(hikingMap [][]byte) int {
	var count = 0
	for _, row := range hikingMap {
		for _, col := range row {
			if col != '#' {
				count++
			}
		}
	}
	return count
}

type BitsetKey struct {
	BitSet *bitset.BitSet
}

func (bk BitsetKey) Equals(other BitsetKey) bool {
	return bk.BitSet.Equal(other.BitSet)
}

func (bk BitsetKey) Hash() int {
	h := fnv.New32a()
	bytes, _ := bk.BitSet.MarshalBinary()
	h.Write(bytes)
	return int(h.Sum32())
}

/*
at this point I realized that bloom doesn't make sense at because i can uniquely identify each position.
which means I could create a bitset which takes up less space.
*/
func findLongestPathIgnoreSlopes(hikingMap [][]byte) uint {
	startPosn := [2]int{0, 1}
	numRows := len(hikingMap)
	numCols := len(hikingMap[0])
	endPosn := [2]int{numRows - 1, numCols - 2}
	
	var queue []HikingPath
	startPath := HikingPath{}
	startPath.CurrentPosition = startPosn
	startPath.Visited = bitset.New(64)
	posnToId := mapPosnToId(hikingMap)
	startPath.Visited.Set(posnToId[startPosn])
	startPath.HashOfPath = hash(hikingMap, startPosn)
	queue = append(queue, startPath)

	longestSoFar := map[[2]int]uint{}
	// try to deduplicate paths using the hash of the path so far
	dedupePaths := map[uint]bool{}

	for len(queue) > 0 {
		currentPath := queue[0]
		queue = queue[1:]

		if DEBUG {
			fmt.Println(currentPath)
		}

		if longestSoFar[currentPath.CurrentPosition] > currentPath.Visited.Count() {
			longestSoFar[currentPath.CurrentPosition] = currentPath.Visited.Count()
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
		if !dedupePaths[nextHash] &&
			nextPosn[0] >= 0 &&
			nextPosn[0] < len(hikingMap) &&
			nextPosn[1] >= 0 &&
			nextPosn[1] < len(hikingMap[0]) &&
			hikingMap[nextPosn[0]][nextPosn[1]] != '#' &&
			!currentPath.Visited.Test(posnToId[nextPosn]) {
				newBitSet := currentPath.Visited.Clone()
				newBitSet.Set(posnToId[nextPosn])
				nextPath := HikingPath{
					CurrentPosition: nextPosn,
					HashOfPath: currentPath.HashOfPath + nextHash,
					Visited: newBitSet,
				}
				queue = append(queue, nextPath)
		}

		nextPosn = [2]int{currentPath.CurrentPosition[0] + 1, currentPath.CurrentPosition[1]}
		nextHash = currentPath.HashOfPath + hash(hikingMap, nextPosn)
		if !dedupePaths[nextHash] &&
			nextPosn[0] >= 0 &&
			nextPosn[0] < len(hikingMap) &&
			nextPosn[1] >= 0 &&
			nextPosn[1] < len(hikingMap[0]) &&
			hikingMap[nextPosn[0]][nextPosn[1]] != '#' &&
			!currentPath.Visited.Test(posnToId[nextPosn]) {
				newBitSet := currentPath.Visited.Clone()
				newBitSet.Set(posnToId[nextPosn])
				nextPath := HikingPath{
					CurrentPosition: nextPosn,
					HashOfPath: currentPath.HashOfPath + nextHash,
					Visited: newBitSet,
				}
				queue = append(queue, nextPath)
		}

		nextPosn = [2]int{currentPath.CurrentPosition[0], currentPath.CurrentPosition[1] - 1}
		nextHash = currentPath.HashOfPath + hash(hikingMap, nextPosn)
		if !dedupePaths[nextHash] &&
			nextPosn[0] >= 0 &&
			nextPosn[0] < len(hikingMap) &&
			nextPosn[1] >= 0 &&
			nextPosn[1] < len(hikingMap[0]) &&
			hikingMap[nextPosn[0]][nextPosn[1]] != '#'  &&
			!currentPath.Visited.Test(posnToId[nextPosn]){
				newBitSet := currentPath.Visited.Clone()
				newBitSet.Set(posnToId[nextPosn])
				nextPath := HikingPath{
					CurrentPosition: nextPosn,
					HashOfPath: currentPath.HashOfPath + nextHash,
					Visited: newBitSet,
				}
				queue = append(queue, nextPath)
		}

		nextPosn = [2]int{currentPath.CurrentPosition[0] - 1, currentPath.CurrentPosition[1]}
		nextHash = currentPath.HashOfPath + hash(hikingMap, nextPosn)
		if !dedupePaths[nextHash] &&
			nextPosn[0] >= 0 &&
			nextPosn[0] < len(hikingMap) &&
			nextPosn[1] >= 0 &&
			nextPosn[1] < len(hikingMap[0]) &&
			hikingMap[nextPosn[0]][nextPosn[1]] != '#' &&
			!currentPath.Visited.Test(posnToId[nextPosn]) {
				newBitSet := currentPath.Visited.Clone()
				newBitSet.Set(posnToId[nextPosn])
				nextPath := HikingPath{
					CurrentPosition: nextPosn,
					HashOfPath: currentPath.HashOfPath + nextHash,
					Visited: newBitSet,
				}
				queue = append(queue, nextPath)
		}
	}

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
