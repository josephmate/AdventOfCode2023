package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func getNextPosition(beam [3]int) [3]int {
	r := beam[0]
	c := beam[1]
	direction := beam[2]
	if direction == UP {
		return [3]int{r - 1, c, direction}
	}
	if direction == DOWN {
		return [3]int{r + 1, c, direction}
	}
	if direction == LEFT {
		return [3]int{r, c - 1, direction}
	}
	if direction == RIGHT {
		return [3]int{r, c + 1, direction}
	}
	fmt.Println("getNextPosition encountered unexpected direction=", direction, "r=", r, "c=", c)
	panic("getNextPosition encountered unexpected direction")
}

func getNextPositionByParams(r int, c int, direction int) [3]int {
	return getNextPosition([3]int{r, c, direction})
}

func parseMirrorMap(input string) [][]byte {
	var byteLines [][]byte

	for _, line := range strings.Split(input, "\n") {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			byteLines = append(byteLines, []byte(trimmedLine))
		}
	}

	return byteLines
}

func simulateBeams(mirrorMap [][]byte, startBeam [3]int) int {
	// r, c, direction
	var beamQueue [][3]int
	engergized := map[[2]int]bool{}
	visited := map[[3]int]bool{}
	// The beam enters in the top-left corner from the left and heading to the right.
	beamQueue = append(beamQueue, startBeam)

	for len(beamQueue) > 0 {
		currentBeam := beamQueue[0]
		beamQueue = beamQueue[1:]

		if visited[currentBeam] {
			continue
		}
		r := currentBeam[0]
		c := currentBeam[1]

		// out of bounds
		if r < 0 || r >= len(mirrorMap) {
			continue
		}
		if c < 0 || c >= len(mirrorMap[r]) {
			continue
		}

		direction := currentBeam[2]
		encounter := mirrorMap[r][c]
		visited[currentBeam] = true
		engergized[[2]int{r, c}] = true

		if encounter == '.' {
			// If the beam encounters empty space (.), it continues in the same direction.
			beamQueue = append(beamQueue, getNextPosition(currentBeam))
		} else if encounter == '/' {
			// If the beam encounters a mirror  (/ or \), the beam is reflected 90 degrees
			// depending  on the  angle of the  mirror. For  instance, a  rightward-moving
			// beam that encounters a / mirror would continue upward in the mirror's
			// column, while a rightward-moving beam that encounters a \ mirror would
			// continue downward from the mirror's column.
			if direction == UP {
				beamQueue = append(beamQueue, getNextPositionByParams(r, c, RIGHT))
			} else if direction == DOWN {
				beamQueue = append(beamQueue, getNextPositionByParams(r, c, LEFT))
			} else if direction == LEFT {
				beamQueue = append(beamQueue, getNextPositionByParams(r, c, DOWN))
			} else if direction == RIGHT {
				beamQueue = append(beamQueue, getNextPositionByParams(r, c, UP))
			} else {
				fmt.Println("simulateBeams encountered unexpected direction=", direction, "r=", r, "c=", c)
				panic("simulateBeams encountered unexpected direction")
			}
		} else if encounter == '\\' {
			// If the beam encounters a mirror  (/ or \), the beam is reflected 90 degrees
			// depending  on the  angle of the  mirror. For  instance, a  rightward-moving
			// beam that encounters a / mirror would continue upward in the mirror's
			// column, while a rightward-moving beam that encounters a \ mirror would
			// continue downward from the mirror's column.
			if direction == UP {
				beamQueue = append(beamQueue, getNextPositionByParams(r, c, LEFT))
			} else if direction == DOWN {
				beamQueue = append(beamQueue, getNextPositionByParams(r, c, RIGHT))
			} else if direction == LEFT {
				beamQueue = append(beamQueue, getNextPositionByParams(r, c, UP))
			} else if direction == RIGHT {
				beamQueue = append(beamQueue, getNextPositionByParams(r, c, DOWN))
			} else {
				fmt.Println("simulateBeams encountered unexpected direction=", direction, "r=", r, "c=", c)
				panic("simulateBeams encountered unexpected direction")
			}
		} else if encounter == '|' {
			if direction == UP || direction == DOWN {
				// If the beam encounters the pointy end of a splitter (| or -), the beam
				// passes through the splitter as if the splitter were empty space. For
				// instance, a rightward-moving beam that encounters a - splitter would
				// continue in the same direction.
				beamQueue = append(beamQueue, getNextPosition(currentBeam))
			} else if direction == LEFT || direction == RIGHT {
				// If the beam encounters the flat side of a splitter (| or -), the beam is
				// split into two beams going in each of the two directions the splitter's
				// pointy ends are pointing. For instance, a rightward-moving beam that
				// encounters a | splitter would split into two beams: one that continues
				// upward from the splitter's column and one that continues downward from the
				// splitter's column.
				beamQueue = append(beamQueue, getNextPositionByParams(r, c, UP))
				beamQueue = append(beamQueue, getNextPositionByParams(r, c, DOWN))
			} else {
				fmt.Println("simulateBeams encountered unexpected direction=", direction, "r=", r, "c=", c)
				panic("simulateBeams encountered unexpected direction")
			}
		} else if encounter == '-' {
			if direction == LEFT || direction == RIGHT {
				// If the beam encounters the pointy end of a splitter (| or -), the beam
				// passes through the splitter as if the splitter were empty space. For
				// instance, a rightward-moving beam that encounters a - splitter would
				// continue in the same direction.
				beamQueue = append(beamQueue, getNextPosition(currentBeam))
			} else if direction == UP || direction == DOWN {
				// If the beam encounters the flat side of a splitter (| or -), the beam is
				// split into two beams going in each of the two directions the splitter's
				// pointy ends are pointing. For instance, a rightward-moving beam that
				// encounters a | splitter would split into two beams: one that continues
				// upward from the splitter's column and one that continues downward from the
				// splitter's column.
				beamQueue = append(beamQueue, getNextPositionByParams(r, c, LEFT))
				beamQueue = append(beamQueue, getNextPositionByParams(r, c, RIGHT))
			} else {
				fmt.Println("simulateBeams encountered unexpected direction=", direction, "r=", r, "c=", c)
				panic("simulateBeams encountered unexpected direction")
			}
		} else {
			fmt.Println("simulateBeams unexecpted encounter=", encounter, "r=", r, "c=", c)
			panic("simulateBeams unexpected encounter")
		}
	}

	return len(engergized)
}

func simulateExampleBeam(mirrorMap [][]byte) int {
	return simulateBeams(mirrorMap, [3]int{0, 0, RIGHT})
}

func maximizeBeamEnergy(mirrorMap [][]byte) int {
	var max = 0

	// enter from all the tops and bottoms
	for c, _ := range mirrorMap[0] {
		max = Max(max, simulateBeams(mirrorMap, [3]int{0, c, DOWN}))
		max = Max(max, simulateBeams(mirrorMap, [3]int{len(mirrorMap) - 1, c, UP}))
	}

	// enter from all the lefts and right
	for r, _ := range mirrorMap {
		max = Max(max, simulateBeams(mirrorMap, [3]int{r, 0, RIGHT}))
		max = Max(max, simulateBeams(mirrorMap, [3]int{r, len(mirrorMap[r]) - 1, LEFT}))
	}

	return max
}

func Day16() {

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
	mirrorMap := parseMirrorMap(text)
	if DEBUG {
		fmt.Println(mirrorMap)
	}
	fmt.Println("Part 1:")
	fmt.Println(simulateExampleBeam(mirrorMap))
	fmt.Println("Part 2:")
	fmt.Println(maximizeBeamEnergy(mirrorMap))
}
