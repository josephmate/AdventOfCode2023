package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func parseDay10(input string) []string {
	lines := strings.Split(input, "\n")

	var result []string

	for _, line := range lines {
		trimmedStr := strings.TrimSpace(line)
		if len(trimmedStr) > 0 {
			result = append(result, trimmedStr)
		}
	}

	return result
}

func findStart(maze []string) (int, int) {
	for r, row := range maze {
		for c, col := range row {
			if col == 'S' {
				return r, c
			}
		}
	}
	return -1, -1
}

const NORTH = 0
const EAST = 1
const SOUTH = 2
const WEST = 3

func expandMoves(maze []string, r int, c int) [][3]int {
	var moves [][3]int

	if r-1 >= 0 {
		moves = append(moves, [3]int{r - 1, c, NORTH})
	}
	if r+1 < len(maze) {
		moves = append(moves, [3]int{r + 1, c, SOUTH})
	}
	if c-1 >= 0 {
		moves = append(moves, [3]int{r, c - 1, WEST})
	}
	if c+1 < len(maze[r]) {
		moves = append(moves, [3]int{r, c + 1, EAST})
	}

	return moves
}

/*
| is a vertical pipe connecting north and south.
- is a horizontal pipe connecting east and west.
L is a 90-degree bend connecting north and east.
J is a 90-degree bend connecting north and west.
7 is a 90-degree bend connecting south and west.
F is a 90-degree bend connecting south and east.
. is ground; there is no pipe in this tile.
S is the starting position
*/

func findFurthest(maze []string) int {
	startR, startC := findStart(maze)
	if DEBUG {
		fmt.Println("findFurthest", "startR", startR, "startC", startC)
	}

	var longestSoFar = 0
	visited := map[[2]int]bool{}
	visited[[2]int{startR, startC}] = true
	var queue [][4]int
	for _, startingMove := range expandMoves(maze, startR, startC) {
		queue = append(queue, [4]int{startingMove[0], startingMove[1], startingMove[2], 1})
	}

	for len(queue) > 0 {
		if DEBUG {
			fmt.Println("findFurthest", "longestSoFar", longestSoFar, "queue", queue)
		}
		nextEntry := queue[0]
		queue = queue[1:]
		r := nextEntry[0]
		c := nextEntry[1]
		direction := nextEntry[2]
		distance := nextEntry[3]

		if visited[[2]int{r, c}] {
			continue
		}

		pipe := maze[r][c]
		/*
		   | is a vertical pipe connecting north and south.
		   - is a horizontal pipe connecting east and west.
		   L is a 90-degree bend connecting north and east.
		   J is a 90-degree bend connecting north and west.
		   7 is a 90-degree bend connecting south and west.
		   F is a 90-degree bend connecting south and east.
		   . is ground; there is no pipe in this tile.
		   S is the starting position
		*/
		if pipe == '.' || pipe == 'S' {
			continue
		} else if pipe == '|' && direction == NORTH {
			// | is a vertical pipe connecting north and south.
			if distance > longestSoFar {
				longestSoFar = distance
			}
			visited[[2]int{r, c}] = true
			if r-1 >= 0 {
				queue = append(queue, [4]int{r - 1, c, NORTH, distance + 1})
			}
		} else if pipe == '|' && direction == SOUTH {
			if distance > longestSoFar {
				longestSoFar = distance
			}
			visited[[2]int{r, c}] = true
			if r+1 < len(maze) {
				queue = append(queue, [4]int{r + 1, c, SOUTH, distance + 1})
			}
		} else if pipe == '-' && direction == EAST {
			// - is a horizontal pipe connecting east and west.
			if distance > longestSoFar {
				longestSoFar = distance
			}
			visited[[2]int{r, c}] = true
			if c+1 < len(maze[r]) {
				queue = append(queue, [4]int{r, c + 1, EAST, distance + 1})
			}
		} else if pipe == '-' && direction == WEST {
			if distance > longestSoFar {
				longestSoFar = distance
			}
			visited[[2]int{r, c}] = true
			if c-1 >= 0 {
				queue = append(queue, [4]int{r, c - 1, WEST, distance + 1})
			}
		} else if pipe == 'L' && direction == SOUTH {
			// L is a 90-degree bend connecting north and east.
			if distance > longestSoFar {
				longestSoFar = distance
			}
			visited[[2]int{r, c}] = true
			if c+1 < len(maze[r]) {
				queue = append(queue, [4]int{r, c + 1, EAST, distance + 1})
			}
		} else if pipe == 'L' && direction == WEST {
			// L is a 90-degree bend connecting north and east.
			if distance > longestSoFar {
				longestSoFar = distance
			}
			visited[[2]int{r, c}] = true
			if r-1 >= 0 {
				queue = append(queue, [4]int{r - 1, c, NORTH, distance + 1})
			}
		} else if pipe == 'J' && direction == EAST {
			// J is a 90-degree bend connecting north and west.
			if distance > longestSoFar {
				longestSoFar = distance
			}
			visited[[2]int{r, c}] = true
			if r-1 >= 0 {
				queue = append(queue, [4]int{r - 1, c, NORTH, distance + 1})
			}
		} else if pipe == 'J' && direction == SOUTH {
			// J is a 90-degree bend connecting north and west.
			if distance > longestSoFar {
				longestSoFar = distance
			}
			visited[[2]int{r, c}] = true
			if c-1 >= 0 {
				queue = append(queue, [4]int{r, c - 1, WEST, distance + 1})
			}
		} else if pipe == '7' && direction == NORTH {
			// 7 is a 90-degree bend connecting south and west.
			if distance > longestSoFar {
				longestSoFar = distance
			}
			visited[[2]int{r, c}] = true
			if c-1 >= 0 {
				queue = append(queue, [4]int{r, c - 1, WEST, distance + 1})
			}
		} else if pipe == '7' && direction == EAST {
			// 7 is a 90-degree bend connecting south and west.
			if distance > longestSoFar {
				longestSoFar = distance
			}
			visited[[2]int{r, c}] = true
			if r+1 < len(maze) {
				queue = append(queue, [4]int{r + 1, c, SOUTH, distance + 1})
			}
		} else if pipe == 'F' && direction == NORTH {
			// F is a 90-degree bend connecting south and east.
			if distance > longestSoFar {
				longestSoFar = distance
			}
			visited[[2]int{r, c}] = true
			if c+1 < len(maze[r]) {
				queue = append(queue, [4]int{r, c + 1, EAST, distance + 1})
			}
		} else if pipe == 'F' && direction == WEST {
			// F is a 90-degree bend connecting south and east.
			if distance > longestSoFar {
				longestSoFar = distance
			}
			visited[[2]int{r, c}] = true
			if r+1 < len(maze) {
				queue = append(queue, [4]int{r + 1, c, SOUTH, distance + 1})
			}
		} else {
			//fmt.Println("findFurthest unexpected pipe ", string(pipe), " at r=", r, "c=", c)
			//return 0
		}

	}

	return longestSoFar
}

func Day10() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: aoc 10 <part 1 input>")
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
	maze := parseDay10(text)
	if DEBUG {
		fmt.Println(maze)
	}
	fmt.Println("Part 1:")
	fmt.Println(findFurthest(maze))
	fmt.Println("Part 2:")
}
