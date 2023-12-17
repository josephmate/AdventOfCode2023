package main

import (
	"container/heap"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
)

func parseCrucibleMap(input string) [][]int {
	var crucibleMap [][]int

	for r, line := range strings.Split(input, "\n") {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			var values []int
			for c, val := range trimmedLine {
				if val < '0' || val > '9' {
					fmt.Println("parseCrucibleMap hit unexpected value=", val, " at r=", r, "c=", c)
					panic("parseCrucibleMap hit unexpected value")
				}
				values = append(values, int(val-'0'))
			}
			crucibleMap = append(crucibleMap, values)
		}
	}

	return crucibleMap
}

type CruciblePositionWithCost struct {
	Row             int
	Col             int
	StraightCounter int
	Distance        int
	Cost            int
	Direction       int
}

type CruciblePosition struct {
	Row             int
	Col             int
	StraightCounter int
	Direction       int
}

func toCruciblePosition(cruciblePosition CruciblePositionWithCost) CruciblePosition {
	return CruciblePosition{
		Row:             cruciblePosition.Row,
		Col:             cruciblePosition.Col,
		StraightCounter: cruciblePosition.StraightCounter,
		Direction:       cruciblePosition.Direction,
	}
}

type PriorityQueue []CruciblePositionWithCost

func (h PriorityQueue) Len() int { return len(h) }

func (h PriorityQueue) Less(i, j int) bool {
	if h[i].Cost != h[j].Cost {
		return h[i].Cost < h[j].Cost
	}
	if h[i].Distance != h[j].Distance {
		return h[i].Distance < h[j].Distance
	}
	return h[i].StraightCounter < h[j].StraightCounter
}

func (h PriorityQueue) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *PriorityQueue) Push(x interface{}) {
	*h = append(*h, x.(CruciblePositionWithCost))
}

func (h *PriorityQueue) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func goStraight(crucibleMap [][]int, cruciblePositionWithCost CruciblePositionWithCost, targetPosn [2]int) (CruciblePositionWithCost, bool) {
	direction := cruciblePositionWithCost.Direction
	currentRow := cruciblePositionWithCost.Row
	currentCol := cruciblePositionWithCost.Col
	if direction == UP {
		newRow := currentRow - 1
		newCol := currentCol
		if newRow >= 0 {
			return CruciblePositionWithCost{
				Row:             newRow,
				Col:             newCol,
				StraightCounter: cruciblePositionWithCost.StraightCounter + 1,
				Distance:        ManhattanDistance([2]int{newRow, newCol}, targetPosn),
				Cost:            cruciblePositionWithCost.Cost + crucibleMap[newRow][newCol],
				Direction:       direction,
			}, true
		}
	} else if direction == DOWN {
		newRow := currentRow + 1
		newCol := currentCol
		if newRow < len(crucibleMap) {
			return CruciblePositionWithCost{
				Row:             newRow,
				Col:             newCol,
				StraightCounter: cruciblePositionWithCost.StraightCounter + 1,
				Distance:        ManhattanDistance([2]int{newRow, newCol}, targetPosn),
				Cost:            cruciblePositionWithCost.Cost + crucibleMap[newRow][newCol],
				Direction:       direction,
			}, true
		}
	} else if direction == LEFT {
		newRow := currentRow
		newCol := currentCol - 1
		if newCol >= 0 {
			return CruciblePositionWithCost{
				Row:             newRow,
				Col:             newCol,
				StraightCounter: cruciblePositionWithCost.StraightCounter + 1,
				Distance:        ManhattanDistance([2]int{newRow, newCol}, targetPosn),
				Cost:            cruciblePositionWithCost.Cost + crucibleMap[newRow][newCol],
				Direction:       direction,
			}, true
		}
	} else if direction == RIGHT {
		newRow := currentRow
		newCol := currentCol + 1
		if newCol < len(crucibleMap[0]) {
			return CruciblePositionWithCost{
				Row:             newRow,
				Col:             newCol,
				StraightCounter: cruciblePositionWithCost.StraightCounter + 1,
				Distance:        ManhattanDistance([2]int{newRow, newCol}, targetPosn),
				Cost:            cruciblePositionWithCost.Cost + crucibleMap[newRow][newCol],
				Direction:       direction,
			}, true
		}
	}
	return CruciblePositionWithCost{}, false
}

func goLeftRight(crucibleMap [][]int, cruciblePositionWithCost CruciblePositionWithCost, targetPosn [2]int) []CruciblePositionWithCost {
	direction := cruciblePositionWithCost.Direction
	currentRow := cruciblePositionWithCost.Row
	currentCol := cruciblePositionWithCost.Col
	if direction == UP || direction == DOWN {
		newRowA := currentRow
		newColA := currentCol - 1
		newRowB := currentRow
		newColB := currentCol + 1

		var result []CruciblePositionWithCost
		if newColA >= 0 && newColA < len(crucibleMap[0]) {
			result = append(result, CruciblePositionWithCost{
				Row:             newRowA,
				Col:             newColA,
				StraightCounter: 1,
				Distance:        ManhattanDistance([2]int{newRowA, newColA}, targetPosn),
				Cost:            cruciblePositionWithCost.Cost + crucibleMap[newRowA][newColA],
				Direction:       LEFT,
			})
		}
		if newColB >= 0 && newColB < len(crucibleMap[0]) {
			result = append(result, CruciblePositionWithCost{
				Row:             newRowB,
				Col:             newColB,
				StraightCounter: 1,
				Distance:        ManhattanDistance([2]int{newRowB, newColB}, targetPosn),
				Cost:            cruciblePositionWithCost.Cost + crucibleMap[newRowB][newColB],
				Direction:       RIGHT,
			})
		}
		return result
	} else { // LEFT/RIGHT
		newRowA := currentRow - 1
		newColA := currentCol
		newRowB := currentRow + 1
		newColB := currentCol

		var result []CruciblePositionWithCost
		if newRowA >= 0 && newRowA < len(crucibleMap) {
			result = append(result, CruciblePositionWithCost{
				Row:             newRowA,
				Col:             newColA,
				StraightCounter: 1,
				Distance:        ManhattanDistance([2]int{newRowA, newColA}, targetPosn),
				Cost:            cruciblePositionWithCost.Cost + crucibleMap[newRowA][newColA],
				Direction:       UP,
			})
		}
		if newRowB >= 0 && newRowB < len(crucibleMap) {
			result = append(result, CruciblePositionWithCost{
				Row:             newRowB,
				Col:             newColB,
				StraightCounter: 1,
				Distance:        ManhattanDistance([2]int{newRowB, newColB}, targetPosn),
				Cost:            cruciblePositionWithCost.Cost + crucibleMap[newRowB][newColB],
				Direction:       DOWN,
			})
		}
		return result
	}
}

func goLeftRightWith4(crucibleMap [][]int, cruciblePositionWithCost CruciblePositionWithCost, targetPosn [2]int) []CruciblePositionWithCost {
	direction := cruciblePositionWithCost.Direction
	currentRow := cruciblePositionWithCost.Row
	currentCol := cruciblePositionWithCost.Col
	if direction == UP || direction == DOWN {
		var newRowA = currentRow
		var newColA = currentCol
		var newCostA = cruciblePositionWithCost.Cost
		var newRowB = currentRow
		var newColB = currentCol
		var newCostB = cruciblePositionWithCost.Cost
		for i := 0; i < 4; i++ {

			newColA -= 1
			if newColA >= 0 && newColA < len(crucibleMap[0]) {
				newCostA += crucibleMap[newRowA][newColA]
			}

			newColB += 1
			if newColB >= 0 && newColB < len(crucibleMap[0]) {
				newCostB += crucibleMap[newRowB][newColB]
			}
		}

		var result []CruciblePositionWithCost
		if newColA >= 0 && newColA < len(crucibleMap[0]) {
			result = append(result, CruciblePositionWithCost{
				Row:             newRowA,
				Col:             newColA,
				StraightCounter: 4,
				Distance:        ManhattanDistance([2]int{newRowA, newColA}, targetPosn),
				Cost:            newCostA,
				Direction:       LEFT,
			})
		}
		if newColB >= 0 && newColB < len(crucibleMap[0]) {
			result = append(result, CruciblePositionWithCost{
				Row:             newRowB,
				Col:             newColB,
				StraightCounter: 4,
				Distance:        ManhattanDistance([2]int{newRowB, newColB}, targetPosn),
				Cost:            newCostB,
				Direction:       RIGHT,
			})
		}
		return result
	} else { // LEFT/RIGHT
		var newRowA = currentRow
		var newColA = currentCol
		var newCostA = cruciblePositionWithCost.Cost
		var newRowB = currentRow
		var newColB = currentCol
		var newCostB = cruciblePositionWithCost.Cost
		for i := 0; i < 4; i++ {

			newRowA -= 1
			if newRowA >= 0 && newRowA < len(crucibleMap) {
				newCostA += crucibleMap[newRowA][newColA]
			}

			newRowB += 1
			if newRowB >= 0 && newRowB < len(crucibleMap) {
				newCostB += crucibleMap[newRowB][newColB]
			}
		}

		var result []CruciblePositionWithCost
		if newRowA >= 0 && newRowA < len(crucibleMap) {
			result = append(result, CruciblePositionWithCost{
				Row:             newRowA,
				Col:             newColA,
				StraightCounter: 4,
				Distance:        ManhattanDistance([2]int{newRowA, newColA}, targetPosn),
				Cost:            newCostA,
				Direction:       UP,
			})
		}
		if newRowB >= 0 && newRowB < len(crucibleMap) {
			result = append(result, CruciblePositionWithCost{
				Row:             newRowB,
				Col:             newColB,
				StraightCounter: 4,
				Distance:        ManhattanDistance([2]int{newRowB, newColB}, targetPosn),
				Cost:            newCostB,
				Direction:       DOWN,
			})
		}
		return result
	}
}

func findCoolestPath(crucibleMap [][]int) int {
	// use priority queue to find shortest path
	// unfortunately, go doesn't have priority queue
	// so we need to use a heap
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	targetPosn := [2]int{len(crucibleMap) - 1, len(crucibleMap[0]) - 1}

	// going right as first step
	heap.Push(&pq, CruciblePositionWithCost{
		Row:             0,
		Col:             1,
		StraightCounter: 1,
		Distance:        ManhattanDistance([2]int{0, 1}, targetPosn),
		Cost:            crucibleMap[0][1],
		Direction:       RIGHT,
	})
	// going down as first step
	heap.Push(&pq, CruciblePositionWithCost{
		Row:             1,
		Col:             0,
		StraightCounter: 1,
		Distance:        ManhattanDistance([2]int{1, 0}, targetPosn),
		Cost:            crucibleMap[1][0],
		Direction:       DOWN,
	})

	// r, c, straightCounter, Direction
	minCostSoFar := map[CruciblePosition]int{}
	for len(pq) > 0 {
		currentPositionWithCost := heap.Pop(&pq).(CruciblePositionWithCost)
		if DEBUG {
			fmt.Println("findCoolestPath", currentPositionWithCost)
		}
		currentPosition := toCruciblePosition(currentPositionWithCost)
		// check if we already visited this position with a lower cost
		minCost := minCostSoFar[currentPosition]
		if (minCost > 0 && currentPositionWithCost.Cost >= minCost) ||
			currentPositionWithCost.Row < 0 ||
			currentPositionWithCost.Row >= len(crucibleMap) ||
			currentPositionWithCost.Col < 0 ||
			currentPositionWithCost.Col >= len(crucibleMap[0]) {
			// we already visited this scenario with less cost
			continue
		}

		minCostSoFar[currentPosition] = currentPositionWithCost.Cost

		if currentPositionWithCost.StraightCounter < 3 {
			nextMove, hasIt := goStraight(crucibleMap, currentPositionWithCost, targetPosn)
			if hasIt {
				heap.Push(&pq, nextMove)
			}
		}
		// turn left and right
		for _, nextMove := range goLeftRight(crucibleMap, currentPositionWithCost, targetPosn) {
			heap.Push(&pq, nextMove)
		}
	}

	var min = math.MaxInt32
	for straightCounter := 1; straightCounter <= 3; straightCounter++ {
		for direction := 0; direction <= 3; direction++ {
			val := minCostSoFar[CruciblePosition{
				Row:             targetPosn[0],
				Col:             targetPosn[1],
				StraightCounter: straightCounter,
				Direction:       direction,
			}]
			if val > 0 {
				min = Min(min, val)
			}
		}
	}

	return min
}

func findCoolestPathForUltraCrucible(crucibleMap [][]int) int {
	// use priority queue to find shortest path
	// unfortunately, go doesn't have priority queue
	// so we need to use a heap
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	targetPosn := [2]int{len(crucibleMap) - 1, len(crucibleMap[0]) - 1}

	// going right as first step
	heap.Push(&pq, CruciblePositionWithCost{
		Row:             0,
		Col:             4,
		StraightCounter: 4,
		Distance:        ManhattanDistance([2]int{0, 4}, targetPosn),
		Cost:            crucibleMap[0][1] + crucibleMap[0][2] + crucibleMap[0][3] + crucibleMap[0][4],
		Direction:       RIGHT,
	})
	// going down as first step
	heap.Push(&pq, CruciblePositionWithCost{
		Row:             4,
		Col:             0,
		StraightCounter: 4,
		Distance:        ManhattanDistance([2]int{4, 0}, targetPosn),
		Cost:            crucibleMap[1][0] + crucibleMap[2][0] + crucibleMap[3][0] + crucibleMap[4][0],
		Direction:       DOWN,
	})

	// r, c, straightCounter, Direction
	minCostSoFar := map[CruciblePosition]int{}
	for len(pq) > 0 {
		currentPositionWithCost := heap.Pop(&pq).(CruciblePositionWithCost)
		if DEBUG {
			fmt.Println("findCoolestPathForUltraCrucible", currentPositionWithCost)
		}
		currentPosition := toCruciblePosition(currentPositionWithCost)
		// check if we already visited this position with a lower cost
		minCost := minCostSoFar[currentPosition]
		if (minCost > 0 && currentPositionWithCost.Cost >= minCost) ||
			currentPositionWithCost.Row < 0 ||
			currentPositionWithCost.Row >= len(crucibleMap) ||
			currentPositionWithCost.Col < 0 ||
			currentPositionWithCost.Col >= len(crucibleMap[0]) {
			// we already visited this scenario with less cost
			continue
		}

		minCostSoFar[currentPosition] = currentPositionWithCost.Cost

		if currentPositionWithCost.StraightCounter < 10 {
			nextMove, hasIt := goStraight(crucibleMap, currentPositionWithCost, targetPosn)
			if hasIt {
				heap.Push(&pq, nextMove)
			}
		}
		// turn left and right
		for _, nextMove := range goLeftRightWith4(crucibleMap, currentPositionWithCost, targetPosn) {
			heap.Push(&pq, nextMove)
		}
	}

	var min = math.MaxInt32
	for straightCounter := 1; straightCounter <= 3; straightCounter++ {
		for direction := 0; direction <= 3; direction++ {
			val := minCostSoFar[CruciblePosition{
				Row:             targetPosn[0],
				Col:             targetPosn[1],
				StraightCounter: straightCounter,
				Direction:       direction,
			}]
			if val > 0 {
				min = Min(min, val)
			}
		}
	}

	return min
}

func Day17() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: aoc 13 <input>")
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
	crucibleMap := parseCrucibleMap(text)
	if DEBUG {
		fmt.Println(crucibleMap)
	}
	fmt.Println("Part 1:")
	fmt.Println(findCoolestPath(crucibleMap))
	fmt.Println("Part 2:")
	fmt.Println(findCoolestPathForUltraCrucible(crucibleMap))
}
