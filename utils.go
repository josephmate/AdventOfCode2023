package main

const DEBUG = true

const UP = 0
const RIGHT = 1
const DOWN = 2
const LEFT = 3

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
