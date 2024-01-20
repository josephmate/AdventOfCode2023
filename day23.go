package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
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

type CompressedHikingPath struct {
	CurrentPosition uint
	HashOfPath      uint64
	CostSoFar       uint
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
			//fmt.Println(currentPath)
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
	return hashId(x)
}

func hashId(x uint) uint {
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

func generateConnectedIds(hikingMap [][]byte, posnToId map[[2]int]uint, posn [2]int) []uint {
	var result []uint
	var nextRow, nextCol int

	nextRow = posn[0] + 1
	nextCol = posn[1]
	if nextRow >= 0 &&
		nextRow < len(hikingMap) && 
		nextCol >= 0 &&
		nextCol < len(hikingMap[0]) &&
		hikingMap[nextRow][nextCol] != '#' {
		result = append(result, posnToId[[2]int{nextRow, nextCol}])
	}
	nextRow = posn[0] - 1
	nextCol = posn[1]
	if nextRow >= 0 &&
		nextRow < len(hikingMap) && 
		nextCol >= 0 &&
		nextCol < len(hikingMap[0]) &&
		hikingMap[nextRow][nextCol] != '#' {
		result = append(result, posnToId[[2]int{nextRow, nextCol}])
	}
	nextRow = posn[0]
	nextCol = posn[1] + 1
	if nextRow >= 0 &&
		nextRow < len(hikingMap) && 
		nextCol >= 0 &&
		nextCol < len(hikingMap[0]) &&
		hikingMap[nextRow][nextCol] != '#' {
		result = append(result, posnToId[[2]int{nextRow, nextCol}])
	}
	nextRow = posn[0]
	nextCol = posn[1] - 1
	if nextRow >= 0 &&
		nextRow < len(hikingMap) && 
		nextCol >= 0 &&
		nextCol < len(hikingMap[0]) &&
		hikingMap[nextRow][nextCol] != '#' {
		result = append(result, posnToId[[2]int{nextRow, nextCol}])
	}
	return result
}

func printAsGraphviz(hikingMap [][]byte) {
	if !DEBUG {
		return
	}

	posnToId := mapPosnToId(hikingMap)

	edges := map[[2]uint]bool{}
	for posn := range posnToId {
		startId := posnToId[posn]
		for _, connectedId := range generateConnectedIds(hikingMap, posnToId, posn) {
			var low = startId
			var high = connectedId
			if high < low {
				var tmp = low
				low = high
				high = tmp
			}
			edges[[2]uint{low, high}] = true
		}
	}

	file, err := os.Create("dot.dot")
	if err != nil {
		fmt.Println("Could not open dot.dot")
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	_, err = fmt.Fprintln(writer, "graph Connected_Components {")
	if err != nil {
		fmt.Println("Could not write to dot.dot")
		return
	}

	for edge := range edges {
		_, err := fmt.Fprintf(writer, "    %d -- %d [tooltip=\"%d<->%d\"]\n", edge[0], edge[1], edge[0], edge[1])
		if err != nil {
			fmt.Println("Could not write to dot.dot")
			return
		}
	}

	_, err = fmt.Fprintln(writer, "}")
	if err != nil {
		fmt.Println("Could not write to dot.dot")
		return
	}
}

func printCompressedAsGraphviz(nodes  map[uint][]HikingMapEdge) {
	if !DEBUG {
		return
	}

	file, err := os.Create("dot_compressed.dot")
	if err != nil {
		fmt.Println("Could not open dot.dot")
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	_, err = fmt.Fprintln(writer, "graph Connected_Components {")
	if err != nil {
		fmt.Println("Could not write to dot.dot")
		return
	}

	for nodeId := range nodes {
		for _, edge := range nodes[nodeId] {
			if nodeId < edge.Destination {
				_, err := fmt.Fprintf(writer, "    %d -- %d [label=\"%d\"]\n", nodeId, edge.Destination, edge.Cost)
				if err != nil {
					fmt.Println("Could not write to dot.dot")
					return
				}
			}
		}
	}

	_, err = fmt.Fprintln(writer, "}")
	if err != nil {
		fmt.Println("Could not write to dot.dot")
		return
	}
}


type HikingMapEdge struct {
	Destination uint
	Cost uint
}


func hikingMapToGraph(hikingMap [][]byte) (map[uint][]HikingMapEdge, map[[2]int]uint) {
	result :=  map[uint][]HikingMapEdge{}
	
	posnToId := mapPosnToId(hikingMap)
	visited := map[[2]uint]bool{}
	for posn := range posnToId {
		startId := posnToId[posn]
		for _, connectedId := range generateConnectedIds(hikingMap, posnToId, posn) {

			
			var low = startId
			var high = connectedId
			if high < low {
				low = connectedId
				high = startId
			}
			visitedPair := [2]uint{low, high}
			if visited[visitedPair] {
				continue
			}

			result[startId] = append(result[startId], HikingMapEdge{
				Destination: connectedId,
				Cost: 1,
			})
			result[connectedId] = append(result[connectedId], HikingMapEdge{
				Destination: startId,
				Cost: 1,
			})

			visited[visitedPair] = true
		}
	}
	return result, posnToId
}

type EdgeQueueEntry struct {
	Current uint
	Prev uint
}

func removeEdge(edges []HikingMapEdge, id uint) []HikingMapEdge {
	var result []HikingMapEdge
	for _, edge := range edges {
		if edge.Destination != id {
			result = append(result, edge)
		}
	}
	return result
}

func compressGraph(hikingMapAsGraph map[uint][]HikingMapEdge) map[uint][]HikingMapEdge {

	// BFS looking for edges of degree two
	visited := map[uint]bool{}
	var queue []EdgeQueueEntry
	queue = append(queue, EdgeQueueEntry{0, 0})

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current.Current] {
			continue
		}
		visited[current.Current] = true

		neighbours := hikingMapAsGraph[current.Current]
		if len(neighbours) == 2 {
			var next = neighbours[0]
			var prev = neighbours[1]
			if next.Destination == current.Prev {
				next = neighbours[1]
				prev = neighbours[0]
			}

			// calculate new weight
			newWeight := prev.Cost + 1

			// remove current from prev
			hikingMapAsGraph[prev.Destination] = removeEdge(hikingMapAsGraph[prev.Destination], current.Current)

			// remove current from next
			hikingMapAsGraph[next.Destination] = removeEdge(hikingMapAsGraph[next.Destination], current.Current)

			// delete current from map
			delete(hikingMapAsGraph, current.Current)

			// connect prev to next with new cost
			hikingMapAsGraph[prev.Destination] = append(hikingMapAsGraph[prev.Destination], HikingMapEdge{
				Destination: next.Destination,
				Cost: newWeight,
			})

			// connect next to prev with new cost
			hikingMapAsGraph[next.Destination] = append(hikingMapAsGraph[next.Destination], HikingMapEdge{
				Destination: prev.Destination,
				Cost: newWeight,
			})
		}

		for _, neighbourId := range neighbours {
			queue = append(queue, EdgeQueueEntry{neighbourId.Destination, current.Current})
		}
	}

	return hikingMapAsGraph
}


func compressGraphIds(hikingMapAsGraph map[uint][]HikingMapEdge) (map[uint][]HikingMapEdge, map[uint]uint) {
	var id uint = 1
	oldIdToNewId := map[uint]uint{}
	oldIdToNewId[0] = 0
	for nodeId := range hikingMapAsGraph {
		oldIdToNewId[nodeId] = id
		id++
	}
	
	compressedGraph := map[uint][]HikingMapEdge{}
	for nodeId := range hikingMapAsGraph {
		var newEdges []HikingMapEdge
		for _, edge := range hikingMapAsGraph[nodeId] {
			newEdges = append(newEdges, HikingMapEdge{
				Destination: oldIdToNewId[edge.Destination],
				Cost: edge.Cost,
			})
		}
		compressedGraph[oldIdToNewId[nodeId]] = newEdges
	}

	return compressedGraph, oldIdToNewId
}

func bitSetToUint64(bitset *bitset.BitSet) uint64 {
	var buf bytes.Buffer
	_, err := bitset.WriteTo(&buf)
	if err != nil {
		fmt.Println("Failed to write to buffer", err)
		os.Exit(-1)
	}
	var value uint64
	err2 := binary.Read(&buf, binary.LittleEndian, &value)
	if err2 != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	return value
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
	
	hikingMapAsGraph, posnToId := hikingMapToGraph(hikingMap)
	if DEBUG {
		fmt.Println("findLongestPathIgnoreSlopes", "hikingMapToGraph", hikingMapAsGraph)
	}
	compressedGraph, idToCompressedId := compressGraphIds(compressGraph(hikingMapAsGraph))
	if DEBUG {
		fmt.Println("findLongestPathIgnoreSlopes", "compressGraph", compressedGraph)
	}
	printCompressedAsGraphviz(compressedGraph)


	var queue []CompressedHikingPath
	startPath := CompressedHikingPath{}
	startCompressedId := idToCompressedId[posnToId[startPosn]]
	startPath.CurrentPosition = startCompressedId
	startPath.Visited = bitset.New(64) // will never exceed 64 bits (1 long)
	startPath.Visited.Set(posnToId[startPosn])
	startPath.HashOfPath = bitSetToUint64(startPath.Visited)
	startPath.CostSoFar = 0
	queue = append(queue, startPath)

	longestSoFar := map[uint]uint{}
	// try to deduplicate paths using the hash of the path so far
	dedupePaths := map[uint]bool{}

	for len(queue) > 0 {
		currentPath := queue[0]
		queue = queue[1:]

		if DEBUG {
			//fmt.Println(currentPath)
		}

		if dedupePaths[currentPath.HashOfPath] {
			continue
		}
		dedupePaths[currentPath.HashOfPath] = true
		if longestSoFar[currentPath.CurrentPosition] < currentPath.CostSoFar {
			longestSoFar[currentPath.CurrentPosition] = currentPath.CostSoFar
		}

		for _, edge := range compressedGraph[currentPath.CurrentPosition] {
			nextHash := currentPath.HashOfPath + hashId(edge.Destination)
			if !dedupePaths[nextHash] &&
				!currentPath.Visited.Test(edge.Destination) {
					newBitSet := currentPath.Visited.Clone()
					newBitSet.Set(edge.Destination)
					nextPath := CompressedHikingPath{
						CurrentPosition: edge.Destination,
						HashOfPath: nextHash,
						CostSoFar: currentPath.CostSoFar + edge.Cost,
						Visited: newBitSet,
					}
					queue = append(queue, nextPath)
			}
		}
	}

	endId := idToCompressedId[posnToId[endPosn]]
	return longestSoFar[endId] - 1
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
	printAsGraphviz(hikingMap)
	if DEBUG {
		fmt.Println(hikingMap)
	}
	fmt.Println(findLongestPath(hikingMap))

	fmt.Println("Part 2:")
	fmt.Println(findLongestPathIgnoreSlopes(hikingMap))
}
