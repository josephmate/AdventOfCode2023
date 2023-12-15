package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func parseInitSeq(input string) []string {
	var result []string
	for _, column := range strings.Split(input, ",") {
		trimmed := strings.TrimSpace(column)
		if column != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func hashStep(step string) int {
	ascii := []byte(step)

	var currentHash = 0
	for _, num := range ascii {
		currentHash += int(num)
		currentHash *= 17
		currentHash = currentHash % 256
	}

	return currentHash
}

func hashInitSeq(initSeq []string) int {
	var sum = 0
	for _, step := range initSeq {
		sum += hashStep(step)
	}

	return sum
}

type HashMapEntry struct {
	Key   string
	Value int
}

func addHashMapEntry(box []HashMapEntry, key string, value int) []HashMapEntry {
	// try to replace first
	for idx, entry := range box {
		if entry.Key == key {
			if DEBUG {
				fmt.Println("addHashMapEntry replaced key ", key, ". had value ", entry.Value, " now has ", value)
			}
			box[idx] = HashMapEntry{
				Key:   key,
				Value: value,
			}
			return box
		}
	}

	// did not find an existing entry. make a new one
	return append(box, HashMapEntry{
		Key:   key,
		Value: value,
	})
}

func removeHashMapEntry(box []HashMapEntry, key string) []HashMapEntry {
	var newBox []HashMapEntry
	for _, entry := range box {
		if entry.Key != key {
			newBox = append(newBox, entry)
		}
	}
	return newBox
}

func hashmapInitSeq(initSeq []string) int {

	var hashmap [256][]HashMapEntry

	for lineNum, step := range initSeq {
		if strings.Contains(step, "=") {
			key := step[0:strings.Index(step, "=")]
			hashVal := hashStep(key)
			value := step[len(step)-1]
			if value >= '1' && value <= '9' {
				hashmap[hashVal] = addHashMapEntry(hashmap[hashVal], key, int(value-'0'))
			} else {
				fmt.Println("hashmapInitSeq: invalid = value ", string(value), " on line ", lineNum+1)
			}
		} else if strings.Contains(step, "-") {
			key := step[0:strings.Index(step, "-")]
			hashVal := hashStep(key)
			hashmap[hashVal] = removeHashMapEntry(hashmap[hashVal], key)
		} else {
			fmt.Println("hashmapInitSeq: invalid operation ", step, " on line ", lineNum+1)
		}
		if DEBUG {
			fmt.Println(hashmap)
		}
	}

	var sum = 0
	for boxIdx, box := range hashmap {
		for entryIdx, entry := range box {
			entryCalc := (boxIdx + 1) * (entryIdx + 1) * entry.Value
			if DEBUG {
				// rn: 1 (box 0) * 1 (first slot) * 1 (focal length) = 1
				fmt.Println("hashmapInitSeq", "calc", entry.Key, ":", (boxIdx + 1), "(box", boxIdx, ") *", (entryIdx + 1), "(", (entryIdx + 1), "slot) *", entry.Value, "(focal length) = ", entryCalc)
			}
			sum += entryCalc
		}
	}
	return sum
}

func Day15() {

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
	initSeq := parseInitSeq(text)
	if DEBUG {
		fmt.Println(initSeq)
	}
	fmt.Println("Part 1:")
	fmt.Println(hashInitSeq(initSeq))
	fmt.Println("Part 2:")
	fmt.Println(hashmapInitSeq(initSeq))
}
