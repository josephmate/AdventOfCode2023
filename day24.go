package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type HailRecord struct {
	Position [3]int64
	Velocity [3]int64
}

func parseHailRecords(input string) []HailRecord {
	lines := strings.Split(input, "\n")
	records := make([]HailRecord, len(lines))

	for i, line := range lines {
		fields := strings.Split(line, "@")
		posFields := strings.Split(strings.TrimSpace(fields[0]), ",")
		velFields := strings.Split(strings.TrimSpace(fields[1]), ",")

		record := HailRecord{}

		for j := 0; j < 3; j++ {
			posVal, _ := strconv.ParseInt(strings.TrimSpace(posFields[j]), 10, 64)
			velVal, _ := strconv.ParseInt(strings.TrimSpace(velFields[j]), 10, 64)

			record.Position[j] = posVal
			record.Velocity[j] = velVal
		}

		records[i] = record
	}

	return records
}

/*
x(t) = x0 + xv * t
y(t) = y0 + yv * t

Where
t is time
x(t) is the x location at time t
x0 is the x location at t=0
xv is the x velocity
y(t) is the y location at time t
y0 is the y location at t=0
yv is the x velocity

insert a and b for the two lines

when x's intersect
xa0 + xav * t = xb0 + xbv * t
(xav - xbv)*t = xb0 - xa0
t = xb0 - xa0

	    ---------
			xav - xbv

when y's intersect
ya0 + yav + t = yb0 + ybv * t
t = yb0 - ya0

	    ---------
			yav - ybv
*/
func calculateTime(a HailRecord, b HailRecord, index int) float64 {
	return float64(b.Position[index]-a.Position[index]) /
		float64(a.Velocity[index]-b.Velocity[index])
}

const DELTA float64 = 0.000001

func calcPosn(a HailRecord, time float64, index int) float64 {
	return float64(a.Position[index]) + float64(a.Velocity[index])*time
}

func absFloat64(a float64) float64 {
	if a < 0 {
		return a * -1.0
	}
	return a
}

func has2DCollision(a HailRecord, b HailRecord, lowerBound int64, upperBound int64) bool {
	timeUsingXs := calculateTime(a, b, 0)
	timeUsingYs := calculateTime(a, b, 1)
	x := calcPosn(a, timeUsingXs, 0)
	y := calcPosn(a, timeUsingXs, 1)

	if DEBUG {
		fmt.Println("has2DCollision", "a=", a, "b=", b, "timeUsingXs=", timeUsingXs, "timeUsingYs=", timeUsingYs, "x=", x, "y=", y)
	}

	return timeUsingXs >= 0 &&
		timeUsingYs >= 0 &&
		absFloat64(timeUsingXs-timeUsingYs) <= DELTA &&
		x >= float64(lowerBound) && x <= float64(upperBound) &&
		y >= float64(lowerBound) && y <= float64(upperBound)
}

func count2DCollisions(hailRecords []HailRecord, lowerBound int64, upperBound int64) int64 {
	var count int64 = 0

	for i := 0; i < len(hailRecords); i++ {
		for j := i + 1; j < len(hailRecords); j++ {
			if has2DCollision(hailRecords[i], hailRecords[j], lowerBound, upperBound) {
				count++
			}
		}
	}

	return count
}

func Day24() {

	if len(os.Args) < 5 {
		fmt.Println("Usage: aoc 24 lowerbound upperbound <input>")
		os.Exit(1)
	}

	lowerbound := ParseInt64OrExit(os.Args[2])
	upperbound := ParseInt64OrExit(os.Args[3])
	text := ReadFileOrExit(os.Args[4])

	fmt.Println("Part 1:")
	if DEBUG {
		fmt.Println(text)
		fmt.Println(lowerbound)
		fmt.Println(upperbound)
	}
	hailRecords := parseHailRecords(text)
	if DEBUG {
		fmt.Println(hailRecords)
		fmt.Println(lowerbound)
		fmt.Println(upperbound)
	}
	fmt.Println(count2DCollisions(hailRecords, lowerbound, upperbound))

	fmt.Println("Part 2:")
}
