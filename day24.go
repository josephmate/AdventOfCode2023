package main

import (
	"fmt"
	"math"
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
ya0 + yav * t = yb0 + ybv * t
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
	y := calcPosn(a, timeUsingYs, 1)

	if DEBUG {
		fmt.Println("has2DCollision", "a=", a, "b=", b, "timeUsingXs=", timeUsingXs, "timeUsingYs=", timeUsingYs, "x=", x, "y=", y)
	}

	return timeUsingXs >= 0 &&
		timeUsingYs >= 0 &&
		absFloat64(timeUsingXs-timeUsingYs) <= DELTA &&
		x >= float64(lowerBound) && x <= float64(upperBound) &&
		y >= float64(lowerBound) && y <= float64(upperBound)
}

/*
Need to find the collision of paths/lines, not objects. Ooooops!

equation of line is

y = mx + b

where y is the y coord given some value on the x axis

m is the slope of the line (rise over run)

b is the y intercept (when x = 0)

m is vy/vx (rise over run)

where

vy is the y velocity

vx is the x velocity

b = y - mx

b = py - (vy/vx)px

where

py is the initial y

px is the initial x

putting that altogether is:

y = (vy/vx) * x + (py - (vy/vx)px)

now we have two equations of lines and we want to know the _x_ value where the ys are equal

(v0y/v0x) * _x_ + (p0y - (v0y/v0x)p0x) = (v1y/v1x) * _x_ + (p1y - (v1y/v1x)p1x)

solve for x

(v0y/v0x - v1y/v1x)_x_ = (p1y - (v1y/v1x)p1x) - (p0y - (v0y/v0x)p0x)

_x_ = (p1y - (v1y/v1x)p1x) - (p0y - (v0y/v0x)p0x)

	      -------------------------------------------
				(v0y/v0x - v1y/v1x)

we can use the y=mx+b to find the y value

just need to check both this x and y values are within the bounding rectangle

we also need to check that the time is not negative

recall from an early attempt the relationship between x and time (x(t) and t)

x(t) = x0 + xv * t

solving for t

t = x(t) - x0

	    ---------
			   xv

need to check if that is positive
*/
func has2DPathCollision(a HailRecord, b HailRecord, lowerBound int64, upperBound int64) bool {
	/*
		_x_ = (p1y - (v1y/v1x)p1x) - (p0y - (v0y/v0x)p0x)

		      -------------------------------------------
					(v0y/v0x - v1y/v1x)
	*/
	x := ((float64(b.Position[1]) - float64(b.Velocity[1])/float64(b.Velocity[0])*float64(b.Position[0])) -
		float64(a.Position[1]) - float64(a.Velocity[1])/float64(a.Velocity[0])*float64(a.Position[0])) /
		(float64(a.Velocity[1])/float64(a.Velocity[0]) - float64(b.Velocity[1])/float64(b.Velocity[0]))

	/*
		y = (vy/vx) * x + (py - (vy/vx)px)
	*/
	y := float64(a.Velocity[1])/float64(a.Velocity[0])*x + (float64(a.Position[1]) - float64(a.Velocity[1])/float64(a.Velocity[0])*float64(a.Position[0]))

	/*
		t = x(t) - x0

					---------
						xv
	*/
	t := (x - float64(a.Position[0])) / float64(a.Velocity[0])

	if DEBUG {
		fmt.Println("has2DPathCollision", "a:", a, "b:", b, "x:", x, "y:", y, "t:", t)
	}

	if math.IsInf(x, 1) || math.IsInf(x, -1) {
		return false
	}

	if !(x >= float64(lowerBound) && x <= float64(upperBound)) ||
		!(y >= float64(lowerBound) && y <= float64(upperBound)) {
		// out of bounds
		return false
	}

	return t >= 0
}

func count2DCollisions(hailRecords []HailRecord, lowerBound int64, upperBound int64) int64 {
	var count int64 = 0

	for i := 0; i < len(hailRecords); i++ {
		for j := i + 1; j < len(hailRecords); j++ {
			if has2DPathCollision(hailRecords[i], hailRecords[j], lowerBound, upperBound) {
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
