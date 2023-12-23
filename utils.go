package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const DEBUG = false

const UP = 0
const RIGHT = 1
const DOWN = 2
const LEFT = 3

func ParseIntOrExit(number string) int {
	steps, err := strconv.Atoi(number)
	if err != nil {
		fmt.Println("Error parsing number:", number, err)
		os.Exit(1)
	}
	return steps
}

func ReadFileOrExit(path string) string {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", path, err)
		os.Exit(1)
	}
	defer file.Close()
	reader := io.Reader(file)
	textBytes, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println("Error opening file:", path, err)
		os.Exit(1)
	}
	return string(textBytes)
}

func FindChar(char byte, matrix [][]byte) (int, int, bool) {
	for r, row := range matrix {
		for c, col := range row {
			if col == char {
				return r, c, true
			}
		}
	}
	return -1, -1, false
}

func Contains(arr []int, num int) bool {
	for _, value := range arr {
		if value == num {
			return true
		}
	}
	return false
}

func ManhattanDistance(a [2]int, b [2]int) int {
	return abs(a[0]-b[0]) + abs(a[1]-b[1])
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func ParseAs2DMatrix(input string) [][]byte {
	lines := strings.Split(input, "\n")
	result := make([][]byte, 0)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			result = append(result, []byte(trimmed))
		}
	}

	return result
}
