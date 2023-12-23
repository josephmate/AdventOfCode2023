package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type CommModRecord struct {
	IsFlipFlop    bool
	IsConjunction bool
	Name          string
	Destinations  []string
}

func parseCommModRecords(input string) []CommModRecord {
	lines := strings.Split(input, "\n")
	records := make([]CommModRecord, 0)

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}
		if DEBUG {
			fmt.Println("parseCommModRecords parsing line", trimmedLine)
		}
		parts := strings.Split(trimmedLine, " -> ")

		name := parts[0]
		destinations := strings.Split(parts[1], ", ")

		var isFlipFlop, isConjunction bool
		if strings.HasPrefix(name, "%") {
			isFlipFlop = true
			name = strings.TrimPrefix(name, "%")
		} else if strings.HasPrefix(name, "&") {
			isConjunction = true
			name = strings.TrimPrefix(name, "&")
		}

		record := CommModRecord{
			IsFlipFlop:    isFlipFlop,
			IsConjunction: isConjunction,
			Name:          name,
			Destinations:  destinations,
		}

		records = append(records, record)
	}

	return records
}

type Signal struct {
	Signalee string
	Signaler string
	IsHigh   bool
}

type CommMod struct {
	IsFlipFlop    bool
	IsFlipFlopOn  bool
	IsConjunction bool
	Name          string
	Outputs       []string
	Inputs        map[string]bool
}

func broadcast(queue *[]Signal, outputSignal bool, signaler string, outputs []string) {
	for _, output := range outputs {
		if DEBUG {
			var lowHigh = "low"
			if outputSignal {
				lowHigh = "high"
			}
			fmt.Println("simulateCommModStep", signaler, "--", lowHigh, "-->", output)
		}
		*queue = append(*queue, Signal{
			Signaler: signaler,
			Signalee: output,
			IsHigh:   outputSignal,
		})
	}
}

func simulateCommModStep(lookupMap *map[string]*CommMod) (int, int) {
	var lowCount = 0
	var highCount = 0
	var queue []Signal
	queue = append(queue, Signal{
		Signalee: "broadcaster",
		Signaler: "button",
		IsHigh:   false,
	})

	for len(queue) > 0 {
		currentSignal := queue[0]
		queue = queue[1:]

		if currentSignal.IsHigh {
			highCount++
		} else {
			lowCount++
		}

		commMod, hasIt := (*lookupMap)[currentSignal.Signalee]
		if !hasIt {
			continue
		}

		if commMod.IsFlipFlop {
			if !currentSignal.IsHigh {
				commMod.IsFlipFlopOn = !commMod.IsFlipFlopOn
				outputSignal := commMod.IsFlipFlopOn
				broadcast(&queue, outputSignal, currentSignal.Signalee, commMod.Outputs)
			}
		} else if commMod.IsConjunction {
			// modify the incoming signal
			commMod.Inputs[currentSignal.Signaler] = currentSignal.IsHigh
			// check if all low
			var allTrue = true
			for key := range commMod.Inputs {
				allTrue = allTrue && commMod.Inputs[key]
			}
			broadcast(&queue, !allTrue, currentSignal.Signalee, commMod.Outputs)
		} else { // broadcast
			broadcast(&queue, false, currentSignal.Signalee, commMod.Outputs)
		}
	}

	return lowCount, highCount
}

func simulate1000(commModRecords []CommModRecord) int {
	// preprocessing for efficient lookup
	lookUpMap := map[string]*CommMod{}
	for _, commModRec := range commModRecords {
		lookUpMap[commModRec.Name] = &CommMod{
			IsFlipFlop:    commModRec.IsFlipFlop,
			IsConjunction: commModRec.IsConjunction,
			Name:          commModRec.Name,
			Inputs:        map[string]bool{},
			Outputs:       commModRec.Destinations,
		}
	}
	if DEBUG {
		fmt.Println("simulate1000", "lookUpMap", lookUpMap)
	}

	// fill in Outpus
	for _, commModRec := range commModRecords {
		commMod := lookUpMap[commModRec.Name]
		for _, output := range commModRec.Destinations {
			outputCommMod, hasIt := lookUpMap[output]
			if hasIt {
				if outputCommMod.IsConjunction {
					outputCommMod.Inputs[commMod.Name] = false
				}
			}
		}
	}

	var lowPulses = 0
	var highPulses = 0
	for i := 0; i < 1000; i++ {
		lows, highs := simulateCommModStep(&lookUpMap)
		lowPulses += lows
		highPulses += highs
		fmt.Println("simulate1000", "step", i, "lowPulses", lowPulses, "highPulses", highPulses)
	}

	return lowPulses * highPulses
}

func Day20() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: aoc 20 <input>")
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
	commModRecords := parseCommModRecords(text)
	if DEBUG {
		fmt.Println(commModRecords)
	}
	fmt.Println("Part 1:")
	fmt.Println(simulate1000(commModRecords))
	fmt.Println("Part 2:")
}
