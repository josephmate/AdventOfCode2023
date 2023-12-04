package main

var DEBUG = true

func Contains(arr []int, num int) bool {
	for _, value := range arr {
		if value == num {
			return true
		}
	}
	return false
}
